package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func clearExportEnv(t *testing.T) {
	t.Helper()
	t.Setenv("OTEL_SDK_DISABLED", "")
	for _, key := range exportEnvs {
		t.Setenv(key, "")
	}
}

func TestEnabled(t *testing.T) {
	clearExportEnv(t)
	if Enabled() {
		t.Fatal("Enabled() = true without OTEL_* export settings")
	}
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://collector:4318")
	if !Enabled() {
		t.Fatal("Enabled() = false with OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	t.Setenv("OTEL_SDK_DISABLED", "true")
	if Enabled() {
		t.Fatal("Enabled() = true with OTEL_SDK_DISABLED=true")
	}
}

func TestSetupDisabledIsNoop(t *testing.T) {
	clearExportEnv(t)
	shutdown, err := Setup(context.Background())
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown() error = %v", err)
	}
}

func TestHTTPHandlerNamesSpanAfterRoute(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() { otel.SetTracerProvider(previous) })

	router := chi.NewRouter()
	router.Use(ChiRoute)
	router.Get("/api/v1/teams/{teamID}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	router.Get(healthPath, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := HTTPHandler(router)

	for _, path := range []string{"/api/v1/teams/42", healthPath} {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("recorded %d spans, want 1 (health checks are filtered)", len(spans))
	}
	if got, want := spans[0].Name(), "GET /api/v1/teams/{teamID}"; got != want {
		t.Fatalf("span name = %q, want %q", got, want)
	}
	var route string
	for _, kv := range spans[0].Attributes() {
		if kv.Key == attribute.Key("http.route") {
			route = kv.Value.AsString()
		}
	}
	if route != "/api/v1/teams/{teamID}" {
		t.Fatalf("http.route = %q, want route pattern", route)
	}
}
