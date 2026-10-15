package discord

import "net/http"

type guildMembershipError struct{ cause error }

func (e *guildMembershipError) Error() string {
	return ErrGuildMembership.Error() + ": " + e.cause.Error()
}
func (e *guildMembershipError) Unwrap() []error        { return []error{ErrGuildMembership, e.cause} }
func (e *guildMembershipError) OAuthErrorKind() string { return "guild_membership" }

func (e *HTTPError) OAuthErrorKind() string {
	if e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError && e.StatusCode != http.StatusTooManyRequests {
		return "authentication_rejected"
	}
	return "upstream_unavailable"
}
