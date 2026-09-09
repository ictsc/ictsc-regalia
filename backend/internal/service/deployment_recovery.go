package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func definitiveDeploymentRejection(err error) bool {
	var classified interface{ DefinitiveDeploymentRejection() bool }
	return errors.As(err, &classified) && classified.DefinitiveDeploymentRejection()
}

// FailQueuedDeployment resolves a definitively rejected local queue entry using
// only transitions allowed by the deployment contract. It is also used by the
// operator-only sync path after SState confirms that no request exists.
func (s *Service) FailQueuedDeployment(ctx context.Context, deployment core.Deployment, message string) error {
	transitionID, err := randomUUID()
	if err != nil {
		return core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create deployment recovery event", err)
	}
	failedID, err := randomUUID()
	if err != nil {
		return core.WrapError(http.StatusInternalServerError, "internal_error", "Could not create deployment failure event", err)
	}
	occurredAt := s.Now().UTC().Truncate(time.Microsecond)
	deploying, _, err := s.Store.AppendDeploymentEvent(ctx, deployment.RequestID, core.DeploymentEvent{
		EventID: transitionID, OccurredAt: occurredAt, Status: core.DeploymentDeploying,
	}, false)
	if err != nil {
		return err
	}
	_ = s.publishDeployment(ctx, deploying)
	failed, _, err := s.Store.AppendDeploymentEvent(ctx, deployment.RequestID, core.DeploymentEvent{
		EventID: failedID, OccurredAt: occurredAt.Add(time.Microsecond), Status: core.DeploymentFailed, Message: &message,
	}, false)
	if err != nil {
		return err
	}
	_ = s.publishDeployment(ctx, failed)
	return nil
}
