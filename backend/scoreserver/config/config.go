package config

import (
	"net/netip"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	AdminAPI      AdminAPI
	ContestantAPI ContestantAPI

	PgConfig pgx.ConnConfig
	Redis    redis.Options
}

type (
	AdminAPI struct {
		Address netip.AddrPort
		Authn   AdminAuthn
		Authz   AdminAuthz
		Growi   Growi
	}
	AdminAuthn struct {
		Issuers []Issuer `yaml:"issuers"`
	}
	Issuer struct {
		Issuer            string `yaml:"issuer"`
		InsecureIssuerURL string `yaml:"insecure_issuer_url"`
		ClientID          string `yaml:"client_id"`
		CAFile            string `yaml:"ca_file"`

		UsernameKey string   `yaml:"username_key"`
		GroupKeys   []string `yaml:"group_keys"`
	}
	AdminAuthz struct {
		Policy string
	}
	Growi struct {
		BaseURL *url.URL
		Token   string
	}
)

type (
	ContestantAPI struct {
		Address        netip.AddrPort
		Auth           ContestantAuth
		AllowedOrigins []string
	}
	ContestantAuth struct {
		BaseURL             *url.URL
		DiscordClientID     string
		DiscordClientSecret string
	}
)
