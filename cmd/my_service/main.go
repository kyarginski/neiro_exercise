// My Service for Neiro.
//
// # Description of the REST API of the service for Neiro.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http, https
// Host: localhost
// Version: 1.0.0
//
// swagger:meta
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"neiro/internal/app"
	"neiro/internal/config"
	"neiro/internal/lib/logger/sl"

	_ "github.com/lib/pq"
)

const (
	serviceName = "my_service"
)

func main() {
	cfg := config.MustLoad(serviceName)
	log := sl.SetupLogger(cfg.Env)
	log.Info(
		"starting server "+serviceName,
		slog.String("env", cfg.Env),
		slog.String("version", cfg.Version),
		slog.Int("port", cfg.Port),
		slog.Int("ttl", cfg.TTL),
		slog.Int("clear_interval", cfg.ClearInterval),
		slog.Bool("use_tracing", cfg.UseTracing),
		slog.String("tracing_address", cfg.TracingAddress),
	)

	if err := run(log, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
}

func run(log *slog.Logger, cfg *config.Config) error {
	application, err := app.NewMyService(log, cfg.Port, cfg.UseTracing, cfg.TracingAddress, cfg.TTL, time.Duration(cfg.ClearInterval), serviceName)
	if err != nil {
		return err
	}
	defer application.Stop()

	application.Start()

	return nil
}
