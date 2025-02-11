package pg

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type teamMemberRow struct {
	User userRow `db:"u"`
	Team teamRow `db:"t"`
}

func (r teamMemberRow) toData() *domain.TeamMemberData {
	return &domain.TeamMemberData{
		User: (*domain.UserData)(&r.User),
		Team: (*domain.TeamData)(&r.Team),
	}
}

var _ domain.TeamMemberGetter = (*repo)(nil)

func (r *repo) GetTeamMemberByID(ctx context.Context, userID uuid.UUID) (*domain.TeamMemberData, error) {
	var row teamMemberRow
	if err := sqlx.GetContext(ctx, r.ext, &row, `
		SELECT
			u.id AS "u.id", u.name AS "u.name",
			t.id AS "t.id", t.code AS "t.code", t.name AS "t.name", t.organization AS "t.organization"
		FROM team_members AS tm
		LEFT JOIN users AS u ON u.id = tm.user_id
		LEFT JOIN teams AS t ON t.id = tm.team_id
		WHERE tm.user_id = $1
	`, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("team member", nil)
		}
		return nil, errors.Wrap(err, "failed to select team member")
	}

	return row.toData(), nil
}

func (r *repo) CountTeamMembers(ctx context.Context, teamID uuid.UUID) (uint, error) {
	var count uint
	if err := sqlx.GetContext(ctx, r.ext, &count, "SELECT COUNT(*) FROM team_members WHERE team_id = $1", teamID); err != nil {
		return 0, errors.Wrap(err, "failed to count team members")
	}
	return count, nil
}

var _ domain.TeamMemberManager = (*RepositoryTx)(nil)

func (r *RepositoryTx) AddTeamMember(ctx context.Context, userID, invitationCodeID, teamID uuid.UUID) error {
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO team_members (user_id, invitation_code_id, team_id)
		VALUES ($1, $2, $3)
	`, userID, invitationCodeID, teamID); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewAlreadyExistsError("team member", nil)
		}
		return errors.Wrap(err, "failed to insert team member")
	}
	return nil
}
