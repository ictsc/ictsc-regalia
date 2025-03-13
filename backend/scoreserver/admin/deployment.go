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

type DeploymentServiceHandler struct {
	adminv1connect.UnimplementedDeploymentServiceHandler

	Enforcer   *auth.Enforcer
	ListEffect domain.DeploymentReader
}

var _ adminv1connect.DeploymentServiceHandler = (*DeploymentServiceHandler)(nil)

func newDeploymentServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *DeploymentServiceHandler {
	return &DeploymentServiceHandler{
		Enforcer:   enforcer,
		ListEffect: repo,
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
