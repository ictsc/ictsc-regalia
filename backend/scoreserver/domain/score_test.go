package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

// ============================================================
// TestNewScoreUpdatePolicy
// ============================================================

func TestNewScoreUpdatePolicy(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		mode          domain.ScoreUpdateMode
		inContest     bool
		rankingFrozen bool
		want          domain.ScoreUpdatePolicy
	}{
		"Normal/outOfContest/notFrozen": {
			mode:          domain.ScoreUpdateModeNormal,
			inContest:     false,
			rankingFrozen: false,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true},
		},
		"Normal/outOfContest/frozen": {
			mode:          domain.ScoreUpdateModeNormal,
			inContest:     false,
			rankingFrozen: true,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true},
		},
		"Normal/inContest/notFrozen": {
			mode:          domain.ScoreUpdateModeNormal,
			inContest:     true,
			rankingFrozen: false,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true, UpdatePublic: true},
		},
		"Normal/inContest/frozen": {
			mode:          domain.ScoreUpdateModeNormal,
			inContest:     true,
			rankingFrozen: true,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true},
		},
		"RevealFinal/outOfContest/notFrozen": {
			mode:          domain.ScoreUpdateModeRevealFinal,
			inContest:     false,
			rankingFrozen: false,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true, UpdatePublic: true, BypassVisibilityDelay: true},
		},
		"RevealFinal/outOfContest/frozen": {
			mode:          domain.ScoreUpdateModeRevealFinal,
			inContest:     false,
			rankingFrozen: true,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true, UpdatePublic: true, BypassVisibilityDelay: true},
		},
		"RevealFinal/inContest/notFrozen": {
			mode:          domain.ScoreUpdateModeRevealFinal,
			inContest:     true,
			rankingFrozen: false,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true, UpdatePublic: true, BypassVisibilityDelay: true},
		},
		"RevealFinal/inContest/frozen": {
			mode:          domain.ScoreUpdateModeRevealFinal,
			inContest:     true,
			rankingFrozen: true,
			want:          domain.ScoreUpdatePolicy{UpdatePrivate: true, UpdateTeam: true, UpdatePublic: true, BypassVisibilityDelay: true},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := domain.NewScoreUpdatePolicy(tt.mode, tt.inContest, tt.rankingFrozen)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewScoreUpdatePolicy(%v, inContest=%v, rankingFrozen=%v) mismatch (-want +got):\n%s",
					tt.mode, tt.inContest, tt.rankingFrozen, diff)
			}
		})
	}
}

// ============================================================
// TestAnswerUpdateScore
// ============================================================

