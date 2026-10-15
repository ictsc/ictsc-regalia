package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func TestOpenAPIResponseValidation(t *testing.T) {
	spec, err := api.GetSwagger()
	if err != nil {
		t.Fatal(err)
	}
	spec.Servers = nil
	tests := []struct {
		name       string
		handler    http.Handler
		wantStatus int
		wantCode   string
	}{
		{
			name: "valid health response",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			}),
			wantStatus: http.StatusOK,
		},
		{
			name: "undeclared response status",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			}),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
		},
		{
			name: "schema-invalid response body",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"wrong"}`))
			}),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validated, err := validateOpenAPIResponses(spec, test.handler)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			validated.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if test.wantCode != "" {
				assertProblem(t, recorder, test.wantStatus, test.wantCode)
			}
		})
	}
}
