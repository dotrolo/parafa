package keys

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

type Seed struct {
	bytes []byte
}

const seedSize = 32

func Load(seedPath string) (*Seed, error) {
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

	// read seed
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return nil, err
	}
	if len(data) != seedSize {
		return nil, fmt.Errorf("seed file %q is %d bytes, expected %d", seedPath, len(data), seedSize)
	}

	return &Seed{bytes: data}, nil
}

func Create(seedPath string) error {
	seed := make([]byte, seedSize)
	rand.Read(seed)

	dir := filepath.Dir(seedPath)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	// rwx------
	if err := os.Chmod(dir, 0700); err != nil {
		return err
	}

	// temp file (f) now then rename it to seedPath
	f, err := os.CreateTemp(dir, "seed.*")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name()) // cleanup in case of failure

	if _, err := f.Write(seed); err != nil {
		return err
	}

	// rw-------
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
