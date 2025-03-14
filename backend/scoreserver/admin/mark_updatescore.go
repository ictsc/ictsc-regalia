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
}

func newUpdateScoreEffect(repo *pg.Repository) UpdateScoreEffect {
	return repo
}

type UpdateScoreResult struct{}

func UpdateScore(
	ctx context.Context,
	eff UpdateScoreEffect,
	now time.Time,
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
	}{
		MarkingResultReader:         cachedMarkReader,
		DeploymentReader:            cachedDeployReader,
		MarkingResultPenaltyUpdator: eff,
		ScoreWriter:                 eff,

		AnswerReader:      eff,
		TeamProblemLister: eff,
	}
	return updateScore(ctx, innerEff, now)
}

func updateScore(ctx context.Context, eff UpdateScoreEffect, now time.Time) (*UpdateScoreResult, error) {
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
		if err := answer.UpdateScore(ctx, eff, now); err != nil {
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
		if err := teamProblem.UpdateScore(ctx, eff); err != nil {
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
