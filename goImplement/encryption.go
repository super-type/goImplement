package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
)

func gcmEncrypt(plaintext []byte, key string, iv []byte, additionalData []byte) (cipherText []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	cipherText = aesgcm.Seal(nil, iv, plaintext, additionalData)
	return cipherText, nil
}

func Encrypt(message string, pubKey *ecdsa.PublicKey) (cipherText []byte, capsule *Capsule, err error) {
	capsule, keyBytes, err := EncryptKeyGen(pubKey)
	if err != nil {
		return nil, nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	cipherText, err = gcmEncrypt([]byte(message), key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, nil, err
	}
	return cipherText, capsule, nil
}

func EncryptKeyGen(pubKey *ecdsa.PublicKey) (capsule *Capsule, keyBytes []byte, err error) {
	s := new(big.Int)
	// generate E,V key-pairs
	priE, pubE, err := GenerateKeys()
	priV, pubV, err := GenerateKeys()
	if err != nil {
		return nil, nil, err
	}
	// get H2(E || V)
	h := HashToCurve(
		ConcatBytes(
			PointToBytes(pubE),
			PointToBytes(pubV)))
	// get s = v + e * H2(E || V)
	s = BigIntAdd(priV.D, BigIntMul(priE.D, h))
	// get (pk_A)^{e+v}
	point := PointScalarMul(pubKey, BigIntAdd(priE.D, priV.D))
	// generate aes key
	keyBytes, err = Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, nil, err
	}
	capsule = &Capsule{
		E: pubE,
		V: pubV,
		S: s,
	}
	return capsule, keyBytes, nil
}
