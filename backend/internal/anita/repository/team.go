package repository

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
)

// TeamRepository チームリポジトリ
type TeamRepository interface {
	SelectTeam(ctx context.Context, id value.TeamID) (*domain.Team, error)
	SelectTeams(ctx context.Context) ([]*domain.Team, error)
	SelectTeamByInvitationCode(ctx context.Context, code value.TeamInvitationCode) (*domain.Team, error)
	UpsertTeam(ctx context.Context, team *domain.Team) error // Bastionまでは挿入・更新するが、Memberの挿入・更新はしない
	DeleteTeam(ctx context.Context, id value.TeamID) error
}
