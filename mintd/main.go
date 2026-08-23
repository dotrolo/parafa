package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotrolo/parafa/mintd/internal/admin"
	"github.com/dotrolo/parafa/mintd/internal/api"
	"github.com/dotrolo/parafa/mintd/internal/config"
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
