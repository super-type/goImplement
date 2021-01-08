package goimplement

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
@param userKey the user's unique AES encryption key
*/
func Produce(data string, attribute string, supertypeID string, apiKey string, userKey string) (*ObservationRequest, error) {
	// Encrypt data using basic AES encryption
	ciphertext, iv, err := Encrypt(data, userKey)
	if err != nil {
		return nil, err
	}

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  *ciphertext,
		SupertypeID: supertypeID,
		IV:          *iv,
	}

	// Produce (upload) data to DynamoDB
	requestBody, err := json.Marshal(obs)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://supertype.io/produce", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	return &obs, nil
}

/*
Consume receives data from the Supertype data network, re-encrypts, and decrypts it
This data is source-agnostic, and encrypted end-to-end
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param apiKey the vendor's secret key
@param userKey the user's unique AES encryption key

@return plaintext the decrypted observation the vendor is requesting
*/
func Consume(attribute string, supertypeID string, apiKey string, userKey string) (*ConsumeResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
	})
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	// req, err := http.NewRequest("POST", "https://supertype.io/consume", bytes.NewReader(requestBody))
	req, err := http.NewRequest("POST", "http://localhost:5000/consume", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var observation ObservationResponse
	json.Unmarshal(body, &observation)
	plaintext, _, err := Decrypt(observation.Ciphertext, userKey)

	response := ConsumeResponse{
		Plaintext:   *plaintext,
		DateAdded:   observation.DateAdded,
		PublicKey:   observation.PublicKey,
		SupertypeID: observation.SupertypeID,
	}

	return &response, nil
}

/*
ValidateWebhookRequest validates incoming Webhook requests to ensure they're coming from Supertype
@param apiKey your secret API Key
@param supertypeSignature the signature Supertype sends as header X-Supertype-Signature

@return error if invalid signature, else nil
*/
func ValidateWebhookRequest(apiKey string, supertypeSignature string) error {
	h := sha256.New()
	h.Write([]byte(apiKey))
	apiKeyHash := hex.EncodeToString(h.Sum(nil))

	if apiKeyHash != supertypeSignature {
		return errors.New("Signatures do not match. Request does not come from Supertype")
	}

	return nil
}
