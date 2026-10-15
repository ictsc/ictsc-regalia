package postgres

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/jackc/pgx/v5"
)

func deploymentEventsFrom(ctx context.Context, queryer interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, requestID string) ([]core.DeploymentEvent, error) {
	rows, err := queryer.Query(ctx, `
		SELECT event_id,occurred_at,status,message FROM deployment_events
		WHERE request_id=$1 ORDER BY occurred_at,event_id`, requestID)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	events := make([]core.DeploymentEvent, 0)
	for rows.Next() {
		var event core.DeploymentEvent
		if err := rows.Scan(&event.EventID, &event.OccurredAt, &event.Status, &event.Message); err != nil {
			return nil, dbError(err)
		}
		events = append(events, event)
	}
	return events, dbError(rows.Err())
}

func sameDeploymentEvent(left, right core.DeploymentEvent) bool {
	if left.EventID != right.EventID || !left.OccurredAt.Equal(right.OccurredAt) || left.Status != right.Status {
		return false
	}
	if left.Message == nil || right.Message == nil {
		return left.Message == nil && right.Message == nil
	}
	return *left.Message == *right.Message
}
