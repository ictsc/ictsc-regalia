package batch

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/sstate"
)

type DeploymentSync struct {
	period           time.Duration
	deploymentClient adminv1connect.DeploymentServiceClient
	sstateClient     *sstate.SStateClient
}

func NewDeploymentSync(
	cfg config.DeploySync,
	deploymentClient adminv1connect.DeploymentServiceClient,
	sstateClient *sstate.SStateClient,
) *DeploymentSync {
	return &DeploymentSync{
		period:           cfg.Period,
		deploymentClient: deploymentClient,
		sstateClient:     sstateClient,
	}
}

func (d *DeploymentSync) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Start sync deployments")
	ticker := time.NewTicker(d.period)
	for {
		if err := d.sync(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to sync deployments", "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (d *DeploymentSync) sync(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "DeploymentSync.sync")
	defer span.End()
	slog.InfoContext(ctx, "Sync deployments")

	ctx, stop := context.WithTimeout(ctx, d.period)
	defer stop()

	deploymentsResp, err := d.deploymentClient.ListDeployments(
		ctx, connect.NewRequest(&adminv1.ListDeploymentsRequest{}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to list deployments")
	}

	for _, deployment := range deploymentsResp.Msg.GetDeployments() {
		logger := slog.With(
			"team", deployment.GetTeamCode(),
			"problem", deployment.GetProblemCode(),
			"revision", deployment.GetRevision(),
		)
		switch deployment.GetLatestEvent() {
		case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED:
			if err := d.queue(ctx, logger, deployment); err != nil {
				logger.ErrorContext(ctx, "failed to queue deployment", "error", err)
			}
		case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING,
			adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_ERROR:
			if err := d.updateState(ctx, logger, deployment); err != nil {
				logger.ErrorContext(ctx, "failed to check deployment state", "error", err)
			}
		case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_UNSPECIFIED:
			fallthrough
		case adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_FINISHED:
			fallthrough
		default:
			continue
		}
	}

	return nil
}

func (d *DeploymentSync) queue(ctx context.Context, logger *slog.Logger, deployment *adminv1.Deployment) error {
	ctx, span := tracer.Start(ctx, "DeploymentSync.queue")
	defer span.End()

	logger.InfoContext(ctx, "Queue deployment")

	var nextStatus adminv1.DeploymentEventType
	err := d.sstateClient.Redeploy(
		ctx, deployment.GetTeamCode(), deployment.GetProblemCode(),
	)
	switch {
	case err != nil && errors.Is(err, sstate.ErrProblemNotFound):
		logger.WarnContext(ctx, "Not managed by sstate")
		return nil
	case err != nil && errors.Is(err, sstate.ErrAlreadyDeploying):
		logger.WarnContext(ctx, "Already queued")
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING
	case err != nil:
		logger.ErrorContext(ctx, "Failed to redeploy", "error", err)
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED
	default:
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING
	}

	if nextStatus != deployment.GetLatestEvent() {
		if _, err := d.deploymentClient.UpdateDeploymentStatus(ctx,
			connect.NewRequest(&adminv1.UpdateDeploymentStatusRequest{
				TeamCode:    deployment.GetTeamCode(),
				ProblemCode: deployment.GetProblemCode(),
				Revision:    uint32(deployment.GetRevision()), //nolint:gosec
				Status:      nextStatus,
			}),
		); err != nil {
			return errors.Wrap(err, "failed to update deployment status")
		}
	}

	return nil
}

func (d *DeploymentSync) updateState(ctx context.Context, logger *slog.Logger, deployment *adminv1.Deployment) error {
	ctx, span := tracer.Start(ctx, "DeploymentSync.updateState")
	defer span.End()
	logger.InfoContext(ctx, "Update deployment state")

	status, _, err := d.sstateClient.GetStatus(ctx, deployment.GetTeamCode(), deployment.GetProblemCode())
	if err != nil {
		if errors.Is(err, sstate.ErrProblemNotFound) {
			logger.WarnContext(ctx, "Not managed by sstate")
			return nil
		}
		return err
	}

	var nextStatus adminv1.DeploymentEventType
	switch status {
	case domain.DeploymentStatusCompleted:
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_FINISHED
	case domain.DeploymentStatusCreating:
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_CREATING
	case domain.DeploymentStatusQueued:
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_QUEUED
	case domain.DeploymentStatusFailed:
		nextStatus = adminv1.DeploymentEventType_DEPLOYMENT_EVENT_TYPE_ERROR
	case domain.DeploymentStatusUnknown:
		fallthrough
	default:
	}

	if nextStatus != deployment.GetLatestEvent() {
		if _, err := d.deploymentClient.UpdateDeploymentStatus(ctx,
			connect.NewRequest(&adminv1.UpdateDeploymentStatusRequest{
				TeamCode:    deployment.GetTeamCode(),
				ProblemCode: deployment.GetProblemCode(),
				Revision:    uint32(deployment.GetRevision()), //nolint:gosec
				Status:      nextStatus,
			}),
		); err != nil {
			return errors.Wrap(err, "failed to update deployment status")
		}
	}

	return nil
}
