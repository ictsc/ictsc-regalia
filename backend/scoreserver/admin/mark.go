package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MarkServiceHandler struct {
	adminv1connect.UnimplementedMarkServiceHandler

	Enforcer         *auth.Enforcer
	ListAnswerEffect domain.AnswerReader
	GetAnswerEffect  domain.AnswerReader
}

var _ adminv1connect.MarkServiceHandler = (*MarkServiceHandler)(nil)

func newMarkServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *MarkServiceHandler {
	return &MarkServiceHandler{
		Enforcer:         enforcer,
		ListAnswerEffect: repo,
		GetAnswerEffect:  repo,
	}
}

func (h *MarkServiceHandler) ListAnswers(
	ctx context.Context,
	req *connect.Request[adminv1.ListAnswersRequest],
) (*connect.Response[adminv1.ListAnswersResponse], error) {
	if err := enforce(ctx, h.Enforcer, "answers", "list"); err != nil {
		return nil, err
	}

	answers, err := domain.ListAnswers(ctx, h.ListAnswerEffect)
	if err != nil {
		return nil, err
	}

	protoAnswers := make([]*adminv1.Answer, 0, len(answers))
	for _, answer := range answers {
		protoAnswers = append(protoAnswers, &adminv1.Answer{
			Id:   answer.Number(),
			Team: convertTeam(answer.Team()),
			Author: &adminv1.Contestant{
				Name: string(answer.Author().Name()),
				Team: convertTeam(answer.Author().Team()),
			},
			Problem:   convertProblem(answer.Problem()),
			CreatedAt: timestamppb.New(answer.CreatedAt()),
		})
	}

	return connect.NewResponse(&adminv1.ListAnswersResponse{
		Answers: protoAnswers,
	}), nil
}

func (h *MarkServiceHandler) GetAnswer(
	ctx context.Context,
	req *connect.Request[adminv1.GetAnswerRequest],
) (*connect.Response[adminv1.GetAnswerResponse], error) {
	if err := enforce(ctx, h.Enforcer, "answers", "get"); err != nil {
		return nil, err
	}

	protoTeamCode := req.Msg.GetTeamCode()
	if protoTeamCode == 0 {
		return nil, domain.NewInvalidArgumentError("team_code is required", nil)
	}
	protoProblemCode := req.Msg.GetProblemCode()
	if protoProblemCode == "" {
		return nil, domain.NewInvalidArgumentError("problem_code is required", nil)
	}
	protoID := req.Msg.GetId()
	if protoID == 0 {
		return nil, domain.NewInvalidArgumentError("id is required", nil)
	}

	teamCode, err := domain.NewTeamCode(int64(protoTeamCode))
	if err != nil {
		return nil, err
	}

	problemCode, err := domain.NewProblemCode(protoProblemCode)
	if err != nil {
		return nil, err
	}

	answer, err := domain.GetAnswerDetail(ctx, h.GetAnswerEffect, teamCode, problemCode, protoID)
	if err != nil {
		return nil, err
	}
	protoAnser := &adminv1.Answer{
		Id:   answer.Number(),
		Team: convertTeam(answer.Team()),
		Author: &adminv1.Contestant{
			Name: string(answer.Author().Name()),
			Team: convertTeam(answer.Author().Team()),
		},
		Problem:   convertProblem(answer.Problem()),
		CreatedAt: timestamppb.New(answer.CreatedAt()),
	}
	switch answer.Problem().Type() {
	case domain.ProblemTypeDescriptive:
		desc, err := answer.Body().Descriptive()
		if err != nil {
			return nil, err
		}
		protoAnser.Body = &adminv1.AnswerBody{
			Body: &adminv1.AnswerBody_Descriptive{
				Descriptive: &adminv1.DescriptiveAnswer{
					Body: desc.Body(),
				},
			},
		}
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
	}

	return connect.NewResponse(&adminv1.GetAnswerResponse{
		Answer: protoAnser,
	}), nil
}
