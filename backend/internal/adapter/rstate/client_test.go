package rstate

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const testPassword = "rstate-admin-password"

func TestQueuePostsJobFormWithBasicAuth(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost || r.URL.Path != "/admin/postJob" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertBasicAuth(t, r)
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("Content-Type = %q", got)
		}
		if r.PostFormValue("team_id") != "team02" || r.PostFormValue("prob_id") != "a01" {
			t.Fatalf("form = %v", r.PostForm)
		}
		_, _ = w.Write([]byte(`{"id":"8ed8f0d9-8eb9-11ed-9bc1-befd30ae56b9","state":0,"team_id":"team02","prob_id":"a01"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server)
	jobID, err := client.Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: "A01", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	if jobID != "8ed8f0d9-8eb9-11ed-9bc1-befd30ae56b9" {
		t.Fatalf("job ID = %q", jobID)
	}
	if requests != 1 {
		t.Fatalf("Queue performed %d requests", requests)
	}
}

func TestQueueClassifiesRejections(t *testing.T) {
	for name, tc := range map[string]struct {
		status     int
		body       string
		definitive bool
		reason     string
	}{
		"bad parameters":   {http.StatusBadRequest, "BadRequest!", true, "rejected the team or problem ID"},
		"preparatory":      {http.StatusBadRequest, "https://rstate.example/status/0b6f-secret-token", true, "preparatory phase"},
		"wrong password":   {http.StatusUnauthorized, "", false, ""},
		"upstream failure": {http.StatusBadGateway, "", false, ""},
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/admin/list" {
					// A finished job does not explain the refusal.
					_, _ = w.Write([]byte(`[{"id":"old","state":"終了","created_time":1,"team_id":"02","prob_id":"a01"}]`))
					return
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			_, err := newTestClient(t, server).Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: "A01", Revision: 1})
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) || httpErr.StatusCode != tc.status {
				t.Fatalf("error = %v, want HTTPError status %d", err, tc.status)
			}
			if httpErr.DefinitiveDeploymentRejection() != tc.definitive {
				t.Fatalf("DefinitiveDeploymentRejection = %v, want %v", !tc.definitive, tc.definitive)
			}
			if !strings.Contains(err.Error(), tc.reason) || strings.Contains(err.Error(), "secret-token") {
				t.Fatalf("error = %q, want reason %q without the status page URL", err, tc.reason)
			}
		})
	}
}

func TestQueueTakesARefusalExplainedByAnActiveJobAsAccepted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r)
		switch r.URL.Path {
		case "/admin/postJob":
			// A retry after the first response was lost.
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("https://rstate.example/status/0b6f-secret-token"))
		case "/admin/list":
			_, _ = w.Write([]byte(`[
				{"id":"other-team","state":"実行中","created_time":2,"team_id":"03","prob_id":"a01"},
				{"id":"taken","state":"開始待ち","created_time":2,"team_id":"02","prob_id":"a01"}
			]`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	jobID, err := newTestClient(t, server).Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: "A01", Revision: 1})
	if err != nil || jobID != "taken" {
		t.Fatalf("job ID = %q, err = %v, want the active job", jobID, err)
	}
}

func TestQueueRejectsProblemCodesRstateRefuses(t *testing.T) {
	client, err := New(Config{BaseURL: "https://rstate.example", AdminPassword: testPassword})
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"A1", "../x", ""} {
		_, err := client.Queue(context.Background(), service.DeploymentRequest{TeamCode: 2, ProblemCode: code})
		var rejection interface{ DefinitiveDeploymentRejection() bool }
		if !errors.As(err, &rejection) || !rejection.DefinitiveDeploymentRejection() {
			t.Fatalf("problem code %q: error = %v, want a definitive rejection", code, err)
		}
	}
}

func TestQueueAcceptsTeamOne(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/postJob" || r.PostFormValue("team_id") != "team01" {
			t.Fatalf("unexpected request: %s %s form=%v", r.Method, r.URL.Path, r.PostForm)
		}
		_, _ = w.Write([]byte(`{"id":"team-one-job","state":0,"team_id":"team01","prob_id":"a01"}`))
	}))
	defer server.Close()
	jobID, err := newTestClient(t, server).Queue(context.Background(), service.DeploymentRequest{TeamCode: 1, ProblemCode: "A01", Revision: 1})
	if err != nil || jobID != "team-one-job" {
		t.Fatalf("job ID = %q, err = %v", jobID, err)
	}
}

func TestQueueRejectsTeamCodesOutsideRstateWorkspaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("rstate was asked for %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()
	client := newTestClient(t, server)
	for _, team := range []int64{0, 100} {
		_, err := client.Queue(context.Background(), service.DeploymentRequest{TeamCode: team, ProblemCode: "A01", Revision: 1})
		var targetErr *TargetError
		if !errors.As(err, &targetErr) {
			t.Fatalf("team %d: error = %v, want a TargetError", team, err)
		}
	}
}

func TestStatusFindsTeamOneJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"job","state":"終了","created_time":1,"team_id":"01","prob_id":"a01"}]`))
	}))
	defer server.Close()
	event, err := newTestClient(t, server).Status(context.Background(), 1, "A01")
	if err != nil || event.Status != core.DeploymentCompleted {
		t.Fatalf("event = %#v, err = %v", event, err)
	}
}

func TestListJobsReadsEveryJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/admin/list" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertBasicAuth(t, r)
		_, _ = w.Write([]byte(`[
			{"id":"running","state":"実行中","created_time":2,"end_time":0,"team_id":"02","prob_id":"a01"},
			{"id":"failed","state":"エラー","created_time":1,"end_time":0,"team_id":"12","prob_id":"b01"},
			{"id":"bad-team","state":"終了","created_time":1,"end_time":1,"team_id":"xx","prob_id":"a01"},
			{"id":"bad-state","state":"Unknown State","created_time":1,"end_time":1,"team_id":"02","prob_id":"a01"}
		]`))
	}))
	defer server.Close()

	jobs, err := newTestClient(t, server).ListJobs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 {
		t.Fatalf("jobs = %#v, want the two readable entries", jobs)
	}
	if got := jobs[0]; got.ID != "running" || got.TeamCode != 2 || got.ProblemCode != "a01" || got.Status != core.DeploymentDeploying || got.Message != nil {
		t.Fatalf("running job = %#v", got)
	}
	if got := jobs[1]; got.ID != "failed" || got.TeamCode != 12 || got.Status != core.DeploymentFailed || got.Message == nil {
		t.Fatalf("failed job = %#v", got)
	}
}

func TestListJobsMapsStates(t *testing.T) {
	for state, want := range map[string]core.DeploymentStatus{
		"開始待ち": core.DeploymentQueued, "実行中": core.DeploymentDeploying, "終了": core.DeploymentCompleted,
		"エラー": core.DeploymentFailed, "Limit": core.DeploymentFailed,
	} {
		t.Run(state, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`[{"id":"job","state":"` + state + `","created_time":1,"end_time":1,"team_id":"02","prob_id":"a01"}]`))
			}))
			defer server.Close()
			jobs, err := newTestClient(t, server).ListJobs(context.Background())
			if err != nil || len(jobs) != 1 || jobs[0].Status != want {
				t.Fatalf("jobs = %#v, err = %v, want status %s", jobs, err, want)
			}
		})
	}
}

func TestStatusPicksNewestJobAndReportsMissingOnes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[
			{"id":"old","state":"エラー","created_time":1,"team_id":"02","prob_id":"a01"},
			{"id":"new","state":"終了","created_time":2,"team_id":"02","prob_id":"a01"}
		]`))
	}))
	defer server.Close()
	client := newTestClient(t, server)
	event, err := client.Status(context.Background(), 2, "A01")
	if err != nil || event.Status != core.DeploymentCompleted {
		t.Fatalf("event = %#v, err = %v", event, err)
	}
	_, err = client.Status(context.Background(), 3, "A01")
	var notFound interface{ DeploymentStatusNotFound() bool }
	if !errors.As(err, &notFound) || !notFound.DeploymentStatusNotFound() {
		t.Fatalf("error = %v, want not found", err)
	}
}

func TestStatusReportsWaitingJobAsDeploying(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"job","state":"開始待ち","created_time":1,"team_id":"02","prob_id":"a01"}]`))
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

func TestAdminUserCanBeOverridden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, password, ok := r.BasicAuth(); !ok || user != "proxy-user" || password != testPassword {
			t.Fatalf("Basic auth = %q/%q (%v)", user, password, ok)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, AdminUser: "proxy-user", AdminPassword: testPassword, HTTPClient: server.Client(), AllowInsecureHTTP: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.ListJobs(context.Background()); err != nil {
		t.Fatal(err)
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
	if !ok || user != defaultAdminUser || password != testPassword {
		t.Fatalf("Basic auth = %q/%q (%v)", user, password, ok)
	}
}
