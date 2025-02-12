package domain

import (
	"context"
	"strconv"

	"github.com/gofrs/uuid/v5"
)

type (
	DiscordUserID   int64
	DiscordIdentity struct {
		id         DiscordUserID
		username   string
		globalName string
	}
	DiscordIdentityData struct {
		ID         string
		Username   string
		GlobalName string
	}
)

func (d *DiscordIdentity) ID() DiscordUserID {
	return d.id
}

func (d *DiscordIdentity) Username() string {
	return d.username
}

func (d *DiscordIdentity) GlobalName() string {
	return d.globalName
}

func (d *DiscordIdentity) Data() *DiscordIdentityData {
	return &DiscordIdentityData{
		ID:         strconv.FormatInt(int64(d.id), 10),
		Username:   d.username,
		GlobalName: d.globalName,
	}
}

func GetDiscordIdentity(ctx context.Context, eff DiscordIdentityGetter, token string) (*DiscordIdentity, error) {
	userData, err := eff.GetIdentity(ctx, token)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get discord user")
	}
	return userData.parse()
}

func (u *User) LinkDiscord(ctx context.Context, eff DiscordUserLinker, id DiscordUserID) error {
	if err := eff.LinkDiscordUser(ctx, uuid.UUID(u.userID), int64(id)); err != nil {
		return WrapAsInternal(err, "failed to link discord user")
	}
	return nil
}

func (id DiscordUserID) User(ctx context.Context, eff DiscordLinkedUserGetter) (*User, error) {
	userData, err := eff.GetDiscordLinkedUser(ctx, int64(id))
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get linked user")
	}
	user, err := userData.parse()
	if err != nil {
		return nil, err
	}
	return user, nil
}

type (
	DiscordIdentityGetter interface {
		GetIdentity(ctx context.Context, token string) (*DiscordIdentityData, error)
	}
	DiscordUserLinker interface {
		LinkDiscordUser(ctx context.Context, userID uuid.UUID, discordUserID int64) error
	}
	DiscordLinkedUserGetter interface {
		GetDiscordLinkedUser(ctx context.Context, discordUserID int64) (*UserData, error)
	}
)

func NewDiscordID(id string) (DiscordUserID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, NewInvalidArgumentError("invalid discord user id", err)
	}
	return DiscordUserID(i), nil
}

func (d *DiscordIdentityData) parse() (*DiscordIdentity, error) {
	id, err := NewDiscordID(d.ID)
	if err != nil {
		return nil, err
	}

	return &DiscordIdentity{
		id:         id,
		username:   d.Username,
		globalName: d.GlobalName,
	}, nil
}
