package service

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/pkg/optional"
)

// CreateTeamArgs チーム作成引数
type CreateTeamArgs struct {
	Code          value.TeamCode
	Name          value.TeamName
	Org           value.TeamOrganization
	CodeRemaining value.TeamCodeRemaining
}

// UpdateTeamArgs チーム更新引数
type UpdateTeamArgs struct {
	ID            value.TeamID
	Code          optional.Of[value.TeamCode]
	Name          optional.Of[value.TeamName]
	Org           optional.Of[value.TeamOrganization]
	CodeRemaining optional.Of[value.TeamCodeRemaining]
	Bastion       optional.Of[value.Bastion]
}

// TeamService チームサービス
type TeamService interface {
	ReadTeam(ctx context.Context, id value.TeamID) (*domain.Team, error)
	ReadTeams(ctx context.Context) ([]*domain.Team, error)
	ReadTeamByInvitationCode(ctx context.Context, code value.TeamInvitationCode) (*domain.Team, error)
	CreateTeam(ctx context.Context, args CreateTeamArgs) (*domain.Team, error)
	UpdateTeam(ctx context.Context, args UpdateTeamArgs) (*domain.Team, error)
	DeleteTeam(ctx context.Context, id value.TeamID) error
	MoveMember(ctx context.Context, to value.TeamID, memberID value.UserID) error
}
