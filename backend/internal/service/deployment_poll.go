package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// polledEventNamespace scopes the deterministic IDs of polled events.
var polledEventNamespace = uuid.MustParse("6f1d7a52-3c0e-4b8a-9d55-2a0f4c1e9b37")

var activeDeploymentStatuses = []core.DeploymentStatus{core.DeploymentQueued, core.DeploymentDeploying}

// PollDeployments moves every active deployment forward from one job list of
// a DeploymentProber. It is a no-op for gateways that call back.
func (s *Service) PollDeployments(ctx context.Context) error {
	prober, ok := s.Deployments.(DeploymentProber)
	if !ok {
		return nil
	}
	deployments, err := s.Store.ListDeployments(ctx, core.DeploymentFilter{Statuses: activeDeploymentStatuses})
	if err != nil || len(deployments) == 0 {
		return err
	}
	jobs, err := prober.ListJobs(ctx)
	if err != nil {
		return err
	}
	var errs []error
	for _, deployment := range deployments {
		if _, err := s.advancePolledDeployment(ctx, deployment, jobs); err != nil {
			errs = append(errs, fmt.Errorf("poll deployment %d/%s revision %d: %w", deployment.TeamCode, deployment.ProblemCode, deployment.Revision, err))
		}
	}
	return errors.Join(errs...)
}

// SyncPolledDeployment is the operator sync for a DeploymentProber: one job
// list read applied through the same staged transitions as polling, so a
// job that finished before the first tick or a repeated sync never asks the
// store for a transition it refuses. A finished deployment is left as is.
func (s *Service) SyncPolledDeployment(ctx context.Context, prober DeploymentProber, deployment core.Deployment) error {
	if deployment.LatestStatus != core.DeploymentQueued && deployment.LatestStatus != core.DeploymentDeploying {
		return nil
	}
	jobs, err := prober.ListJobs(ctx)
	if err != nil {
		return core.WrapError(http.StatusBadGateway, "upstream_unavailable", "Deployment job lookup failed", err)
	}
	found, err := s.advancePolledDeployment(ctx, deployment, jobs)
	if err != nil || found {
		return err
	}
	if deployment.LatestStatus == core.DeploymentQueued {
		if err := s.FailQueuedDeployment(ctx, deployment, "The redeploy service has no job for this request"); err != nil {
			return err
		}
	}
	return core.NewError(http.StatusNotFound, "resource_not_found", "The redeploy service has no job for this deployment")
}

// advancePolledDeployment applies the state of the deployment's own job and
// reports whether that job was found.
func (s *Service) advancePolledDeployment(ctx context.Context, deployment core.Deployment, jobs []DeploymentJob) (bool, error) {
	job, ok := deploymentJob(deployment, jobs)
	if !ok {
		return false, nil
	}
	if deployment.GatewayJobID == "" {
		if err := s.Store.SetDeploymentJobID(ctx, deployment.RequestID, job.ID); err != nil {
			return true, err
		}
	}
	var steps []core.DeploymentStatus
	switch job.Status {
	case core.DeploymentDeploying:
		if deployment.LatestStatus == core.DeploymentQueued {
			steps = []core.DeploymentStatus{core.DeploymentDeploying}
		}
	case core.DeploymentCompleted, core.DeploymentFailed:
		if deployment.LatestStatus == core.DeploymentQueued {
			steps = append(steps, core.DeploymentDeploying)
		}
		if deployment.LatestStatus != job.Status {
			steps = append(steps, job.Status)
		}
	}
	for _, status := range steps {
		// The gateway keeps no event times we can trust against ours, and
		// the store refuses events older than the latest one.
		event := core.DeploymentEvent{
			EventID:    polledEventID(deployment.RequestID, status),
			OccurredAt: s.nextDeploymentEventTime(deployment),
			Status:     status,
		}
		if status == job.Status {
			event.Message = job.Message
		}
		updated, _, err := s.ApplyDeploymentEvent(ctx, deployment.TeamCode, deployment.ProblemCode, deployment.Revision, event, false)
		if err != nil {
			return true, err
		}
		deployment = updated
	}
	return true, nil
}

// deploymentJob finds the job by the ID recorded for the deployment. Before
// an ID is recorded, only a waiting or running job of the same team and
// problem can be its job: the gateway runs one at a time per problem, and a
// finished one may be from an earlier request.
func deploymentJob(deployment core.Deployment, jobs []DeploymentJob) (DeploymentJob, bool) {
	for _, job := range jobs {
		if deployment.GatewayJobID != "" {
			if job.ID == deployment.GatewayJobID {
				return job, true
			}
			continue
		}
		if job.TeamCode == deployment.TeamCode && strings.EqualFold(job.ProblemCode, deployment.ProblemCode) &&
			(job.Status == core.DeploymentQueued || job.Status == core.DeploymentDeploying) {
			return job, true
		}
	}
	return DeploymentJob{}, false
}

func (s *Service) nextDeploymentEventTime(deployment core.Deployment) time.Time {
	next := s.Now().UTC().Truncate(time.Microsecond)
	for _, event := range deployment.Events {
		if !next.After(event.OccurredAt) {
			next = event.OccurredAt.UTC().Truncate(time.Microsecond).Add(time.Microsecond)
		}
	}
	return next
}

// polledEventID is stable per deployment and status, so a status applied by
// an overlapping poll is recognized by its event ID.
func polledEventID(requestID string, status core.DeploymentStatus) string {
	return uuid.NewSHA1(polledEventNamespace, []byte(requestID+":"+string(status))).String()
}
