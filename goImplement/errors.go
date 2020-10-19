package goimplement

import "errors"

// ErrEncryptingData ... error encrypting data
var ErrEncryptingData = errors.New("Error encrypting data")

// ErrMarshaling ... error marshaling request
var ErrMarshaling = errors.New("Error marshaling request")

// ErrHTML ... error with HTML request
var ErrHTML = errors.New("Error with HTML request")

// ErrIORead ... error reading response body
var ErrIORead = errors.New("Error reading response body")

// ErrDecrypting ... error decrypting
var ErrDecrypting = errors.New("Error decrypting")
