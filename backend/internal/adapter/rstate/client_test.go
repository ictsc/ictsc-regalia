package rstate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const testPassword = "rstate-admin-password"

func TestQueuePostsToBackendRouteWithTeamID(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost || r.URL.Path != "/backend/team02/A01" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertBasicAuth(t, r)
		_, _ = w.Write([]byte(`{"id":"8ed8f0d9-8eb9-11ed-9bc1-befd30ae56b9","state":0,"team_id":"team02","prob_id":"A01"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server)
	if err := client.Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: "A01", Revision: 1}); err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("Queue performed %d requests", requests)
	}
}

func TestQueueClassifiesRejections(t *testing.T) {
	for status, definitive := range map[int]bool{
		http.StatusBadRequest: true, http.StatusConflict: true,
		http.StatusUnauthorized: false, http.StatusBadGateway: false,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()
			err := newTestClient(t, server).Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: "A01", Revision: 1})
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) || httpErr.StatusCode != status {
				t.Fatalf("error = %v, want HTTPError status %d", err, status)
			}
			if httpErr.DefinitiveDeploymentRejection() != definitive {
				t.Fatalf("DefinitiveDeploymentRejection = %v, want %v", !definitive, definitive)
			}
		})
	}
}

func TestQueueRejectsProblemCodesRstateRefuses(t *testing.T) {
	client, err := New(Config{BaseURL: "https://rstate.example", AdminPassword: testPassword})
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"A1", "../x", ""} {
		if err := client.Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: code}); err == nil {
			t.Fatalf("problem code %q was accepted", code)
		}
	}
}

func TestObservePicksNewestJobForTeamAndProblem(t *testing.T) {
	created := time.Date(2026, 11, 4, 1, 2, 3, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/admin/list" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertBasicAuth(t, r)
		_, _ = w.Write([]byte(`[
			{"id":"old","state":"終了","created_time":1,"end_time":2,"team_id":"02","prob_id":"A01"},
			{"id":"new","state":"実行中","created_time":` + itoa(created.UnixNano()) + `,"end_time":0,"team_id":"02","prob_id":"A01"},
			{"id":"other-team","state":"エラー","created_time":` + itoa(created.UnixNano()+1) + `,"end_time":0,"team_id":"12","prob_id":"A01"},
			{"id":"other-problem","state":"エラー","created_time":` + itoa(created.UnixNano()+1) + `,"end_time":0,"team_id":"02","prob_id":"B01"}
		]`))
	}))
	defer server.Close()

	observation, err := newTestClient(t, server).Observe(context.Background(), 2, "A01")
	if err != nil {
		t.Fatal(err)
	}
	if observation.JobID != "new" || observation.Status != core.DeploymentDeploying || !observation.CreatedAt.Equal(created) {
		t.Fatalf("unexpected observation: %#v", observation)
	}
}

func TestObserveMapsStates(t *testing.T) {
	for state, want := range map[string]core.DeploymentStatus{
		"開始待ち": core.DeploymentQueued, "実行中": core.DeploymentDeploying, "終了": core.DeploymentCompleted,
		"エラー": core.DeploymentFailed, "Limit": core.DeploymentFailed,
	} {
		t.Run(state, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`[{"id":"job","state":"` + state + `","created_time":1,"end_time":1,"team_id":"02","prob_id":"A01"}]`))
			}))
			defer server.Close()
			observation, err := newTestClient(t, server).Observe(context.Background(), 2, "A01")
			if err != nil {
				t.Fatal(err)
			}
			if observation.Status != want || (want == core.DeploymentFailed) != (observation.Message != nil) {
				t.Fatalf("observation = %#v, want status %s", observation, want)
			}
		})
	}
}

func TestObserveReportsMissingJobAndUnknownState(t *testing.T) {
	for name, body := range map[string]string{"null list": `null`, "other team": `[{"id":"job","state":"終了","created_time":1,"team_id":"03","prob_id":"A01"}]`} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(body)) }))
			defer server.Close()
			_, err := newTestClient(t, server).Status(context.Background(), 2, "A01")
			var notFound interface{ DeploymentStatusNotFound() bool }
			if !errors.As(err, &notFound) || !notFound.DeploymentStatusNotFound() {
				t.Fatalf("error = %v, want not found", err)
			}
		})
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"job","state":"Unknown State","created_time":1,"team_id":"02","prob_id":"A01"}]`))
	}))
	defer server.Close()
	if _, err := newTestClient(t, server).Observe(context.Background(), 2, "A01"); err == nil {
		t.Fatal("unknown rstate state was accepted")
	}
}

func TestStatusReportsWaitingJobAsDeploying(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"job","state":"開始待ち","created_time":1,"team_id":"02","prob_id":"A01"}]`))
	}))
	defer server.Close()
	event, err := newTestClient(t, server).Status(context.Background(), 2, "A01")
	if err != nil {
		t.Fatal(err)
	}
	if event.Status != core.DeploymentDeploying || event.EventID != "" || !event.OccurredAt.IsZero() {
		t.Fatalf("unexpected event: %#v", event)
	}
}

func TestNewRequiresHTTPSUnlessAllowed(t *testing.T) {
	if _, err := New(Config{BaseURL: "http://rstate:8089", AdminPassword: testPassword}); err == nil {
		t.Fatal("plain HTTP was accepted without AllowInsecureHTTP")
	}
	if _, err := New(Config{BaseURL: "https://rstate.example"}); err == nil {
		t.Fatal("empty admin password was accepted")
	}
	if _, err := New(Config{BaseURL: "http://rstate:8089", AdminPassword: testPassword, AllowInsecureHTTP: true}); err != nil {
		t.Fatal(err)
	}
}

func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	client, err := New(Config{BaseURL: server.URL, AdminPassword: testPassword, HTTPClient: server.Client(), AllowInsecureHTTP: true})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func assertBasicAuth(t *testing.T, request *http.Request) {
	t.Helper()
	user, password, ok := request.BasicAuth()
	if !ok || user != adminUser || password != testPassword {
		t.Fatalf("Basic auth = %q/%q (%v)", user, password, ok)
	}
}

func itoa(value int64) string {
	return fmt.Sprint(value)
}
