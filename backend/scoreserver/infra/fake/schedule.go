package fake

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type FakeScheduler config.FakeSchedule

var _ domain.ScheduleReader = (*FakeScheduler)(nil)

func (f *FakeScheduler) GetSchedule(context.Context) ([]*domain.ScheduleData, error) {
	return []*domain.ScheduleData{
		{
			ID:      uuid.FromStringOrNil("f179260e-5bc3-4d9e-8022-d87392711a21"),
			Phase:   f.Phase,
			StartAt: f.StartAt,
			EndAt:   f.EndAt,
		},
		{
			ID:      uuid.FromStringOrNil("dcedec19-f42a-46ca-a4c1-497cdf13085d"),
			Phase:   f.NextPhase,
			StartAt: f.EndAt,
		},
	}, nil
}
