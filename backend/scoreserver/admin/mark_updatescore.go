package admin

import (
	"context"
	"log/slog"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type UpdateScoreEffect interface {
	domain.MarkingResultReader
	domain.MarkingResultPenaltyUpdator
	domain.DeploymentReader
	domain.AnswerReader
	domain.UpdateAnswerScoreEffect
	domain.TeamProblemLister
	domain.UpdateProblemScoreEffect
	domain.ScheduleReader
	domain.ScoreVisibilitySettingsReader
}

func newUpdateScoreEffect(repo *pg.Repository, scheduleReader domain.ScheduleReader) UpdateScoreEffect {
	return struct {
		*pg.Repository
		domain.ScheduleReader
	}{
		Repository:     repo,
		ScheduleReader: scheduleReader,
	}
}

type UpdateScoreResult struct{}

func UpdateScore(
	ctx context.Context,
	eff UpdateScoreEffect,
	now time.Time,
	mode domain.ScoreUpdateMode,
) (*UpdateScoreResult, error) {
	deploymentsList, err := eff.ListDeployments(ctx)
	if err != nil {
		return nil, domain.WrapAsInternal(err, "failed to list deployments")
	}
	cachedDeployReader := &cachedDeploymentReader{list: deploymentsList}

	markingResultsList, err := eff.ListMarkingResults(ctx)
	if err != nil {
		return nil, domain.WrapAsInternal(err, "failed to list marking results")
	}
	cachedMarkReader := &cachedMarkingResultReader{list: markingResultsList}

	innerEff := struct {
		domain.MarkingResultReader
		domain.DeploymentReader
		domain.MarkingResultPenaltyUpdator
		domain.ScoreWriter

		domain.AnswerReader
		domain.TeamProblemLister
		domain.ScheduleReader
		domain.ScoreVisibilitySettingsReader
	}{
		MarkingResultReader:         cachedMarkReader,
		DeploymentReader:            cachedDeployReader,
		MarkingResultPenaltyUpdator: eff,
		ScoreWriter:                 eff,

		AnswerReader:                  eff,
		TeamProblemLister:             eff,
		ScheduleReader:                eff,
		ScoreVisibilitySettingsReader: eff,
	}
	return updateScore(ctx, innerEff, now, mode)
}

func updateScore(ctx context.Context, eff UpdateScoreEffect, now time.Time, mode domain.ScoreUpdateMode) (*UpdateScoreResult, error) {
	inContest := false
	if mode == domain.ScoreUpdateModeNormal {
		schedules, err := domain.GetSchedule(ctx, eff)
		if err != nil {
			return nil, err
		}
		inContest = schedules.Current(now) != nil
	}

	settings, err := domain.GetScoreVisibilitySettings(ctx, eff)
	if err != nil {
		return nil, err
	}
	policy := domain.NewScoreUpdatePolicy(mode, inContest, settings.IsRankingFrozenAt(now))

	slog.InfoContext(ctx, "Update marking results")
	markingResults, err := domain.ListAllMarkingResults(ctx, eff)
	if err != nil {
		return nil, err
	}
	for _, markingResult := range markingResults {
		if err := markingResult.UpdatePenalty(ctx, eff); err != nil {
			return nil, err
		}
	}

	slog.InfoContext(ctx, "Update answer scores")
	answers, err := domain.ListAnswersForAdmin(ctx, eff)
	if err != nil {
		return nil, err
	}
	for _, answer := range answers {
		if err := answer.UpdateScore(ctx, eff, now, policy); err != nil {
			return nil, err
		}
	}

	slog.InfoContext(ctx, "Update problem scores")
	teamProblems, err := domain.ListTeamProblemsForAdmin(ctx, eff)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "Update team problems", "count", len(teamProblems))
	for _, teamProblem := range teamProblems {
		if err := teamProblem.UpdateScore(ctx, eff, policy); err != nil {
			return nil, err
		}
	}

	return &UpdateScoreResult{}, nil
}

type cachedDeploymentReader struct {
	list []*domain.DeploymentData
}

func (r *cachedDeploymentReader) ListDeployments(context.Context) ([]*domain.DeploymentData, error) {
	return r.list, nil
}

type cachedMarkingResultReader struct {
	list []*domain.MarkingResultData
}

func (r *cachedMarkingResultReader) ListMarkingResults(context.Context) ([]*domain.MarkingResultData, error) {
	return r.list, nil
}
