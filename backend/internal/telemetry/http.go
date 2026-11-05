package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// healthPath is polled by probes and would otherwise dominate traces.
const healthPath = "/api/v1/health"

// HTTPHandler records a server span and request metrics for every request
// except health checks.
func HTTPHandler(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool { return r.URL.Path != healthPath }),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string { return r.Method }),
	)
}

// ChiRoute names the server span after the matched chi route pattern and adds
// http.route to the span and request metrics. It must run inside a chi router
// that is itself wrapped by HTTPHandler.
func ChiRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil {
			SetRoute(r, routeCtx.RoutePattern())
		}
	})
}

// SetRoute records pattern as the request's http.route. Use it where a request
// is answered before reaching the chi router, e.g. OpenAPI validation errors.
func SetRoute(r *http.Request, pattern string) {
	if pattern == "" {
		return
	}
	route := semconv.HTTPRoute(pattern)
	span := trace.SpanFromContext(r.Context())
	span.SetName(r.Method + " " + pattern)
	span.SetAttributes(route)
	if labeler, ok := otelhttp.LabelerFromContext(r.Context()); ok {
		labeler.Add(route)
	}
}

// HTTPTransport wraps base (http.DefaultTransport when nil) so outgoing
// requests produce client spans and carry trace context.
func HTTPTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return otelhttp.NewTransport(base)
}
