package domain

import (
	"context"
	"testing"
	"time"
)

type scheduleWriterStub struct {
	saved []*ScheduleData
}

func (s *scheduleWriterStub) GetSchedule(context.Context) ([]*ScheduleData, error) {
	return nil, nil
}

func (s *scheduleWriterStub) SaveSchedule(_ context.Context, data []*ScheduleData) error {
	s.saved = data
	return nil
}

func TestSaveScheduleRejectsDuplicateNames(t *testing.T) {
	t.Parallel()

	writer := &scheduleWriterStub{}
	err := SaveSchedule(t.Context(), writer, []*UpdateScheduleInput{
		{
			Name:    "sched-1",
			StartAt: time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC),
			EndAt:   time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC),
		},
		{
			Name:    "sched-1",
			StartAt: time.Date(2026, 1, 1, 3, 0, 0, 0, time.UTC),
			EndAt:   time.Date(2026, 1, 1, 4, 0, 0, 0, time.UTC),
		},
	})
	if err == nil {
		t.Fatal("expected duplicated schedule name to fail")
	}
	if len(writer.saved) != 0 {
		t.Fatal("unexpected write on invalid input")
	}
}
