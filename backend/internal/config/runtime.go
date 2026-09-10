package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Runtime struct {
	Address           string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
	DatabaseURL       string
	RedisURL          string
	SecureCookies     bool
	AllowedOrigins    []string
	DevFakeMode       bool

	DiscordClientID          string
	DiscordClientSecret      string
	ContestantRedirectURI    string
	AdminRedirectURI         string
	DiscordAdminGuildID      string
	DiscordContestantGuildID string
	DiscordAdminRoleIDs      []string

	GitHubAPIBaseURL   string
	GitHubAPIToken     string
	GitHubOIDCIssuer   string
	GitHubOIDCAudience string
	GitHubRepositoryID string
	GitHubRepository   string
	GitHubRef          string
	GitHubWorkflowRef  string
	GitHubManifestPath string
	GitHubDiscoveryURL string

	SStateBaseURL       string
	SStateRequestToken  string
	SStateCallbackToken string
	CallbackBaseURL     string
	AllowInsecureHTTP   bool
}

func LoadRuntime() (Runtime, error) {
	dev, err := boolEnv("ICTSC_DEV_FAKE_MODE", false)
	if err != nil {
		return Runtime{}, err
	}
	secure, err := boolEnv("ICTSC_SECURE_COOKIES", !dev)
	if err != nil {
		return Runtime{}, err
	}
	insecure, err := boolEnv("ICTSC_ALLOW_INSECURE_UPSTREAMS", false)
	if err != nil {
		return Runtime{}, err
	}
	if insecure && !dev {
		return Runtime{}, fmt.Errorf("ICTSC_ALLOW_INSECURE_UPSTREAMS is only allowed with ICTSC_DEV_FAKE_MODE")
	}
	cfg := Runtime{
		Address:                  envOrDefault("ICTSC_API_ADDRESS", defaultAddress),
		ReadHeaderTimeout:        defaultReadHeaderTimeout,
		ShutdownTimeout:          15 * time.Second,
		DatabaseURL:              strings.TrimSpace(os.Getenv("ICTSC_DATABASE_URL")),
		RedisURL:                 strings.TrimSpace(os.Getenv("ICTSC_REDIS_URL")),
		SecureCookies:            secure,
		AllowedOrigins:           csvEnv("ICTSC_ALLOWED_ORIGINS"),
		DevFakeMode:              dev,
		DiscordClientID:          strings.TrimSpace(os.Getenv("ICTSC_DISCORD_CLIENT_ID")),
		DiscordClientSecret:      strings.TrimSpace(os.Getenv("ICTSC_DISCORD_CLIENT_SECRET")),
		ContestantRedirectURI:    strings.TrimSpace(os.Getenv("ICTSC_DISCORD_CONTESTANT_REDIRECT_URI")),
		AdminRedirectURI:         strings.TrimSpace(os.Getenv("ICTSC_DISCORD_ADMIN_REDIRECT_URI")),
		DiscordAdminGuildID:      strings.TrimSpace(os.Getenv("ICTSC_DISCORD_ADMIN_GUILD_ID")),
		DiscordContestantGuildID: strings.TrimSpace(os.Getenv("ICTSC_DISCORD_CONTESTANT_GUILD_ID")),
		DiscordAdminRoleIDs:      csvEnv("ICTSC_DISCORD_ADMIN_ROLE_IDS"),
		GitHubAPIBaseURL:         envOrDefault("ICTSC_GITHUB_API_BASE_URL", "https://api.github.com"),
		GitHubAPIToken:           strings.TrimSpace(os.Getenv("ICTSC_GITHUB_API_TOKEN")),
		GitHubOIDCIssuer:         envOrDefault("ICTSC_GITHUB_OIDC_ISSUER", "https://token.actions.githubusercontent.com"),
		GitHubOIDCAudience:       strings.TrimSpace(os.Getenv("ICTSC_GITHUB_OIDC_AUDIENCE")),
		GitHubRepositoryID:       strings.TrimSpace(os.Getenv("ICTSC_GITHUB_REPOSITORY_ID")),
		GitHubRepository:         strings.TrimSpace(os.Getenv("ICTSC_GITHUB_REPOSITORY")),
		GitHubRef:                strings.TrimSpace(os.Getenv("ICTSC_GITHUB_REF")),
		GitHubWorkflowRef:        strings.TrimSpace(os.Getenv("ICTSC_GITHUB_WORKFLOW_REF")),
		GitHubManifestPath:       envOrDefault("ICTSC_GITHUB_MANIFEST_PATH", "content/manifest.yaml"),
		GitHubDiscoveryURL:       strings.TrimSpace(os.Getenv("ICTSC_GITHUB_OIDC_DISCOVERY_URL")),
		SStateBaseURL:            strings.TrimSpace(os.Getenv("ICTSC_SSTATE_BASE_URL")),
		SStateRequestToken:       strings.TrimSpace(os.Getenv("ICTSC_SSTATE_REQUEST_TOKEN")),
		SStateCallbackToken:      strings.TrimSpace(os.Getenv("ICTSC_SSTATE_CALLBACK_TOKEN")),
		CallbackBaseURL:          strings.TrimSpace(os.Getenv("ICTSC_CALLBACK_BASE_URL")),
		AllowInsecureHTTP:        insecure,
	}
	if raw := os.Getenv("ICTSC_API_READ_HEADER_TIMEOUT"); raw != "" {
		cfg.ReadHeaderTimeout, err = positiveDuration("ICTSC_API_READ_HEADER_TIMEOUT", raw)
		if err != nil {
			return Runtime{}, err
		}
	}
	if raw := os.Getenv("ICTSC_API_SHUTDOWN_TIMEOUT"); raw != "" {
		cfg.ShutdownTimeout, err = positiveDuration("ICTSC_API_SHUTDOWN_TIMEOUT", raw)
		if err != nil {
			return Runtime{}, err
		}
	}
	if err := cfg.Validate(); err != nil {
		return Runtime{}, err
	}
	return cfg, nil
}

