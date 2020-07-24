package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
)

// Decrypt by my own private key
func DecryptOnMyPriKey(aPriKey *ecdsa.PrivateKey, capsule *Capsule, cipherText []byte) (plainText []byte, err error) {
	keyBytes, err := RecreateAESKeyByMyPriKey(capsule, aPriKey)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	plainText, err = GCMDecrypt(cipherText, key[:32], keyBytes[:12], nil)
	return plainText, err
}

// Recreate aes key
func RecreateAESKeyByMyPriKey(capsule *Capsule, aPriKey *ecdsa.PrivateKey) (keyBytes []byte, err error) {
	point1 := PointScalarAdd(capsule.E, capsule.V)
	point := PointScalarMul(point1, aPriKey.D)
	// generate aes key
	keyBytes, err = Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

func decryptKeyGen(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey) (keyBytes []byte, err error) {
	// S = X_A^{sk_B}
	S := PointScalarMul(pubX, bPriKey.D)
	// recreate d = H3(X_A || pk_B || S)
	d := HashToCurve(
		ConcatBytes(
			ConcatBytes(
				PointToBytes(pubX),
				PointToBytes(&bPriKey.PublicKey)),
			PointToBytes(S)))
	point := PointScalarMul(
		PointScalarAdd(capsule.E, capsule.V), d)
	keyBytes, err = Sha3Hash(PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// Recreate the aes key then decrypt the cipherText
func Decrypt(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey, cipherText []byte) (plainText []byte, err error) {
	keyBytes, err := decryptKeyGen(bPriKey, capsule, pubX)
	if err != nil {
		return nil, err
	}
	// recreate aes key = G((E' * V')^d)
	key := hex.EncodeToString(keyBytes)
	// use aes gcm to decrypt
	// mark keyBytes[:12] as nonce
	plainText, err = GCMDecrypt(cipherText, key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("pt: %v\n", plainText)
	return plainText, nil
}

func GCMDecrypt(cipherText []byte, key string, iv []byte, additionalData []byte) (plainText []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plainText, err = aesgcm.Open(nil, iv, cipherText, additionalData)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil, err
	}
	return plainText, nil
}
