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

type teamRow struct {
	ID           uuid.UUID `db:"id"`
	Code         int64     `db:"code"`
	Name         string    `db:"name"`
	Organization string    `db:"organization"`
	MaxMembers   uint      `db:"max_members"`
}

var _ domain.TeamsLister = (*repo)(nil)

func (r *repo) ListTeams(ctx context.Context) ([]*domain.TeamData, error) {
	rows, err := r.ext.QueryxContext(ctx, "SELECT id, code, name, organization, max_members FROM teams ORDER BY code ASC")
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

func (r *repo) GetTeamByCode(ctx context.Context, code int64) (*domain.TeamData, error) {
	var row teamRow
	if err := sqlx.GetContext(ctx, r.ext,
		&row, "SELECT id, code, name, organization, max_members FROM teams WHERE code = $1 LIMIT 1",
		code,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrTypeNotFound, errors.New("team not found"))
		}
		return nil, errors.Wrap(err, "failed to select team")
	}
	return (*domain.TeamData)(&row), nil
}

var _ domain.TeamCreator = (*repo)(nil)

func (r *repo) CreateTeam(ctx context.Context, team *domain.TeamData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO teams (id, code, name, organization, max_members, created_at, updated_at)
		VALUES (:id, :code, :name, :organization, :max_members, NOW(), NOW())`,
		(*teamRow)(team),
	); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			// 一意性制約違反
			if pgErr.Code == "23505" {
				return domain.NewError(domain.ErrTypeAlreadyExists, errors.New("team already exists"))
			}
		}
		return errors.Wrap(err, "failed to insert team")
	}

	return nil
}

var _ domain.TeamUpdater = (*repo)(nil)

func (r *repo) UpdateTeam(ctx context.Context, team *domain.TeamData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		UPDATE teams SET
			code = :code, name = :name, organization = :organization,
			max_members = :max_members, updated_at = NOW()
		WHERE id = :id`,
		(*teamRow)(team),
	); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.NewError(domain.ErrTypeAlreadyExists, errors.New("team already exists"))
			}
		}
		return errors.Wrap(err, "failed to update team")
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
