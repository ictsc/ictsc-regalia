package main

import (
	"context"
	"errors"
	"net/url"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

var errDevDependencyDisabled = errors.New("external dependency is disabled in development fake mode")

type disabledDiscord struct{}

func (disabledDiscord) AuthorizationURL(_, _, _ string, _ bool) string { return "/" }
func (disabledDiscord) Exchange(context.Context, string, string, string, bool) (service.DiscordResult, error) {
	return service.DiscordResult{}, errDevDependencyDisabled
}

type disabledContent struct{}

func (disabledContent) VerifyActionsToken(context.Context, string) (service.ActionsClaims, error) {
	return service.ActionsClaims{}, errDevDependencyDisabled
}
func (disabledContent) Fetch(context.Context, string, string, string) (core.ContentSnapshot, error) {
	return core.ContentSnapshot{}, errDevDependencyDisabled
}
func (disabledContent) IsAncestor(context.Context, string, string, string) (bool, error) {
	return false, errDevDependencyDisabled
}

type disabledDeployments struct{}

func (disabledDeployments) Queue(context.Context, service.DeploymentRequest) error {
	return errDevDependencyDisabled
}
func (disabledDeployments) Status(context.Context, int64, string) (core.DeploymentEvent, error) {
	return core.DeploymentEvent{}, errDevDependencyDisabled
}

// previewDiscord is used only in fake mode without Discord credentials.
// It requires an authenticated external gateway when exposed publicly.
type previewDiscord struct {
	guildID string
	roleIDs []string
}

func (d previewDiscord) AuthorizationURL(state, _, redirectURI string, _ bool) string {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return "/"
	}
	q := u.Query()
	q.Set("state", state)
	q.Set("code", "development-preview")
	u.RawQuery = q.Encode()
	return u.String()
}
func (d previewDiscord) Exchange(_ context.Context, code, _, _ string, admin bool) (service.DiscordResult, error) {
	if code != "development-preview" {
		return service.DiscordResult{}, errDevDependencyDisabled
	}
	identity := core.DiscordIdentity{ID: "900000000000000001", Username: "preview-contestant", DisplayName: "Preview Contestant"}
	if admin {
		identity = core.DiscordIdentity{ID: "900000000000000002", Username: "preview-admin", DisplayName: "Preview Admin"}
	}
	return service.DiscordResult{Identity: identity, GuildID: d.guildID, RoleIDs: d.roleIDs}, nil
}
