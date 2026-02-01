package contestant

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnswerServiceHandler struct {
	contestantv1connect.UnimplementedAnswerServiceHandler

	ListEffect   ListAnswersEffect
	GetEffect    GetAnswerEffect
	SubmitEffect SubmitAnswerEffect
}

var _ contestantv1connect.AnswerServiceHandler = (*AnswerServiceHandler)(nil)

func newAnswerServiceHandler(repo *pg.Repository) *AnswerServiceHandler {
	return &AnswerServiceHandler{
		ListEffect: repo,
		GetEffect:  repo,
		SubmitEffect: struct {
			domain.TeamMemberGetter
			domain.ProblemReader
			domain.ScheduleReader
			domain.Tx[domain.AnswerWriter]
		}{
			TeamMemberGetter: repo,
			ProblemReader:    repo,
			ScheduleReader:   repo,
			Tx:               pg.Tx(repo, func(rt *pg.RepositoryTx) domain.AnswerWriter { return rt }),
		},
	}
}

type ListAnswersEffect interface {
	domain.TeamMemberGetter
	domain.AnswerReader
}

func (h *AnswerServiceHandler) ListAnswers(
	ctx context.Context,
	req *connect.Request[contestantv1.ListAnswersRequest],
) (*connect.Response[contestantv1.ListAnswersResponse], error) {
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

	protoProblemCode := req.Msg.GetProblemCode()
	if protoProblemCode == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("problem_code is required"))
	}
	problemCode, err := domain.NewProblemCode(protoProblemCode)
	if err != nil {
		return nil, err
	}

	answers, err := domain.ListAnswersByTeamProblemForPublic(ctx, h.ListEffect, teamMember.Team().Code(), problemCode)
	if err != nil {
		return nil, err
	}

	protoAnswers := make([]*contestantv1.Answer, 0, len(answers))
	latestSubmitTime := time.Time{}
	for _, answer := range answers {
		if latestSubmitTime == (time.Time{}) || answer.CreatedAt().After(latestSubmitTime) {
			latestSubmitTime = answer.CreatedAt()
		}
		protoAnswers = append(protoAnswers, convertAnswer(answer))
	}

	resp := &contestantv1.ListAnswersResponse{
		Answers:        protoAnswers,
		SubmitInterval: durationpb.New(domain.AnswerInterval),
	}
	if latestSubmitTime != (time.Time{}) {
		resp.LastSubmittedAt = timestamppb.New(latestSubmitTime)
	}

	return connect.NewResponse(resp), nil
}

type GetAnswerEffect interface {
	domain.TeamMemberGetter
	domain.AnswerReader
}

func (h *AnswerServiceHandler) GetAnswer(
	ctx context.Context,
	req *connect.Request[contestantv1.GetAnswerRequest],
) (*connect.Response[contestantv1.GetAnswerResponse], error) {
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

	protoProblemCode := req.Msg.GetProblemCode()
	if protoProblemCode == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("problem_code is required"))
	}
	problemCode, err := domain.NewProblemCode(protoProblemCode)
	if err != nil {
		return nil, err
	}

	answerDetail, err := domain.GetAnswerDetailForPublic(ctx, h.ListEffect, teamMember.Team().Code(), problemCode, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	protoAnswer, err := convertAnswerDetail(answerDetail)
	if err != nil {
		return nil, err
	}

	resp := &contestantv1.GetAnswerResponse{
		Answer: protoAnswer,
	}

	return connect.NewResponse(resp), nil
}

type SubmitAnswerEffect interface {
	domain.TeamMemberGetter
	domain.ProblemReader
	domain.ScheduleReader
	domain.Tx[domain.AnswerWriter]
}

func (h *AnswerServiceHandler) SubmitAnswer(
	ctx context.Context,
	req *connect.Request[contestantv1.SubmitAnswerRequest],
) (*connect.Response[contestantv1.SubmitAnswerResponse], error) {
	userSess, err := session.UserSessionStore.Get(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.SubmitEffect)
	if err != nil {
		return nil, err
	}

	protoProblemCode := req.Msg.GetProblemCode()
	if protoProblemCode == "" {
		return nil, domain.NewInvalidArgumentError("problem_code is required", nil)
	}
	problemCode, err := domain.NewProblemCode(protoProblemCode)
	if err != nil {
		return nil, err
	}
	problem, err := problemCode.Problem(ctx, h.SubmitEffect)
	if err != nil {
		return nil, err
	}

	// スケジュールベースの提出可否チェック
	now := time.Now()
	isSubmittable, err := problem.IsSubmittableAt(ctx, now, h.SubmitEffect)
	if err != nil {
		return nil, err
	}
	if !isSubmittable {
		return nil, connect.NewError(
			connect.CodeFailedPrecondition,
			errors.New("この問題は現在提出できません"),
		)
	}

	body := req.Msg.GetBody()
	if body == "" {
		return nil, domain.NewInvalidArgumentError("body is required", nil)
	}
	answer, err := domain.RunTx(ctx, h.SubmitEffect, func(tx domain.AnswerWriter) (*domain.AnswerDetail, error) {
		return teamMember.SubmitDescriptiveAnswer(ctx, now, tx, problem, body)
	})
	if err != nil {
		if errors.Is(err, domain.ErrTooEarlyToSubmitAnswer) {
			return nil, connect.NewError(connect.CodeFailedPrecondition, domain.ErrTooEarlyToSubmitAnswer)
		}
		return nil, err
	}

	protoAnswer, err := convertAnswerDetail(answer)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&contestantv1.SubmitAnswerResponse{
		Answer: protoAnswer,
	}), nil
}

func convertAnswer(answer *domain.Answer) *contestantv1.Answer {
	protoAnswer := &contestantv1.Answer{
		Id:          answer.Number(),
		SubmittedAt: timestamppb.New(answer.CreatedAt()),
	}
	switch answer.Problem().Type() {
	case domain.ProblemTypeDescriptive:
		protoAnswer.Body = &contestantv1.AnswerBody{
			Type: contestantv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE,
		}
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
		protoAnswer.Body = &contestantv1.AnswerBody{
			Type: contestantv1.ProblemType_PROBLEM_TYPE_UNSPECIFIED,
		}
	}
	if score := answer.Score(); score != nil {
		protoAnswer.Score = &contestantv1.Score{
			MarkedScore: score.MarkedScore(),
			Penalty:     score.Penalty(),
			Score:       score.TotalScore(),
			MaxScore:    score.MaxScore(),
		}
	}
	return protoAnswer
}

func convertAnswerDetail(answer *domain.AnswerDetail) (*contestantv1.Answer, error) {
	proto := convertAnswer(answer.Answer())
	switch answer.Body().Type() {
	case domain.ProblemTypeDescriptive:
		desc, err := answer.Body().Descriptive()
		if err != nil {
			return nil, err
		}
		proto.Body = &contestantv1.AnswerBody{
			Type: contestantv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE,
			Body: &contestantv1.AnswerBody_Descriptive{
				Descriptive: &contestantv1.DescriptiveAnswer{
					Body: desc.Body(),
				},
			},
		}
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
	}
	return proto, nil
}
