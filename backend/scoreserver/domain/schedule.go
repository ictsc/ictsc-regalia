package domain

import (
	"context"
	"slices"
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

type (
	Schedule      []*ScheduleEntry
	ScheduleEntry struct {
		id      uuid.UUID
		phase   Phase
		startAt time.Time
		endAt   time.Time
	}
)

func (s Schedule) Current(now time.Time) *ScheduleEntry {
	idx := s.currentIndex(now)
	if idx < 0 {
		return s.overTimedSchedule(now)
	}
	return s[idx]
}

func (s Schedule) Next(now time.Time) *ScheduleEntry {
	idx := s.currentIndex(now)
	if idx < 0 || idx == len(s)-1 {
		return s.overTimedSchedule(now)
	}
	return s[idx+1]
}

func (s Schedule) currentIndex(now time.Time) int {
	return slices.IndexFunc(s, func(e *ScheduleEntry) bool {
		// now <@ [startAt, endAt)
		return (now.Equal(e.startAt) || now.After(e.startAt)) && now.Before(e.endAt)
	})
}

func (s Schedule) overTimedSchedule(now time.Time) *ScheduleEntry {
	var minTime, maxTime time.Time
	if len(s) > 0 {
		minTime, maxTime = s[0].startAt, s[len(s)-1].endAt
	}
	if now.Before(minTime) {
		return &ScheduleEntry{phase: PhaseOutOfContest, endAt: minTime}
	} else {
		return &ScheduleEntry{phase: PhaseAfterContest, startAt: maxTime}
	}
}

func (s *ScheduleEntry) Phase() Phase {
	return s.phase
}

func (s *ScheduleEntry) StartAt() time.Time {
	return s.startAt
}

func (s *ScheduleEntry) EndAt() time.Time {
	return s.endAt
}

func GetSchedule(ctx context.Context, eff ScheduleReader) (Schedule, error) {
	scheduleData, err := eff.GetSchedule(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get schedule")
	}
	entries := make([]*ScheduleEntry, 0, len(scheduleData))
	for _, data := range scheduleData {
		entries = append(entries, data.parse())
	}
	slices.SortFunc(entries, func(i, j *ScheduleEntry) int {
		return i.startAt.Compare(j.startAt)
	})
	return entries, nil
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

func (d *ScheduleData) parse() *ScheduleEntry {
	return &ScheduleEntry{
		id:      d.ID,
		phase:   d.Phase,
		startAt: d.StartAt,
		endAt:   d.EndAt,
	}
}
