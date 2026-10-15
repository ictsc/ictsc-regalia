package sstate

import "net/http"

func (e *HTTPError) DeploymentStatusNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *HTTPError) DefinitiveDeploymentRejection() bool {
	return e.StatusCode >= http.StatusBadRequest && e.StatusCode < http.StatusInternalServerError &&
		e.StatusCode != http.StatusTooManyRequests
}
