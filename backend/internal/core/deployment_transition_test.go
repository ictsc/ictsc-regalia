package core

import "testing"

func TestCanTransitionDeploymentMatrix(t *testing.T) {
	statuses := []DeploymentStatus{DeploymentQueued, DeploymentDeploying, DeploymentCompleted, DeploymentFailed}
	for _, recovery := range []bool{false, true} {
		for _, from := range statuses {
			for _, to := range statuses {
				want := from == DeploymentQueued && to == DeploymentDeploying ||
					from == DeploymentDeploying && (to == DeploymentCompleted || to == DeploymentFailed) ||
					recovery && from == DeploymentFailed && to == DeploymentDeploying
				if got := CanTransitionDeployment(from, to, recovery); got != want {
					t.Fatalf("transition %s -> %s (recovery=%v) = %v, want %v", from, to, recovery, got, want)
				}
			}
		}
	}
}
