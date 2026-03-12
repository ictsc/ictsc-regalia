package domain

import (
	"context"
	"time"
)

type ScoreVisibility string

const (
	ScoreVisibilityPrivate ScoreVisibility = "PRIVATE"
	ScoreVisibilityTeam    ScoreVisibility = "TEAM"
	ScoreVisibilityPublic  ScoreVisibility = "PUBLIC"
)

type ScoreUpdateMode int

const (
	ScoreUpdateModeNormal ScoreUpdateMode = iota
	ScoreUpdateModeRevealFinal
)

type ScoreVisibilitySettings struct {
	rankingFreezeAt *time.Time
}

func (s *ScoreVisibilitySettings) RankingFreezeAt() *time.Time {
	return s.rankingFreezeAt
}

func (s *ScoreVisibilitySettings) IsRankingFrozenAt(now time.Time) bool {
	return s.rankingFreezeAt != nil && !now.Before(*s.rankingFreezeAt)
}

func GetScoreVisibilitySettings(
	ctx context.Context,
	eff ScoreVisibilitySettingsReader,
) (*ScoreVisibilitySettings, error) {
	data, err := eff.GetScoreVisibilitySettings(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get score visibility settings")
	}
	return data.parse(), nil
}

func SaveScoreVisibilitySettings(
	ctx context.Context,
	eff ScoreVisibilitySettingsWriter,
	input *UpdateScoreVisibilitySettingsInput,
) error {
	if err := eff.SaveScoreVisibilitySettings(ctx, (&ScoreVisibilitySettingsData{
		RankingFreezeAt: input.RankingFreezeAt,
	})); err != nil {
		return WrapAsInternal(err, "failed to save score visibility settings")
	}
	return nil
}

type (
	ScoreVisibilitySettingsData struct {
		RankingFreezeAt *time.Time
	}

	UpdateScoreVisibilitySettingsInput struct {
		RankingFreezeAt *time.Time
	}

	ScoreVisibilitySettingsReader interface {
		GetScoreVisibilitySettings(ctx context.Context) (*ScoreVisibilitySettingsData, error)
	}

	ScoreVisibilitySettingsWriter interface {
		ScoreVisibilitySettingsReader
		SaveScoreVisibilitySettings(ctx context.Context, data *ScoreVisibilitySettingsData) error
	}
)

func (d *ScoreVisibilitySettingsData) parse() *ScoreVisibilitySettings {
	if d == nil {
		return &ScoreVisibilitySettings{}
	}
	return &ScoreVisibilitySettings{
		rankingFreezeAt: d.RankingFreezeAt,
	}
}
