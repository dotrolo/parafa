package keys

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

type Seed struct {
	bytes []byte
}

const (
	seedSize     = 32
	saltSize     = 16
	argonTime    = 3
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func Load(seedPath string, passphrase []byte) (*Seed, error) {
	// seed fileinfo
	fi, err := os.Stat(seedPath)
	if err != nil {
		return nil, err
	}

	// file's permissions
	fip := fi.Mode().Perm()

	// AND bitwise to check whether group/any user has access
	// normally 0600 & 0077 = 0
	if fip&0o077 != 0 {
		return nil, fmt.Errorf("dangerous perms on seed file %q: %04o", seedPath, fip)
	}

	// also check dir's perms so seed can't be replaced by another user
	dir := filepath.Dir(seedPath)

	di, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	dip := di.Mode().Perm()
	if dip&0o077 != 0 {
		return nil, fmt.Errorf("dangerous perms on seed parent dir %q: %04o", dir, dip)
	}

	// read encrypted data
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return nil, err
	}

	if len(data) < saltSize {
		return nil, fmt.Errorf("seed file %q is only %d bytes", seedPath, len(data))
	}

	salt := data[:saltSize]

	key := argon2.IDKey(passphrase, salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceEnd := saltSize + gcm.NonceSize()

	if len(data) < nonceEnd {
		return nil, fmt.Errorf("seed file %q is only %d bytes", seedPath, len(data))
	}

	nonce := data[saltSize:nonceEnd]

	seed, err := gcm.Open(nil, nonce, data[nonceEnd:], nil)
	if err != nil {
		return nil, err
	}

	return &Seed{bytes: seed}, nil
}

func Create(seedPath string, passphrase []byte) error {
	salt := make([]byte, saltSize)
	rand.Read(salt)

	// generate hash from salt & passphrase using argon2
	// it gives a fixed sized key back and is secure against brute-force
	key := argon2.IDKey(passphrase, salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// turn it into a block aes can work on
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// adds tag, so we can verify passphrase
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	seed := make([]byte, seedSize)
	rand.Read(seed)

	sealed := gcm.Seal(nil, nonce, seed, nil)

	fileData := bytes.Join([][]byte{salt, nonce, sealed}, nil)

	dir := filepath.Dir(seedPath)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return err
	}

	// temp file (f) now then rename it to seedPath
	f, err := os.CreateTemp(dir, "seed.*")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name()) // cleanup in case of failure

	if _, err := f.Write(fileData); err != nil {
		return err
	}
	if err := os.Chmod(f.Name(), 0600); err != nil {
		return err
	}

	// sync file
	if err := f.Sync(); err != nil {
		return err
	}

	f.Close()

	if err := os.Rename(f.Name(), seedPath); err != nil {
		return err
	}

	// sync parent dir
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	if err := d.Sync(); err != nil {
		return err
	}

	d.Close()

	return nil
}

// helper used in main
func ReadPassphrase() ([]byte, error) {
	fd := int(os.Stdin.Fd())

	var pass []byte
	var err error

	if term.IsTerminal(fd) {
		fmt.Print("Enter seed passphrase: ")
		pass, err = term.ReadPassword(fd)
		fmt.Println()
	} else {
		// if password goes through a pipe
		pass, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		return nil, err
	}

	pass = bytes.TrimRight(pass, "\r\n")

	if len(pass) == 0 {
		return nil, errors.New("empty passphrase")
	}

	return pass, nil
}
