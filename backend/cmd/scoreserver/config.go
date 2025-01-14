package main

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

func newConfig(opts *CLIOption) (*config.Config, error) {
	var err error

	cfg, err := pgx.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB_DSN")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse REDIS_URL")
	}

	adminCfg, err := parseAdminConfig(opts)
	if err != nil {
		return nil, err
	}

	contestantCfg, err := parseContestantConfig(opts)
	if err != nil {
		return nil, err
	}

	return &config.Config{
		PgConfig:      *cfg,
		Redis:         *redisOptions,
		AdminAPI:      *adminCfg,
		ContestantAPI: *contestantCfg,
	}, nil
}

func parseAdminConfig(opts *CLIOption) (*config.AdminAPI, error) {
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

	return &config.AdminAPI{
		Address: opts.AdminHTTPAddr,
		Authn:   adminAuthnConfig,
		Authz: config.AdminAuthz{
			Policy: adminAuthPolicy,
		},
	}, nil
}

func parseContestantConfig(opts *CLIOption) (*config.ContestantAPI, error) {
	return &config.ContestantAPI{
		Address: opts.ContestantHTTPAddr,
		Auth: config.ContestantAuth{
			BaseURL: *opts.ContestantBaseURL.Value(),
		},
	}, nil
}
