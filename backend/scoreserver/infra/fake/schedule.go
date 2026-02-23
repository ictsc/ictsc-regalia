package fake

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type FakeScheduler config.FakeSchedule

var _ domain.ScheduleReader = (*FakeScheduler)(nil)

func (f *FakeScheduler) GetSchedule(context.Context) ([]*domain.ScheduleData, error) {
	return []*domain.ScheduleData{
		{
			Name:    f.Name,
			StartAt: f.StartAt,
			EndAt:   f.EndAt,
		},
	}, nil
}
