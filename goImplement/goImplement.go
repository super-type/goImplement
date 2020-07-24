package goimplement

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/crypto/sha3"
)

var CURVE = elliptic.P256()
var P = CURVE.Params().P
var N = CURVE.Params().N

type CurvePoint = ecdsa.PublicKey

func stringToPrivateKey(skString *string, pk ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	n := new(big.Int)
	n, ok := n.SetString(*skString, 10)
	if !ok {
		return nil, errors.New("SetString error")
	}

	sk := ecdsa.PrivateKey{
		PublicKey: pk,
		D:         n,
	}

	return &sk, nil
}

func stringToPublicKey(pkString *string) (*ecdsa.PublicKey, error) {
	pkTempBytes, err := hex.DecodeString(*pkString)
	if err != nil {
		return nil, errors.New("Error decoding bytes from string")
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pkTempBytes)
	publicKeyFinal := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return &publicKeyFinal, nil
}

func generateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	sk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return sk, &sk.PublicKey, nil
}

func pointScalarMul(a *CurvePoint, k *big.Int) *CurvePoint {
	x, y := a.ScalarMult(a.X, a.Y, k.Bytes())
	return &CurvePoint{CURVE, x, y}
}

func pointToBytes(point *CurvePoint) (res []byte) {
	res = elliptic.Marshal(CURVE, point.X, point.Y)
	return
}

func hashToCurve(hash []byte) *big.Int {
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.Mod(hashInt, N)
}

func concatBytes(a, b []byte) []byte {
	var buf bytes.Buffer
	buf.Write(a)
	buf.Write(b)
	return buf.Bytes()
}

func bigIntAdd(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Add(a, b)
	res.Mod(res, N)
	return
}

func bigIntMul(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Mul(a, b)
	res.Mod(res, N)
	return
}

func sha3Hash(message []byte) ([]byte, error) {
	sha := sha3.New256()
	_, err := sha.Write(message)
	if err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}

func encryptKeyGen(pubKey *ecdsa.PublicKey) (capsule *Capsule, keyBytes []byte, err error) {
	s := new(big.Int)
	// generate E,V key-pairs
	priE, pubE, err := generateKeys()
	priV, pubV, err := generateKeys()
	if err != nil {
		return nil, nil, err
	}
	// get H2(E || V)
	h := hashToCurve(
		concatBytes(
			pointToBytes(pubE),
			pointToBytes(pubV)))
	// get s = v + e * H2(E || V)
	s = bigIntAdd(priV.D, bigIntMul(priE.D, h))
	// get (pk_A)^{e+v}
	point := pointScalarMul(pubKey, bigIntAdd(priE.D, priV.D))
	// generate aes key
	keyBytes, err = sha3Hash(pointToBytes(point))
	if err != nil {
		return nil, nil, err
	}
	capsule = &Capsule{
		E: pubE,
		V: pubV,
		S: s,
	}
	return capsule, keyBytes, nil
}

func gcmEncrypt(plaintext []byte, key string, iv []byte, additionalData []byte) (cipherText []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	cipherText = aesgcm.Seal(nil, iv, plaintext, additionalData)
	return cipherText, nil
}

