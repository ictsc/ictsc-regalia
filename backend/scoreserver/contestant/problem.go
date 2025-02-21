package contestant

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type ProblemServiceHandler struct {
	contestantv1connect.UnimplementedProblemServiceHandler

	ListEffect ProblemListEffect
	GetEffect  ProblemGetEffect
}

var _ contestantv1connect.ProblemServiceHandler = (*ProblemServiceHandler)(nil)

func newProblemServiceHandler(repo *pg.Repository) *ProblemServiceHandler {
	return &ProblemServiceHandler{
		ListEffect: repo,
		GetEffect:  repo,
	}
}

type ProblemListEffect interface {
	domain.TeamMemberGetter
	domain.ProblemReader
}

func (h *ProblemServiceHandler) ListProblems(
	ctx context.Context,
	req *connect.Request[contestantv1.ListProblemsRequest],
) (*connect.Response[contestantv1.ListProblemsResponse], error) {
	userSess, err := session.UserSessionStore.Get(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.ListEffect)
	if err != nil {
		return nil, err
	}

	problems, err := teamMember.Team().Problems(ctx, h.ListEffect)
	if err != nil {
		return nil, err
	}

	protoProblems := make([]*contestantv1.Problem, 0, len(problems))
	for _, problem := range problems {
		protoProblems = append(protoProblems, &contestantv1.Problem{
			Code:     string(problem.Problem().Code()),
			Title:    problem.Problem().Title(),
			MaxScore: problem.Problem().MaxScore(),
		})
	}

	return connect.NewResponse(&contestantv1.ListProblemsResponse{
		Problems: protoProblems,
	}), nil
}

type ProblemGetEffect = ProblemListEffect

func (h *ProblemServiceHandler) GetProblem(
	ctx context.Context,
	req *connect.Request[contestantv1.GetProblemRequest],
) (*connect.Response[contestantv1.GetProblemResponse], error) {
	userSess, err := session.UserSessionStore.Get(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetCode())
	if err != nil {
		return nil, err
	}

	tp, err := teamMember.Team().ProblemDetailByCode(ctx, h.GetEffect, code)
	if err != nil {
		return nil, err
	}
	detail := tp.ProblemDetail()

	proto := &contestantv1.Problem{
		Code:     string(detail.Code()),
		Title:    detail.Title(),
		MaxScore: detail.MaxScore(),
		Body: &contestantv1.ProblemBody{
			Type: contestantv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE,
			Body: &contestantv1.ProblemBody_Descriptive{
				Descriptive: &contestantv1.DescriptiveProblem{
					Body: detail.Body(),
				},
			},
		},
	}
	return connect.NewResponse(&contestantv1.GetProblemResponse{
		Problem: proto,
	}), nil
}
