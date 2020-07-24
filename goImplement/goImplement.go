package goimplement

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
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
		return errors.New("Error converting to public key")
	}

	// TODO we will use this for re-encryption with new vendors
	// sk, err := stringToPrivateKey(&skVendor, *pk)
	// if err != nil {
	// 	return errors.New("Error converting to private key")
	// }

	// Encrypt data
	cipherText, capsule, err := Encrypt(data, pk)
	if err != nil {
		return errors.New("Error encrypt data")
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
		return errors.New("Error marshaling request")
	}

	_, err = http.Post("http://localhost:8080/produce", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("Error posting data")
	}

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
func Consume(attribute string, supertypeID string, skVendor string, pkVendor string) error {
	// Get data from server
	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
		"pk":          pkVendor,
	})
	if err != nil {
		return errors.New("Error marshaling request")
	}

	resp, err := http.Post("http://localhost:8080/consume", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("Error posting data")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error reading response")
	}

	var observations []ObservationResponse
	json.Unmarshal(body, &observations) // todo figure out a way to not use json.unmarshal in favor of decoder

	// Get public and private keys in usable form
	pk, err := StringToPublicKey(&pkVendor)
	if err != nil {
		return errors.New("Error converting to public key")
	}

	sk, err := StringToPrivateKey(&skVendor, *pk)
	if err != nil {
		return errors.New("Error converting to private key")
	}

	// Iterate through each observation
	for _, obs := range observations {
		capsuleE, err := StringToPublicKey(&obs.CapsuleE)
		if err != nil {
			return errors.New("Error decoding capsule")
		}

		capsuleV, err := StringToPublicKey(&obs.CapsuleV)
		if err != nil {
			return errors.New("Error decoding capsule")
		}

		capsuleS := new(big.Int)
		capsuleS, ok := capsuleS.SetString(obs.CapsuleS, 10)
		if !ok {
			return errors.New("Error decoding capsule")
		}

		decodedCapsule := Capsule{
			E: capsuleE,
			V: capsuleV,
			S: capsuleS,
		}

		// ciphertextAsBytes, err := base64.RawStdEncoding.DecodeString(obs.Ciphertext)
		ciphertextAsBytes, err := hex.DecodeString(obs.Ciphertext)
		if err != nil {
			return errors.New("Error decoding cipehrtext")
		}

		rekey := new(big.Int)
		rekey, ok = rekey.SetString(obs.ReencryptionMetadata[0], 10)
		if !ok {
			return errors.New("Error setting rekey")
		}

		fmt.Printf("rekey: %v\n", rekey)

		pkX, err := StringToPublicKey(&(obs.ReencryptionMetadata[1]))
		if err != nil {
			return errors.New("Error decoding pkX")
		}

		fmt.Printf("pkX string: %v\n", obs.ReencryptionMetadata[1])

		newCapsule, err := ReEncryption(rekey, &decodedCapsule)
		if err != nil {
			return errors.New("Error re-encrypting")
		}

		fmt.Printf("ciphertextAsBytes: %v\n", ciphertextAsBytes)
		plainText, err := Decrypt(sk, newCapsule, pkX, ciphertextAsBytes)
		if err != nil {
			fmt.Printf("Error decrypting... %v\n", err)
		}
		fmt.Printf("plaintext: %v\n", string(plainText))
	}

	return nil
}
