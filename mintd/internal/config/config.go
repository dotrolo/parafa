package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"errors"
	"io/fs"
	"strings"
	"io"
)

type Config struct {
	SeedPath   string // where to create/find seed
	APIAddr    string
	AdminAddr  string
	Passphrase string // comes from stdin
}

func Load(args []string) (Config, []string, error) {
	
	// initialize empty flag sets
	conFS := flag.NewFlagSet("mintd", flag.ContinueOnError)
	conFS.SetOutput(io.Discard) // silence output on conFS, so mainFS will be written out
	
	mainFS := flag.NewFlagSet("mintd", flag.ContinueOnError)

	// set only config file path
	configPath := conFS.String(
		"config",
		getConfEnv("PARAFA_CONFIG", filepath.Join("/etc", "parafa", "parafa.conf")), // default value
		"",
	)
	
	// keeping these empty so "flag" package just ignores these instead of giving an error
	conFS.String("seed-path", "", "")
	conFS.String("api-addr", "", "")
	conFS.String("admin-addr", "", "")
	
	if err := conFS.Parse(args); err != nil && !errors.Is(err, flag.ErrHelp) {
		return Config{}, nil, err
	}
	
	values, err := readConfig(*configPath)
	if err != nil {
		return Config{}, nil, err
	}
	
	// set config, flag overwrites env overwrites hardcoded
	seedPath := mainFS.String(
		"seed-path",
		selectConfig(values, "seed_path", "PARAFA_SEED_PATH", filepath.Join("/var", "lib", "parafa", "seed")), // default value
		"Set path where the program can find/write your (super secret) seed file.",
	)
	apiAddr := mainFS.String("api-addr",
		selectConfig(values, "api_addr", "PARAFA_API_ADDRESS", "127.0.0.1:8080"),
		"Set public api address (where wallets can reach your server).",
	)
	adminAddr := mainFS.String(
		"admin-addr",
		selectConfig(values, "admin_addr", "PARAFA_ADMIN_ADDRESS", "127.0.0.1:8081"),
		"Set admin api address (that is only used locally).",
	)
	// only for --help
	mainFS.String(
		"config",
		getConfEnv("PARAFA_CONFIG", filepath.Join("/etc", "parafa", "parafa.conf")),
		"Set configuration file.",
	)

	// apply values to flag set
	if err := mainFS.Parse(args); err != nil {
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

// helper: returns env var / config file / hardcoded.
func selectConfig(values map[string]string, key, envKey, fallback string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	
	if v, ok := values[key]; ok {
		return v
	}

	return fallback
}

// helper for config file path
func getConfEnv(envKey, fallback string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}

	return fallback
}

// read config file into a map
func readConfig(path string) (map[string]string, error) {
	values := make(map[string]string)

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return values, nil
	}
	if err != nil {
		return nil, err
	}

	// read every line & toml key value pair
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return values, nil
}
