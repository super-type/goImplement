package goimplement

import (
	"crypto/ecdsa"
	"math/big"
)

type Capsule struct {
	E *ecdsa.PublicKey
	V *ecdsa.PublicKey
	S *big.Int
}
