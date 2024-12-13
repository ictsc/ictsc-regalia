package scoreserver

import (
	"net/netip"
	"net/url"
)

type Config struct {
	AdminAPI AdminAPIConfig

	ContestantHTTPAddress netip.AddrPort
	ContestantBaseURLs    []url.URL

	DBDSN string
}

type AdminAPIConfig struct {
	Address netip.AddrPort
}
