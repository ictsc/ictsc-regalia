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
		name    string // スケジュール名 (例: "day1-am")
		startAt time.Time
		endAt   time.Time
	}
)

// Current は現在時刻に該当するスケジュールエントリを返す
// 該当がない場合はnilを返す
func (s Schedule) Current(now time.Time) *ScheduleEntry {
	for _, entry := range s {
		if entry.Contains(now) {
			return entry
		}
	}
	return nil
}

// Contains は指定時刻がこのスケジュールの範囲内かどうかを返す
func (s *ScheduleEntry) Contains(t time.Time) bool {
	// t <@ [startAt, endAt)
	afterStart := t.Equal(s.startAt) || t.After(s.startAt)
	beforeEnd := t.Before(s.endAt)
	return afterStart && beforeEnd
}

func (s *ScheduleEntry) ID() uuid.UUID {
	return s.id
}

func (s *ScheduleEntry) Name() string {
	return s.name
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
			ID:      uuid.Nil, // Repository層で適切なIDが設定される
			Name:    schedule.Name,
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
		Name    string // スケジュール名
		StartAt time.Time
		EndAt   time.Time
	}
	UpdateScheduleInput struct {
		Name    string // スケジュール名
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
		name:    d.Name,
		startAt: d.StartAt,
		endAt:   d.EndAt,
	}
}
