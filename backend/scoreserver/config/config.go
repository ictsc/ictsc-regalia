package config

import (
	"net/netip"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	AdminAPI AdminAPIConfig

	ContestantHTTPAddress netip.AddrPort
	ContestantBaseURLs    []url.URL

	PgConfig pgx.ConnConfig
}

type AdminAPIConfig struct {
	Address netip.AddrPort
	Authn   AdminAuthnConfig
}

type AdminAuthnConfig struct {
	Issuers []IssuerConfig `yaml:"issuers"`
}

type IssuerConfig struct {
	Issuer            string `yaml:"issuer"`
	InsecureIssuerURL string `yaml:"insecure_issuer_url"`
	ClientID          string `yaml:"client_id"`
	CAFile            string `yaml:"ca_file"`

	UsernameKey string   `yaml:"username_key"`
	GroupKeys   []string `yaml:"group_keys"`
}
