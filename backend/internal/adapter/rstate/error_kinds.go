package rstate

import "net/http"

func (e *NotFoundError) DeploymentStatusNotFound() bool { return true }

func (e *TargetError) DefinitiveDeploymentRejection() bool { return true }

func (e *HTTPError) DeploymentStatusNotFound() bool { return false }

// DefinitiveDeploymentRejection covers postJob's 400, which means a bad team
// or problem ID, a job already waiting or running, or the preparatory phase.
func (e *HTTPError) DefinitiveDeploymentRejection() bool {
	return e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError &&
		e.StatusCode != http.StatusTooManyRequests && e.StatusCode != http.StatusUnauthorized
}
