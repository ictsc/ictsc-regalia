package scoreserver

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
}
