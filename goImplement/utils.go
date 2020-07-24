package goimplement

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"golang.org/x/crypto/sha3"
)

var CURVE = elliptic.P256()
var P = CURVE.Params().P
var N = CURVE.Params().N

type Capsule struct {
	E *ecdsa.PublicKey
	V *ecdsa.PublicKey
	S *big.Int
}

type CurvePoint = ecdsa.PublicKey

func PointScalarAdd(a, b *CurvePoint) *CurvePoint {
	x, y := CURVE.Add(a.X, a.Y, b.X, b.Y)
	return &CurvePoint{CURVE, x, y}
}

func PointScalarMul(a *CurvePoint, k *big.Int) *CurvePoint {
	x, y := a.ScalarMult(a.X, a.Y, k.Bytes())
	return &CurvePoint{CURVE, x, y}
}

func PointToBytes(point *CurvePoint) (res []byte) {
	res = elliptic.Marshal(CURVE, point.X, point.Y)
	return
}

func HashToCurve(hash []byte) *big.Int {
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.Mod(hashInt, N)
}

func ConcatBytes(a, b []byte) []byte {
	var buf bytes.Buffer
	buf.Write(a)
	buf.Write(b)
	return buf.Bytes()
}

func BigIntAdd(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Add(a, b)
	res.Mod(res, N)
	return
}

func BigIntMul(a, b *big.Int) (res *big.Int) {
	res = new(big.Int).Mul(a, b)
	res.Mod(res, N)
	return
}

func Sha3Hash(message []byte) ([]byte, error) {
	sha := sha3.New256()
	_, err := sha.Write(message)
	if err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}