func (c Runtime) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("ICTSC_API_ADDRESS is required")
	}
	if len(c.AllowedOrigins) == 0 {
		return fmt.Errorf("ICTSC_ALLOWED_ORIGINS must contain at least one origin")
	}
	for _, origin := range c.AllowedOrigins {
		if err := validateAbsoluteURL("allowed origin", origin, false, c.DevFakeMode); err != nil {
			return err
		}
	}
	if !c.DevFakeMode {
		if !c.SecureCookies {
			return fmt.Errorf("ICTSC_SECURE_COOKIES must be true outside development fake mode")
		}
		for name, value := range map[string]string{
			"ICTSC_DATABASE_URL": c.DatabaseURL, "ICTSC_REDIS_URL": c.RedisURL,
			"ICTSC_DISCORD_CLIENT_ID": c.DiscordClientID, "ICTSC_DISCORD_CLIENT_SECRET": c.DiscordClientSecret,
			"ICTSC_DISCORD_CONTESTANT_REDIRECT_URI": c.ContestantRedirectURI, "ICTSC_DISCORD_ADMIN_REDIRECT_URI": c.AdminRedirectURI,
			"ICTSC_DISCORD_ADMIN_GUILD_ID": c.DiscordAdminGuildID, "ICTSC_GITHUB_API_TOKEN": c.GitHubAPIToken,
			"ICTSC_GITHUB_OIDC_AUDIENCE": c.GitHubOIDCAudience, "ICTSC_GITHUB_REPOSITORY_ID": c.GitHubRepositoryID,
			"ICTSC_GITHUB_REPOSITORY": c.GitHubRepository, "ICTSC_GITHUB_REF": c.GitHubRef,
			"ICTSC_GITHUB_WORKFLOW_REF": c.GitHubWorkflowRef, "ICTSC_SSTATE_BASE_URL": c.SStateBaseURL,
			"ICTSC_SSTATE_REQUEST_TOKEN": c.SStateRequestToken, "ICTSC_SSTATE_CALLBACK_TOKEN": c.SStateCallbackToken,
			"ICTSC_CALLBACK_BASE_URL": c.CallbackBaseURL,
		} {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("%s is required outside development fake mode", name)
			}
		}
		if len(c.DiscordAdminRoleIDs) == 0 {
			return fmt.Errorf("ICTSC_DISCORD_ADMIN_ROLE_IDS must contain at least one role")
		}
	}
	for name, value := range map[string]string{
		"contestant redirect URI": c.ContestantRedirectURI, "admin redirect URI": c.AdminRedirectURI,
		"GitHub API base URL": c.GitHubAPIBaseURL, "GitHub OIDC issuer": c.GitHubOIDCIssuer,
		"SState base URL": c.SStateBaseURL, "callback base URL": c.CallbackBaseURL,
	} {
		if value != "" {
			if err := validateAbsoluteURL(name, value, true, c.DevFakeMode && c.AllowInsecureHTTP); err != nil {
				return err
			}
		}
	}
	if c.DatabaseURL != "" {
		if err := validateConnectionURL("ICTSC_DATABASE_URL", c.DatabaseURL, "postgres", "postgresql"); err != nil {
			return err
		}
	}
	if c.RedisURL != "" {
		if err := validateConnectionURL("ICTSC_REDIS_URL", c.RedisURL, "redis", "rediss"); err != nil {
			return err
		}
	}
	if c.GitHubDiscoveryURL != "" {
		if err := validateAbsoluteURL("GitHub discovery URL", c.GitHubDiscoveryURL, true, c.DevFakeMode && c.AllowInsecureHTTP); err != nil {
			return err
		}
	}
	return nil
}

func positiveDuration(name, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", name)
	}
	return duration, nil
}

func boolEnv(name string, fallback bool) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", name, err)
	}
	return value, nil
}

func csvEnv(name string) []string {
	raw := strings.Split(os.Getenv(name), ",")
	values := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		values = append(values, item)
	}
	return values
}

func validateConnectionURL(name, raw string, schemes ...string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("%s must be an absolute connection URL", name)
	}
	for _, scheme := range schemes {
		if parsed.Scheme == scheme {
			return nil
		}
	}
	return fmt.Errorf("%s uses unsupported scheme %q", name, parsed.Scheme)
}

func validateAbsoluteURL(name, raw string, allowPath, allowHTTP bool) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("%s must be an absolute URL without credentials, query, or fragment", name)
	}
	if parsed.Scheme != "https" && !(allowHTTP && parsed.Scheme == "http") {
		return fmt.Errorf("%s must use HTTPS", name)
	}
	if !allowPath && parsed.Path != "" && parsed.Path != "/" {
		return fmt.Errorf("%s may not contain a path", name)
	}
	return nil
}
