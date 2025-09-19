package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestSubmitDescriptiveAnswer(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		at      time.Time
		eff     answerWriter
		member  *domain.TeamMember
		problem *domain.Problem
		body    string

		wantErr    error
		wantAnswer *domain.AnswerDetailData
	}{
		"new answer": {
			at: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
			eff: answerWriter{
				getLatestAnswerByTeamProblemForPublicFunc: func(context.Context, uuid.UUID, uuid.UUID) (*domain.AnswerData, error) {
					return nil, domain.ErrNotFound
				},
			},
			member:  domain.FixTeamMember1(t, nil),
			problem: domain.FixDescriptiveProblem1(t, nil).Problem(),
			body:    "answer",

			wantAnswer: &domain.AnswerDetailData{
				Answer: &domain.AnswerData{
					Number:    1,
					Team:      domain.FixTeam1(t, nil).Data(),
					Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
					Author:    domain.FixUser1(t, nil).Data(),
					CreatedAt: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
					Interval:  1 * time.Millisecond,
				},
				Body: &domain.AnswerBodyData{
					Descriptive: &domain.DescriptiveAnswerBodyData{
						Body: "answer",
					},
				},
			},
		},
		"additional answer": {
			at: time.Date(2021, 1, 1, 1, 30, 0, 0, time.UTC),
			eff: answerWriter{
				getLatestAnswerByTeamProblemForPublicFunc: func(
					ctx context.Context, teamID, problemID uuid.UUID,
				) (*domain.AnswerData, error) {
					return &domain.AnswerData{
						ID:        uuid.Must(uuid.NewV4()),
						Number:    1,
						Team:      domain.FixTeam1(t, nil).Data(),
						Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
						Author:    domain.FixUser1(t, nil).Data(),
						CreatedAt: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
						Interval:  20 * time.Minute,
					}, nil
				},
			},
			member:  domain.FixTeamMember1(t, nil),
			problem: domain.FixDescriptiveProblem1(t, nil).Problem(),
			body:    "answer",

			wantAnswer: &domain.AnswerDetailData{
				Answer: &domain.AnswerData{
					Number:    2,
					Team:      domain.FixTeam1(t, nil).Data(),
					Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
					Author:    domain.FixUser1(t, nil).Data(),
					CreatedAt: time.Date(2021, 1, 1, 1, 30, 0, 0, time.UTC),
					Interval:  1 * time.Millisecond,
				},
				Body: &domain.AnswerBodyData{
					Descriptive: &domain.DescriptiveAnswerBodyData{
						Body: "answer",
					},
				},
			},
		},
		"too early": {
			at: time.Date(2021, 1, 1, 1, 10, 0, 0, time.UTC),
			eff: answerWriter{
				getLatestAnswerByTeamProblemForPublicFunc: func(
					ctx context.Context, teamID, problemID uuid.UUID,
				) (*domain.AnswerData, error) {
					return &domain.AnswerData{
						ID:        uuid.Must(uuid.NewV4()),
						Number:    1,
						Team:      domain.FixTeam1(t, nil).Data(),
						Problem:   domain.FixDescriptiveProblem1(t, nil).Problem().Data(),
						Author:    domain.FixUser1(t, nil).Data(),
						CreatedAt: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
						Interval:  20 * time.Minute,
					}, nil
				},
			},
			member:  domain.FixTeamMember1(t, nil),
			problem: domain.FixDescriptiveProblem1(t, nil).Problem(),
			body:    "answer",

			wantErr: domain.ErrTooEarlyToSubmitAnswer,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var createdAnswer *domain.AnswerDetailData
			eff := answerWriter{
				getLatestAnswerByTeamProblemForPublicFunc: tt.eff.getLatestAnswerByTeamProblemForPublicFunc,
				createAnswerFunc: func(ctx context.Context, data *domain.AnswerDetailData) error {
					if f := tt.eff.createAnswerFunc; f != nil {
						if err := f(ctx, data); err != nil {
							return err
						}
					}
					createdAnswer = data
					return nil
				},
			}

			answer, err := tt.member.SubmitDescriptiveAnswer(t.Context(), tt.at, eff, tt.problem, tt.body)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err: %v, want err: %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if diff := cmp.Diff(
				createdAnswer, answer.Data(),
			); diff != "" {
				t.Errorf("answer mismatch (-created +got):\n%s", diff)
			}

			if diff := cmp.Diff(
				tt.wantAnswer, answer.Data(),
				cmpopts.IgnoreFields(domain.AnswerData{}, "ID"),
			); diff != "" {
				t.Errorf("answer mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type (
	listAnswersForAdminFunc                   func(ctx context.Context) ([]*domain.AnswerData, error)
	listAnswersByTeamProblemForAdminFunc      func(ctx context.Context, teamCode int64, problemCode string) ([]*domain.AnswerData, error)
	listAnswersByTeamProblemForPublicFunc     func(ctx context.Context, teamCode int64, problemCode string) ([]*domain.AnswerData, error)
	getLatestAnswerByTeamProblemForPublicFunc func(ctx context.Context, teamID, problemID uuid.UUID) (*domain.AnswerData, error)

	getAnswerDetailForAdminFunc  func(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*domain.AnswerDetailData, error)
	getAnswerDetailForPublicFunc func(ctx context.Context, teamCode int64, problemCode string, answerNumber uint32) (*domain.AnswerDetailData, error)
	createAnswerFunc             func(ctx context.Context, data *domain.AnswerDetailData) error

	answerReader struct {
		listAnswersForAdminFunc
		listAnswersByTeamProblemForAdminFunc
		listAnswersByTeamProblemForPublicFunc
		getAnswerDetailForPublicFunc
		getAnswerDetailForAdminFunc
	}
	answerWriter struct {
		getLatestAnswerByTeamProblemForPublicFunc
		createAnswerFunc
	}
)

var (
	_ domain.AnswerReader = answerReader{}
	_ domain.AnswerWriter = answerWriter{}
)

func (f listAnswersForAdminFunc) ListAnswersForAdmin(ctx context.Context) ([]*domain.AnswerData, error) {
	return f(ctx)
}

func (f listAnswersByTeamProblemForAdminFunc) ListAnswersByTeamProblemForAdmin(
	ctx context.Context, teamCode int64, problemCode string,
) ([]*domain.AnswerData, error) {
	return f(ctx, teamCode, problemCode)
}

func (f listAnswersByTeamProblemForPublicFunc) ListAnswersByTeamProblemForPublic(
	ctx context.Context, teamCode int64, problemCode string,
) ([]*domain.AnswerData, error) {
	return f(ctx, teamCode, problemCode)
}

func (f getLatestAnswerByTeamProblemForPublicFunc) GetLatestAnswerByTeamProblemForPublic(
	ctx context.Context, teamID, problemID uuid.UUID,
) (*domain.AnswerData, error) {
	return f(ctx, teamID, problemID)
}

func (f getAnswerDetailForPublicFunc) GetAnswerDetailForPublic(
	ctx context.Context,
	teamCode int64, problemCode string, answerNumber uint32,
) (*domain.AnswerDetailData, error) {
	return f(ctx, teamCode, problemCode, answerNumber)
}

func (f getAnswerDetailForAdminFunc) GetAnswerDetailForAdmin(
	ctx context.Context,
	teamCode int64, problemCode string, answerNumber uint32,
) (*domain.AnswerDetailData, error) {
	return f(ctx, teamCode, problemCode, answerNumber)
}

func (f createAnswerFunc) CreateAnswer(ctx context.Context, data *domain.AnswerDetailData) error {
	return f(ctx, data)
}
