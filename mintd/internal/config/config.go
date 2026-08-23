package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

type Config struct {
	SeedPath   string // where to create/find seed
	APIAddr    string
	AdminAddr  string
	Passphrase string // comes from stdin
}

func Load(args []string) (Config, []string, error) {

	fs := flag.NewFlagSet("mintd", flag.ContinueOnError)

	// set config, flag overwrites env overwrites hardcoded
	seedPath := fs.String(
		"seed-path",
		getEnv("PARAFA_SEED_PATH", filepath.Join("/var", "lib", "parafa", "seed")),
		"Set path where the program can find/write your (super secret) seed file.",
	)
	apiAddr := fs.String("api-addr",
		getEnv("PARAFA_API_ADDRESS", "127.0.0.1:8080"),
		"Set public api address (where wallets can reach your server).",
	)
	adminAddr := fs.String("admin-addr",
		getEnv("PARAFA_ADMIN_ADDRESS", "127.0.0.1:8081"),
		"Set admin api address (that is only used locally).",
	)

	// apply values
	if err := fs.Parse(args); err != nil {
		return Config{}, nil, err
	}

	config := Config{
		SeedPath:  *seedPath,
		APIAddr:   *apiAddr,
		AdminAddr: *adminAddr,
	}

	if err := config.validate(); err != nil {
		return Config{}, nil, err
	}

	warns := config.warnings()

	return config, warns, nil
}

// basic input validation on config
func (c Config) validate() error {
	_, _, err := net.SplitHostPort(c.APIAddr)
	if err != nil {
		return fmt.Errorf("invalid api-addr %q: %w", c.APIAddr, err)
	}

	_, _, err = net.SplitHostPort(c.AdminAddr)
	if err != nil {
		return fmt.Errorf("invalid admin-addr %q: %w", c.AdminAddr, err)
	}

	return nil
}

// find warnings on config
func (c Config) warnings() []string {
	var warns []string
	h, _, _ := net.SplitHostPort(c.AdminAddr)
	ip := net.ParseIP(h)

	switch {
	case h == "":
		warns = append(warns, "admin API is listening on all network interfaces")
	case h == "localhost":
		// loopback
	case ip == nil:
		warns = append(warns, "admin API host is not an IP, cannot verify it is local")
	case !ip.IsLoopback():
		warns = append(warns, "admin API is not bound to loopback")
	}

	return warns
}

// helper: return env var instead if exists
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
