package goimplement

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
)

func StringToPrivateKey(skString *string, pk ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	n := new(big.Int)
	n, ok := n.SetString(*skString, 10)
	if !ok {
		return nil, errors.New("SetString error")
	}

	sk := ecdsa.PrivateKey{
		PublicKey: pk,
		D:         n,
	}

	return &sk, nil
}

func StringToPublicKey(pkString *string) (*ecdsa.PublicKey, error) {
	pkTempBytes, err := hex.DecodeString(*pkString)
	if err != nil {
		return nil, errors.New("Error decoding bytes from string")
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pkTempBytes)
	publicKeyFinal := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return &publicKeyFinal, nil
}

func GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	sk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return sk, &sk.PublicKey, nil
}

// convert private key to string
func PrivateKeyToString(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

// convert public key to string
func PublicKeyToString(publicKey *ecdsa.PublicKey) string {
	pubKeyBytes := PointToBytes(publicKey)
	return hex.EncodeToString(pubKeyBytes)
}
