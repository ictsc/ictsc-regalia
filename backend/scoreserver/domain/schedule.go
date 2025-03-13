package domain

import (
	"context"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

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
	entry, _ := s.getEntry(now)
	return entry
}

func (s Schedule) Next(now time.Time) *ScheduleEntry {
	entry, idx := s.getEntry(now)
	if idx+1 < len(s) {
		return s[idx+1]
	} else if idx+1 == len(s) {
		return &ScheduleEntry{phase: PhaseAfterContest, startAt: s[idx].endAt}
	}
	return entry
}

func (s Schedule) getEntry(now time.Time) (*ScheduleEntry, int) {
	if len(s) == 0 {
		return &ScheduleEntry{phase: PhaseOutOfContest}, 0
	}
	idx := slices.IndexFunc(s, func(entry *ScheduleEntry) bool {
		// now <@ [startAt, endAt)
		var lowerBound, upperBound bool
		if entry.startAt.IsZero() {
			lowerBound = true
		} else {
			lowerBound = now.Equal(entry.startAt) || now.After(entry.startAt)
		}
		if entry.endAt.IsZero() {
			upperBound = true
		} else {
			upperBound = now.Before(entry.endAt)
		}
		return lowerBound && upperBound
	})
	if idx >= 0 {
		return s[idx], idx
	}
	if now.Before(s[0].startAt) {
		return &ScheduleEntry{phase: PhaseOutOfContest, endAt: s[0].startAt}, -1
	} else {
		return &ScheduleEntry{phase: PhaseAfterContest, startAt: s[len(s)-1].endAt}, len(s)
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
