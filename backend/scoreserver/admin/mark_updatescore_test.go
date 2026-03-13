package admin

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestCachedScoreWriterUpdateAnswerScoreSkipsUnchanged(t *testing.T) {
	t.Parallel()

	answerID := uuid.Must(uuid.NewV4())
	currentMarkingResultID := uuid.Must(uuid.NewV4())
	nextMarkingResultID := uuid.Must(uuid.NewV4())

	reader := scoreCacheReaderStub{
		listAnswersFunc: func(_ context.Context, visibility domain.ScoreVisibility) ([]*domain.AnswerData, error) {
			if visibility != domain.ScoreVisibilityPrivate {
				return nil, nil
			}
			return []*domain.AnswerData{{
				ID:    answerID,
				Score: &domain.ScoreData{MarkingResultID: currentMarkingResultID},
			}}, nil
		},
	}

	var updateCount int
	writer, err := newCachedScoreWriter(t.Context(), reader, scoreWriterStub{
		updateAnswerScoreFunc: func(_ context.Context, input *domain.UpdateAnswerScoreInput) error {
			updateCount++
			if input.MarkingResultID != nextMarkingResultID {
				t.Fatalf("unexpected marking result id: %v", input.MarkingResultID)
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("newCachedScoreWriter() error: %v", err)
	}

	if err := writer.UpdateAnswerScore(t.Context(), &domain.UpdateAnswerScoreInput{
		AnswerID:        answerID,
		MarkingResultID: currentMarkingResultID,
		Visibility:      domain.ScoreVisibilityPrivate,
	}); err != nil {
		t.Fatalf("UpdateAnswerScore(skip current) error: %v", err)
	}
	if updateCount != 0 {
		t.Fatalf("UpdateAnswerScore(skip current) called fallback %d times", updateCount)
	}

	input := &domain.UpdateAnswerScoreInput{
		AnswerID:        answerID,
		MarkingResultID: nextMarkingResultID,
		Visibility:      domain.ScoreVisibilityPrivate,
	}
	if err := writer.UpdateAnswerScore(t.Context(), input); err != nil {
		t.Fatalf("UpdateAnswerScore(change) error: %v", err)
	}
	if err := writer.UpdateAnswerScore(t.Context(), input); err != nil {
		t.Fatalf("UpdateAnswerScore(skip cached) error: %v", err)
	}
	if updateCount != 1 {
		t.Fatalf("UpdateAnswerScore(change) called fallback %d times, want 1", updateCount)
	}
}

func TestCachedScoreWriterUpdateProblemScoreSkipsUnchanged(t *testing.T) {
	t.Parallel()

	teamID := uuid.Must(uuid.NewV4())
	problemID := uuid.Must(uuid.NewV4())
	currentMarkingResultID := uuid.Must(uuid.NewV4())
	nextMarkingResultID := uuid.Must(uuid.NewV4())

	reader := scoreCacheReaderStub{
		listTeamProblemScoresFunc: func(_ context.Context, visibility domain.ScoreVisibility) ([]*domain.TeamProblemScoreData, error) {
			if visibility != domain.ScoreVisibilityTeam {
				return nil, nil
			}
			return []*domain.TeamProblemScoreData{{
				TeamID:    teamID,
				ProblemID: problemID,
				Score:     domain.ScoreData{MarkingResultID: currentMarkingResultID},
			}}, nil
		},
	}

	var updateCount int
	writer, err := newCachedScoreWriter(t.Context(), reader, scoreWriterStub{
		updateProblemScoreFunc: func(_ context.Context, input *domain.UpdateProblemScoreInput) error {
			updateCount++
			if input.MarkingResultID != nextMarkingResultID {
				t.Fatalf("unexpected marking result id: %v", input.MarkingResultID)
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("newCachedScoreWriter() error: %v", err)
	}

	if err := writer.UpdateProblemScore(t.Context(), &domain.UpdateProblemScoreInput{
		TeamID:          teamID,
		ProblemID:       problemID,
		MarkingResultID: currentMarkingResultID,
		UpdateSubmitAt:  time.Unix(0, 0),
		Visibility:      domain.ScoreVisibilityTeam,
	}); err != nil {
		t.Fatalf("UpdateProblemScore(skip current) error: %v", err)
	}
	if updateCount != 0 {
		t.Fatalf("UpdateProblemScore(skip current) called fallback %d times", updateCount)
	}

	input := &domain.UpdateProblemScoreInput{
		TeamID:          teamID,
		ProblemID:       problemID,
		MarkingResultID: nextMarkingResultID,
		UpdateSubmitAt:  time.Unix(0, 0),
		Visibility:      domain.ScoreVisibilityTeam,
	}
	if err := writer.UpdateProblemScore(t.Context(), input); err != nil {
		t.Fatalf("UpdateProblemScore(change) error: %v", err)
	}
	if err := writer.UpdateProblemScore(t.Context(), input); err != nil {
		t.Fatalf("UpdateProblemScore(skip cached) error: %v", err)
	}
	if updateCount != 1 {
		t.Fatalf("UpdateProblemScore(change) called fallback %d times, want 1", updateCount)
	}
}

type scoreCacheReaderStub struct {
	listAnswersFunc           func(ctx context.Context, visibility domain.ScoreVisibility) ([]*domain.AnswerData, error)
	listTeamProblemScoresFunc func(ctx context.Context, visibility domain.ScoreVisibility) ([]*domain.TeamProblemScoreData, error)
}

func (s scoreCacheReaderStub) ListAnswers(ctx context.Context, visibility domain.ScoreVisibility) ([]*domain.AnswerData, error) {
	if s.listAnswersFunc != nil {
		return s.listAnswersFunc(ctx, visibility)
	}
	return nil, nil
}

func (scoreCacheReaderStub) ListAnswersByTeamProblem(context.Context, domain.ScoreVisibility, int64, string) ([]*domain.AnswerData, error) {
	return nil, nil
}

func (scoreCacheReaderStub) GetAnswerDetail(context.Context, domain.ScoreVisibility, int64, string, uint32) (*domain.AnswerDetailData, error) {
	return nil, nil
}

func (s scoreCacheReaderStub) ListTeamProblemScores(ctx context.Context, visibility domain.ScoreVisibility) ([]*domain.TeamProblemScoreData, error) {
	if s.listTeamProblemScoresFunc != nil {
		return s.listTeamProblemScoresFunc(ctx, visibility)
	}
	return nil, nil
}

func (scoreCacheReaderStub) GetTeamProblemScore(context.Context, domain.ScoreVisibility, uuid.UUID, uuid.UUID) (*domain.ScoreData, error) {
	return nil, nil
}

func (scoreCacheReaderStub) ListTeamProblemScoresByTeamID(context.Context, domain.ScoreVisibility, uuid.UUID) ([]*domain.TeamProblemScoreData, error) {
	return nil, nil
}

type scoreWriterStub struct {
	updateAnswerScoreFunc  func(ctx context.Context, input *domain.UpdateAnswerScoreInput) error
	updateProblemScoreFunc func(ctx context.Context, input *domain.UpdateProblemScoreInput) error
}

func (s scoreWriterStub) UpdateAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput) error {
	if s.updateAnswerScoreFunc != nil {
		return s.updateAnswerScoreFunc(ctx, input)
	}
	return nil
}

func (s scoreWriterStub) UpdateProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput) error {
	if s.updateProblemScoreFunc != nil {
		return s.updateProblemScoreFunc(ctx, input)
	}
	return nil
}
