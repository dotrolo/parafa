package keys

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// create a private key from seed with params: denomination and epoch
func (s *Seed) derive(denom uint64, epoch string) *secp256k1.ModNScalar {
	// our scalar (private key)
	var k secp256k1.ModNScalar

	// in case of overflow we need to make a new one with a slightly different input (counter)
	counter := uint64(0)
	for {
		// hash by "adding" each element to the mix
		h := sha256.New()
		h.Write(s.bytes)
		h.Write([]byte(epoch))
		// convert uint64 to bytes
		denomBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(denomBytes, denom)
		h.Write(denomBytes)

		counterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(counterBytes, counter)
		h.Write(counterBytes)

		hash := h.Sum(nil)

		// SetBytes needs type [32]byte
		var arr [32]byte
		copy(arr[:], hash)

		overflow := k.SetBytes(&arr)
		if overflow == 0 && !k.IsZero() {
			return &k
		}
		counter++
	}
}

func (s *Seed) DerivePublic(denom uint64, epoch string) *secp256k1.JacobianPoint {
	k := s.derive(denom, epoch)

	var pub secp256k1.JacobianPoint
	secp256k1.ScalarBaseMultNonConst(k, &pub)
	// convert to normal coordinate format
	pub.ToAffine()

	return &pub
}
