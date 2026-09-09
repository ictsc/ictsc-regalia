package config

import (
	"strings"
	"testing"
)

var runtimeEnvKeys = []string{
	"ICTSC_DEV_FAKE_MODE", "ICTSC_SECURE_COOKIES", "ICTSC_ALLOW_INSECURE_UPSTREAMS",
	"ICTSC_ALLOWED_ORIGINS", "ICTSC_DATABASE_URL", "ICTSC_REDIS_URL",
	"ICTSC_DISCORD_CLIENT_ID", "ICTSC_DISCORD_CLIENT_SECRET", "ICTSC_DISCORD_CONTESTANT_REDIRECT_URI",
	"ICTSC_DISCORD_ADMIN_REDIRECT_URI", "ICTSC_DISCORD_ADMIN_GUILD_ID", "ICTSC_DISCORD_ADMIN_ROLE_IDS",
	"ICTSC_GITHUB_API_TOKEN", "ICTSC_GITHUB_OIDC_AUDIENCE", "ICTSC_GITHUB_REPOSITORY_ID",
	"ICTSC_GITHUB_REPOSITORY", "ICTSC_GITHUB_REF", "ICTSC_GITHUB_WORKFLOW_REF",
	"ICTSC_SSTATE_BASE_URL", "ICTSC_SSTATE_REQUEST_TOKEN", "ICTSC_SSTATE_CALLBACK_TOKEN",
	"ICTSC_CALLBACK_BASE_URL", "ICTSC_API_READ_HEADER_TIMEOUT", "ICTSC_API_SHUTDOWN_TIMEOUT",
}

func clearRuntimeEnv(t *testing.T) {
	t.Helper()
	for _, key := range runtimeEnvKeys {
		t.Setenv(key, "")
	}
}

func setValidProductionRuntimeEnv(t *testing.T) {
	t.Helper()
	clearRuntimeEnv(t)
	values := map[string]string{
		"ICTSC_ALLOWED_ORIGINS":                 "https://score.example,https://admin.example",
		"ICTSC_DATABASE_URL":                    "postgres://ictsc:secret@postgres:5432/ictscore",
		"ICTSC_REDIS_URL":                       "redis://redis:6379/0",
		"ICTSC_DISCORD_CLIENT_ID":               "123456789012345678",
		"ICTSC_DISCORD_CLIENT_SECRET":           "discord-secret",
		"ICTSC_DISCORD_CONTESTANT_REDIRECT_URI": "https://score.example/api/v1/auth/discord/callback",
		"ICTSC_DISCORD_ADMIN_REDIRECT_URI":      "https://score.example/api/v1/admin/auth/discord/callback",
		"ICTSC_DISCORD_ADMIN_GUILD_ID":          "223456789012345678",
		"ICTSC_DISCORD_ADMIN_ROLE_IDS":          "323456789012345678,423456789012345678",
		"ICTSC_GITHUB_API_TOKEN":                "github-token",
		"ICTSC_GITHUB_OIDC_AUDIENCE":            "ictsc-score",
		"ICTSC_GITHUB_REPOSITORY_ID":            "123456789",
		"ICTSC_GITHUB_REPOSITORY":               "ictsc/content",
		"ICTSC_GITHUB_REF":                      "refs/heads/main",
		"ICTSC_GITHUB_WORKFLOW_REF":             "ictsc/content/.github/workflows/publish.yml@refs/heads/main",
		"ICTSC_SSTATE_BASE_URL":                 "https://sstate.example",
		"ICTSC_SSTATE_REQUEST_TOKEN":            "request-token",
		"ICTSC_SSTATE_CALLBACK_TOKEN":           "callback-token",
		"ICTSC_CALLBACK_BASE_URL":               "https://score.example",
	}
	for key, value := range values {
		t.Setenv(key, value)
	}
}

func TestLoadRuntimeProductionValidatesAndDefaultsSecureCookies(t *testing.T) {
	setValidProductionRuntimeEnv(t)
	cfg, err := LoadRuntime()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.SecureCookies || cfg.DevFakeMode {
		t.Fatalf("unexpected security mode: %#v", cfg)
	}
	if len(cfg.AllowedOrigins) != 2 || len(cfg.DiscordAdminRoleIDs) != 2 {
		t.Fatalf("CSV settings were not parsed: %#v", cfg)
	}
}

func TestLoadRuntimeRejectsMissingProductionSecret(t *testing.T) {
	setValidProductionRuntimeEnv(t)
	t.Setenv("ICTSC_SSTATE_CALLBACK_TOKEN", "")
	_, err := LoadRuntime()
	if err == nil || !strings.Contains(err.Error(), "ICTSC_SSTATE_CALLBACK_TOKEN") {
		t.Fatalf("error = %v", err)
	}
}

func TestLoadRuntimeDevelopmentFakeAllowsLocalHTTP(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("ICTSC_DEV_FAKE_MODE", "true")
	t.Setenv("ICTSC_ALLOWED_ORIGINS", "http://localhost:3000")
	cfg, err := LoadRuntime()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SecureCookies || !cfg.DevFakeMode {
		t.Fatalf("unexpected fake mode: %#v", cfg)
	}
}

func TestLoadRuntimeRejectsInsecureUpstreamsOutsideDev(t *testing.T) {
	setValidProductionRuntimeEnv(t)
	t.Setenv("ICTSC_ALLOW_INSECURE_UPSTREAMS", "true")
	if _, err := LoadRuntime(); err == nil {
		t.Fatal("expected insecure upstream validation error")
	}
}

func TestLoadRuntimeRejectsOriginPathAndBadDuration(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "origin path", key: "ICTSC_ALLOWED_ORIGINS", value: "https://score.example/path"},
		{name: "duration", key: "ICTSC_API_READ_HEADER_TIMEOUT", value: "0s"},
		{name: "database scheme", key: "ICTSC_DATABASE_URL", value: "mysql://db/ictscore"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidProductionRuntimeEnv(t)
			t.Setenv(test.key, test.value)
			if _, err := LoadRuntime(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
