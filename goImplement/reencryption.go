package goimplement

import (
	"fmt"
	"math/big"
)

// Server executes Re-Encryption method
func ReEncryption(rk *big.Int, capsule *Capsule) (*Capsule, error) {
	// check g^s == V * E^{H2(E || V)}
	x1, y1 := CURVE.ScalarBaseMult(capsule.S.Bytes())
	tempX, tempY := CURVE.ScalarMult(capsule.E.X, capsule.E.Y,
		HashToCurve(
			ConcatBytes(
				PointToBytes(capsule.E),
				PointToBytes(capsule.V))).Bytes())
	x2, y2 := CURVE.Add(capsule.V.X, capsule.V.Y, tempX, tempY)
	// if check failed return error
	if x1.Cmp(x2) != 0 || y1.Cmp(y2) != 0 {
		return nil, fmt.Errorf("%s", "Capsule not match")
	}
	// E' = E^{rk}, V' = V^{rk}
	newCapsule := &Capsule{
		E: PointScalarMul(capsule.E, rk),
		V: PointScalarMul(capsule.V, rk),
		S: capsule.S,
	}
	return newCapsule, nil
}
