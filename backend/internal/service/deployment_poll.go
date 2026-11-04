package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// polledJobClockSkew lets a job whose creation time reads slightly before the
// request still count as that request's job, since the deployer's clock is
// not ours. Older jobs belong to an earlier redeploy and are ignored.
const polledJobClockSkew = 2 * time.Minute

// polledEventNamespace scopes the deterministic IDs of polled events.
var polledEventNamespace = uuid.MustParse("6f1d7a52-3c0e-4b8a-9d55-2a0f4c1e9b37")

// PollDeployments moves every active deployment forward from what a
// DeploymentProber reports. It is a no-op for gateways that call back.
func (s *Service) PollDeployments(ctx context.Context) error {
	prober, ok := s.Deployments.(DeploymentProber)
	if !ok {
		return nil
	}
	deployments, err := s.Store.ListDeployments(ctx, core.DeploymentFilter{})
	if err != nil {
		return err
	}
	type target struct {
		team    int64
		problem string
	}
	latest := make(map[target]core.Deployment, len(deployments))
	for _, deployment := range deployments {
		key := target{deployment.TeamCode, deployment.ProblemCode}
		if current, ok := latest[key]; !ok || deployment.Revision > current.Revision {
			latest[key] = deployment
		}
	}
	var errs []error
	for _, deployment := range latest {
		if deployment.LatestStatus != core.DeploymentQueued && deployment.LatestStatus != core.DeploymentDeploying {
			continue
		}
		if err := s.pollDeployment(ctx, prober, deployment); err != nil {
			errs = append(errs, fmt.Errorf("poll deployment %d/%s revision %d: %w", deployment.TeamCode, deployment.ProblemCode, deployment.Revision, err))
		}
	}
	return errors.Join(errs...)
}

func (s *Service) pollDeployment(ctx context.Context, prober DeploymentProber, deployment core.Deployment) error {
	observation, err := prober.Observe(ctx, deployment.TeamCode, deployment.ProblemCode)
	if err != nil {
		var notFound interface{ DeploymentStatusNotFound() bool }
		if errors.As(err, &notFound) && notFound.DeploymentStatusNotFound() {
			return nil
		}
		return err
	}
	if observation.CreatedAt.Before(deployment.RequestedAt.Add(-polledJobClockSkew)) {
		return nil
	}
	var steps []core.DeploymentStatus
	switch observation.Status {
	case core.DeploymentDeploying:
		if deployment.LatestStatus == core.DeploymentQueued {
			steps = []core.DeploymentStatus{core.DeploymentDeploying}
		}
	case core.DeploymentCompleted, core.DeploymentFailed:
		if deployment.LatestStatus == core.DeploymentQueued {
			steps = append(steps, core.DeploymentDeploying)
		}
		steps = append(steps, observation.Status)
	}
	for _, status := range steps {
		// The deployer's timestamps may run behind ours and the store refuses
		// events older than the latest one, so polled events use our clock.
		event := core.DeploymentEvent{
			EventID:    polledEventID(deployment.RequestID, status),
			OccurredAt: s.nextDeploymentEventTime(deployment),
			Status:     status,
		}
		if status == observation.Status {
			event.Message = observation.Message
		}
		updated, _, err := s.ApplyDeploymentEvent(ctx, deployment.TeamCode, deployment.ProblemCode, deployment.Revision, event, false)
		if err != nil {
			return err
		}
		deployment = updated
	}
	return nil
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
