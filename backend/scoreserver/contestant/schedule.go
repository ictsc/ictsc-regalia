package contestant

import (
	"context"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type ScheduleEnforcer struct {
	ScheduleReader domain.ScheduleReader
}

func (h *ScheduleEnforcer) Enforce(ctx context.Context, allowedPhases ...domain.Phase) error {
	schedule, err := domain.GetSchedule(ctx, h.ScheduleReader)
	if err != nil {
		return err
	}

	current := schedule.Current(time.Now())
	if !slices.Contains(allowedPhases, current.Phase()) {
		return connect.NewError(connect.CodePermissionDenied, errors.New("not allowed in the current phase"))
	}

	return nil
}
