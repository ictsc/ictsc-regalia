package service

import "errors"

func discordOAuthErrorKind(err error) string {
	var classified interface{ OAuthErrorKind() string }
	if errors.As(err, &classified) {
		return classified.OAuthErrorKind()
	}
	return "upstream_unavailable"
}
