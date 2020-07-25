package goimplement

type ObservationRequest struct {
	Attribute  string `json:"attribute"`
	Ciphertext string `json:"ciphertext"`
	// Capsule     string `json:"capsule"`
	CapsuleE    string `json:"capsuleE"`
	CapsuleV    string `json:"capsuleV"`
	CapsuleS    string `json:"capsuleS"`
	SupertypeID string `json:"supertypeID"`
	PublicKey   string `json:"pk"`
}

// ObservationResponse is returned from the server
// todo should this be public?
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
