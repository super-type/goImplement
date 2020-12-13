package goimplement

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
@param userKey the user's unique AES encryption key
*/
func Produce(data string, attribute string, supertypeID string, skVendor string, pkVendor string, userKey string) error {
	// Encrypt data using basic AES encryption
	ciphertext, iv, err := Encrypt(data, userKey)
	if err != nil {
		return err
	}

	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  *ciphertext,
		SupertypeID: supertypeID,
		PublicKey:   pkVendor,
		SkHash:      skHash,
		IV:          *iv,
	}

	// Produce (upload) data to DynamoDB
	requestBody, err := json.Marshal(obs)
	if err != nil {
		return err
	}

	_, err = http.Post("https://supertype.io/produce", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	return nil
}

/*
Consume receives data from the Supertype data network, re-encrypts, and decrypts it
This data is source-agnostic, and encrypted end-to-end
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param skVendor the vendor's secret key
@param pkVendor the vendor's public key
@param userKey the user's unique AES encryption key

@return plaintext the decrypted observation the vendor is requesting
*/
func Consume(attribute string, supertypeID string, skVendor string, pkVendor string, userKey string) (plaintext *[]string, err error) {
	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
		"pk":          pkVendor,
		"skHash":      skHash,
	})
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	var result []string

	resp, err = http.Post("https://supertype.io/consume", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var observations []ObservationResponse
	json.Unmarshal(body, &observations)

	// Iterate through each observation
	for _, obs := range observations {
		plaintext, _, err := Decrypt(obs.Ciphertext, userKey)
		if err != nil {
			return nil, err
		}
		result = append(result, *plaintext)
	}

	return &result, nil
}
