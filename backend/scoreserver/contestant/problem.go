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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProblemServiceHandler struct {
	contestantv1connect.UnimplementedProblemServiceHandler

	ListProblemsEffect    ProblemListEffect
	GetProblemEffect      ProblemGetEffect
	ListDeploymentsEffect DeploymentsListEffect
	DeployEffect          DeployEffect
}

var _ contestantv1connect.ProblemServiceHandler = (*ProblemServiceHandler)(nil)

func newProblemServiceHandler(repo *pg.Repository) *ProblemServiceHandler {
	return &ProblemServiceHandler{
		ListProblemsEffect:    repo,
		GetProblemEffect:      repo,
		ListDeploymentsEffect: repo,
		DeployEffect: struct {
			domain.TeamMemberGetter
			domain.TeamProblemReader
			domain.Tx[domain.DeploymentWriter]
		}{
			TeamMemberGetter:  repo,
			TeamProblemReader: repo,
			Tx:                pg.Tx(repo, func(rt *pg.RepositoryTx) domain.DeploymentWriter { return rt }),
		},
	}
}

type ProblemListEffect interface {
	domain.TeamMemberGetter
	domain.TeamProblemReader
	domain.ScheduleReader
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

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.ListProblemsEffect)
	if err != nil {
		return nil, err
	}

	problems, err := teamMember.Team().ProblemsForPublic(ctx, h.ListProblemsEffect)
	if err != nil {
		return nil, err
	}

	// 可視性フィルタリング: 過去に一度でも提出可能だった問題のみ表示
	now := time.Now()
	visibleProblems := make([]*domain.TeamProblem, 0, len(problems))
	for _, problem := range problems {
		isVisible, err := problem.Problem().IsVisibleAt(ctx, now, h.ListProblemsEffect)
		if err != nil {
			return nil, err
		}
		if isVisible {
			visibleProblems = append(visibleProblems, problem)
		}
	}

	protoProblems := make([]*contestantv1.Problem, 0, len(visibleProblems))
	for _, problem := range visibleProblems {
		// 提出状態を計算
		submissionStatus, err := h.calculateSubmissionStatus(ctx, problem.Problem(), now)
		if err != nil {
			return nil, err
		}

		proto := convertTeamProblem(problem)
		proto.SubmissionStatus = submissionStatus
		protoProblems = append(protoProblems, proto)
	}

	return connect.NewResponse(&contestantv1.ListProblemsResponse{
		Problems: protoProblems,
	}), nil
}

func (h *ProblemServiceHandler) calculateSubmissionStatus(
	ctx context.Context,
	problem *domain.Problem,
	now time.Time,
) (*contestantv1.SubmissionStatus, error) {
	schedules, err := domain.GetSchedule(ctx, h.ListProblemsEffect)
	if err != nil {
		return nil, err
	}

	status := &contestantv1.SubmissionStatus{}

	var currentWindow, nextWindow *domain.ScheduleEntry

	for _, scheduleName := range problem.SubmissionableScheduleNames() {
		for _, entry := range schedules {
			if entry.Name() != scheduleName {
				continue
			}

			if !now.Before(entry.StartAt()) && now.Before(entry.EndAt()) {
				currentWindow = entry
				break
			}

			if now.Before(entry.StartAt()) {
				if nextWindow == nil || entry.StartAt().Before(nextWindow.StartAt()) {
					nextWindow = entry
				}
			}
		}
		if currentWindow != nil {
			break
		}
	}

	if currentWindow != nil {
		status.IsSubmittable = true
		status.SubmittableUntil = timestamppb.New(currentWindow.EndAt())
	} else if nextWindow != nil {
		status.IsSubmittable = false
		status.SubmittableFrom = timestamppb.New(nextWindow.StartAt())
	} else {
		status.IsSubmittable = false
	}

	return status, nil
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

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.GetProblemEffect)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(req.Msg.GetCode())
	if err != nil {
		return nil, err
	}

	teamProblem, err := teamMember.Team().ProblemDetailByCodeForPublic(ctx, h.GetProblemEffect, code)
	if err != nil {
		return nil, err
	}

	// 可視性チェック: まだ開始されていないスケジュールの問題はアクセス不可
	now := time.Now()
	isVisible, err := teamProblem.TeamProblem().Problem().IsVisibleAt(ctx, now, h.GetProblemEffect)
	if err != nil {
		return nil, err
	}
	if !isVisible {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}

	detail := teamProblem.ProblemDetail()
	submissionStatus, err := h.calculateSubmissionStatus(ctx, teamProblem.TeamProblem().Problem(), now)
	if err != nil {
		return nil, err
	}

	proto := convertTeamProblem(teamProblem.TeamProblem())
	proto.SubmissionStatus = submissionStatus
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
	domain.TeamProblemReader
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

	teamProblem, err := teamMember.Team().ProblemByCodeForPublic(ctx, h.ListDeploymentsEffect, code)
	if err != nil {
		return nil, err
	}

	deployments, err := teamProblem.Deployments(ctx, h.ListDeploymentsEffect)
	if err != nil {
		return nil, err
	}

	protoDeployments := make([]*contestantv1.DeploymentRequest, 0, len(deployments))
	for _, deployment := range deployments {
		protoDeployments = append(protoDeployments, convertDeployment(teamProblem.Problem(), deployment))
	}

	return connect.NewResponse(&contestantv1.ListDeploymentsResponse{
		Deployments: protoDeployments,
	}), nil
}

type DeployEffect interface {
	domain.TeamMemberGetter
	domain.TeamProblemReader
	domain.Tx[domain.DeploymentWriter]
}

func (h *ProblemServiceHandler) Deploy(
	ctx context.Context,
	req *connect.Request[contestantv1.DeployRequest],
) (*connect.Response[contestantv1.DeployResponse], error) {
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

	teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.DeployEffect)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewProblemCode(protoCode)
	if err != nil {
		return nil, err
	}

	teamProblem, err := teamMember.Team().ProblemByCodeForPublic(ctx, h.DeployEffect, code)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	deployment, err := domain.RunTx(ctx, h.DeployEffect,
		func(eff domain.DeploymentWriter) (*domain.Deployment, error) {
			return teamProblem.Deploy(ctx, eff, now)
		},
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&contestantv1.DeployResponse{
		Deployment: convertDeployment(teamProblem.Problem(), deployment),
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

func convertDeployment(problem *domain.Problem, deployment *domain.Deployment) *contestantv1.DeploymentRequest {
	return &contestantv1.DeploymentRequest{
		Revision:            deployment.Revision(),
		Status:              convertDeploymentStatus(deployment.Status()),
		RequestedAt:         timestamppb.New(deployment.CreatedAt()),
		AllowedRequestCount: problem.RemainingDeployments(deployment.Revision()),
		Penalty:             problem.Penalty(deployment.Revision()),
	}
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
