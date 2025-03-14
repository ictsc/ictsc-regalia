package main

import (
	"flag"
	"log/slog"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/pkg/slogutil"
)

type CLIOption struct {
	GracefulPeriod time.Duration
	LogFormat      slogutil.Format
	LogLevel       slog.Level

	APIURL       string
	APITokenFile string

	DeploymentSyncPeriod time.Duration
	ScoreUpdatePeriod    time.Duration

	SStateCAFile        string
	SStateTLSSkipVerify bool
}

//nolint:varnamelen,mnd
func NewOption(fs *flag.FlagSet) *CLIOption {
	opt := CLIOption{}

	fs.DurationVar(&opt.GracefulPeriod, "graceful-period", 30*time.Second, "Graceful period for shutdown")
	fs.TextVar(&opt.LogFormat, "log-format", slogutil.FormatJSON, "Log format (json, console, pretty)")
	fs.TextVar(&opt.LogLevel, "log-level", slog.LevelInfo, "Log level")

	fs.StringVar(&opt.APIURL, "api-url", "http://localhost:8081", "API URL")
	fs.StringVar(&opt.APITokenFile, "api-token-file", "", "API token file")

	fs.DurationVar(&opt.DeploymentSyncPeriod, "deployment-sync-period", 15*time.Second, "Deployment sync period")
	fs.DurationVar(&opt.ScoreUpdatePeriod, "score-update-period", time.Minute, "Score update period")

	fs.StringVar(&opt.SStateCAFile, "sstate-ca-file", "", "SState CA file")
	fs.BoolVar(&opt.SStateTLSSkipVerify, "sstate-tls-skip-verify", false, "SState TLS skip verify")

	fs.BoolFunc("dev", "Use development mode", func(_ string) error {
		opt.LogFormat = slogutil.FormatPretty
		opt.LogLevel = slog.LevelDebug
		return nil
	})

	return &opt
}
