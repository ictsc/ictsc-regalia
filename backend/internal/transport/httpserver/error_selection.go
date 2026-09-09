package httpserver

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// preferredDomainError selects the most actionable error when an OpenAPI
// security requirement contains alternatives. In particular, an identity
// provider outage must not be hidden by the missing-cookie 401 from another
// authentication alternative.
func preferredDomainError(err error) *core.Error {
	candidates := make([]*core.Error, 0, 2)
	collectDomainErrors(err, &candidates)
	if len(candidates) == 0 {
		var domainErr *core.Error
		if errors.As(err, &domainErr) {
			return domainErr
		}
		return nil
	}
	selected := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.Status > selected.Status {
			selected = candidate
		}
	}
	return selected
}

func collectDomainErrors(err error, target *[]*core.Error) {
	if err == nil {
		return
	}
	if domainErr, ok := err.(*core.Error); ok {
		*target = append(*target, domainErr)
		return
	}
	switch value := err.(type) {
	case openapi3.MultiError:
		for _, nested := range value {
			collectDomainErrors(nested, target)
		}
		return
	case *openapi3.MultiError:
		for _, nested := range *value {
			collectDomainErrors(nested, target)
		}
		return
	case interface{ Unwrap() []error }:
		for _, nested := range value.Unwrap() {
			collectDomainErrors(nested, target)
		}
	case interface{ Unwrap() error }:
		collectDomainErrors(value.Unwrap(), target)
	}
}
