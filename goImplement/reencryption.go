package goimplement

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
)

// ReEncrypt re-encrypts the capsule into a new capsule, able to be decrypted with the rekey
func ReEncrypt(rk *big.Int, capsule *Capsule) (*Capsule, error) {
	// check g^s == V * E^{H2(E || V)}
	x1, y1 := CURVE.ScalarBaseMult(capsule.S.Bytes())
	tempX, tempY := CURVE.ScalarMult(capsule.E.X, capsule.E.Y, HashToCurve(ConcatBytes(PointToBytes(capsule.E), PointToBytes(capsule.V))).Bytes())
	x2, y2 := CURVE.Add(capsule.V.X, capsule.V.Y, tempX, tempY)
	// if check failed return error
	if x1.Cmp(x2) != 0 || y1.Cmp(y2) != 0 {
		return nil, fmt.Errorf("%s", "Capsule not match")
	}
	// E' = E^{rk}, V' = V^{rk}
	newCapsule := &Capsule{
		E: PointScalarMul(capsule.E, rk),
		V: PointScalarMul(capsule.V, rk),
		S: capsule.S,
	}
	return newCapsule, nil
}

// CreateReencryptionKeys creates re-encryption keys for public keys returned from DynamoDB
func CreateReencryptionKeys(vendor string, sk *ecdsa.PrivateKey) (*[2]string, error) {
	var connections [2]string

	pk, err := StringToPublicKey(&vendor)
	if err != nil {
		return nil, ErrStringToPublicKey
	}

	rekey, pkX, err := ReKeyGen(sk, pk)
	if err != nil {
		return nil, err
	}

	rekeyStr := rekey.String()
	pkXStr := PublicKeyToString(pkX)

	connections[0] = rekeyStr
	connections[1] = pkXStr

	return &connections, nil
}

// generate re-encryption key and sends it to Server
// rk = sk_A * d^{-1}
func ReKeyGen(aPriKey *ecdsa.PrivateKey, bPubKey *ecdsa.PublicKey) (*big.Int, *ecdsa.PublicKey, error) {
	// generate x,X key-pair
	priX, pubX, err := GenerateKeys()
	if err != nil {
		return nil, nil, err
	}
	// get d = H3(X_A || pk_B || pk_B^{x_A})
	point := PointScalarMul(bPubKey, priX.D)
	d := HashToCurve(
		ConcatBytes(
			ConcatBytes(
				PointToBytes(pubX),
				PointToBytes(bPubKey)),
			PointToBytes(point)))
	// rk = sk_A * d^{-1}
	rk := BigIntMul(aPriKey.D, GetInvert(d))
	rk.Mod(rk, N)
	return rk, pubX, nil
}
