package rstate

import "net/http"

func (e *NotFoundError) DeploymentStatusNotFound() bool { return true }

func (e *HTTPError) DeploymentStatusNotFound() bool { return false }

// DefinitiveDeploymentRejection covers 400 (bad team or problem) and 409 (a
// job for the same team and problem is already waiting or running).
func (e *HTTPError) DefinitiveDeploymentRejection() bool {
	return e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError &&
		e.StatusCode != http.StatusTooManyRequests && e.StatusCode != http.StatusUnauthorized
}
