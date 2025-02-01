package domain_test

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type teamCodeGetterFunc func(ctx context.Context, code domain.TeamCode) (*domain.Team, error)

func (f teamCodeGetterFunc) GetTeamByCode(ctx context.Context, code domain.TeamCode) (*domain.Team, error) {
	return f(ctx, code)
}
