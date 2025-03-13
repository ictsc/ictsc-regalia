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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProblemServiceHandler struct {
	contestantv1connect.UnimplementedProblemServiceHandler

	Enforcer *ScheduleEnforcer

	ListProblemsEffect    ProblemListEffect
	GetProblemEffect      ProblemGetEffect
	ListDeploymentsEffect DeploymentsListEffect
}

var _ contestantv1connect.ProblemServiceHandler = (*ProblemServiceHandler)(nil)

func newProblemServiceHandler(enforcer *ScheduleEnforcer, repo *pg.Repository) *ProblemServiceHandler {
	return &ProblemServiceHandler{
		Enforcer:           enforcer,
		ListProblemsEffect: repo,
		GetProblemEffect:   repo,
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
	if err := h.Enforcer.Enforce(ctx, domain.PhaseInContest); err != nil {
		return nil, err
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.ListProblemsEffect)
	if err != nil {
		return nil, err
	}

	problems, err := teamMember.Team().Problems(ctx, h.ListProblemsEffect)
	if err != nil {
		return nil, err
	}

	protoProblems := make([]*contestantv1.Problem, 0, len(problems))
	for _, problem := range problems {
		protoProblems = append(protoProblems, convertTeamProblem(problem))
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
	if err := h.Enforcer.Enforce(ctx, domain.PhaseInContest); err != nil {
		return nil, err
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.GetProblemEffect)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetCode())
	if err != nil {
		return nil, err
	}

	teamProblem, err := teamMember.Team().ProblemDetailByCode(ctx, h.GetProblemEffect, code)
	if err != nil {
		return nil, err
	}
	detail := teamProblem.ProblemDetail()

	proto := convertTeamProblem(teamProblem.TeamProblem())
	proto.Body = &contestantv1.ProblemBody{
		Type: contestantv1.ProblemType_PROBLEM_TYPE_DESCRIPTIVE,
		Body: &contestantv1.ProblemBody_Descriptive{
			Descriptive: &contestantv1.DescriptiveProblem{
				Body: detail.Body(),
			},
		},
	}

	return connect.NewResponse(&contestantv1.GetProblemResponse{
		Problem: proto,
	}), nil
}

type DeploymentsListEffect interface {
	domain.TeamMemberGetter
	domain.ProblemReader
	domain.DeploymentReader
}

func (h *ProblemServiceHandler) ListDeployments(
	ctx context.Context,
	req *connect.Request[contestantv1.ListDeploymentsRequest],
) (*connect.Response[contestantv1.ListDeploymentsResponse], error) {
	userSess, err := session.UserSessionStore.Get(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	protoCode := req.Msg.GetCode()
	if protoCode == "" {
		return nil, domain.NewInvalidArgumentError("code is required", nil)
	}

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.ListDeploymentsEffect)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(protoCode)
	if err != nil {
		return nil, err
	}

	teamProblem, err := teamMember.Team().ProblemDetailByCode(ctx, h.ListDeploymentsEffect, code)
	if err != nil {
		return nil, err
	}

	deployments, err := teamProblem.TeamProblem().Deployments(ctx, h.ListDeploymentsEffect)
	if err != nil {
		return nil, err
	}

	protoDeployments := make([]*contestantv1.DeploymentRequest, 0, len(deployments))
	for _, deployment := range deployments {
		protoDeployments = append(protoDeployments, &contestantv1.DeploymentRequest{
			Revision:            deployment.Revision(),
			Status:              convertDeploymentStatus(deployment.Status()),
			RequestedAt:         timestamppb.New(deployment.CreatedAt()),
			AllowedRequestCount: teamProblem.TeamProblem().RemainingDeployments(deployment.Revision()),
			Penalty:             teamProblem.TeamProblem().Penalty(deployment.Revision()),
		})
	}

	return connect.NewResponse(&contestantv1.ListDeploymentsResponse{
		Deployments: protoDeployments,
	}), nil
}

func convertTeamProblem(problem *domain.TeamProblem) *contestantv1.Problem {
	proto := &contestantv1.Problem{
		Code:     string(problem.Code()),
		Title:    problem.Title(),
		MaxScore: problem.MaxScore(),
		Category: problem.Category(),
		Deployment: &contestantv1.Deployment{
			Redeployable: problem.Redeployable(),
		},
	}
	if score := problem.Score(); score != nil {
		proto.Score = &contestantv1.Score{
			MarkedScore: score.MarkedScore(),
			Penalty:     score.Penalty(),
			Score:       score.TotalScore(),
			MaxScore:    score.MaxScore(),
		}
	}
	if problem.RedeployRule() == domain.RedeployRulePercentagePenalty {
		proto.Deployment.PenaltyThreashold = problem.PercentagePenalty().Threshold
	}
	return proto
}

func convertDeploymentStatus(status domain.DeploymentStatus) contestantv1.DeploymentStatus {
	switch status {
	case domain.DeploymentStatusQueued:
		fallthrough
	case domain.DeploymentStatusCreating:
		return contestantv1.DeploymentStatus_DEPLOYMENT_STATUS_DEPLOYING
	case domain.DeploymentStatusCompleted:
		return contestantv1.DeploymentStatus_DEPLOYMENT_STATUS_DEPLOYED
	case domain.DeploymentStatusFailed:
		return contestantv1.DeploymentStatus_DEPLOYMENT_STATUS_FAILED
	case domain.DeploymentStatusUnknown:
		fallthrough
	default:
		return contestantv1.DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED
	}
}
