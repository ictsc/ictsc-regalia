package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// validateOpenAPIResponses buffers finite HTTP responses until kin-openapi has
// checked their status, headers, and body against the matched operation. SSE
// responses are intentionally excluded because buffering an unbounded stream
// would prevent EventSource clients from receiving events; their generated
// response type and dedicated contract tests cover the streaming envelope.
func validateOpenAPIResponses(spec *openapi3.T, next http.Handler) (http.Handler, error) {
	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		return nil, fmt.Errorf("build OpenAPI response router: %w", err)
	}
	validator := openapi3filter.NewValidator(
		router,
		openapi3filter.Strict(true),
		openapi3filter.ValidationOptions(openapi3filter.Options{
			// Request authentication was already performed by the outer request
			// validator. Repeat structural request validation only to resolve the
			// operation and path parameters used for response validation.
			AuthenticationFunc:    func(context.Context, *openapi3filter.AuthenticationInput) error { return nil },
			MultiError:            true,
			IncludeResponseStatus: true,
		}),
		openapi3filter.OnLog(func(_ context.Context, message string, err error) {
			if strings.Contains(message, "failed to write response") && errors.Is(err, http.ErrBodyNotAllowed) {
				return
			}
			log.Printf("OpenAPI response validation: %s: %v", message, err)
		}),
		openapi3filter.OnErr(func(ctx context.Context, w http.ResponseWriter, _ int, code openapi3filter.ErrCode, validationErr error) {
			request := requestFrom(ctx)
			if request == nil {
				request = &http.Request{URL: &url.URL{Path: "/"}}
			}
			switch code {
			case openapi3filter.ErrCodeCannotFindRoute:
				writeProblem(w, request, core.WrapError(http.StatusNotFound, "resource_not_found", "API operation was not found", validationErr))
			case openapi3filter.ErrCodeRequestInvalid:
				writeProblem(w, request, core.WrapError(http.StatusUnprocessableEntity, "validation_error", "Request does not match the OpenAPI contract", validationErr))
			default:
				writeProblem(w, request, core.WrapError(http.StatusInternalServerError, "internal_error", "Response does not match the OpenAPI contract", validationErr))
			}
		}),
	)
	validated := validator.Middleware(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDeploymentSSEPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), requestKey{}, r)
		validated.ServeHTTP(w, r.WithContext(ctx))
	}), nil
}

func isDeploymentSSEPath(path string) bool {
	if path == "/api/v1/admin/deployments/stream" {
		return true
	}
	return strings.HasPrefix(path, "/api/v1/contestant/problems/") && strings.HasSuffix(path, "/deployments/stream")
}
