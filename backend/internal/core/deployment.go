package core

func CanTransitionDeployment(from, to DeploymentStatus, recovery bool) bool {
	switch from {
	case DeploymentQueued:
		return to == DeploymentDeploying
	case DeploymentDeploying:
		return to == DeploymentCompleted || to == DeploymentFailed
	case DeploymentFailed:
		return recovery && to == DeploymentDeploying
	default:
		return false
	}
}
