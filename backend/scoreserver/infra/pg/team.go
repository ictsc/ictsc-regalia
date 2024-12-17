package pg

import (
	"context"
	"database/sql"
	"time"

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
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (t *teamRow) asDomain() (*domain.Team, error) {
	return domain.NewTeam(domain.TeamInput{
		ID:           t.ID,
		Code:         int(t.Code),
		Name:         t.Name,
		Organization: t.Organization,
	})
}

var _ domain.TeamsLister = (*repo)(nil)

func (r *repo) ListTeams(ctx context.Context) ([]*domain.Team, error) {
	var rows []teamRow
	if err := sqlx.SelectContext(ctx, r.ext,
		&rows, "SELECT * FROM teams ORDER BY code ASC",
	); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to select teams")
	}

	teams := make([]*domain.Team, 0, len(rows))
	for _, row := range rows {
		team, err := row.asDomain()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to validate: code=%d", row.Code)
		}
		teams = append(teams, team)
	}

	return teams, nil
}

var _ domain.TeamGetter = (*repo)(nil)

func (r *repo) GetTeamByCode(ctx context.Context, code domain.TeamCode) (*domain.Team, error) {
	var row teamRow
	if err := sqlx.GetContext(ctx, r.ext,
		&row, "SELECT * FROM teams WHERE code = $1 LIMIT 1",
		int64(code),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrTypeNotFound, errors.New("team not found"))
		}
		return nil, errors.Wrap(err, "failed to select team")
	}
	return row.asDomain()
}

var _ domain.TeamCreator = (*repo)(nil)

func (r *repo) CreateTeam(ctx context.Context, team *domain.Team) error {
	now := time.Now()

	row := teamRow{
		ID:           team.ID(),
		Code:         int64(team.Code()),
		Name:         team.Name(),
		Organization: team.Organization(),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := sqlx.NamedExecContext(ctx, r.ext,
		`INSERT INTO teams (id, code, name, organization, created_at, updated_at)
		 VALUES (:id, :code, :name, :organization, :created_at, :updated_at)`,
		row,
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

func (r *repo) UpdateTeam(ctx context.Context, team *domain.Team) error {
	now := time.Now()
	if _, err := r.ext.ExecContext(ctx,
		`UPDATE teams
		 SET code = $1, name = $2, organization = $3, updated_at = $4
		 WHERE id = $5`,
		team.Code(), team.Name(), team.Organization(), now, team.ID(),
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

func (r *repo) DeleteTeam(ctx context.Context, team *domain.Team) error {
	if _, err := r.ext.ExecContext(ctx, "DELETE FROM teams WHERE id = $1", team.ID()); err != nil {
		return errors.Wrap(err, "failed to delete team")
	}

	return nil
}
