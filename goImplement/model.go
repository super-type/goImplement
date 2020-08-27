package goimplement

import (
	"crypto/ecdsa"
	"math/big"
)

// Capsule is essentially the re-encryption key
type Capsule struct {
	E *ecdsa.PublicKey
	V *ecdsa.PublicKey
	S *big.Int
}

// ObservationRequest is an incoming observation from a vendor
type ObservationRequest struct {
	Attribute   string `json:"attribute"`
	Ciphertext  string `json:"ciphertext"`
	CapsuleE    string `json:"capsuleE"`
	CapsuleV    string `json:"capsuleV"`
	CapsuleS    string `json:"capsuleS"`
	SupertypeID string `json:"supertypeID"`
	PublicKey   string `json:"pk"`
	SkHash      string `json:"skHash"`
}

// ObservationResponse is returned from the server
type ObservationResponse struct {
	Ciphertext string `json:"ciphertext"`
	// Capsule              string    `json:"capsule"`
	CapsuleE             string    `json:"capsuleE"`
	CapsuleV             string    `json:"capsuleV"`
	CapsuleS             string    `json:"capsuleS"`
	DateAdded            string    `json:"dateAdded"`
	PublicKey            string    `json:"pk"`
	SupertypeID          string    `json:"supertypeID"`
	ReencryptionMetadata [2]string `json:"reencryptionMetadata"`
}

// Observation is a decrypted observation containing plaintext
type Observation struct {
	DateAdded string `json:"dateAdded"`
	PublicKey string `json:"pk"`
	Plaintext string `json:"plaintext"`
}

// MetadataResponse is the vendor metadata used to sync vendors on data production
type MetadataResponse struct {
	VendorConnections []string `json:"connections"`
	Vendors           []string `json:"vendors"`
}

// ReencryptionKeysRequest is sent when adding new re-encryption keys to a pre-existing vendor on produce
type ReencryptionKeysRequest struct {
	Connections map[string][2]string `json:"connections"`
	PublicKey   string               `json:"pk"`
}
