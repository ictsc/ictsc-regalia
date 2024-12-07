package scoreserver

import (
	"fmt"
	"net/netip"
	"net/url"
)

type Config struct {
	AdminAPI AdminAPIConfig

	ContestantHTTPAddress netip.AddrPort
	ContestantBaseURLs    []url.URL

	DBURL string
}

type AdminAPIConfig struct {
	Address        netip.AddrPort
}

type DBConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Name     string
}

func (dbc *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", url.QueryEscape(dbc.User), url.QueryEscape(dbc.Password), dbc.Host, dbc.Name)
}
