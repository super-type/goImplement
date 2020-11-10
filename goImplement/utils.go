package goimplement

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
)

// Encrypt creates the AES key, encrypts it it with GCM, and returns ciphertext and capsule
// https://play.golang.org/p/4FQBAeHgRs
func Encrypt(message string, userKey string) (*string, *string, error) {
	plaintext := []byte(message)
	userKeyBytes := []byte(userKey[0:32])
	fmt.Printf("userKey: %v\n", userKey)
	fmt.Printf("userkeyBytes: %v\n", userKeyBytes)

	block, err := aes.NewCipher(userKeyBytes)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
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
		panic(err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	stream.XORKeyStream(plaintext, ciphertextBytes)

	res := string(plaintext)

	return &res, &attribute, nil
}

// GetSecretKeyHash returns the hashed value of the secret key
func GetSecretKeyHash(skVendor string) string {
	h := sha256.New()
	h.Write([]byte(skVendor))
	return hex.EncodeToString(h.Sum(nil))
}
