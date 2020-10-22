package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strings"
)

// Encrypt creates the AES key, encrypts it it with GCM, and returns ciphertext and capsule
// https://play.golang.org/p/4FQBAeHgRs
func Encrypt(message string, userKey string) (*string, *string, error) {
	userKeyBytes := []byte(userKey)
	messageBytes := []byte(message)
	block, err := aes.NewCipher(userKeyBytes[0:32])
	if err != nil {
		return nil, nil, err
	}

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}
	ivString := base64.StdEncoding.EncodeToString(iv)

	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(messageBytes))
	cfb.XORKeyStream(ciphertext, messageBytes)
	ciphertextString := base64.StdEncoding.EncodeToString(ciphertext)
	return &ciphertextString, &ivString, nil
}

// Decrypt recreates the AES key, then decrypts the encrypted data
// https://play.golang.org/p/4FQBAeHgRs
func Decrypt(ciphertext string, userKey string) (*string, *string, error) {
	// Split ciphertext from iv
	metadata := strings.Split(ciphertext, "|")
	ciphertext = metadata[0]
	attribute := metadata[2]
	iv, err := base64.StdEncoding.DecodeString(metadata[1])
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher([]byte(userKey[0:32]))
	if err != nil {
		return nil, nil, err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, []byte(iv))
	plaintext := make([]byte, len(data))
	cfb.XORKeyStream(plaintext, data)
	res := string(plaintext)
	return &res, &attribute, nil
}

// GetSecretKeyHash returns the hashed value of the secret key
func GetSecretKeyHash(skVendor string) string {
	h := sha256.New()
	h.Write([]byte(skVendor))
	return hex.EncodeToString(h.Sum(nil))
}
