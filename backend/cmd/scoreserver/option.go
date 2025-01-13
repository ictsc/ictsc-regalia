package main

import (
	"flag"
	"time"
)

type CLIOption struct {
	Dev     bool
	Verbose bool

	GracefulPeriod time.Duration

	AdminHTTPAddr         string
	AdminAuthConfig       string
	AdminAuthConfigInline string
}

// NewOption creates a new CLIOption combined with the given flag.FlagSet.
//
//nolint:varnamelen
func NewOption(fs *flag.FlagSet) *CLIOption {
	opt := CLIOption{}

	fs.BoolVar(&opt.Dev, "dev", false, "Development mode")
	fs.BoolVar(&opt.Verbose, "v", false, "Verbose logging")

	fs.DurationVar(&opt.GracefulPeriod, "graceful-period", 30*time.Second, "Graceful period for shutdown")

	fs.StringVar(&opt.AdminHTTPAddr, "admin.http-addr", "0.0.0.0:8080", "Admin HTTP server address")
	fs.StringVar(&opt.AdminAuthConfig, "admin.auth-config", "", "Admin API authentication config file")
	fs.StringVar(&opt.AdminAuthConfigInline, "admin.auth-config-inline", "", "Admin API authentication config (inline)")

	return &opt
}
