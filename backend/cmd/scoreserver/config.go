package main

import (
	"net/netip"
	"net/url"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver"
	"github.com/jackc/pgx/v5"
)

func newConfig() (*scoreserver.Config, error) {
	adminHTTPAddr, err := netip.ParseAddrPort(*flagAdminHTTPAddr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid admin HTTP address")
	}

	cfg, err := pgx.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB_DSN")
	}

	return &scoreserver.Config{
		AdminAPI: scoreserver.AdminAPIConfig{
			Address: adminHTTPAddr,
		},

		ContestantHTTPAddress: netip.AddrPort{},
		ContestantBaseURLs:    []url.URL{},

		PgConfig: *cfg,
	}, nil
}
