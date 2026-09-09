package httpserver

import "github.com/ictsc/ictsc-regalia/backend/internal/core"

func latestDeploymentEventID(deployment core.Deployment) string {
	latest := core.DeploymentEvent{}
	for _, event := range deployment.Events {
		if latest.EventID == "" || event.OccurredAt.After(latest.OccurredAt) ||
			(event.OccurredAt.Equal(latest.OccurredAt) && event.EventID > latest.EventID) {
			latest = event
		}
	}
	if latest.EventID != "" {
		return latest.EventID
	}
	return deployment.RequestID
}
