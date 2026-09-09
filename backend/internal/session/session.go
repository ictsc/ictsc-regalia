package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type Kind string

const (
	KindOAuth      Kind = "oauth"
	KindAdminOAuth Kind = "admin-oauth"
	KindSignup     Kind = "signup"
	KindContestant Kind = "contestant"
	KindAdmin      Kind = "admin"
)

type Data struct {
	Kind           Kind                 `json:"kind"`
	OAuthState     string               `json:"oauth_state,omitempty"`
	PKCEVerifier   string               `json:"pkce_verifier,omitempty"`
	Next           string               `json:"next,omitempty"`
	GuildID        string               `json:"guild_id,omitempty"`
	RoleIDs        []string             `json:"role_ids,omitempty"`
	Discord        core.DiscordIdentity `json:"discord,omitempty"`
	ContestantName string               `json:"contestant_name,omitempty"`
	AdminName      string               `json:"admin_name,omitempty"`
	ImpersonatedBy string               `json:"impersonated_by,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

var ErrNotFound = errors.New("session not found")

type Store interface {
	Create(ctx context.Context, data Data, ttl time.Duration) (string, error)
	Get(ctx context.Context, token string, expected Kind) (Data, error)
	Consume(ctx context.Context, token string, expected Kind) (Data, error)
	Delete(ctx context.Context, token string) error
}

func Token() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}
