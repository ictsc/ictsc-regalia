package main

import (
	"flag"
	"log/slog"
	"net/netip"
	"time"

	"github.com/cockroachdb/errors"
)

type CLIOption struct {
	Dev     bool

	GracefulPeriod time.Duration
	LogLevel       slog.Level

	AdminHTTPAddr         AddrPortValue
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

	opt.AdminHTTPAddr = AddrPortValue(netip.MustParseAddrPort("127.0.0.1:8081"))
	fs.Var(&opt.AdminHTTPAddr, "admin.http-addr", "Admin HTTP server address")
	fs.StringVar(&opt.AdminAuthConfig, "admin.auth-config", "", "Admin API authentication config file")
	fs.StringVar(&opt.AdminAuthConfigInline, "admin.auth-config-inline", "", "Admin API authentication config (inline)")

	fs.BoolFunc("v", "Verbose logging", func(string) error {
		opt.LogLevel = slog.LevelDebug
		return nil
	})

	return &opt
}

type AddrPortValue netip.AddrPort

func (v *AddrPortValue) Set(s string) error {
	addrPort, err := netip.ParseAddrPort(s)
	if err != nil {
		return errors.WithStack(err)
	}
	*v = AddrPortValue(addrPort)
	return nil
}
func (v *AddrPortValue) String() string {
	return netip.AddrPort(*v).String()
}
