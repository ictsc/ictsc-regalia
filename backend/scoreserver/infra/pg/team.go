package pg

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type teamRow struct {
	ID           uuid.UUID `db:"id"`
	Code         int64     `db:"code"`
	Name         string    `db:"name"`
	Organization string    `db:"organization"`
	MaxMembers   uint      `db:"max_members"`
}

var teamColumns = columns([]string{"id", "code", "name", "organization", "max_members"})

var _ domain.TeamsLister = (*repo)(nil)

var listTeamsQuery = "SELECT " + teamColumns.String("") + " FROM teams ORDER BY code ASC"

func (r *repo) ListTeams(ctx context.Context) ([]*domain.TeamData, error) {
	rows, err := r.ext.QueryxContext(ctx, listTeamsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select teams")
	}
	defer rows.Close() //nolint:errcheck

	var teams []*domain.TeamData
	for rows.Next() {
		var row teamRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan team")
		}
		teams = append(teams, (*domain.TeamData)(&row))
	}

	return teams, nil
}

var _ domain.TeamGetter = (*repo)(nil)

var getTeamQuery = "SELECT " + teamColumns.String("") + " FROM teams WHERE code = $1 LIMIT 1"

func (r *repo) GetTeamByCode(ctx context.Context, code int64) (*domain.TeamData, error) {
	var row teamRow
	if err := sqlx.GetContext(ctx, r.ext, &row, getTeamQuery, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("team", nil)
		}
		return nil, errors.Wrap(err, "failed to select team")
	}
	return (*domain.TeamData)(&row), nil
}

var (
	_ domain.TeamCreator = (*repo)(nil)
	_ domain.TeamUpdater = (*repo)(nil)
)

func (r *repo) CreateTeam(ctx context.Context, team *domain.TeamData) error {
	return r.saveTeam(ctx, team)
}

func (r *repo) UpdateTeam(ctx context.Context, team *domain.TeamData) error {
	return r.saveTeam(ctx, team)
}

func (r *RepositoryTx) SaveTeam(ctx context.Context, team *domain.TeamData) error {
	return r.saveTeam(ctx, team)
}

func (r *repo) saveTeam(ctx context.Context, team *domain.TeamData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO teams (id, code, name, organization, max_members, created_at, updated_at)
		VALUES (:id, :code, :name, :organization, :max_members, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code, name = EXCLUDED.name, organization = EXCLUDED.organization,
			max_members = EXCLUDED.max_members, updated_at = NOW()`,
		(*teamRow)(team),
	); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.NewAlreadyExistsError("team", nil)
			}
		}
		return errors.Wrap(err, "failed to save team")
	}
	return nil
}

var _ domain.TeamDeleter = (*repo)(nil)

func (r *repo) DeleteTeam(ctx context.Context, teamID uuid.UUID) error {
	if _, err := r.ext.ExecContext(ctx, "DELETE FROM teams WHERE id = $1", teamID); err != nil {
		return errors.Wrap(err, "failed to delete team")
	}

	return nil
}
