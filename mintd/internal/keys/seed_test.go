package keys

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Skip("TODO: update for encrypted seed format")
	tests := []struct {
		name      string
		setup     func(t *testing.T, dir string) string
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "valid seed",
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "seed")
				if err := os.WriteFile(path, make([]byte, 32), 0o600); err != nil {
					t.Fatalf("setup: %v", err)
				}
				return path
			},
			wantErr: false,
		},
		{
			name: "invalid seed wrong len of bytes",
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "seed")
				if err := os.WriteFile(path, make([]byte, 14), 0o600); err != nil {
					t.Fatalf("setup: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		{
			name: "invalid seed wrong open permissions",
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "seed")
				if err := os.WriteFile(path, make([]byte, 32), 0o644); err != nil {
					t.Fatalf("setup: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		{
			name: "invalid wrong permission on parent dir",
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "seed")
				if err := os.WriteFile(path, make([]byte, 32), 0o600); err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.Chmod(dir, 0o744); err != nil {
					t.Fatalf("setup: %v", err)
				}
				return path
			},
			wantErr: true,
		},
		{
			name: "invalid seed doesnt exist",
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "seed")
				return path
			},
			wantErr:   true,
			wantErrIs: fs.ErrNotExist,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.Chmod(dir, 0o700); err != nil {
				t.Fatalf("setup: %v", err)
			}

			path := test.setup(t, dir)

			_, err := Load(path)

			if (err != nil) != test.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, test.wantErr)
			}

			if test.wantErrIs != nil && !errors.Is(err, test.wantErrIs) {
				t.Errorf("Load() error = %v, want errors.Is(..., %v)", err, test.wantErrIs)
			}

			//t.Logf("got error: %v", err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Skip("TODO: update for encrypted seed format")
	t.Run("writes a valid seed file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "seed")

		if err := Create(path); err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}

		fi, err := os.Stat(path)
		if err != nil {
			t.Fatalf("seed file was not created: %v", err)
		}

		if fi.Size() != seedSize {
			t.Errorf("seed file is %d bytes, want %d", fi.Size(), seedSize)
		}

		if perm := fi.Mode().Perm(); perm != 0o600 {
			t.Errorf("seed file mode is %04o, want 0600", perm)
		}

		di, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("could not stat parent dir: %v", err)
		}

		if perm := di.Mode().Perm(); perm != 0o700 {
			t.Errorf("parent dir mode is %04o, want 0700", perm)
		}
	})

	t.Run("creates missing parent directories", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "a", "b", "seed")

		if err := Create(path); err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}

		if _, err := os.Stat(path); err != nil {
			t.Errorf("seed file was not created: %v", err)
		}
	})

	t.Run("produces a different seed each time", func(t *testing.T) {
		dir := t.TempDir()
		first := filepath.Join(dir, "first")
		second := filepath.Join(dir, "second")

		if err := Create(first); err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}
		if err := Create(second); err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}

		a, err := os.ReadFile(first)
		if err != nil {
			t.Fatalf("could not read first seed: %v", err)
		}
		b, err := os.ReadFile(second)
		if err != nil {
			t.Fatalf("could not read second seed: %v", err)
		}

		if bytes.Equal(a, b) {
			t.Error("identical seeds!!!!")
		}
	})

	t.Run("leaves no temporary files behind", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "seed")

		if err := Create(path); err != nil {
			t.Fatalf("Create() error = %v, want nil", err)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("could not read dir: %v", err)
		}

		if len(entries) != 1 {
			t.Errorf("directory has %d entries, want 1", len(entries))
		}
	})
}
