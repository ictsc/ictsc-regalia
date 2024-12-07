package main

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver"
)

func newConfig() (*scoreserver.Config, error) {
	var errs []error

	adminHTTPAddr, err := netip.ParseAddrPort(*flagAdminHTTPAddr)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid admin HTTP address"))
	}

	dsn, err := dsnFromEnv()
	if err != nil {
		errs = append(errs, err)
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

		DBURL: dsn,
	}, nil
}

func dsnFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "5432"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", errors.Wrap(err, "invalid DB_PORT")
	}
	if port <= 0 || port > 65535 {
		return "", errors.New("invalid DB_PORT")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "ictsc"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		if password = os.Getenv("PGPASSWORD"); password == "" {
			return "", errors.New("DB_PASSWORD or PGPASSWORD must be set")
		}
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "ictscore"
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s", url.QueryEscape(user), url.QueryEscape(password), net.JoinHostPort(host, strconv.Itoa(port)), name), nil
}
