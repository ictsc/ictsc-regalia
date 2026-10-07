package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func TestAdminAnswerWorkflowAndSubmissionTimeGrader(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	ctx := context.Background()
	initial, err := fixture.store.ActiveContent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assigned := initial
	assigned.CommitSHA = strings.Repeat("b", 40)
	assigned.Manifest.Problems = append([]core.Problem(nil), initial.Manifest.Problems...)
	assigned.Manifest.Problems[0].DefaultGraderDiscordID = "123456789012345678"
	if _, err := fixture.store.ActivateContent(ctx, assigned, initial.CommitSHA); err != nil {
		t.Fatal(err)
	}
	answer, _, err := fixture.store.SubmitAnswer(ctx, core.Answer{
		TeamCode: contestant.TeamCode, ProblemCode: "P1", AuthorName: contestant.Name,
		Body: "answer\nsecond line\nthird line", SubmittedAt: testNow, ContentCommit: assigned.CommitSHA,
		MaxScore: 100, RedeployRule: assigned.Manifest.Problems[0].Redeploy,
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	changed := assigned
	changed.CommitSHA = strings.Repeat("c", 40)
	changed.Manifest.Problems = append([]core.Problem(nil), assigned.Manifest.Problems...)
	changed.Manifest.Problems[0].DefaultGraderDiscordID = "999999999999999999"
	if _, err := fixture.store.ActivateContent(ctx, changed, assigned.CommitSHA); err != nil {
		t.Fatal(err)
	}
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge", Discord: core.DiscordIdentity{ID: "222222222222222222"}}, service.AdminTTL)
	path := "/api/v1/admin/answers/2/P1/1"
	read := func() api.AdminAnswer {
		t.Helper()
		response := fixture.request(t, http.MethodGet, path, "", admin)
		if response.Code != http.StatusOK {
			t.Fatalf("GET status %d: %s", response.Code, response.Body.String())
		}
		var body api.AdminAnswerResponse
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		return body.Answer
	}
	current := read()
	if current.Workflow.AssigneeDiscordId == nil || *current.Workflow.AssigneeDiscordId != "123456789012345678" || current.Workflow.Status != "WAITING" {
		t.Fatalf("initial workflow: %#v", current.Workflow)
	}
	mutate := func(method, target, body string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(method, target, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(admin)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		return response
	}
	workflowPath := path + "/workflow"
	if response := mutate(http.MethodPatch, workflowPath, `{"expected_revision":0,"status":"COMPLETED"}`); response.Code != http.StatusConflict {
		t.Fatalf("unmarked completion: %d %s", response.Code, response.Body.String())
	}
	if response := mutate(http.MethodPatch, workflowPath, `{"expected_revision":0,"status":"IN_PROGRESS","assignment":"CLAIM_SELF"}`); response.Code != http.StatusOK {
		t.Fatalf("claim: %d %s", response.Code, response.Body.String())
	}
	current = read()
	if current.Workflow.AssigneeDiscordId == nil || *current.Workflow.AssigneeDiscordId != "222222222222222222" || current.Workflow.Revision != 1 {
		t.Fatalf("claim workflow: %#v", current.Workflow)
	}
	if response := mutate(http.MethodPatch, workflowPath, `{"expected_revision":0,"status":"WAITING"}`); response.Code != http.StatusConflict {
		t.Fatalf("stale update: %d %s", response.Code, response.Body.String())
	}
	if response := mutate(http.MethodPost, "/api/v1/admin/marking-results", `{"answer":{"team_code":2,"problem_code":"P1","answer_number":1},"score":80,"rationale":"ok","line_comments":[{"line_number":2,"end_line_number":4,"body":"invalid"}]}`); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("out-of-range line: %d %s", response.Code, response.Body.String())
	}
	if response := mutate(http.MethodPost, "/api/v1/admin/marking-results", `{"answer":{"team_code":2,"problem_code":"P1","answer_number":1},"score":80,"rationale":"ok","line_comments":[{"line_number":3,"end_line_number":2,"body":"invalid"}]}`); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("reversed line range: %d %s", response.Code, response.Body.String())
	}
	markBody := `{"answer":{"team_code":2,"problem_code":"P1","answer_number":1},"score":80,"rationale":"ok","line_comments":[{"line_number":1,"end_line_number":3,"body":"check this range"}]}`
	if response := mutate(http.MethodPost, "/api/v1/admin/marking-results", markBody); response.Code != http.StatusCreated {
		t.Fatalf("mark: %d %s", response.Code, response.Body.String())
	}
	current = read()
	if current.Workflow.Status != "COMPLETED" || current.Workflow.Revision != 2 || current.Score == nil {
		t.Fatalf("marked workflow: %#v", current)
	}
	markings := fixture.request(t, http.MethodGet, "/api/v1/admin/marking-results", "", admin)
	if markings.Code != http.StatusOK {
		t.Fatalf("marking history: %d %s", markings.Code, markings.Body.String())
	}
	var history api.MarkingResultsResponse
	if err := json.Unmarshal(markings.Body.Bytes(), &history); err != nil {
		t.Fatal(err)
	}
	if len(history.MarkingResults) != 1 || len(history.MarkingResults[0].LineComments) != 1 || history.MarkingResults[0].LineComments[0].Body != "check this range" || history.MarkingResults[0].LineComments[0].EndLineNumber == nil || *history.MarkingResults[0].LineComments[0].EndLineNumber != 3 {
		t.Fatalf("line comment history: %#v", history.MarkingResults)
	}
	deletePath := "/api/v1/admin/marking-results/" + history.MarkingResults[0].Id.String() + "/line-comments/0"
	if response := mutate(http.MethodDelete, deletePath, ""); response.Code != http.StatusOK {
		t.Fatalf("delete line comment: %d %s", response.Code, response.Body.String())
	}
	if response := mutate(http.MethodDelete, deletePath, ""); response.Code != http.StatusConflict {
		t.Fatalf("repeat deletion: %d %s", response.Code, response.Body.String())
	}
	if response := mutate(http.MethodDelete, "/api/v1/admin/marking-results/"+history.MarkingResults[0].Id.String()+"/line-comments/1", ""); response.Code != http.StatusNotFound {
		t.Fatalf("missing comment: %d %s", response.Code, response.Body.String())
	}
	markings = fixture.request(t, http.MethodGet, "/api/v1/admin/marking-results", "", admin)
	if err := json.Unmarshal(markings.Body.Bytes(), &history); err != nil {
		t.Fatal(err)
	}
	deleted := history.MarkingResults[0].LineComments[0]
	if deleted.DeletedAt == nil || deleted.DeletedBy == nil || *deleted.DeletedBy != "Judge" || deleted.Body != "check this range" || history.MarkingResults[0].Score != 80 {
		t.Fatalf("soft deleted line comment: %#v", history.MarkingResults[0])
	}
	if response := mutate(http.MethodPatch, workflowPath, `{"expected_revision":2,"status":"WAITING","assignment":"RESET_TO_DEFAULT"}`); response.Code != http.StatusOK {
		t.Fatalf("reopen: %d %s", response.Code, response.Body.String())
	}
	current = read()
	if current.Workflow.Status != "WAITING" || current.Score == nil || current.Workflow.AssigneeDiscordId == nil || *current.Workflow.AssigneeDiscordId != "123456789012345678" {
		t.Fatalf("reopened workflow: %#v", current)
	}
	_ = answer
}
