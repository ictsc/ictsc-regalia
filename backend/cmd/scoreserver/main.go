package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/slogutil"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver"
	"golang.org/x/sys/unix"
)

// nolint:gochecknoglobals
var (
	// flagContestantHTTPAddr = flag.String("contestant-http-addr", "0.0.0.0:8080", "Contestant API HTTP Address")
	flagAdminHTTPAddr         = flag.String("admin-http-addr", "0.0.0.0:8081", "Admin API HTTP Address")
	flagAdminAuthConfig       = flag.String("admin-auth-config", "", "Admin API authentication config file")
	flagAdminAuthConfigInline = flag.String("admin-auth-config-inline", "", "Admin API authentication config (inline)")

	flagGracefulPeriod = flag.String("graceful-period", "30s", "Graceful period before shutting down the server")
	flagDev            = flag.Bool("dev", false, "Run in development mode")
	flagVerbose        = flag.Bool("v", false, "Verbose logging")
)

func main() {
	os.Exit(start())
}

func start() int {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, unix.SIGTERM)
	defer stop()

	logLevel := slog.LevelInfo
	if *flagVerbose {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slogutil.NewHandler(*flagDev, logLevel)))

	gracefulPeriod, err := parseGracefulPeriod()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse graceful period", "error", err)
		return 1
	}

	shutdownOTel, err := setupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup OpenTelemetry", "error", err)
		return 1
	}

	cfg, err := newConfig()
	if err != nil {
		slog.Error("Failed to create app config", "error", err)
		return 1
	}

	server, err := scoreserver.New(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to init score server", "error", err)
		return 1
	}

	if err := server.Start(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to start score server", "error", err)
		return 1
	}

	<-ctx.Done()

	slog.Info("Shutting down the server gracefully", "graceful_period", gracefulPeriod)
	ctx, cancel := context.WithTimeout(context.Background(), gracefulPeriod)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown the server gracefully", "error", err)
		return 1
	}

	if err := shutdownOTel(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown OpenTelemetry", "error", err)
		return 1
	}

	return 0
}

func parseGracefulPeriod() (time.Duration, error) {
	gracefulPeriod, err := time.ParseDuration(*flagGracefulPeriod)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if gracefulPeriod < 0 {
		return 0, errors.New("graceful period must be positive")
	}
	return gracefulPeriod, nil
}
