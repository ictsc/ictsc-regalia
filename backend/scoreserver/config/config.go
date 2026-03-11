package config

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type Config struct {
	AdminAPI      AdminAPI
	ContestantAPI ContestantAPI

	PgConfig     pgx.ConnConfig
	Redis        redis.Options
	FakeSchedule *FakeSchedule
}

// FakeSchedule はDBを使わずにスケジュールをシミュレートするための設定。
// 指定した名前・期間のスケジュールが常に返される。
type FakeSchedule struct {
	Name    string
	StartAt time.Time
	EndAt   time.Time
}

type (
	Batch struct {
		APIURL         string
		APITokenSource oauth2.TokenSource

		DeploymentSync *DeploySync
		ScoreUpdate    *ScoreUpdate
	}
	DeploySync struct {
		Period time.Duration
		SState SState
	}
	SState struct {
		URL                string
		CA                 string
		InsecureSkipVerify bool
		User               string
		Password           string
	}
	ScoreUpdate struct {
		Period time.Duration
	}
)

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
		Address netip.AddrPort
		Auth    ContestantAuth
	}
	ContestantAuth struct {
		ExternalURL *url.URL

		DiscordClientID     string
		DiscordClientSecret string
	}
)

