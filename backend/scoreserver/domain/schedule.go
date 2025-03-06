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
		schedule, err := data.parse()
		if err != nil {
			return nil, WrapAsInternal(err, "failed to parse schedule")
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func SaveSchedule(ctx context.Context, eff ScheduleWriter, input []*UpdateScheduleInput) error {
	data := make([]*ScheduleData, 0, len(input))
	for _, schedule := range input {
		phase, err := schedule.Phase.toStringPhase()
		if err != nil {
			return WrapAsInternal(err, "failed to convert phase")
		}
		data = append(data, &ScheduleData{
			ID:          uuid.Must(uuid.NewV4()),
			StringPhase: phase,
			StartAt:     schedule.StartAt,
			EndAt:       schedule.EndAt,
		})
	}
	if err := eff.SaveSchedule(ctx, data); err != nil {
		return WrapAsInternal(err, "failed to save schedules")
	}
	return nil
}

type StringPhase string

func (s *StringPhase) String() (Phase, error) {
	switch *s {
	case "OUT_OF_CONTEST":
		return PhaseOutOfContest, nil
	case "IN_CONTEST":
		return PhaseInContest, nil
	case "BREAK":
		return PhaseBreak, nil
	case "AFTER_CONTEST":
		return PhaseAfterContest, nil
	case "UNSPECIFIED":
		fallthrough
	default:
		return PhaseUnspecified, NewInvalidArgumentError("phase is required", nil)
	}
}

type (
	ScheduleData struct {
		ID          uuid.UUID
		StringPhase StringPhase
		StartAt     time.Time
		EndAt       time.Time
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

func (d *ScheduleData) parse() (*Schedule, error) {
	phase, err := d.StringPhase.String()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to parse phase")
	}
	return &Schedule{
			id:      d.ID,
			phase:   phase,
			startAt: d.StartAt,
			endAt:   d.EndAt,
		},
		nil
}

func (s *Schedule) Data() (*ScheduleData, error) {
	phase, err := s.phase.toStringPhase()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to convert phase")
	}
	return &ScheduleData{
			ID:          s.id,
			StringPhase: phase,
			StartAt:     s.startAt,
			EndAt:       s.endAt,
		},
		nil
}

func (p Phase) toStringPhase() (StringPhase, error) {
	switch p {
	case PhaseOutOfContest:
		return "OUT_OF_CONTEST", nil
	case PhaseInContest:
		return "IN_CONTEST", nil
	case PhaseBreak:
		return "BREAK", nil
	case PhaseAfterContest:
		return "AFTER_CONTEST", nil
	case PhaseUnspecified:
		fallthrough
	default:
		return "", NewInvalidArgumentError("phase is required", nil)
	}
}
