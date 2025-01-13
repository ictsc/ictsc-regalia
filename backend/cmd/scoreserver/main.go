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

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opts := NewOption(fs)
	_ = fs.Parse(os.Args[1:])

	os.Exit(start(opts))
}

func start(opts *CLIOption) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, unix.SIGTERM)
	defer stop()

	logLevel := slog.LevelInfo
	if opts.Verbose {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slogutil.NewHandler(opts.Dev, logLevel)))

	gracefulPeriod, err := parseGracefulPeriod(opts.GracefulPeriod)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse graceful period", "error", err)
		return 1
	}

	shutdownOTel, err := setupOpenTelemetry(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup OpenTelemetry", "error", err)
		return 1
	}

	cfg, err := newConfig(opts)
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

func parseGracefulPeriod(s string) (time.Duration, error) {
	gracefulPeriod, err := time.ParseDuration(s)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if gracefulPeriod < 0 {
		return 0, errors.New("graceful period must be positive")
	}
	return gracefulPeriod, nil
}
