package main

import (
	"net/netip"
	"net/url"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver"
)

func newConfig() (*scoreserver.Config, error) {
	var errs []error

	adminHTTPAddr, err := netip.ParseAddrPort(*flagAdminHTTPAddr)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid admin HTTP address"))
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		errs = append(errs, errors.New("DB_DSN must be set"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &scoreserver.Config{
		AdminAPI: scoreserver.AdminAPIConfig{
			Address: adminHTTPAddr,
		},

		ContestantHTTPAddress: netip.AddrPort{},
		ContestantBaseURLs:    []url.URL{},

		DBDSN: dsn,
	}, nil
}
