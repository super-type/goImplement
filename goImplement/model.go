package goimplement

// ObservationRequest is an incoming observation from a vendor
type ObservationRequest struct {
	Attribute   string `json:"attribute"`
	Ciphertext  string `json:"ciphertext"`
	SupertypeID string `json:"supertypeID"`
	IV          string `json:"iv"`
}

// ObservationResponse is returned from the server
type ObservationResponse struct {
	Ciphertext  string `json:"ciphertext"`
	DateAdded   string `json:"dateAdded"`
	PublicKey   string `json:"pk"`
	SupertypeID string `json:"supertypeID"`
}

// ConsumeResponse is the value returned to user
type ConsumeResponse struct {
	Plaintext   string `json:"plaintext"`
	DateAdded   string `json:"dateAdded"`
	PublicKey   string `json:"pk"`
	SupertypeID string `json:"supertypeID"`
}
