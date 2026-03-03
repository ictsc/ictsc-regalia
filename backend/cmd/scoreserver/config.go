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
	var errs []error

	adminAPI, err := newAdminConfig(opts)
	if err != nil {
		errs = append(errs, err)
	}

	contestantAPI, err := newContestantConfig(opts)
	if err != nil {
		errs = append(errs, err)
	}

	cfg, err := pgx.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid DB_DSN"))
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid REDIS_URL"))
	}

	var fakeSchedule *config.FakeSchedule
	if opts.UseFakeSchedule {
		fakeSchedule = &config.FakeSchedule{
			Name:    opts.FakeScheduleName,
			StartAt: opts.FakeScheduleStartAt,
			EndAt:   opts.FakeScheduleEndAt,
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &config.Config{
		AdminAPI:      *adminAPI,
		ContestantAPI: *contestantAPI,
		PgConfig:      *cfg,
		Redis:         *redisOpts,
		FakeSchedule:  fakeSchedule,
	}, nil
}
func newAdminConfig(opts *CLIOption) (*config.AdminAPI, error) {
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

func newContestantConfig(opts *CLIOption) (*config.ContestantAPI, error) {
	externalURL := &opts.ContestantExternalURL
	if opts.ContestantRoutePrefix != "" {
		if externalURL.Path == "" {
			externalURL.Path = "/"
		}
		externalURL = externalURL.JoinPath(opts.ContestantRoutePrefix)
	}

	var errs []error
	discordClientID := os.Getenv("DISCORD_CLIENT_ID")
	if discordClientID == "" {
		errs = append(errs, errors.New("DISCORD_CLIENT_ID is not set"))
	}
	discordClientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	if discordClientSecret == "" {
		errs = append(errs, errors.New("DISCORD_CLIENT_SECRET is not set"))
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &config.ContestantAPI{
		Address: opts.ContestantHTTPAddr,
		Auth: config.ContestantAuth{
			ExternalURL:         externalURL,
			DiscordClientID:     discordClientID,
			DiscordClientSecret: discordClientSecret,
		},
	}, nil
}
