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

// ErrDecoding ... error decoding string to bytes
var ErrDecoding = errors.New("Error decoding string to byte array")

// ErrStringToBigInt ... error converting from string to biginteger
var ErrStringToBigInt = errors.New("Error converting from string to big integer")

// ErrReencrypting ... error re-encrypting
var ErrReencrypting = errors.New("Error re-encrypting")

// ErrDecrypting ... error decrypting
var ErrDecrypting = errors.New("Error decrypting")
