package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Phase int32

const (
	PhaseUnspecified Phase = iota
	PhaseOutOfContest
	PhaseInContest
	PhaseBreak
	PhaseAfterContest
)

func (p Phase) String() string {
	switch p {
	case PhaseOutOfContest:
		return "OUT_OF_CONTEST"
	case PhaseInContest:
		return "IN_CONTEST"
	case PhaseBreak:
		return "BREAK"
	case PhaseAfterContest:
		return "AFTER_CONTEST"
	case PhaseUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}

type Schedule struct {
	id      uuid.UUID
	phase   Phase
	startAt time.Time
	endAt   time.Time
}

func (s *Schedule) Phase() Phase {
	return s.phase
}

func (s *Schedule) StartAt() time.Time {
	return s.startAt
}

func (s *Schedule) EndAt() time.Time {
	return s.endAt
}

func GetSchedule(ctx context.Context, eff ScheduleReader) ([]*Schedule, error) {
	scheduleData, err := eff.GetSchedule(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get schedule")
	}
	schedules := make([]*Schedule, 0, len(scheduleData))
	for _, data := range scheduleData {
		schedules = append(schedules, data.parse())
	}
	return schedules, nil
}

func SaveSchedule(ctx context.Context, eff ScheduleWriter, input []*UpdateScheduleInput) error {
	data := make([]*ScheduleData, 0, len(input))
	for _, schedule := range input {
		data = append(data, &ScheduleData{
			ID:      uuid.Must(uuid.NewV4()),
			Phase:   schedule.Phase,
			StartAt: schedule.StartAt,
			EndAt:   schedule.EndAt,
		})
	}
	if err := eff.SaveSchedule(ctx, data); err != nil {
		return WrapAsInternal(err, "failed to save schedules")
	}
	return nil
}

type (
	ScheduleData struct {
		ID      uuid.UUID
		Phase   Phase
		StartAt time.Time
		EndAt   time.Time
	}
	UpdateScheduleInput struct {
		Phase   Phase
		StartAt time.Time
		EndAt   time.Time
	}
	ScheduleReader interface {
		GetSchedule(ctx context.Context) ([]*ScheduleData, error)
	}
	ScheduleWriter interface {
		ScheduleReader
		SaveSchedule(ctx context.Context, data []*ScheduleData) error
	}
)

func (d *ScheduleData) parse() *Schedule {
	return &Schedule{
		id:      d.ID,
		phase:   d.Phase,
		startAt: d.StartAt,
		endAt:   d.EndAt,
	}
}

func (s *Schedule) Data() *ScheduleData {
	return &ScheduleData{
		ID:      s.id,
		Phase:   s.phase,
		StartAt: s.startAt,
		EndAt:   s.endAt,
	}
}
