package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ictsc/ictsc-regalia/backend/pkg/otelutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/slogutil"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/batch"
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

	slog.SetDefault(
		slog.New(slogutil.NewHandler(os.Stdout, opts.LogFormat, opts.LogLevel)),
	)

	cfg, err := newConfig(opts)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create config", "error", err)
		return 1
	}

	shutdownOTel, err := otelutil.Setup(ctx, "scoreserver-batch")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup OpenTelemetry", "error", err)
		return 1
	}

	batchApp, err := batch.NewBatch(*cfg)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize batch", "error", err)
		return 1
	}

	if err := batchApp.Run(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to start batch", "error", err)
	}

	slog.Info("Shutting down app gracefully", "graceful_period", opts.GracefulPeriod)
	ctx, cancel := context.WithTimeout(context.Background(), opts.GracefulPeriod)
	defer cancel()

	if err := shutdownOTel(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown OpenTelemetry", "error", err)
		return 1
	}

	return 0
}
