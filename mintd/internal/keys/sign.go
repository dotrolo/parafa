package keys

import "github.com/decred/dcrd/dcrec/secp256k1/v4"

func (s *Seed) Sign(blinded *secp256k1.JacobianPoint, denom uint64, epoch string) *secp256k1.JacobianPoint {
	k := s.derive(denom, epoch)

	var result secp256k1.JacobianPoint

	// multiply blinded serial with k
	secp256k1.ScalarMultNonConst(k, blinded, &result)

	return &result
}

func (s *Seed) Verify(serial []byte, stamp *secp256k1.JacobianPoint, denom uint64, epoch string) bool {
	point := hashToCurve(serial)

	k := s.derive(denom, epoch)

	// sign the provided serial
	var expected secp256k1.JacobianPoint
	secp256k1.ScalarMultNonConst(k, point, &expected)

	// return whether stamp & signed serial is the same
	return expected.EquivalentNonConst(stamp)
}
