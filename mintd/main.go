package main

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotrolo/parafa/mintd/internal/admin"
	"github.com/dotrolo/parafa/mintd/internal/api"
	"github.com/dotrolo/parafa/mintd/internal/config"
	"github.com/dotrolo/parafa/mintd/internal/keys"
)

func main() {
	// load and validate configuration
	cfg, warns, err := config.Load(os.Args[1:])
	if err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
	// log settings then warnings
	slog.Info("configuration loaded",
		"seed_path", cfg.SeedPath,
		"api_addr", cfg.APIAddr,
		"admin_addr", cfg.AdminAddr,
	)
	for _, v := range warns {
		slog.Warn(v, "addr", cfg.AdminAddr)
	}

	// ask for password to deal with seed file
	passphrase, err := keys.ReadPassphrase()
	if err != nil {
		slog.Error("reading passphrase failed", "err", err)
		os.Exit(1)
	}
	// load seed file, if doesn't exist generate one
	seed, err := keys.Load(cfg.SeedPath, passphrase)
	if errors.Is(err, fs.ErrNotExist) {
		if err := keys.Create(cfg.SeedPath, passphrase); err != nil {
			slog.Error("seed generation failed", "err", err, "seed_path", cfg.SeedPath)
			os.Exit(1)
		}

		slog.Warn("seed successfully generated, BACK IT UP before restarting", "seed_path", cfg.SeedPath)

		os.Exit(0)
	} else if err != nil {
		slog.Error("loading seed failed", "err", err, "seed_path", cfg.SeedPath)
		os.Exit(1)
	}

	_ = seed // temporary, todo: use seed

	// public api used by wallets
	pub := &http.Server{
		Addr:         cfg.APIAddr,
		Handler:      api.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// admin api, only used locally
	adm := &http.Server{
		Addr:         cfg.AdminAddr,
		Handler:      admin.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errs := make(chan error, 2)

	// run both admin and public servers
	go func() {
		slog.Info("server listening", "role", "api", "addr", pub.Addr)
		errs <- pub.ListenAndServe()
	}()

	go func() {
		slog.Info("server listening", "role", "admin", "addr", adm.Addr)
		errs <- adm.ListenAndServe()
	}()

	// catch os signals such as interrupt
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	exitCode := 0

	// continue with the first incoming signal/error
	select {
	case sig := <-sigs:
		slog.Info("shutting down", "signal", sig)
	case err := <-errs:
		slog.Error("server failed", "err", err)
		exitCode = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown to make sure active requests finish without interruption
	if err := pub.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "err", err, "role", "api", "addr", pub.Addr)
	} else {
		slog.Info("shutdown done", "role", "api", "addr", pub.Addr)
	}

	if err := adm.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "err", err, "role", "admin", "addr", adm.Addr)
	} else {
		slog.Info("shutdown done", "role", "admin", "addr", adm.Addr)
	}

	os.Exit(exitCode)
}
