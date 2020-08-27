package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
)

// GCMEncrypt uses the GCM algorithm to encrypt AES key
func GCMEncrypt(plaintext []byte, key string, iv []byte, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText := aesgcm.Seal(nil, iv, plaintext, additionalData)

	return cipherText, nil
}

// Encrypt creates the AES key, encrypts it it with GCM, and returns ciphertext and capsule
func Encrypt(message string, pk *ecdsa.PublicKey) (cipherText []byte, capsule *Capsule, err error) {
	capsule, aesKeyBytes, err := EncryptKeyGen(pk)
	if err != nil {
		return nil, nil, err
	}

	key := hex.EncodeToString(aesKeyBytes)
	cipherText, err = GCMEncrypt([]byte(message), key[:32], aesKeyBytes[:12], nil)
	if err != nil {
		return nil, nil, err
	}

	return cipherText, capsule, nil
}

// EncryptKeyGen generates the capsule and AES key
func EncryptKeyGen(pk *ecdsa.PublicKey) (*Capsule, []byte, error) {
	s := new(big.Int)
	skE, pkE, err := GenerateKeys()
	skV, pkV, err := GenerateKeys()
	if err != nil {
		return nil, nil, err
	}

	h := HashToCurve(
		ConcatBytes(
			PointToBytes(pkE),
			PointToBytes(pkV)))
	s = AddBigInteger(skV.D, MultiplyBigInteger(skV.D, h))
	point := PointScalarMul(pk, AddBigInteger(skE.D, skV.D))

	aesKeyBytes, err := Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, nil, err
	}

	capsule := &Capsule{
		E: pkE,
		V: pkV,
		S: s,
	}

	return capsule, aesKeyBytes, nil
}
