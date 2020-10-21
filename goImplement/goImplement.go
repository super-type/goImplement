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
*/
func Produce(data string, attribute string, supertypeID string, skVendor string, pkVendor string) error {
	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  data, // TODO change this so that it's not just plaintext!!
		SupertypeID: supertypeID,
		PublicKey:   pkVendor,
		SkHash:      skHash,
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
*/
func Consume(attribute string, supertypeID string, skVendor string, pkVendor string) (*[]string, error) {
	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	requestBody, err := json.Marshal(map[string]string{
		"attribute":   attribute,
		"supertypeID": supertypeID,
		"pk":          pkVendor,
		"skHash":      skHash,
	})
	if err != nil {
		return nil, ErrMarshaling
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

	fmt.Printf("observations: %v\n", observations)

	var result []string

	// Iterate through each observation
	// TODO eventually we'll do AES decryption here... right now it's just plaintext
	for _, obs := range observations {
		result = append(result, obs.Ciphertext)
	}

	return &result, nil
}

/*
ConsumeWS subscribes this node to the specified attribute(s)
@param attribute to consume data from
@param supertypeID the vendor's Supertype ID
@param skVendor the vendor's secret key
@param pkVendor the vendor's public key
*/
func ConsumeWS(attribute string, supertypeID string, skVendor string, pkVendor string) error {
	// TODO first make a POST reqeust to the server to add attributes to Redis, then subscribe via WebSocket...
	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	// Establish WebSocket connection between device <-> server
	interrupt := make(chan os.Signal, 1)

	var addr = flag.String("addr", "localhost:8081", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/consume"}
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
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			// todo we should listen to something better than "Subscribed" - maybe write a specific message type
			if string(message) == "Subscribed" {
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
			log.Println("interrupt")

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
