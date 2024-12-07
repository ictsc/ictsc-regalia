package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/pkg/slogutil"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver"
	"golang.org/x/sys/unix"
)

// nolint:gochecknoglobals
var (
	flagContestantHTTPAddr = flag.String("contestant-http-addr", "0.0.0.0:8080", "Contestant API HTTP Address")
	flagAdminHTTPAddr      = flag.String("admin-http-addr", "0.0.0.0:8081", "Admin API HTTP Address")
	flagTelemetryAddr      = flag.String("telemetry-addr", "0.0.0.0:9090", "Telemetry HTTP Address")
	flagGracefulPeriod     = flag.String("graceful-period", "30s", "Graceful period before shutting down the server")
	flagDev                = flag.Bool("dev", false, "Run in development mode")
	flagVerbose            = flag.Bool("v", false, "Verbose logging")
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

	gracefulPeriod, err := time.ParseDuration(*flagGracefulPeriod)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse graceful period", "error", err)
		return 1
	}
	if gracefulPeriod < 0 {
		slog.ErrorContext(ctx, "Graceful period must be positive.")
		return 1
	}

	cfg, err := newConfig()
	if err != nil {
		slog.Error("Failed to create app config", "error", err)
		return 1
	}

	server, err := scoreserver.New(cfg)
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

	return 0
}
