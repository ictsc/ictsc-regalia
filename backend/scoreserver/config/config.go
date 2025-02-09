package config

import (
	"net/netip"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	AdminAPI      AdminAPI
	ContestantAPI ContestantAPI

	PgConfig pgx.ConnConfig
}

type AdminAPI struct {
	Address netip.AddrPort
	Authn   AdminAuthn
	Authz   AdminAuthz
}

type AdminAuthn struct {
	Issuers []Issuer `yaml:"issuers"`
}

type Issuer struct {
	Issuer            string `yaml:"issuer"`
	InsecureIssuerURL string `yaml:"insecure_issuer_url"`
	ClientID          string `yaml:"client_id"`
	CAFile            string `yaml:"ca_file"`

	UsernameKey string   `yaml:"username_key"`
	GroupKeys   []string `yaml:"group_keys"`
}

type AdminAuthz struct {
	Policy string
}

type ContestantAPI struct {
	Address netip.AddrPort
}
