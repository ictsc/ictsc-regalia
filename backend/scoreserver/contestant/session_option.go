package contestant

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

func externalURL(base *url.URL, r *http.Request) *url.URL {
	url := &url.URL{}

	if base != nil {
		if url.Scheme == "" && base.Scheme != "" {
			url.Scheme = base.Scheme
		}
		if url.Host == "" && base.Host != "" {
			url.Host = base.Host
		}
		if url.Path == "" && base.Path != "" {
			url.Path = base.Path
		}
	}

	if url.Scheme == "" && strings.ToLower(r.Header.Get("X-Forwarded-Proto")) == "https" {
		url.Scheme = "https"
	}
	if url.Scheme == "" {
		if r.TLS != nil {
			url.Scheme = "https"
		} else {
			url.Scheme = "http"
		}
	}
	if url.Host == "" {
		url.Host = r.Host
	}
	if url.Path == "" {
		url.Path = "/"
	}

	return url
}

func oauth2SessionOption(base *url.URL, r *http.Request) *sessions.Options {
	opt := sessionOption(base, r)
	opt.MaxAge = int(authCookieAge.Seconds())
	opt.Path = externalURL(base, r).JoinPath("./auth").Path
	opt.SameSite = http.SameSiteLaxMode
	return opt
}

func userSessionOption(base *url.URL, r *http.Request) *sessions.Options {
	opt := sessionOption(base, r)
	opt.MaxAge = int(userCookieAge.Seconds())
	return opt
}

func signUpSessionOption(base *url.URL, r *http.Request) *sessions.Options {
	opt := sessionOption(base, r)
	opt.MaxAge = int(signUpCookieAge.Seconds())
	return opt
}

func sessionOption(base *url.URL, r *http.Request) *sessions.Options {
	externalURL := externalURL(base, r)
	path := externalURL.Path
	if path == "" {
		path = "/"
	}
	return &sessions.Options{
		Path:     path,
		Secure:   externalURL.Scheme == "https",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

const (
	authCookieAge   = 10 * time.Minute
	signUpCookieAge = 10 * time.Minute
	userCookieAge   = 3 * 24 * time.Hour
)
