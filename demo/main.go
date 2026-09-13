package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/dotrolo/parafa/mintd/keys"
)

func main() {
	fmt.Println("starting demo...")
	spent := make(map[[32]byte]bool)

	// temp directory for demo seed
	dir, err := os.MkdirTemp("", "parafa-demo")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	seedPath := filepath.Join(dir, "seed")

	// demo passphrase for seed file
	passphrase := []byte("demo")

	fmt.Println("loading seed...")
	if err := keys.Create(seedPath, passphrase); err != nil {
		log.Fatal(err)
	}

	seed, err := keys.Load(seedPath, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("loaded")

	fmt.Println("\"client\" generates serial (s -> S)...")

	// new serial (client side)
	var s [32]byte
	rand.Read(s[:])

	S := keys.HashToCurve(s[:])

	fmt.Println("generate blinding factor (r)...")
	// blinding factor
	r := randomScalar()

	fmt.Println("blinding it (B = r * G + S)...")

	// make blinded serial (B)
	var rG secp256k1.JacobianPoint
	secp256k1.ScalarBaseMultNonConst(r, &rG) // r * G

	var B secp256k1.JacobianPoint
	secp256k1.AddNonConst(S, &rG, &B) // r*G + S

	fmt.Println("done, blinded serial (B):", B)

	fmt.Println("mintd signing B with epoch & denom...")
	// mintd signs it with their current epoch, we get k * B = k(S + r*G) = k*S + k*r*G
	kB := seed.Sign(&B, 1, "2026Q3")

	fmt.Println("done, signed:", kB)

	fmt.Println("client unblinding it...")

	// we need the public key, the client (wallet) can simply request this
	K := seed.DerivePublic(1, "2026Q3")

	// multiply r with K (K = k * G), we get k*r*G
	var rK secp256k1.JacobianPoint
	secp256k1.ScalarMultNonConst(r, K, &rK)

	var kS secp256k1.JacobianPoint
	// now we need to subtract rK (=krG) from kB (=kS+krG)

	rK.Y.Negate(1) // to subtract we just negate the Y axis and do addition

	secp256k1.AddNonConst(kB, &rK, &kS)
	fmt.Println("unblinded signature:", kS)

	// verify / redeem note
	fmt.Println("redeeming...")
	redeem(seed, s, &kS, spent)

	// try to redeem it twice (should fail)
	fmt.Println("redeeming again...")
	redeem(seed, s, &kS, spent)
}

// client side, this is used to create the blinding factor (r) that we can use to blind our serial
// This step should only be done offline (in a wallet). we only use this for our demo!
func randomScalar() *secp256k1.ModNScalar {
	var b [32]byte

	for {
		rand.Read(b[:])

		var r secp256k1.ModNScalar

		overflow := r.SetBytes(&b)
		if overflow == 0 && !r.IsZero() {
			return &r
		}
	}
}

// checks a note and burns it, if note is already spent, reject
func redeem(seed *keys.Seed, s [32]byte, stamp *secp256k1.JacobianPoint, spent map[[32]byte]bool) bool {
	if spent[s] {
		fmt.Println("rejected: serial already spent")
		return false
	}

	if !seed.Verify(s[:], stamp, 1, "2026Q3") {
		fmt.Println("rejected: invalid signature")
		return false
	}

	spent[s] = true
	fmt.Println("accepted")
	return true
}
