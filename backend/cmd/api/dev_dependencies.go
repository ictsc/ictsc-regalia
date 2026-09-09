package main

import (
	"context"
	"errors"

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
