package memory

import "github.com/ictsc/ictsc-regalia/backend/internal/core"

func sameDeploymentEvent(left, right core.DeploymentEvent) bool {
	if left.EventID != right.EventID || !left.OccurredAt.Equal(right.OccurredAt) || left.Status != right.Status {
		return false
	}
	if left.Message == nil || right.Message == nil {
		return left.Message == nil && right.Message == nil
	}
	return *left.Message == *right.Message
}