func encrypt(message string, pubKey *ecdsa.PublicKey) (cipherText []byte, capsule *Capsule, err error) {
	capsule, keyBytes, err := encryptKeyGen(pubKey)
	if err != nil {
		return nil, nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	cipherText, err = gcmEncrypt([]byte(message), key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, nil, err
	}
	return cipherText, capsule, nil
}

func encodeCapsule(capsule Capsule) (capsuleAsBytes []byte, err error) {
	gob.Register(elliptic.P256())
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(capsule); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func setupAWSSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	return svc
}

// Observation is the data to upload lmao
// todo should this be public?
type Observation struct {
	CipherText string `json:"ciphertext"`
	Capsule    string `json:"capsule"`
}

type ObservationRequest struct {
	Attribute  string `json:"attribute"`
	Ciphertext string `json:"ciphertext"`
	// Capsule     string `json:"capsule"`
	CapsuleE    string `json:"capsuleE"`
	CapsuleV    string `json:"capsuleV"`
	CapsuleS    string `json:"capsuleS"`
	SupertypeID string `json:"supertypeID"`
	PublicKey   string `json:"pk"`
}

// ObservationResponse is returned from the server
// todo should this be public?
type ObservationResponse struct {
	Ciphertext string `json:"ciphertext"`
	// Capsule              string    `json:"capsule"`
	CapsuleE             string    `json:"capsuleE"`
	CapsuleV             string    `json:"capsuleV"`
	CapsuleS             string    `json:"capsuleS"`
	DateAdded            string    `json:"dateAdded"`
	PublicKey            string    `json:"pk"`
	SupertypeID          string    `json:"supertypeID"`
	ReencryptionMetadata [2]string `json:"reencryptionMetadata"`
}

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
	pk, err := stringToPublicKey(&pkVendor)
	if err != nil {
		return errors.New("Error converting to public key")
	}

	// TODO we will use this for re-encryption with new vendors
	// sk, err := stringToPrivateKey(&skVendor, *pk)
	// if err != nil {
	// 	return errors.New("Error converting to private key")
	// }

	// Encrypt data
	cipherText, capsule, err := encrypt(data, pk)
	if err != nil {
		return errors.New("Error encrypt data")
	}

	// capsuleAsBytes, err := encodeCapsule(*capsule)
	// if err != nil {
	// 	return errors.New("Error encoding data")
	// }

	// fmt.Printf("capsule as bytes: %v\n", capsuleAsBytes)
	// fmt.Printf("string version of capsule as bytes: %v\n", string(capsuleAsBytes))
	// s := string(capsuleAsBytes)
	// for _, r := range s {
	// 	fmt.Printf("r: %v\n", r)
	// }
	// rs := []rune(s)
	// fmt.Printf("rs: %v\n", rs)

	// var outstring string
	// for _, v := range rs {
	// 	outstring += string(v)
	// }

	// fmt.Printf("outstring: %v\n", outstring)

	capsuleE := PublicKeyToString(capsule.E)
	capsuleV := PublicKeyToString(capsule.V)
	capsuleS := capsule.S.String()

	ce, err := stringToPublicKey(&capsuleE)
	if err != nil {
		fmt.Println("1")
	}

	cv, err := stringToPublicKey(&capsuleV)
	if err != nil {
		fmt.Println("1")
	}

	x := new(big.Int)
	x, ok := x.SetString(capsuleS, 10)
	if !ok {
		return errors.New("Ooopsie doopsie")
	}

	decodedCapsule := Capsule{
		E: ce,
		V: cv,
		S: x,
	}

	fmt.Printf("capsule: %v\n", capsule.E)
	fmt.Printf("decoded capsule: %v\n", decodedCapsule.E)

	fmt.Printf("capsule: %v\n", capsule.V)
	fmt.Printf("decoded capsule: %v\n", decodedCapsule.V)

	fmt.Printf("capsule: %v\n", capsule.S)
	fmt.Printf("decoded capsule: %v\n", decodedCapsule.S)

	sk, err := stringToPrivateKey(&skVendor, *pk)
	if err != nil {
		fmt.Println("lskdjf")
	}

	pt, err := DecryptOnMyPriKey(sk, &decodedCapsule, cipherText)
	if err != nil {
		fmt.Println("lkj")
	}

	fmt.Printf("pt: %v\n", string(pt))

	fmt.Printf("ciphertextAsBytes: %v\n", cipherText)

	obs := ObservationRequest{
		Attribute:  attribute,
		Ciphertext: hex.EncodeToString(cipherText),
		// Capsule:     base64.RawStdEncoding.EncodeToString(capsuleAsBytes),
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

// Recreate aes key
func RecreateAESKeyByMyPriKey(capsule *Capsule, aPriKey *ecdsa.PrivateKey) (keyBytes []byte, err error) {
	point1 := pointScalarAdd(capsule.E, capsule.V)
	point := pointScalarMul(point1, aPriKey.D)
	// generate aes key
	keyBytes, err = sha3Hash(pointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// Decrypt by my own private key
func DecryptOnMyPriKey(aPriKey *ecdsa.PrivateKey, capsule *Capsule, cipherText []byte) (plainText []byte, err error) {
	keyBytes, err := RecreateAESKeyByMyPriKey(capsule, aPriKey)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	plainText, err = GCMDecrypt(cipherText, key[:32], keyBytes[:12], nil)
	return plainText, err
}

// convert private key to string
func PrivateKeyToString(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

// convert public key to string
func PublicKeyToString(publicKey *ecdsa.PublicKey) string {
	pubKeyBytes := pointToBytes(publicKey)
	return hex.EncodeToString(pubKeyBytes)
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
	pk, err := stringToPublicKey(&pkVendor)
	if err != nil {
		return errors.New("Error converting to public key")
	}

	sk, err := stringToPrivateKey(&skVendor, *pk)
	if err != nil {
		return errors.New("Error converting to private key")
	}

	// Iterate through each observation
	for _, obs := range observations {
		capsuleE, err := stringToPublicKey(&obs.CapsuleE)
		if err != nil {
			return errors.New("Error decoding capsule")
		}

		capsuleV, err := stringToPublicKey(&obs.CapsuleV)
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

		// capsuleAsBytes, err := base64.RawStdEncoding.DecodeString(obs.CapsuleE)
		// if err != nil {
		// 	return errors.New("Error decoding capsule")
		// }
		// decodedCapsule, err := decodeCapsule(capsuleAsBytes)
		// if err != nil {
		// 	return errors.New("Error decoding capsule")
		// }

		rekey := new(big.Int)
		rekey, ok = rekey.SetString(obs.ReencryptionMetadata[0], 10)
		if !ok {
			return errors.New("Error setting rekey")
		}

		fmt.Printf("rekey: %v\n", rekey)

		pkX, err := stringToPublicKey(&(obs.ReencryptionMetadata[1]))
		if err != nil {
			return errors.New("Error decoding pkX")
		}

		fmt.Printf("pkX string: %v\n", obs.ReencryptionMetadata[1])

		newCapsule, err := reEncryption(rekey, &decodedCapsule)
		if err != nil {
			return errors.New("Error re-encrypting")
		}

		fmt.Printf("ciphertextAsBytes: %v\n", ciphertextAsBytes)
		plainText, err := decrypt(sk, newCapsule, pkX, ciphertextAsBytes)
		if err != nil {
			fmt.Printf("Error decrypting... %v\n", err)
		}
		fmt.Printf("plaintext: %v\n", string(plainText))
	}

	return nil
}

// Server executes Re-Encryption method
func reEncryption(rk *big.Int, capsule *Capsule) (*Capsule, error) {
	// check g^s == V * E^{H2(E || V)}
	x1, y1 := CURVE.ScalarBaseMult(capsule.S.Bytes())
	tempX, tempY := CURVE.ScalarMult(capsule.E.X, capsule.E.Y,
		hashToCurve(
			concatBytes(
				pointToBytes(capsule.E),
				pointToBytes(capsule.V))).Bytes())
	x2, y2 := CURVE.Add(capsule.V.X, capsule.V.Y, tempX, tempY)
	// if check failed return error
	if x1.Cmp(x2) != 0 || y1.Cmp(y2) != 0 {
		return nil, fmt.Errorf("%s", "Capsule not match")
	}
	// E' = E^{rk}, V' = V^{rk}
	newCapsule := &Capsule{
		E: pointScalarMul(capsule.E, rk),
		V: pointScalarMul(capsule.V, rk),
		S: capsule.S,
	}
	return newCapsule, nil
}

func decodeCapsule(capsuleAsBytes []byte) (capsule Capsule, err error) {
	capsule = Capsule{}
	gob.Register(elliptic.P256())
	dec := gob.NewDecoder(bytes.NewBuffer(capsuleAsBytes))
	if err = dec.Decode(&capsule); err != nil {
		fmt.Printf("error... %v\n", err)
		return capsule, err
	}
	return capsule, nil
}

func decryptKeyGen(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey) (keyBytes []byte, err error) {
	// S = X_A^{sk_B}
	S := pointScalarMul(pubX, bPriKey.D)
	// recreate d = H3(X_A || pk_B || S)
	d := hashToCurve(
		concatBytes(
			concatBytes(
				pointToBytes(pubX),
				pointToBytes(&bPriKey.PublicKey)),
			pointToBytes(S)))
	point := pointScalarMul(
		pointScalarAdd(capsule.E, capsule.V), d)
	keyBytes, err = sha3Hash(pointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

func pointScalarAdd(a, b *CurvePoint) *CurvePoint {
	x, y := CURVE.Add(a.X, a.Y, b.X, b.Y)
	return &CurvePoint{CURVE, x, y}
}

// Recreate the aes key then decrypt the cipherText
func decrypt(bPriKey *ecdsa.PrivateKey, capsule *Capsule, pubX *ecdsa.PublicKey, cipherText []byte) (plainText []byte, err error) {
	keyBytes, err := decryptKeyGen(bPriKey, capsule, pubX)
	if err != nil {
		return nil, err
	}
	// recreate aes key = G((E' * V')^d)
	key := hex.EncodeToString(keyBytes)
	// use aes gcm to decrypt
	// mark keyBytes[:12] as nonce
	plainText, err = GCMDecrypt(cipherText, key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("pt: %v\n", plainText)
	return plainText, nil
}

func GCMDecrypt(cipherText []byte, key string, iv []byte, additionalData []byte) (plainText []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plainText, err = aesgcm.Open(nil, iv, cipherText, additionalData)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil, err
	}
	return plainText, nil
}
