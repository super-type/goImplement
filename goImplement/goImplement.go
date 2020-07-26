package goimplement

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

/*
Produce produces data to the Supertype data marketplace
You need only encrypt once to send data anywhere within the ecosystem
@param data the message to encrypt
@param attribute the attribute to produce the data to
@param supertypeId the vendor's Supertype ID
@param skVendor the vendor's secret key
@param pkVendor the vendor's public key
*/
// TODO implement signing
func Produce(data string, attribute string, supertypeID string, skVendor string, pkVendor string) error {
	// Get public and private keys in usable form
	pk, err := StringToPublicKey(&pkVendor)
	if err != nil {
		return ErrStringToPublicKey
	}

	// Encrypt data
	cipherText, capsule, err := Encrypt(data, pk)
	if err != nil {
		return ErrEncryptingData
	}

	capsuleE := PublicKeyToString(capsule.E)
	capsuleV := PublicKeyToString(capsule.V)
	capsuleS := capsule.S.String()

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  hex.EncodeToString(cipherText),
		CapsuleE:    capsuleE,
		CapsuleV:    capsuleV,
		CapsuleS:    capsuleS,
		SupertypeID: supertypeID,
		PublicKey:   pkVendor,
	}

	// Upload data to DynamoDB
	requestBody, err := json.Marshal(obs)
	if err != nil {
		return ErrMarshaling
	}

	_, err = http.Post("http://localhost:8080/produce", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ErrHTML
	}

	// Get connections metadata in order to create any necessary re-encryption keys
	resp, err := http.Post("http://localhost:8080/getVendorComparisonMetadata", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ErrHTML
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrIORead
	}

	var metadata MetadataResponse
	json.Unmarshal(body, &metadata)

	connections := metadata.VendorConnections
	vendors := metadata.Vendors

	sk, err := StringToPrivateKey(&skVendor, *pk)
	if err != nil {
		fmt.Printf("error!!!: %v\n", err)
		return ErrStringToPrivateKey
	}

	newConnections := make(map[string][2]string)
	// Check if our vendor is up-to-date
	for _, vendor := range vendors {
		if !Contains(connections, vendor) && vendor != pkVendor {
			// Create re-encryption keys for current vendor to this vendor, and add it to this our current vendor
			rekeys, err := CreateReencryptionKeys(vendor, sk)
			if err != nil {
				return ErrReencrypting
			}

			// Add this re-encryption key-set to our newConnections map
			newConnections[vendor] = [2]string{rekeys[0], rekeys[1]}
		}
	}

	connectionRequest := ReencryptionKeysRequest{
		PublicKey:   pkVendor,
		Connections: newConnections,
	}

	// Marshal new request body
	requestBody, err = json.Marshal(connectionRequest)
	if err != nil {
		return ErrMarshaling
	}

	// Add these new connections for the current vendor
	resp, err = http.Post("http://localhost:8080/addReencryptionKeys", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ErrHTML
	}

	// TODO maybe /produce will return a list of vendors, as well as this vendor's connections list as its ObservationResponse
	// TODO maybe let's call a new endpoint to get this information... we don't want to slow down production, and we can figure out a way to do this in the background, asynchronously from typical Produce functionality...
	// 	1. Go through each vendor pk
	// 	2. If the vendor pk isn't in this vendor's connections list, create re-encryption keys for it
	// 	3. Update the vendor table through a new API endpoint

	return nil
}

/*
Consume receives data from the Supertype data marketplace and decrypt it
This data is source-agnostic, and encrypted end-to-end
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param skVendor the vendor's secret key
@param pkVendor the vendor's public key
*/
func Consume(attribute string, supertypeID string, skVendor string, pkVendor string) (*[]Observation, error) {
	// Get data from server
	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
		"pk":          pkVendor,
	})
	if err != nil {
		return nil, ErrMarshaling
	}

	// Get public and private keys in usable form
	pk, err := StringToPublicKey(&pkVendor)
	if err != nil {
		return nil, ErrStringToPublicKey
	}

	sk, err := StringToPrivateKey(&skVendor, *pk)
	if err != nil {
		return nil, ErrStringToPrivateKey
	}

	resp, err := http.Post("http://localhost:8080/consume", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, ErrHTML
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrIORead
	}

	var observations []ObservationResponse
	json.Unmarshal(body, &observations)

	var result []Observation

	// Iterate through each observation
	for _, obs := range observations {
		capsuleE, err := StringToPublicKey(&obs.CapsuleE)
		if err != nil {
			return nil, ErrStringToPublicKey
		}

		capsuleV, err := StringToPublicKey(&obs.CapsuleV)
		if err != nil {
			return nil, ErrStringToPublicKey
		}

		capsuleS := new(big.Int)
		capsuleS, ok := capsuleS.SetString(obs.CapsuleS, 10)
		if !ok {
			return nil, ErrStringToPublicKey
		}

		decodedCapsule := Capsule{
			E: capsuleE,
			V: capsuleV,
			S: capsuleS,
		}

		ciphertextAsBytes, err := hex.DecodeString(obs.Ciphertext)
		if err != nil {
			return nil, ErrDecoding
		}

		rekey := new(big.Int)
		rekey, ok = rekey.SetString(obs.ReencryptionMetadata[0], 10)
		if !ok {
			return nil, ErrStringToBigInt
		}

		pkX, err := StringToPublicKey(&(obs.ReencryptionMetadata[1]))
		if err != nil {
			return nil, ErrDecoding
		}

		newCapsule, err := ReEncrypt(rekey, &decodedCapsule)
		if err != nil {
			return nil, ErrReencrypting
		}

		plainText, err := Decrypt(sk, newCapsule, pkX, ciphertextAsBytes)
		if err != nil {
			return nil, ErrDecrypting
		}
		fmt.Printf("plaintext: %v\n", string(plainText))

		reObs := Observation{
			DateAdded: obs.DateAdded,
			PublicKey: obs.PublicKey,
			Plaintext: string(plainText),
		}
		result = append(result, reObs)
	}

	return &result, nil
}
