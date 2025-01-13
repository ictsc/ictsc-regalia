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
			return nil, errors.New("both admin.auth-config and admin.auth-config-file are specified")
		}
		adminAuthConfigData, err = os.ReadFile(opts.AdminAuthConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read admin.auth-config-file")
		}
	}
	var adminAuthnConfig config.AdminAuthn
	if err := yaml.Unmarshal(adminAuthConfigData, &adminAuthnConfig); err != nil {
		return nil, errors.Wrap(err, "failed to parse admin.auth-config")
	}

	adminAuthPolicy := opts.AdminAuthPolicy
	if opts.AdminAuthPolicyFile != "" {
		if adminAuthPolicy != "" {
			return nil, errors.New("both admin.auth-policy and admin.auth-policy-file are specified")
		}
		data, err := os.ReadFile(opts.AdminAuthPolicyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read admin.auth-policy-file")
		}
		adminAuthPolicy = string(data)
	}

	cfg, err := pgx.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB_DSN")
	}

	return &config.Config{
		AdminAPI: config.AdminAPI{
			Address: opts.AdminHTTPAddr,
			Authn:   adminAuthnConfig,
			Authz: config.AdminAuthz{
				Policy: adminAuthPolicy,
			},
		},

		ContestantHTTPAddress: netip.AddrPort{},
		ContestantBaseURLs:    []url.URL{},

		PgConfig: *cfg,
	}, nil
}
