package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var _ domain.DiscordIdentityGetter = (*UserClient)(nil)

func (c *UserClient) GetIdentity(ctx context.Context, token string) (*domain.DiscordIdentityData, error) {
	session, err := discordgo.New("Bearer " + token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discord session")
	}
	if c.HTTPClient != nil {
		session.Client = c.HTTPClient
	}

	user, err := session.User("@me", discordgo.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &domain.DiscordIdentityData{
		ID:         user.ID,
		Username:   user.Username,
		GlobalName: user.GlobalName,
	}, nil
}
