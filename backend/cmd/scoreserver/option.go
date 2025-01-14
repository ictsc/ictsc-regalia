package main

import (
	"flag"
	"log/slog"
	"net/netip"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/slogutil"
)

type CLIOption struct {
	GracefulPeriod time.Duration
	LogFormat      slogutil.Format
	LogLevel       slog.Level

	AdminHTTPAddr       netip.AddrPort
	AdminAuthConfig     string
	AdminAuthConfigFile string
	AdminAuthPolicy     string
	AdminAuthPolicyFile string

	ContestantHTTPAddr netip.AddrPort
	ContestantBaseURL  URLValue
}

// NewOption creates a new CLIOption combined with the given flag.FlagSet.
//
//nolint:varnamelen,mnd
func NewOption(fs *flag.FlagSet) *CLIOption {
	opt := CLIOption{}

	fs.DurationVar(&opt.GracefulPeriod, "graceful-period", 30*time.Second, "Graceful period for shutdown")
	fs.TextVar(&opt.LogFormat, "log-format", slogutil.FormatJSON, "Log format (json, console, pretty)")
	fs.TextVar(&opt.LogLevel, "log-level", slog.LevelInfo, "Log level")

	fs.TextVar(&opt.AdminHTTPAddr, "admin.http-addr", netip.MustParseAddrPort("127.0.0.1:8081"), "Admin HTTP server address")
	fs.StringVar(&opt.AdminAuthConfig, "admin.auth-config", "", "Admin API authentication config")
	fs.StringVar(&opt.AdminAuthConfigFile, "admin.auth-config-file", "", "Admin API authentication config file")
	fs.StringVar(&opt.AdminAuthPolicy, "admin.auth-policy", "", "Admin API authorization policy")
	fs.StringVar(&opt.AdminAuthPolicyFile, "admin.auth-policy-file", "", "Admin API authorization policy file")

	fs.TextVar(&opt.ContestantHTTPAddr, "contestant.http-addr", netip.MustParseAddrPort("127.0.0.1:8080"), "Contestant API server address")
	_ = opt.ContestantBaseURL.Set("http://localhost:8080")
	fs.Var(&opt.ContestantBaseURL, "contestant.base-url", "Contestant API base URL")

	fs.BoolFunc("v", "Verbose logging", func(string) error {
		opt.LogLevel = slog.LevelDebug
		return nil
	})

	fs.BoolFunc("dev", "Development mode", func(string) error {
		opt.LogFormat = slogutil.FormatPretty

		if opt.AdminAuthPolicy == "" && opt.AdminAuthPolicyFile == "" {
			opt.AdminAuthPolicy = "g, system:unauthenticated, role:admin"
		}

		return nil
	})

	return &opt
}

type URLValue url.URL

func (v *URLValue) Set(s string) error {
	url, err := url.Parse(s)
	if err != nil {
		return errors.WithStack(err)
	}
	*v = URLValue(*url)
	return nil
}

func (v *URLValue) String() string {
	return (*url.URL)(v).String()
}

func (v *URLValue) Value() *url.URL {
	return (*url.URL)(v)
}
