package postgres

import (
	"context"
	"errors"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/jackc/pgx/v5"
)

func (s *Store) updateTeamTransactional(ctx context.Context, code int64, patch core.TeamPatch) (core.Team, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return core.Team{}, dbError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var team core.Team
	err = tx.QueryRow(ctx, `SELECT code,name,organization,member_limit,color FROM teams WHERE code=$1 FOR UPDATE`, code).
		Scan(&team.Code, &team.Name, &team.Organization, &team.MemberLimit, &team.Color)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Team{}, core.NewError(http.StatusNotFound, "resource_not_found", "Team was not found")
	}
	if err != nil {
		return core.Team{}, dbError(err)
	}
	var memberCount int32
	if err = tx.QueryRow(ctx, `SELECT count(*)::int FROM contestants WHERE team_code=$1`, code).Scan(&memberCount); err != nil {
		return core.Team{}, dbError(err)
	}
	if patch.Name != nil {
		team.Name = *patch.Name
	}
	if patch.Color != nil {
		team.Color = *patch.Color
	}
	if patch.Organization != nil {
		team.Organization = *patch.Organization
	}
	if patch.MemberLimit != nil {
		if *patch.MemberLimit < memberCount {
			return core.Team{}, core.NewError(http.StatusConflict, "team_full", "Member limit cannot be lower than the current member count")
		}
		team.MemberLimit = *patch.MemberLimit
	}
	err = tx.QueryRow(ctx, `
		UPDATE teams SET name=$2,organization=$3,member_limit=$4,color=$5,updated_at=now()
		WHERE code=$1 RETURNING code,name,organization,member_limit,color`, code, team.Name, team.Organization, team.MemberLimit, core.TeamColor(team.Color)).
		Scan(&team.Code, &team.Name, &team.Organization, &team.MemberLimit, &team.Color)
	if err != nil {
		return core.Team{}, conflictError(err, "conflict", "Team name is already in use")
	}
	if err = tx.Commit(ctx); err != nil {
		return core.Team{}, conflictError(err, "conflict", "Team update conflicted with a concurrent signup")
	}
	return team, nil
}
