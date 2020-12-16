package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

// Encrypt creates the AES key, encrypts it it with GCM, and returns ciphertext and capsule
// https://play.golang.org/p/4FQBAeHgRs
func Encrypt(message string, userKey string) (*string, *string, error) {
	plaintext := []byte(message)
	userKeyBytes := []byte(userKey[0:32])

	block, err := aes.NewCipher(userKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext, plaintext)

	ciphertextString := base64.StdEncoding.EncodeToString(ciphertext)
	ivString := base64.StdEncoding.EncodeToString(iv)
	return &ciphertextString, &ivString, nil
}

// Decrypt recreates the AES key, then decrypts the encrypted data
// https://play.golang.org/p/4FQBAeHgRs
// https://gist.github.com/huyinghuan/e6d3046add5412bcd1129b074dc8b106
// https://gist.github.com/lihnux/2aa4a6f5a9170974f6aa
func Decrypt(ciphertext string, userKey string) (*string, *string, error) {
	keyBytes := []byte(userKey[0:32])
	metadata := strings.Split(ciphertext, "|")
	ciphertext = metadata[0]
	iv, err := base64.StdEncoding.DecodeString(metadata[1])
	if err != nil {
		return nil, nil, err
	}
	attribute := metadata[2]

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	stream.XORKeyStream(plaintext, ciphertextBytes)

	res := string(plaintext)

	return &res, &attribute, nil
}

// GetAPIKeyHash returns the hashed value of the secret key
func GetAPIKeyHash(skVendor string) string {
	h := sha256.New()
	h.Write([]byte(skVendor))
	return hex.EncodeToString(h.Sum(nil))
}
