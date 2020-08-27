package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"encoding/hex"
)

// GCMDecrypt decrypts data using the AES GCM algorithm. aesLocal[:12] is nonce.
func GCMDecrypt(cipherText []byte, key string, iv []byte, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plainText, err := aesgcm.Open(nil, iv, cipherText, additionalData)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// RecreateAESKeyLocalDecrypt creates AES key bytes for local decrypt
func RecreateAESKeyLocalDecrypt(capsule *Capsule, aPriKey *ecdsa.PrivateKey) ([]byte, error) {
	point1 := PointScalarAdd(capsule.E, capsule.V)
	point := PointScalarMul(point1, aPriKey.D)
	aesKey, err := Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return aesKey, nil
}

// LocalDecrypt decrypts data if the data was originally produces by this vendor
func LocalDecrypt(aPriKey *ecdsa.PrivateKey, capsule *Capsule, cipherText []byte) ([]byte, error) {
	aesKeyBytes, err := RecreateAESKeyLocalDecrypt(capsule, aPriKey)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(aesKeyBytes)
	plainText, err := GCMDecrypt(cipherText, key[:32], aesKeyBytes[:12], nil)
	return plainText, err
}

func decryptKeyGen(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey) ([]byte, error) {
	S := PointScalarMul(pubX, bPriKey.D)
	d := HashToCurve(ConcatBytes(ConcatBytes(PointToBytes(pubX), PointToBytes(&bPriKey.PublicKey)), PointToBytes(S)))
	point := PointScalarMul(PointScalarAdd(capsule.E, capsule.V), d)
	keyBytes, err := Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// Decrypt recreates the AES key, then decrypts the encrypted data
func Decrypt(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey, cipherText []byte) ([]byte, error) {
	aesKeyBytes, err := decryptKeyGen(bPriKey, capsule, pubX)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(aesKeyBytes)
	plainText, err := GCMDecrypt(cipherText, key[:32], aesKeyBytes[:12], nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
