package goimplement

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
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
		return ErrEncryptingData
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

	// TODO we are only using this when testing locally. We should only request the real API when publishing to GoDoc
	if len(os.Args) == 1 || os.Args[1] == "local" {
		_, err = http.Post("http://localhost:5000/produce", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
	} else if os.Args[1] == "test" {
		_, err = http.Post("https://supertype.io/produce", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Please use one of: local | test")
		return nil
	}

	// TODO we should probably return something here...
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
	fmt.Println(skHash)

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

	// TODO we are only using this when testing locally. We should only request the real API when publishing to GoDoc
	if len(os.Args) == 1 || os.Args[1] == "local" {
		resp, err = http.Post("http://localhost:5000/consume", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}
	} else if os.Args[1] == "test" {
		resp, err = http.Post("https://supertype.io/consume", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Please use one of: local | test")
		return &result, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrIORead
	}

	var observations []ObservationResponse
	json.Unmarshal(body, &observations)

	// Iterate through each observation
	for _, obs := range observations {
		plaintext, _, err := Decrypt(obs.Ciphertext, userKey)
		if err != nil {
			return nil, ErrDecrypting
		}
		result = append(result, *plaintext)
	}

	return &result, nil
}

/*
ConsumeWS subscribes this node to the specified attribute(s)
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param skVendor the vendor's secret key
@param pkVendor the vendor's public key
@param userKey the user's unique AES encryption key
*/
func ConsumeWS(attribute string, supertypeID string, skVendor string, pkVendor string, userKey string) error {
	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	// Establish WebSocket connection between device <-> server
	interrupt := make(chan os.Signal, 1)

	var addr *string

	// TODO we are only using this when testing locally. We should only request the real API when publishing to GoDoc
	if len(os.Args) == 1 || os.Args[1] == "local" {
		addr = flag.String("addr", "localhost:5001", "http service address")
	} else if os.Args[1] == "test" {
		addr = flag.String("addr", "supertype-demo.us-east-1.elasticbeanstalk.com:8082", "http service address")
	} else {
		fmt.Println("Please use one of: local | test")
		return nil
	}

	// var addr = flag.String("addr", "localhost:5001", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/subscribe"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			switch messageType {
			case 1:
				if string(message) == "Connected" {
					requestBody, err := json.Marshal(map[string]string{
						"attribute":   attribute,
						"supertypeID": supertypeID,
						"pk":          pkVendor,
						"skHash":      skHash,
						"cid":         string(message),
					})
					err = c.WriteMessage(2, requestBody)
					if err != nil {
						return
					}
				} else {
					log.Printf("subscribed to: %s", message)
				}
			case 2:
				var raw map[string]interface{}
				err = json.Unmarshal(message, &raw)

				rawMessage, ok := raw["body"].(string)
				if !ok {
					fmt.Println("Error getting raw body")
				}
				plaintext, attribute, err := Decrypt(string(rawMessage), userKey)
				if err != nil {
					fmt.Printf("ERROR: Error decrypting message: %v\n", err)
					return
				}
				fmt.Printf("Received %v: %v\n", *attribute, *plaintext)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return nil
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
