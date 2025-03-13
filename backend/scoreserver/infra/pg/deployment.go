package pg

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"maps"
	"slices"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var _ domain.DeploymentReader = (*repo)(nil)

type deploymentRow struct {
	ID          uuid.UUID        `db:"id"`
	TeamCode    int64            `db:"team_code"`
	ProblemCode string           `db:"problem_code"`
	Revision    uint32           `db:"revision"`
	Status      deploymentStatus `db:"status"`
	CreatedAt   time.Time        `db:"created_at"`
}

const listDeploymentsQuery = `
SELECT
	req.id AS id, team.code AS team_code, problem.code AS problem_code, req.revision AS revision,
	ev.status AS status, ev.created_at AS created_at
FROM redeployment_events AS ev
JOIN redeployment_requests AS req ON req.id = ev.request_id
JOIN teams AS team ON team.id = req.team_id
JOIN problems AS problem ON problem.id = req.problem_id`

func (r *repo) ListDeployments(ctx context.Context) ([]*domain.DeploymentData, error) {
	ctx, span := tracer.Start(ctx, "pg.ListDeployments")
	defer span.End()

	var rows []deploymentRow
	if err := sqlx.SelectContext(ctx, r.ext, &rows, listDeploymentsQuery); err != nil {
		return nil, errors.Wrap(err, "failed to select deployments")
	}

	deployments := make(map[uuid.UUID]*domain.DeploymentData)
	for _, row := range rows {
		deploy, ok := deployments[row.ID]
		if !ok {
			deploy = &domain.DeploymentData{
				ID:          row.ID,
				TeamCode:    row.TeamCode,
				ProblemCode: row.ProblemCode,
				Revision:    row.Revision,
			}
			deployments[row.ID] = deploy
		}
		deploy.Events = append(deploy.Events, &domain.DeploymentEventData{
			Status:     domain.DeploymentStatus(row.Status),
			OccurredAt: row.CreatedAt,
		})
	}

	return slices.Collect(maps.Values(deployments)), nil
}

var _ domain.DeploymentWriter = (*RepositoryTx)(nil)

func (tx *RepositoryTx) CreateDeployment(ctx context.Context, input *domain.CreateDeploymentInput) error {
	ctx, span := tracer.Start(ctx, "pg.CreateDeployment")
	defer span.End()

	if _, err := tx.ext.ExecContext(ctx, `
		INSERT INTO redeployment_requests (id, team_id, problem_id, revision)
		VALUES ($1, $2, $3, $4)`,
		input.ID, input.TeamID, input.ProblemID, input.Revision,
	); err != nil {
		if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.NewAlreadyExistsError("deployment", nil)
		}
		return errors.Wrap(err, "failed to insert deployment")
	}

	if _, err := tx.ext.ExecContext(ctx, `
		INSERT INTO redeployment_events (request_id, status, created_at)
		VALUES ($1, $2, $3)`,
		input.ID, deploymentStatus(input.Status), input.OccurredAt,
	); err != nil {
		if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.NewAlreadyExistsError("deployment event", nil)
		}
		return errors.Wrap(err, "failed to insert deployment event")
	}

	return nil
}

func (tx *RepositoryTx) UpdateDeploymentStatus(ctx context.Context, input *domain.UpdateDeploymentStatusInput) error {
	ctx, span := tracer.Start(ctx, "pg.UpdateDeploymentStatus")
	defer span.End()

	if _, err := tx.ext.ExecContext(ctx, `
		INSERT INTO redeployment_events (request_id, status, created_at)
		VALUES ($1, $2, $3)`,
		input.ID, deploymentStatus(input.Status), input.OccurredAt,
	); err != nil {
		if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.NewAlreadyExistsError("deployment event", nil)
		}
		return errors.Wrap(err, "failed to insert deployment event")
	}

	return nil
}

type deploymentStatus domain.DeploymentStatus

var (
	_ sql.Scanner   = (*deploymentStatus)(nil)
	_ driver.Valuer = deploymentStatus(domain.DeploymentStatusUnknown)
)

func (s *deploymentStatus) Scan(src any) error {
	*s = deploymentStatus(domain.DeploymentStatusUnknown)
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if !ok {
		return nil
	}
	switch v {
	case "QUEUED":
		*s = deploymentStatus(domain.DeploymentStatusQueued)
	case "DEPLOYING":
		*s = deploymentStatus(domain.DeploymentStatusCreating)
	case "COMPLETED":
		*s = deploymentStatus(domain.DeploymentStatusCompleted)
	case "FAILED":
		*s = deploymentStatus(domain.DeploymentStatusFailed)
	}
	return nil
}

func (s deploymentStatus) Value() (driver.Value, error) {
	switch s {
	case deploymentStatus(domain.DeploymentStatusQueued):
		return "QUEUED", nil
	case deploymentStatus(domain.DeploymentStatusCreating):
		return "DEPLOYING", nil
	case deploymentStatus(domain.DeploymentStatusCompleted):
		return "COMPLETED", nil
	case deploymentStatus(domain.DeploymentStatusFailed):
		return "FAILED", nil
	case deploymentStatus(domain.DeploymentStatusUnknown):
		fallthrough
	default:
		return nil, errors.New("unknown deployment status")
	}
}
