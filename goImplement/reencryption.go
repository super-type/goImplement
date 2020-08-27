package goimplement

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
)

// ReEncrypt re-encrypts the capsule into a new capsule, able to be decrypted with the rekey
func ReEncrypt(rk *big.Int, capsule *Capsule) (*Capsule, error) {
	x1, y1 := CURVE.ScalarBaseMult(capsule.S.Bytes())
	tempX, tempY := CURVE.ScalarMult(capsule.E.X, capsule.E.Y, HashToCurve(ConcatBytes(PointToBytes(capsule.E), PointToBytes(capsule.V))).Bytes())
	x2, y2 := CURVE.Add(capsule.V.X, capsule.V.Y, tempX, tempY)
	if x1.Cmp(x2) != 0 || y1.Cmp(y2) != 0 {
		return nil, fmt.Errorf("%s", "Capsule does not match.")
	}

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

// ReKeyGen generates a re-encryption key and sends it to Server
func ReKeyGen(skA *ecdsa.PrivateKey, pkB *ecdsa.PublicKey) (*big.Int, *ecdsa.PublicKey, error) {
	priX, pubX, err := GenerateKeys()
	if err != nil {
		return nil, nil, err
	}

	point := PointScalarMul(pkB, priX.D)
	d := HashToCurve(ConcatBytes(ConcatBytes(PointToBytes(pubX), PointToBytes(pkB)), PointToBytes(point)))
	rk := MultiplyBigInteger(skA.D, GetInverse(d))
	rk.Mod(rk, N)

	return rk, pubX, nil
}
