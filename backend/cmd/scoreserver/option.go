package main

import (
	"flag"
	"log/slog"
	"net/netip"
	"time"
)

type CLIOption struct {
	Dev bool

	GracefulPeriod time.Duration
	LogLevel       slog.Level

	AdminHTTPAddr         netip.AddrPort
	AdminAuthConfig       string
	AdminAuthConfigInline string
}

// NewOption creates a new CLIOption combined with the given flag.FlagSet.
//
//nolint:varnamelen,mnd
func NewOption(fs *flag.FlagSet) *CLIOption {
	opt := CLIOption{}

	fs.BoolVar(&opt.Dev, "dev", false, "Development mode")

	fs.DurationVar(&opt.GracefulPeriod, "graceful-period", 30*time.Second, "Graceful period for shutdown")
	fs.TextVar(&opt.LogLevel, "log-level", slog.LevelInfo, "Log level")

	fs.TextVar(&opt.AdminHTTPAddr, "admin.http-addr", netip.MustParseAddrPort("127.0.0.1:8081"), "Admin HTTP server address")
	fs.StringVar(&opt.AdminAuthConfig, "admin.auth-config", "", "Admin API authentication config file")
	fs.StringVar(&opt.AdminAuthConfigInline, "admin.auth-config-inline", "", "Admin API authentication config (inline)")

	fs.BoolFunc("v", "Verbose logging", func(string) error {
		opt.LogLevel = slog.LevelDebug
		return nil
	})

	return &opt
}
