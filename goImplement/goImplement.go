package goimplement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
Produce produces data to the Supertype data marketplace
You need only encrypt once to send data anywhere within the ecosystem
@param data the message to encrypt
@param attribute the attribute to produce the data to
@param supertypeId the vendor's Supertype ID
@param apiKey the vendor's secret key
@param pkVendor the vendor's public key
@param userKey the user's unique AES encryption key
*/
func Produce(data string, attribute string, supertypeID string, apiKey string, pkVendor string, userKey string) error {
	// Generate hash of API key to be used as a signing measure for producing/consuming data
	apiKeyHash := GetAPIKeyHash(apiKey)

	// Encrypt data using basic AES encryption
	ciphertext, iv, err := Encrypt(data, userKey)
	if err != nil {
		return err
	}

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  *ciphertext,
		SupertypeID: supertypeID,
		PublicKey:   pkVendor,
		IV:          *iv,
	}

	// Produce (upload) data to DynamoDB
	requestBody, err := json.Marshal(obs)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://supertype.io/produce", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKeyHash)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

/*
Consume receives data from the Supertype data network, re-encrypts, and decrypts it
This data is source-agnostic, and encrypted end-to-end
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param apiKey the vendor's secret key
@param pkVendor the vendor's public key
@param userKey the user's unique AES encryption key

@return plaintext the decrypted observation the vendor is requesting
*/
func Consume(attribute string, supertypeID string, apiKey string, pkVendor string, userKey string) (plaintext *[]string, err error) {
	// Generate hash of API key to be used as a signing measure for producing/consuming data
	apiKeyHash := GetAPIKeyHash(apiKey)

	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
		"pk":          pkVendor,
	})
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://supertype.io/consume", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKeyHash)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []string

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var observations []ObservationResponse
	json.Unmarshal(body, &observations)

	// Iterate through each observation
	for _, obs := range observations {
		fmt.Println("here")
		plaintext, _, err := Decrypt(obs.Ciphertext, userKey)
		if err != nil {
			return nil, err
		}
		result = append(result, *plaintext)
	}

	return &result, nil
}