func TestAnswerUpdateScore(t *testing.T) {
	t.Parallel()

	answerCreatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	// now is >20 minutes after answerCreatedAt so IsPublic(now) returns true
	now := answerCreatedAt.Add(time.Hour)

	// Build an AnswerData with fixtures so that it parses correctly.
	answerData := &domain.AnswerData{
		ID:        uuid.Must(uuid.NewV4()),
		Number:    1,
		Team:      domain.FixTeam1(t, nil).Data(),
		Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
		Author:    domain.FixUser1(t, nil).Data(),
		CreatedAt: answerCreatedAt,
		Interval:  domain.AnswerInterval,
	}

	// Parse the AnswerData into a *Answer via the public reader API.
	answers, err := domain.ListAnswersForAdmin(t.Context(), answerReader{
		listAnswersForAdminFunc: func(context.Context) ([]*domain.AnswerData, error) {
			return []*domain.AnswerData{answerData}, nil
		},
	})
	if err != nil || len(answers) == 0 {
		t.Fatalf("setup ListAnswersForAdmin: %v", err)
	}
	answer := answers[0]

	// A marking result that references the same answer (matched by ID during filtering).
	markData := &domain.MarkingResultData{
		ID:        uuid.Must(uuid.NewV4()),
		Judge:     "judge",
		Answer:    answerData,
		Score:     &domain.ScoreData{MarkedScore: 50},
		Rationale: &domain.MarkingRationaleData{DescriptiveComment: "good"},
		CreatedAt: answerCreatedAt.Add(5 * time.Minute),
	}

	cases := map[string]struct {
		policy           domain.ScoreUpdatePolicy
		wantVisibilities []domain.ScoreVisibility
	}{
		"Normal/outOfContest": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, false, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate},
		},
		"Normal/inContest/notFrozen": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, true, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam, domain.ScoreVisibilityPublic},
		},
		"Normal/inContest/frozen": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, true, true),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam},
		},
		"RevealFinal": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeRevealFinal, false, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam, domain.ScoreVisibilityPublic},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var updatedVisibilities []domain.ScoreVisibility
			eff := updateAnswerScoreEff{
				listMarkingResults: func(context.Context) ([]*domain.MarkingResultData, error) {
					return []*domain.MarkingResultData{markData}, nil
				},
				updateAnswerScore: func(_ context.Context, input *domain.UpdateAnswerScoreInput) error {
					updatedVisibilities = append(updatedVisibilities, input.Visibility)
					return nil
				},
			}

			if err := answer.UpdateScore(t.Context(), eff, now, tt.policy); err != nil {
				t.Fatalf("UpdateScore() error: %v", err)
			}
			if diff := cmp.Diff(tt.wantVisibilities, updatedVisibilities); diff != "" {
				t.Errorf("updated visibilities mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// ============================================================
// TestTeamProblemUpdateScore
// ============================================================

func TestTeamProblemUpdateScore(t *testing.T) {
	t.Parallel()

	// Obtain a *TeamProblem from a parsed answer (no score needed on the answer itself).
	plainAnswerData := &domain.AnswerData{
		ID:        uuid.Must(uuid.NewV4()),
		Number:    1,
		Team:      domain.FixTeam1(t, nil).Data(),
		Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
		Author:    domain.FixUser1(t, nil).Data(),
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Interval:  domain.AnswerInterval,
	}
	plainAnswers, err := domain.ListAnswersForAdmin(t.Context(), answerReader{
		listAnswersForAdminFunc: func(context.Context) ([]*domain.AnswerData, error) {
			return []*domain.AnswerData{plainAnswerData}, nil
		},
	})
	if err != nil || len(plainAnswers) == 0 {
		t.Fatalf("setup ListAnswersForAdmin: %v", err)
	}
	teamProblem := plainAnswers[0].TeamProblem()

	// A scored answer that the ListAnswersByTeamProblem mock will return.
	scoredAnswerData := &domain.AnswerData{
		ID:        uuid.Must(uuid.NewV4()),
		Number:    1,
		Team:      domain.FixTeam1(t, nil).Data(),
		Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
		Author:    domain.FixUser1(t, nil).Data(),
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Interval:  domain.AnswerInterval,
		Score: &domain.ScoreData{
			MarkingResultID: uuid.Must(uuid.NewV4()),
			MarkedScore:     50,
		},
	}

	cases := map[string]struct {
		policy           domain.ScoreUpdatePolicy
		wantVisibilities []domain.ScoreVisibility
	}{
		"Normal/outOfContest": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, false, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate},
		},
		"Normal/inContest/notFrozen": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, true, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam, domain.ScoreVisibilityPublic},
		},
		"Normal/inContest/frozen": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeNormal, true, true),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam},
		},
		"RevealFinal": {
			policy:           domain.NewScoreUpdatePolicy(domain.ScoreUpdateModeRevealFinal, false, false),
			wantVisibilities: []domain.ScoreVisibility{domain.ScoreVisibilityPrivate, domain.ScoreVisibilityTeam, domain.ScoreVisibilityPublic},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var updatedVisibilities []domain.ScoreVisibility
			eff := updateProblemScoreEff{
				listAnswersByTeamProblem: func(_ context.Context, _ domain.ScoreVisibility, _ int64, _ string) ([]*domain.AnswerData, error) {
					return []*domain.AnswerData{scoredAnswerData}, nil
				},
				updateProblemScore: func(_ context.Context, input *domain.UpdateProblemScoreInput) error {
					updatedVisibilities = append(updatedVisibilities, input.Visibility)
					return nil
				},
			}

			if err := teamProblem.UpdateScore(t.Context(), eff, tt.policy); err != nil {
				t.Fatalf("UpdateScore() error: %v", err)
			}
			if diff := cmp.Diff(tt.wantVisibilities, updatedVisibilities); diff != "" {
				t.Errorf("updated visibilities mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// ============================================================
// Effect stubs
// ============================================================

type updateAnswerScoreEff struct {
	listMarkingResults func(ctx context.Context) ([]*domain.MarkingResultData, error)
	updateAnswerScore  func(ctx context.Context, input *domain.UpdateAnswerScoreInput) error
}

var _ domain.UpdateAnswerScoreEffect = updateAnswerScoreEff{}

func (e updateAnswerScoreEff) ListMarkingResults(ctx context.Context) ([]*domain.MarkingResultData, error) {
	if e.listMarkingResults != nil {
		return e.listMarkingResults(ctx)
	}
	return nil, nil
}

func (e updateAnswerScoreEff) UpdateAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput) error {
	if e.updateAnswerScore != nil {
		return e.updateAnswerScore(ctx, input)
	}
	return nil
}

func (e updateAnswerScoreEff) UpdateProblemScore(context.Context, *domain.UpdateProblemScoreInput) error {
	return nil
}

type updateProblemScoreEff struct {
	listAnswersByTeamProblem func(ctx context.Context, visibility domain.ScoreVisibility, teamCode int64, problemCode string) ([]*domain.AnswerData, error)
	updateProblemScore       func(ctx context.Context, input *domain.UpdateProblemScoreInput) error
}

var _ domain.UpdateProblemScoreEffect = updateProblemScoreEff{}

func (e updateProblemScoreEff) ListAnswers(context.Context, domain.ScoreVisibility) ([]*domain.AnswerData, error) {
	return nil, nil
}

func (e updateProblemScoreEff) ListAnswersByTeamProblem(ctx context.Context, visibility domain.ScoreVisibility, teamCode int64, problemCode string) ([]*domain.AnswerData, error) {
	if e.listAnswersByTeamProblem != nil {
		return e.listAnswersByTeamProblem(ctx, visibility, teamCode, problemCode)
	}
	return nil, nil
}

func (e updateProblemScoreEff) GetAnswerDetail(context.Context, domain.ScoreVisibility, int64, string, uint32) (*domain.AnswerDetailData, error) {
	return nil, nil
}

func (e updateProblemScoreEff) ListMarkingResults(context.Context) ([]*domain.MarkingResultData, error) {
	return nil, nil
}

func (e updateProblemScoreEff) UpdateAnswerScore(context.Context, *domain.UpdateAnswerScoreInput) error {
	return nil
}

func (e updateProblemScoreEff) UpdateProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput) error {
	if e.updateProblemScore != nil {
		return e.updateProblemScore(ctx, input)
	}
	return nil
}
