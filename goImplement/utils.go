package goimplement

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"golang.org/x/crypto/sha3"
)

// CURVE is an elliptic curve
var CURVE = elliptic.P256()

// P is the order of the underlying field
var P = CURVE.Params().P

// N is the order of the base point
var N = CURVE.Params().N

// CurvePoint is an ECDSA public key
type CurvePoint = ecdsa.PublicKey

// PointScalarAdd adds two scalars
func PointScalarAdd(a, b *CurvePoint) *CurvePoint {
	x, y := CURVE.Add(a.X, a.Y, b.X, b.Y)
	return &CurvePoint{CURVE, x, y}
}

// PointScalarMul multiplies two scalars
func PointScalarMul(a *CurvePoint, k *big.Int) *CurvePoint {
	x, y := a.ScalarMult(a.X, a.Y, k.Bytes())
	return &CurvePoint{CURVE, x, y}
}

// PointToBytes converts a point to bytes
func PointToBytes(point *CurvePoint) (res []byte) {
	res = elliptic.Marshal(CURVE, point.X, point.Y)
	return
}

// HashToCurve converts a byte hash to a curve
func HashToCurve(hash []byte) *big.Int {
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.Mod(hashInt, N)
}

// ConcatBytes concats b to a
func ConcatBytes(a, b []byte) []byte {
	var buf bytes.Buffer
	buf.Write(a)
	buf.Write(b)
	return buf.Bytes()
}

// AddBigInteger adds a BigInteger
func AddBigInteger(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Add(a, b)
	res.Mod(res, N)
	return
}

// MultiplyBigInteger mulitplies a BigInteger
func MultiplyBigInteger(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Mul(a, b)
	res.Mod(res, N)
	return
}

// Sha3Hash hashes data
func Sha3Hash(message []byte) ([]byte, error) {
	sha := sha3.New256()
	_, err := sha.Write(message)
	if err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}

// Contains checks if a string is contained within a string array
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetInverse gets the inverse of a BigInteger
func GetInverse(a *big.Int) (res *big.Int) {
	res = new(big.Int).ModInverse(a, N)
	return
}
