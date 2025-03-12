package config

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	AdminAPI      AdminAPI
	ContestantAPI ContestantAPI

	PgConfig     pgx.ConnConfig
	Redis        redis.Options
	FakeSchedule *FakeSchedule
}

type (
	AdminAPI struct {
		Address netip.AddrPort
		Authn   AdminAuthn
		Authz   AdminAuthz
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

type FakeSchedule struct {
	Phase     domain.Phase
	NextPhase domain.Phase
	StartAt   time.Time
	EndAt     time.Time
}
