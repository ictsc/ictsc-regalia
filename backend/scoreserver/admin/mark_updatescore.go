package admin

import (
	"context"
	"log/slog"
	"time"

	"github.com/cockroachdb/errors"
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

type problemScoreUpdateEffect interface {
	domain.UpdateProblemScoreEffect
	domain.TeamProblemLister
}

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
	cachedDeployReader, err := domain.NewCachedDeploymentReader(deploymentsList)
	if err != nil {
		return nil, domain.WrapAsInternal(err, "failed to cache deployments")
	}

	markingResultsList, err := eff.ListMarkingResults(ctx)
	if err != nil {
		return nil, domain.WrapAsInternal(err, "failed to list marking results")
	}
	cachedMarkReader, err := domain.NewCachedMarkingResultReader(markingResultsList)
	if err != nil {
		return nil, domain.WrapAsInternal(err, "failed to cache marking results")
	}

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
	problemScoreEff, err := newProblemScoreUpdateEffect(ctx, eff)
	if err != nil {
		return nil, err
	}
	teamProblems, err := domain.ListTeamProblemsForAdmin(ctx, problemScoreEff)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "Update team problems", "count", len(teamProblems))
	for _, teamProblem := range teamProblems {
		if err := teamProblem.UpdateScore(ctx, problemScoreEff, policy); err != nil {
			return nil, err
		}
	}

	return &UpdateScoreResult{}, nil
}

func newProblemScoreUpdateEffect(ctx context.Context, eff UpdateScoreEffect) (problemScoreUpdateEffect, error) {
	cachedAnswerReader, err := newCachedAnswerReader(ctx, eff)
	if err != nil {
		return nil, err
	}
	return struct {
		domain.AnswerReader
		domain.MarkingResultReader
		domain.ScoreWriter
		domain.TeamProblemLister
	}{
		AnswerReader:        cachedAnswerReader,
		MarkingResultReader: eff,
		ScoreWriter:         eff,
		TeamProblemLister:   eff,
	}, nil
}

type cachedAnswerReader struct {
	fallback      domain.AnswerReader
	byVisibility  map[domain.ScoreVisibility][]*domain.AnswerData
	byTeamProblem map[domain.ScoreVisibility]map[answerCacheKey][]*domain.AnswerData
}

type answerCacheKey struct {
	teamCode    int64
	problemCode string
}

func newCachedAnswerReader(ctx context.Context, fallback domain.AnswerReader) (*cachedAnswerReader, error) {
	byVisibility := make(map[domain.ScoreVisibility][]*domain.AnswerData, 3)
	byTeamProblem := make(map[domain.ScoreVisibility]map[answerCacheKey][]*domain.AnswerData, 3)
	for _, visibility := range []domain.ScoreVisibility{
		domain.ScoreVisibilityPrivate,
		domain.ScoreVisibilityTeam,
		domain.ScoreVisibilityPublic,
	} {
		answers, err := fallback.ListAnswers(ctx, visibility)
		if err != nil {
			return nil, domain.WrapAsInternal(err, "failed to list answers")
		}
		byVisibility[visibility] = answers
		grouped := make(map[answerCacheKey][]*domain.AnswerData)
		for _, answer := range answers {
			if answer.Team == nil || answer.Problem == nil {
				return nil, domain.NewInvalidArgumentError("answer cache requires team and problem", nil)
			}
			key := answerCacheKey{
				teamCode:    answer.Team.Code,
				problemCode: answer.Problem.Code,
			}
			grouped[key] = append(grouped[key], answer)
		}
		byTeamProblem[visibility] = grouped
	}
	return &cachedAnswerReader{
		fallback:      fallback,
		byVisibility:  byVisibility,
		byTeamProblem: byTeamProblem,
	}, nil
}

func (r *cachedAnswerReader) ListAnswers(_ context.Context, visibility domain.ScoreVisibility) ([]*domain.AnswerData, error) {
	return r.byVisibility[visibility], nil
}

func (r *cachedAnswerReader) ListAnswersByTeamProblem(
	_ context.Context,
	visibility domain.ScoreVisibility,
	teamCode int64,
	problemCode string,
) ([]*domain.AnswerData, error) {
	return r.byTeamProblem[visibility][answerCacheKey{
		teamCode:    teamCode,
		problemCode: problemCode,
	}], nil
}

func (r *cachedAnswerReader) GetAnswerDetail(
	ctx context.Context,
	visibility domain.ScoreVisibility,
	teamCode int64,
	problemCode string,
	answerNumber uint32,
) (*domain.AnswerDetailData, error) {
	if r.fallback == nil {
		return nil, errors.New("fallback answer reader is nil")
	}
	return r.fallback.GetAnswerDetail(ctx, visibility, teamCode, problemCode, answerNumber)
}
