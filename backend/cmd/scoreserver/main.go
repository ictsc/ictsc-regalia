package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ictsc/ictsc-regalia/backend/pkg/otelutil"
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

	slog.SetDefault(
		slog.New(slogutil.NewHandler(os.Stdout, opts.LogFormat, opts.LogLevel)),
	)

	shutdownOTel, err := otelutil.Setup(ctx, "scoreserver-backend")
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

	slog.Info("Shutting down the server gracefully", "graceful_period", opts.GracefulPeriod)
	ctx, cancel := context.WithTimeout(context.Background(), opts.GracefulPeriod)
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
