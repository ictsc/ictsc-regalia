package main

import (
	"net/netip"
	"net/url"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
)

func newConfig(opts *CLIOption) (*config.Config, error) {
	var err error

	adminAuthConfigData := []byte(opts.AdminAuthConfig)
	if opts.AdminAuthConfigFile != "" {
		if len(adminAuthConfigData) != 0 {
			return nil, errors.New("both admin-auth-config and admin-auth-config-file are specified")
		}
		adminAuthConfigData, err = os.ReadFile(opts.AdminAuthConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read admin-auth-config")
		}
	}
	var adminAuthnConfig config.AdminAuthn
	if err := yaml.Unmarshal(adminAuthConfigData, &adminAuthnConfig); err != nil {
		return nil, errors.Wrap(err, "failed to parse admin-auth-config")
	}

	cfg, err := pgx.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB_DSN")
	}

	return &config.Config{
		AdminAPI: config.AdminAPIConfig{
			Address: opts.AdminHTTPAddr,
			Authn:   adminAuthnConfig,
		},

		ContestantHTTPAddress: netip.AddrPort{},
		ContestantBaseURLs:    []url.URL{},

		PgConfig: *cfg,
	}, nil
}
