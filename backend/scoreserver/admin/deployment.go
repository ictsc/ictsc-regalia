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

type DeploymentServiceHandler struct {
	adminv1connect.UnimplementedDeploymentServiceHandler

	Enforcer     *auth.Enforcer
	ListEffect   domain.DeploymentReader
	UpdateEffect DeploymentStatusUpdateEffect
}

var _ adminv1connect.DeploymentServiceHandler = (*DeploymentServiceHandler)(nil)

func newDeploymentServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *DeploymentServiceHandler {
	return &DeploymentServiceHandler{
		Enforcer:   enforcer,
		ListEffect: repo,
		UpdateEffect: struct {
			domain.TeamGetter
			domain.TeamProblemReader
			domain.DeploymentReader
			domain.Tx[domain.DeploymentWriter]
		}{
			TeamGetter:        repo,
			TeamProblemReader: repo,
			DeploymentReader:  repo,
			Tx:                pg.Tx(repo, func(rt *pg.RepositoryTx) domain.DeploymentWriter { return rt }),
		},
	}
}

func (h *DeploymentServiceHandler) ListDeployments(
	ctx context.Context,
	req *connect.Request[adminv1.ListDeploymentsRequest],
) (*connect.Response[adminv1.ListDeploymentsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "deployments", "list"); err != nil {
		return nil, err
	}

	deployments, err := domain.ListDeployments(ctx, h.ListEffect)
	if err != nil {
		return nil, err
	}

	protoDeployments := make([]*adminv1.Deployment, 0, len(deployments))
	for _, deployment := range deployments {
		protoEvents := make([]*adminv1.DeploymentEvent, 0, len(deployment.Events()))
		for _, event := range deployment.Events() {
			protoEvents = append(protoEvents, &adminv1.DeploymentEvent{
				Type:       convertDeploymentStatus(event.Status()),
				OccurredAt: timestamppb.New(event.OccurredAt()),
			})
		}
		protoDeployments = append(protoDeployments, &adminv1.Deployment{
			TeamCode:    int64(deployment.TeamCode()),
			ProblemCode: string(deployment.ProblemCode()),
			Revision:    int64(deployment.Revision()),
			LatestEvent: convertDeploymentStatus(deployment.Status()),
			Events:      protoEvents,
		})
	}

	return connect.NewResponse(&adminv1.ListDeploymentsResponse{
		Deployments: protoDeployments,
	}), nil
}

type DeploymentStatusUpdateEffect interface {
	domain.TeamGetter
	domain.TeamProblemReader
	domain.DeploymentReader
	domain.Tx[domain.DeploymentWriter]
}

func (h *DeploymentServiceHandler) UpdateDeploymentStatus(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateDeploymentStatusRequest],
) (*connect.Response[adminv1.UpdateDeploymentStatusResponse], error) {
	if err := enforce(ctx, h.Enforcer, "deployments", "update"); err != nil {
		return nil, err
	}

	reqTeamCode := req.Msg.GetTeamCode()
	if reqTeamCode == 0 {
		return nil, domain.NewInvalidArgumentError("team_code is required", nil)
	}
	teamCode, err := domain.NewTeamCode(reqTeamCode)
	if err != nil {
		return nil, err
	}

	reqProblemCode := req.Msg.GetProblemCode()
	if reqProblemCode == "" {
		return nil, domain.NewInvalidArgumentError("problem_code is required", nil)
	}
	problemCode, err := domain.NewProblemCode(reqProblemCode)
	if err != nil {
		return nil, err
	}

	reqRevision := req.Msg.GetRevision()
	if reqRevision == 0 {
		return nil, domain.NewInvalidArgumentError("revision is required", nil)
	}

	status := parseDeploymentStatus(req.Msg.GetStatus())

	team, err := teamCode.Team(ctx, h.UpdateEffect)
	if err != nil {
		return nil, err
	}

	teamProblem, err := team.ProblemByCodeForAdmin(ctx, h.UpdateEffect, problemCode)
	if err != nil {
		return nil, err
	}

	deployment, err := teamProblem.DeploymentByRevision(ctx, h.UpdateEffect, reqRevision)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if err := h.UpdateEffect.RunInTx(ctx, func(eff domain.DeploymentWriter) error {
		return deployment.UpdateStatus(ctx, eff, status, now)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateDeploymentStatusResponse{}), nil
}

func convertDeploymentStatus(status domain.DeploymentStatus) adminv1.DeploymentEventType {
	switch status {
	case domain.DeploymentStatusQueued:
		return adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED
	case domain.DeploymentStatusCreating:
		return adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING
	case domain.DeploymentStatusCompleted:
		return adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_FINISHED
	case domain.DeploymentStatusFailed:
		return adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_ERROR
	case domain.DeploymentStatusUnknown:
		fallthrough
	default:
		return adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED
	}
}

func parseDeploymentStatus(status adminv1.DeploymentEventType) domain.DeploymentStatus {
	switch status {
	case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED:
		return domain.DeploymentStatusQueued
	case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING:
		return domain.DeploymentStatusCreating
	case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_FINISHED:
		return domain.DeploymentStatusCompleted
	case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_ERROR:
		return domain.DeploymentStatusFailed
	case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return domain.DeploymentStatusUnknown
	}
}
