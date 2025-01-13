package main

import (
	"flag"
	"net/netip"
	"time"
)

type CLIOption struct {
	Dev     bool
	Verbose bool

	GracefulPeriod time.Duration

	AdminHTTPAddr         AddrPortValue
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

	opt.AdminHTTPAddr = AddrPortValue(netip.MustParseAddrPort("127.0.0.1:8081"))
	fs.Var(&opt.AdminHTTPAddr, "admin.http-addr", "Admin HTTP server address")
	fs.StringVar(&opt.AdminAuthConfig, "admin.auth-config", "", "Admin API authentication config file")
	fs.StringVar(&opt.AdminAuthConfigInline, "admin.auth-config-inline", "", "Admin API authentication config (inline)")

	return &opt
}

type AddrPortValue netip.AddrPort

func (v *AddrPortValue) Set(s string) error {
	addrPort, err := netip.ParseAddrPort(s)
	if err != nil {
		return err
	}
	*v = AddrPortValue(addrPort)
	return nil
}
func (v *AddrPortValue) String() string {
	return netip.AddrPort(*v).String()
}
