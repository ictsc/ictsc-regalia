package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type (
	DiscordUserID   string
	DiscordIdentity struct {
		ID         DiscordUserID
		Username   string
		GlobalName string
	}
)

func GetDiscordIdentity(ctx context.Context, eff DiscordIdentityGetter, token string) (*DiscordIdentity, error) {
	user, err := eff.GetIdentity(ctx, token)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get discord user")
	}
	return user, nil
}

func (u *User) LinkDiscord(ctx context.Context, eff DiscordUserLinker, id DiscordUserID) error {
	if err := eff.LinkDiscordUser(ctx, u.id, string(id)); err != nil {
		return WrapAsInternal(err, "failed to link discord user")
	}
	return nil
}

func (id DiscordUserID) User(ctx context.Context, eff DiscordLinkedUserGetter) (*UserData, error) {
	user, err := eff.GetDiscordLinkedUser(ctx, string(id))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get linked user")
	}
	return user, nil
}

type (
	DiscordIdentityGetter interface {
		GetIdentity(ctx context.Context, token string) (*DiscordIdentity, error)
	}
	DiscordUserLinker interface {
		LinkDiscordUser(ctx context.Context, userID uuid.UUID, discordUserID string) error
	}
	DiscordLinkedUserGetter interface {
		GetDiscordLinkedUser(ctx context.Context, discordUserID string) (*UserData, error)
	}
)
