package goimplement

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// StringToPrivateKey converts an encoded string into an ECDSA private key
func StringToPrivateKey(skString *string, pk ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	n := new(big.Int)
	n, ok := n.SetString(*skString, 10)
	if !ok {
		return nil, ErrStringToPrivateKey
	}

	sk := ecdsa.PrivateKey{
		PublicKey: pk,
		D:         n,
	}

	return &sk, nil
}

// StringToPublicKey converts an encoded string into an ECDSA public key
func StringToPublicKey(pkString *string) (*ecdsa.PublicKey, error) {
	pkTempBytes, err := hex.DecodeString(*pkString)
	if err != nil {
		return nil, ErrStringToPublicKey
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pkTempBytes)
	publicKeyFinal := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return &publicKeyFinal, nil
}

// PrivateKeyToString encodes an ECDSA private key to a string
func PrivateKeyToString(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

// PublicKeyToString encodes an ECDSA public key to a string
func PublicKeyToString(publicKey *ecdsa.PublicKey) string {
	pubKeyBytes := PointToBytes(publicKey)
	return hex.EncodeToString(pubKeyBytes)
}

// GenerateKeys generates a <private,public> key pair
func GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	sk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return sk, &sk.PublicKey, nil
}
