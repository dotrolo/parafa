package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		args    []string
		want    Config
		wantErr bool
	}{
		{
			name: "defaults when nothing is set",
			args: []string{},
			want: Config{
				SeedPath:  "/var/lib/parafa/seed",
				APIAddr:   "127.0.0.1:8080",
				AdminAddr: "127.0.0.1:8081",
			},
		},
		{
			name: "env overrides defaults",
			env: map[string]string{
				"PARAFA_SEED_PATH":     "/opt/seed",
				"PARAFA_API_ADDRESS":   "127.0.0.1:9000",
				"PARAFA_ADMIN_ADDRESS": "127.0.0.1:9001",
			},
			args: []string{},
			want: Config{
				SeedPath:  "/opt/seed",
				APIAddr:   "127.0.0.1:9000",
				AdminAddr: "127.0.0.1:9001",
			},
		},
		{
			name: "flags override defaults",
			args: []string{
				"--seed-path=/tmp/seed",
				"--api-addr=127.0.0.1:7000",
				"--admin-addr=127.0.0.1:7001",
			},
			want: Config{
				SeedPath:  "/tmp/seed",
				APIAddr:   "127.0.0.1:7000",
				AdminAddr: "127.0.0.1:7001",
			},
		},
		{
			name: "flags override env",
			env: map[string]string{
				"PARAFA_SEED_PATH":     "/opt/seed",
				"PARAFA_API_ADDRESS":   "127.0.0.1:9000",
				"PARAFA_ADMIN_ADDRESS": "127.0.0.1:9001",
			},
			args: []string{
				"--seed-path=/tmp/seed",
				"--api-addr=127.0.0.1:7000",
				"--admin-addr=127.0.0.1:7001",
			},
			want: Config{
				SeedPath:  "/tmp/seed",
				APIAddr:   "127.0.0.1:7000",
				AdminAddr: "127.0.0.1:7001",
			},
		},
		{
			name:    "invalid api address",
			args:    []string{"--api-addr=asdf"},
			wantErr: true,
		},
		{
			name:    "invalid admin address",
			args:    []string{"--admin-addr=fdsa"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"--asdfffdsasdd=1"},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// clear env set before
			t.Setenv("PARAFA_SEED_PATH", "")
			t.Setenv("PARAFA_API_ADDRESS", "")
			t.Setenv("PARAFA_ADMIN_ADDRESS", "")

			for k, v := range test.env {
				t.Setenv(k, v)
			}

			got, _, err := Load(test.args)

			if (err != nil) != test.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}

			if got != test.want {
				t.Errorf("Load() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestWarnings(t *testing.T) {
	tests := []struct {
		name      string
		adminAddr string
		want      int
	}{
		{name: "loopback ipv4", adminAddr: "127.0.0.1:8081", want: 0},
		{name: "loopback ipv6", adminAddr: "[::1]:8081", want: 0},
		{name: "localhost", adminAddr: "localhost:8081", want: 0},
		{name: "all interfaces", adminAddr: ":8081", want: 1},
		{name: "public ip", adminAddr: "1.2.3.4:8081", want: 1},
		{name: "hostname", adminAddr: "example.com:8081", want: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := Config{AdminAddr: test.adminAddr}

			got := cfg.warnings()

			if len(got) != test.want {
				t.Errorf("warnings() returned %d warnings %v, want %d", len(got), got, test.want)
			}
		})
	}
}
