package keys

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// convert a serial to a point on the curve
func hashToCurve(serial []byte) *secp256k1.JacobianPoint {
	counter := uint64(0)

	for {
		// hash serial
		h := sha256.New()
		h.Write(serial)

		counterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(counterBytes, counter)
		h.Write(counterBytes)

		hash := h.Sum(nil)

		var arr [32]byte
		copy(arr[:], hash)

		// turn hash into x
		var x secp256k1.FieldVal
		overflow := x.SetBytes(&arr)
		if overflow == 1 {
			counter++
			continue
		}

		// check if y is on curve
		var y secp256k1.FieldVal
		if !secp256k1.DecompressY(&x, false, &y) { // 'false' is just to tell which y to look at (curve gives 2 y for x)
			counter++
			continue
		}

		// z is set to 1 because we don't want to scale it
		var z secp256k1.FieldVal
		z.SetInt(1)

		p := secp256k1.MakeJacobianPoint(&x, &y, &z)
		return &p
	}
}
