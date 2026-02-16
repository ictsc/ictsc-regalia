package admin

import (
	"context"
	"time"

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

	Enforcer                  *auth.Enforcer
	ListAnswerEffect          domain.AnswerReader
	GetAnswerEffect           domain.AnswerReader
	ListMarkingResultEffect   domain.MarkingResultReader
	CreateMarkingResultEffect domain.Tx[CreateMarkingResultTxEffect]
	UpdateScoreEffect         UpdateScoreEffect
}

var _ adminv1connect.MarkServiceHandler = (*MarkServiceHandler)(nil)

func newMarkServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *MarkServiceHandler {
	return &MarkServiceHandler{
		Enforcer:                  enforcer,
		ListAnswerEffect:          repo,
		GetAnswerEffect:           repo,
		ListMarkingResultEffect:   repo,
		CreateMarkingResultEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) CreateMarkingResultTxEffect { return rt }),
		UpdateScoreEffect:         newUpdateScoreEffect(repo),
	}
}

func (h *MarkServiceHandler) ListAnswers(
	ctx context.Context,
	req *connect.Request[adminv1.ListAnswersRequest],
) (*connect.Response[adminv1.ListAnswersResponse], error) {
	if err := enforce(ctx, h.Enforcer, "answers", "list"); err != nil {
		return nil, err
	}

	answers, err := domain.ListAnswersForAdmin(ctx, h.ListAnswerEffect)
	if err != nil {
		return nil, err
	}

	protoAnswers := make([]*adminv1.Answer, 0, len(answers))
	for _, answer := range answers {
		protoAnswers = append(protoAnswers, convertAnswer(answer))
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

	answer, err := domain.GetAnswerDetailForAdmin(ctx, h.GetAnswerEffect, teamCode, problemCode, protoID)
	if err != nil {
		return nil, err
	}

	protoAnser := convertAnswer(answer.Answer())
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

func (h *MarkServiceHandler) ListMarkingResults(
	ctx context.Context,
	req *connect.Request[adminv1.ListMarkingResultsRequest],
) (*connect.Response[adminv1.ListMarkingResultsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "marking_results", "list"); err != nil {
		return nil, err
	}

	markingResults, err := domain.ListAllMarkingResults(ctx, h.ListMarkingResultEffect)
	if err != nil {
		return nil, err
	}

	protoMarkingResults := make([]*adminv1.MarkingResult, 0, len(markingResults))
	for _, markingResult := range markingResults {
		protoMarkingResults = append(protoMarkingResults, convertMarkingResult(markingResult))
	}

	return connect.NewResponse(&adminv1.ListMarkingResultsResponse{
		MarkingResults: protoMarkingResults,
	}), nil
}

type CreateMarkingResultTxEffect interface {
	domain.AnswerReader
	domain.MarkingResultWriter
}

func (h *MarkServiceHandler) CreateMarkingResult(
	ctx context.Context,
	req *connect.Request[adminv1.CreateMarkingResultRequest],
) (*connect.Response[adminv1.CreateMarkingResultResponse], error) {
	if err := enforce(ctx, h.Enforcer, "marking_results", "create"); err != nil {
		return nil, err
	}

	viewer := auth.GetViewer(ctx)

	reqAnswer := req.Msg.GetMarkingResult().GetAnswer()

	reqTeamCode := reqAnswer.GetTeam().GetCode()
	if reqTeamCode == 0 {
		return nil, domain.NewInvalidArgumentError("team_code is required", nil)
	}
	teamCode, err := domain.NewTeamCode(reqTeamCode)
	if err != nil {
		return nil, err
	}

	reqProblemCode := reqAnswer.GetProblem().GetCode()
	if reqProblemCode == "" {
		return nil, domain.NewInvalidArgumentError("problem_code is required", nil)
	}
	problemCode, err := domain.NewProblemCode(reqProblemCode)
	if err != nil {
		return nil, err
	}

	reqAnswerID := reqAnswer.GetId()
	if reqAnswerID == 0 {
		return nil, domain.NewInvalidArgumentError("answer_id is required", nil)
	}

	reqScore := req.Msg.GetMarkingResult().GetScore()
	reqDescriptiveComment := req.Msg.GetMarkingResult().GetRationale().GetDescriptive().GetComment()

	now := time.Now()

	markingResult, err := domain.RunTx(ctx, h.CreateMarkingResultEffect, func(eff CreateMarkingResultTxEffect) (*domain.MarkingResult, error) {
		answerDetail, err := domain.GetAnswerDetailForAdmin(ctx, eff, teamCode, problemCode, reqAnswerID)
		if err != nil {
			return nil, err
		}

		return answerDetail.Answer().Mark(ctx, eff, now, &domain.MarkInput{
			Score: reqScore,
			Judge: viewer.Name,

			Comment: reqDescriptiveComment,
		})
	})
	if err != nil {
		return nil, err
	}

	if _, err := UpdateScore(ctx, h.UpdateScoreEffect, now); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.CreateMarkingResultResponse{
		MarkingResult: convertMarkingResult(markingResult),
	}), nil
}

func (h *MarkServiceHandler) UpdateScores(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateScoresRequest],
) (*connect.Response[adminv1.UpdateScoresResponse], error) {
	if err := enforce(ctx, h.Enforcer, "scores", "update"); err != nil {
		return nil, err
	}

	if _, err := UpdateScore(ctx, h.UpdateScoreEffect, time.Now()); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateScoresResponse{}), nil
}

func convertAnswer(answer *domain.Answer) *adminv1.Answer {
	proto := &adminv1.Answer{
		Id:        answer.Number(),
		Team:      convertTeam(answer.Team()),
		Problem:   convertProblem(answer.Problem()),
		CreatedAt: timestamppb.New(answer.CreatedAt()),
	}
	if score := answer.Score(); score != nil {
		proto.Score = &adminv1.MarkingScore{
			Total:   score.TotalScore(),
			Marked:  score.MarkedScore(),
			Penalty: score.Penalty(),
			Max:     score.MaxScore(),
		}
	}
	return proto
}

func convertMarkingResult(markingResult *domain.MarkingResult) *adminv1.MarkingResult {
	proto := &adminv1.MarkingResult{
		Answer:    convertAnswer(markingResult.Answer()),
		Judge:     &adminv1.Admin{Name: markingResult.Judge()},
		Score:     markingResult.Score().MarkedScore(),
		CreatedAt: timestamppb.New(markingResult.CreatedAt()),
	}
	rationale := markingResult.Rationale()
	switch rationale.Type() {
	case domain.ProblemTypeDescriptive:
		proto.Rationale = &adminv1.MarkingRationale{
			Type: adminv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE,
			Body: &adminv1.MarkingRationale_Descriptive{
				Descriptive: &adminv1.DescriptiveMarkingRationale{
					Comment: rationale.Descriptive().Comment(),
				},
			},
		}
	case domain.ProblemTypeUnknown:
		fallthrough
	default:
		proto.Rationale = &adminv1.MarkingRationale{Type: adminv1.ProblemType_PROBLEM_TYPE_UNSPECIFIED}
	}
	return proto
}
