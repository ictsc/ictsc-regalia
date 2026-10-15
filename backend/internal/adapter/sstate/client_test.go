package sstate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	testBearer = "sstate-shared-secret"
	testCommit = "0123456789abcdef0123456789abcdef01234567"
)

func TestQueueRequiresAcceptedAndSendsCallbackContract(t *testing.T) {
	var got RedeployRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/redeploy" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertFixedBearer(t, r)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client := newTestClient(t, server)
	request := service.DeploymentRequest{
		RequestID: "123e4567-e89b-12d3-a456-426614174000", TeamCode: 2, ProblemCode: "A01", Revision: 3,
		ContentCommit: testCommit, CallbackURL: "https://score.example/api/v1/admin/deployments/2/A01/3/events",
	}
	if err := client.Queue(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	want := RedeployRequest{
		RequestID: request.RequestID, TeamCode: 2, ProblemCode: "A01", Revision: 3,
		ContentCommit: testCommit, CallbackURL: request.CallbackURL,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("request = %#v, want %#v", got, want)
	}
}

func TestQueueRejectsEveryStatusExceptAccepted(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusBadGateway} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "not accepted", status)
			}))
			defer server.Close()
			client := newTestClient(t, server)
			err := client.Queue(context.Background(), service.DeploymentRequest{
				RequestID: "123e4567-e89b-12d3-a456-426614174000", TeamCode: 2, ProblemCode: "A01", Revision: 1,
				ContentCommit: testCommit, CallbackURL: "https://score.example/callback",
			})
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) || httpErr.StatusCode != status {
				t.Fatalf("error = %v, want HTTPError status %d", err, status)
			}
		})
	}
}

func TestStatusIsOneShotFixedBearerRequest(t *testing.T) {
	requests := 0
	occurredAt := time.Date(2026, 8, 30, 3, 4, 5, 0, time.FixedZone("JST", 9*60*60))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet || r.URL.Path != "/status/02/A01" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertFixedBearer(t, r)
		message := "creating instances"
		if err := json.NewEncoder(w).Encode(StatusResponse{
			EventID: "018f47d2-7f3d-7c6d-8a2a-4d96cd884a01", OccurredAt: occurredAt,
			Status: core.DeploymentDeploying, Message: &message,
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server)
	event, err := client.Status(context.Background(), 2, "A01")
	if err != nil {
		t.Fatal(err)
	}
	if event.EventID != "018f47d2-7f3d-7c6d-8a2a-4d96cd884a01" || event.Status != core.DeploymentDeploying ||
		!event.OccurredAt.Equal(occurredAt) || event.OccurredAt.Location() != time.UTC || event.Message == nil || *event.Message != "creating instances" {
		t.Fatalf("unexpected event: %#v", event)
	}
	if requests != 1 {
		t.Fatalf("Status performed %d requests; background polling is forbidden", requests)
	}
}

func TestStatusRejectsUnknownFieldsAndQueuedState(t *testing.T) {
	for name, response := range map[string]string{
		"unknown field": `{"event_id":"123e4567-e89b-12d3-a456-426614174000","occurred_at":"2026-08-30T03:04:05Z","status":"DEPLOYING","message":null,"extra":true}`,
		"queued":        `{"event_id":"123e4567-e89b-12d3-a456-426614174000","occurred_at":"2026-08-30T03:04:05Z","status":"QUEUED","message":null}`,
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(response))
			}))
			defer server.Close()
			client := newTestClient(t, server)
			if _, err := client.Status(context.Background(), 2, "A01"); err == nil {
				t.Fatal("expected invalid status response to fail")
			}
		})
	}
}

func TestEqualBearerToken(t *testing.T) {
	if !EqualBearerToken("shared-secret", "shared-secret") || EqualBearerToken("shared-secret", "other-secret") || EqualBearerToken("short", "longer") {
		t.Fatal("constant-time bearer comparison returned an unexpected result")
	}
}

func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	client, err := New(Config{
		BaseURL: server.URL, BearerToken: testBearer, HTTPClient: server.Client(), AllowInsecureHTTP: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func assertFixedBearer(t *testing.T, request *http.Request) {
	t.Helper()
	if got := request.Header.Get("Authorization"); got != "Bearer "+testBearer {
		t.Fatalf("Authorization = %q", got)
	}
}
