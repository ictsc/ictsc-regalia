package github

type oidcUpstreamError struct{ cause error }

func (e *oidcUpstreamError) Error() string              { return e.cause.Error() }
func (e *oidcUpstreamError) Unwrap() error              { return e.cause }
func (e *oidcUpstreamError) MachineTokenUpstream() bool { return true }
