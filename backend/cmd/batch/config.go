package main

import (
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"golang.org/x/oauth2"
)

func newConfig(opts *CLIOption) (*config.Batch, error) {
	var tokenSource oauth2.TokenSource
	if opts.APITokenFile != "" {
		tokenSource = oauth2.ReuseTokenSource(nil, &fileTokenSource{Path: opts.APITokenFile})
	} else if token := os.Getenv("ICTSCORE_TOKEN"); token != "" {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	}

	var errs []error
	sstateURL := os.Getenv("SSTATE_URL")
	if sstateURL == "" {
		errs = append(errs, errors.New("SSTATE_URL is required"))
	}
	sstateUser := os.Getenv("SSTATE_USER")
	if sstateUser == "" {
		errs = append(errs, errors.New("SSTATE_USER is required"))
	}
	sstatePassword := os.Getenv("SSTATE_PASSWORD")
	if sstatePassword == "" {
		errs = append(errs, errors.New("SSTATE_PASSWORD is required"))
	}
	var sstateCA string
	if opts.SStateCAFile != "" {
		data, err := os.ReadFile(opts.SStateCAFile)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to read sstate CA file"))
		}
		sstateCA = string(data)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &config.Batch{
		APIURL:         opts.APIURL,
		APITokenSource: tokenSource,

		DeploymentSync: config.DeploySync{
			Period: opts.DeploymentSyncPeriod,
			SState: config.SState{
				URL:                sstateURL,
				CA:                 sstateCA,
				InsecureSkipVerify: opts.SStateTLSSkipVerify,
				User:               sstateUser,
				Password:           sstatePassword,
			},
		},
		ScoreUpdate: config.ScoreUpdate{
			Period: opts.ScoreUpdatePeriod,
		},
	}, nil
}

type fileTokenSource struct {
	Path string
}

var _ oauth2.TokenSource = (*fileTokenSource)(nil)

const tokenExpiry = 5 * time.Second

func (f *fileTokenSource) Token() (*oauth2.Token, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token file")
	}
	token := strings.TrimSpace(string(data))
	return &oauth2.Token{AccessToken: token, Expiry: time.Now().Add(tokenExpiry)}, nil
}
