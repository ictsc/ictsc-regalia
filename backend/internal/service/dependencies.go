package service

import (
	"context"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type DiscordResult struct {
	Identity core.DiscordIdentity
	GuildID  string
	RoleIDs  []string
}

type Discord interface {
	AuthorizationURL(state, codeChallenge, redirectURI string, admin bool) string
	Exchange(ctx context.Context, code, codeVerifier, redirectURI string, admin bool) (DiscordResult, error)
}

type ActionsClaims struct {
	Issuer       string
	Audience     string
	RepositoryID string
	Repository   string
	Ref          string
	WorkflowRef  string
	Subject      string
}

type ContentSource interface {
	VerifyActionsToken(ctx context.Context, token string) (ActionsClaims, error)
	Fetch(ctx context.Context, repository, ref, commit string) (core.ContentSnapshot, error)
	IsAncestor(ctx context.Context, repository, ancestor, descendant string) (bool, error)
}

type DeploymentRequest struct {
	RequestID     string
	TeamCode      int64
	ProblemCode   string
	Revision      int32
	ContentCommit string
	CallbackURL   string
}

type DeploymentGateway interface {
	Queue(ctx context.Context, request DeploymentRequest) error
	Status(ctx context.Context, teamCode int64, problemCode string) (core.DeploymentEvent, error)
}

type Subscription interface {
	Messages() <-chan []byte
	Close() error
}

type EventBus interface {
	Publish(ctx context.Context, event []byte) error
	Subscribe(ctx context.Context) (Subscription, error)
}

type Clock func() time.Time
