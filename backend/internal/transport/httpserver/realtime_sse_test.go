package httpserver

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

type sseEvent struct{ name, data string }

// openSSE connects to path and returns a function that reads the next
// non-keepalive event.
func openSSE(t *testing.T, server *httptest.Server, path string, cookie *http.Cookie) func() sseEvent {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.AddCookie(cookie)
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = response.Body.Close() })
	if response.StatusCode != http.StatusOK || response.Header.Get("Content-Type") != "text/event-stream" {
		t.Fatalf("%s status = %d content-type = %q", path, response.StatusCode, response.Header.Get("Content-Type"))
	}
	reader := bufio.NewReader(response.Body)
	return func() sseEvent {
		t.Helper()
		var event sseEvent
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			line = strings.TrimSuffix(line, "\n")
			switch {
			case line == "" && event.name != "":
				return event
			case strings.HasPrefix(line, "event: "):
				event.name = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				event.data = strings.TrimPrefix(line, "data: ")
			}
		}
	}
}

func TestRealtimeStreamsScopeEventsByTeamAndVisibility(t *testing.T) {
	fixture := newContractFixture(t)
	var now atomic.Pointer[time.Time]
	setNow := func(value time.Time) { now.Store(&value) }
	setNow(testNow)
	fixture.service.Now = func() time.Time { return *now.Load() }
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	contestantCookie := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	adminCookie := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "operator"}, service.AdminTTL)
	server := httptest.NewServer(fixture.handler)
	t.Cleanup(server.Close)

	contestantNext := openSSE(t, server, "/api/v1/contestant/events/stream", contestantCookie)
	adminNext := openSSE(t, server, "/api/v1/admin/events/stream", adminCookie)
	expect := func(next func() sseEvent, name, data string) {
		t.Helper()
		got := next()
		if got.name != name || got.data != data {
			t.Fatalf("event = %q %s, want %q %s", got.name, got.data, name, data)
		}
	}
	expect(contestantNext, "ready", `{}`)
	expect(adminNext, "ready", `{}`)

	mutate := func(method, path, body string, cookie *http.Cookie) {
		t.Helper()
		request := httptest.NewRequest(method, path, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		fixture.handler.ServeHTTP(recorder, request)
		if recorder.Code >= 300 {
			t.Fatalf("%s %s status = %d body = %s", method, path, recorder.Code, recorder.Body.String())
		}
	}

	// Another team's submission reaches operators but never this team.
	fixture.service.PublishRealtime(context.Background(), core.RealtimeEvent{Kind: core.RealtimeAnswerSubmitted, TeamCode: 3, ProblemCode: "P1", AnswerNumber: 1})
	expect(adminNext, "answer", `{"reason":"SUBMITTED","answer":{"answer_number":1,"problem_code":"P1","team_code":3}}`)
	mutate(http.MethodPost, "/api/v1/contestant/problems/P1/answers", `{"body":"answer"}`, contestantCookie)
	expect(contestantNext, "answers", `{"problem_code":"P1"}`)
	expect(adminNext, "answer", `{"reason":"SUBMITTED","answer":{"answer_number":1,"problem_code":"P1","team_code":2}}`)

	// A marking during the publication delay must not reveal its timing.
	mark := `{"answer":{"answer_number":1,"problem_code":"P1","team_code":2},"score":80,"rationale":"ok"}`
	mutate(http.MethodPost, "/api/v1/admin/marking-results", mark, adminCookie)
	expect(adminNext, "answer", `{"reason":"MARKING_UPDATED","answer":{"answer_number":1,"problem_code":"P1","team_code":2}}`)
	fixture.service.PublishRealtime(context.Background(), core.RealtimeEvent{Kind: core.RealtimeTeamsUpdated})
	expect(contestantNext, "teams", `{}`)
	expect(adminNext, "teams", `{}`)

	// The delay's end is announced independently of the marking.
	fixture.service.PublishRealtime(context.Background(), core.RealtimeEvent{Kind: core.RealtimeAnswerPublishDelayElapsed, TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1})
	expect(contestantNext, "scores", `{"problem_code":"P1"}`)
	expect(adminNext, "answer", `{"reason":"PUBLICATION_DELAY_ELAPSED","answer":{"answer_number":1,"problem_code":"P1","team_code":2}}`)

	// Once the score is visible, a regrade is pushed immediately.
	setNow(testNow.Add(core.PublishDelay))
	mutate(http.MethodPost, "/api/v1/admin/marking-results", mark, adminCookie)
	expect(contestantNext, "scores", `{"problem_code":"P1"}`)
	expect(adminNext, "answer", `{"reason":"MARKING_UPDATED","answer":{"answer_number":1,"problem_code":"P1","team_code":2}}`)

	mutate(http.MethodPost, "/api/v1/admin/scores/actions/reveal-final", ``, adminCookie)
	expect(contestantNext, "competition", `{}`)
	expect(adminNext, "competition", `{}`)
}
