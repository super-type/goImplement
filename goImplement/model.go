package goimplement

// ObservationRequest is an incoming observation from a vendor
type ObservationRequest struct {
	Attribute   string `json:"attribute"`
	Ciphertext  string `json:"ciphertext"`
	SupertypeID string `json:"supertypeID"`
	PublicKey   string `json:"pk"`
	SkHash      string `json:"skHash"`
}

// ObservationResponse is returned from the server
type ObservationResponse struct {
	Ciphertext           string    `json:"ciphertext"`
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
