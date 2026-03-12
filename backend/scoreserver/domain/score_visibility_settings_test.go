package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestScoreVisibilitySettingsIsRankingFrozenAt(t *testing.T) {
	t.Parallel()

	freezeAt := time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		settings *domain.ScoreVisibilitySettings
		now      time.Time
		want     bool
	}{
		"nil freeze time": {
			settings: &domain.ScoreVisibilitySettings{},
			now:      freezeAt,
			want:     false,
		},
		"before freeze": {
			settings: settingsWithFreezeAt(freezeAt),
			now:      freezeAt.Add(-time.Second),
			want:     false,
		},
		"at freeze": {
			settings: settingsWithFreezeAt(freezeAt),
			now:      freezeAt,
			want:     true,
		},
		"after freeze": {
			settings: settingsWithFreezeAt(freezeAt),
			now:      freezeAt.Add(time.Second),
			want:     true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.settings.IsRankingFrozenAt(tt.now); got != tt.want {
				t.Fatalf("IsRankingFrozenAt(%v) = %v, want %v", tt.now, got, tt.want)
			}
		})
	}
}

func TestGetScoreVisibilitySettings(t *testing.T) {
	t.Parallel()

	freezeAt := time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		reader  scoreVisibilitySettingsReaderStub
		want    *domain.ScoreVisibilitySettings
		wantErr error
	}{
		"ok with nil settings": {
			reader: scoreVisibilitySettingsReaderStub{
				data: nil,
			},
			want: &domain.ScoreVisibilitySettings{},
		},
		"ok with freeze time": {
			reader: scoreVisibilitySettingsReaderStub{
				data: &domain.ScoreVisibilitySettingsData{
					RankingFreezeAt: &freezeAt,
				},
			},
			want: settingsWithFreezeAt(freezeAt),
		},
		"reader error": {
			reader: scoreVisibilitySettingsReaderStub{
				err: errors.New("boom"),
			},
			wantErr: errors.New("boom"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.GetScoreVisibilitySettings(t.Context(), &tt.reader)
			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetScoreVisibilitySettings() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetScoreVisibilitySettings() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want.RankingFreezeAt(), got.RankingFreezeAt()); diff != "" {
				t.Fatalf("freeze time mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSaveScoreVisibilitySettings(t *testing.T) {
	t.Parallel()

	freezeAt := time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		input   *domain.UpdateScoreVisibilitySettingsInput
		writer  scoreVisibilitySettingsWriterStub
		want    *domain.ScoreVisibilitySettingsData
		wantErr error
	}{
		"ok with nil freeze time": {
			input: &domain.UpdateScoreVisibilitySettingsInput{},
			want:  &domain.ScoreVisibilitySettingsData{},
		},
		"ok with freeze time": {
			input: &domain.UpdateScoreVisibilitySettingsInput{
				RankingFreezeAt: &freezeAt,
			},
			want: &domain.ScoreVisibilitySettingsData{
				RankingFreezeAt: &freezeAt,
			},
		},
		"writer error": {
			input: &domain.UpdateScoreVisibilitySettingsInput{},
			writer: scoreVisibilitySettingsWriterStub{
				err: errors.New("boom"),
			},
			wantErr: errors.New("boom"),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			writer := tt.writer
			err := domain.SaveScoreVisibilitySettings(t.Context(), &writer, tt.input)
			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("SaveScoreVisibilitySettings() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("SaveScoreVisibilitySettings() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, writer.saved); diff != "" {
				t.Fatalf("saved settings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type scoreVisibilitySettingsReaderStub struct {
	data *domain.ScoreVisibilitySettingsData
	err  error
}

func (s *scoreVisibilitySettingsReaderStub) GetScoreVisibilitySettings(context.Context) (*domain.ScoreVisibilitySettingsData, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.data, nil
}

type scoreVisibilitySettingsWriterStub struct {
	scoreVisibilitySettingsReaderStub
	saved *domain.ScoreVisibilitySettingsData
	err   error
}

func (s *scoreVisibilitySettingsWriterStub) SaveScoreVisibilitySettings(_ context.Context, data *domain.ScoreVisibilitySettingsData) error {
	if s.err != nil {
		return s.err
	}
	s.saved = data
	return nil
}

func settingsWithFreezeAt(freezeAt time.Time) *domain.ScoreVisibilitySettings {
	reader := &scoreVisibilitySettingsReaderStub{
		data: &domain.ScoreVisibilitySettingsData{
			RankingFreezeAt: &freezeAt,
		},
	}
	settings, err := domain.GetScoreVisibilitySettings(context.Background(), reader)
	if err != nil {
		panic(err)
	}
	return settings
}
