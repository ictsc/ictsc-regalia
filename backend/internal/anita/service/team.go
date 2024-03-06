package service

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/pkg/optional"
)

// UpdateTeamArgs チーム更新引数
type UpdateTeamArgs struct {
	ID      value.TeamID
	Code    optional.Of[value.TeamCode]
	Name    optional.Of[value.TeamName]
	Org     optional.Of[value.TeamOrganization]
	Bastion optional.Of[value.Bastion]
}

// TeamService チームサービス
type TeamService interface {
	ReadTeam(ctx context.Context, id value.TeamID) (*domain.Team, error)
	ReadTeams(ctx context.Context) ([]*domain.Team, error)
	CreateTeam(ctx context.Context, code value.TeamCode, name value.TeamName, org value.TeamOrganization) (*domain.Team, error)
	UpdateTeam(ctx context.Context, args UpdateTeamArgs) (*domain.Team, error)
	DeleteTeam(ctx context.Context, id value.TeamID) error
	MoveMember(ctx context.Context, to value.TeamID, memberID value.UserID) error
}
