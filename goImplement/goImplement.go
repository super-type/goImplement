package goimplement

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
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
// TODO implement sharding capsule
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

	// Generate hash of secret key to be used as a signing measure for producing/consuming data
	skHash := GetSecretKeyHash(skVendor)

	obs := ObservationRequest{
		Attribute:   attribute,
		Ciphertext:  hex.EncodeToString(cipherText),
		CapsuleE:    capsuleE,
		CapsuleV:    capsuleV,
		CapsuleS:    capsuleS,
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
func Consume(attribute string, supertypeID string, skVendor string, pkVendor string) (*[]Observation, error) {
	// Get public and private keys in usable form
	pk, err := StringToPublicKey(&pkVendor)
	if err != nil {
		return nil, ErrStringToPublicKey
	}

	sk, err := StringToPrivateKey(&skVendor, *pk)
	if err != nil {
		return nil, ErrStringToPrivateKey
	}

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
		// TODO if this is nil, that means it's the vendor's own data, so we need a check & implementation for that
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
