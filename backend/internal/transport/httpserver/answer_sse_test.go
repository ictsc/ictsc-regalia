package httpserver

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

type sseEvent struct {
	name string
	data string
}

// openSSE connects to a stream and returns a function reading one event,
// skipping keepalive comments.
func openSSE(t *testing.T, ctx context.Context, server *httptest.Server, path string, cookie *http.Cookie) func() sseEvent {
	t.Helper()
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
		t.Fatalf("%s: status %d content type %q", path, response.StatusCode, response.Header.Get("Content-Type"))
	}
	reader := bufio.NewReader(response.Body)
	return func() sseEvent {
		t.Helper()
		var event sseEvent
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("%s: read event: %v", path, err)
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

func TestAnswerSSEDeliversSubmissionWorkflowAndMarking(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	team := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	server := httptest.NewServer(fixture.handler)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	contestantEvents := openSSE(t, ctx, server, "/api/v1/contestant/answers/stream", team)
	adminEvents := openSSE(t, ctx, server, "/api/v1/admin/answers/stream", admin)
	for name, next := range map[string]func() sseEvent{"contestant": contestantEvents, "admin": adminEvents} {
		if event := next(); event.name != "ready" || event.data != "{}" {
			t.Fatalf("%s first event = %#v, want ready", name, event)
		}
	}

	// Another team's change reaches admins only.
	other, _ := json.Marshal(core.AnswerChange{TeamCode: 3, ProblemCode: "P1", AnswerNumber: 1, Kind: core.AnswerSubmitted})
	if err := fixture.service.AnswerEvents.Publish(ctx, other); err != nil {
		t.Fatal(err)
	}
	if event := adminEvents(); event.data != `{"team_code":3,"problem_code":"P1","answer_number":1,"kind":"SUBMITTED"}` {
		t.Fatalf("admin other-team event = %#v", event)
	}

	send := func(method, path, body string, cookie *http.Cookie, want int) {
		t.Helper()
		request := httptest.NewRequest(method, path, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(cookie)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		if response.Code != want {
			t.Fatalf("%s %s: %d %s", method, path, response.Code, response.Body.String())
		}
	}
	send(http.MethodPost, "/api/v1/contestant/problems/P1/answers", `{"body":"answer"}`, team, http.StatusCreated)
	send(http.MethodPatch, "/api/v1/admin/answers/2/P1/1/workflow", `{"expected_revision":0,"status":"IN_PROGRESS"}`, admin, http.StatusOK)
	send(http.MethodPost, "/api/v1/admin/marking-results", `{"answer":{"team_code":2,"problem_code":"P1","answer_number":1},"score":80,"rationale":"ok"}`, admin, http.StatusCreated)

	for _, kind := range []string{"SUBMITTED", "WORKFLOW", "MARKED"} {
		if event := adminEvents(); event.name != "answer" || event.data != `{"team_code":2,"problem_code":"P1","answer_number":1,"kind":"`+kind+`"}` {
			t.Fatalf("admin %s event = %#v", kind, event)
		}
		// Contestants get only the identity: no team, kind or score.
		if event := contestantEvents(); event.name != "answer" || event.data != `{"problem_code":"P1","answer_number":1}` {
			t.Fatalf("contestant %s event = %#v", kind, event)
		}
	}
}

func TestAnswerSSERequiresMatchingSession(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	team := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	for _, tc := range []struct {
		name    string
		path    string
		cookies []*http.Cookie
	}{
		{"anonymous contestant stream", "/api/v1/contestant/answers/stream", nil},
		{"anonymous admin stream", "/api/v1/admin/answers/stream", nil},
		{"contestant on admin stream", "/api/v1/admin/answers/stream", []*http.Cookie{team}},
	} {
		response := fixture.request(t, http.MethodGet, tc.path, "", tc.cookies...)
		if response.Code != http.StatusUnauthorized && response.Code != http.StatusForbidden {
			t.Errorf("%s: status %d, want 401 or 403; body %s", tc.name, response.Code, response.Body.String())
		}
		if strings.Contains(response.Header().Get("Content-Type"), "text/event-stream") {
			t.Errorf("%s: opened an event stream", tc.name)
		}
	}
}

type closeTrackingBus struct{ closed chan struct{} }

func (b *closeTrackingBus) Publish(context.Context, []byte) error { return nil }

func (b *closeTrackingBus) Subscribe(context.Context) (service.Subscription, error) {
	return &closeTrackingSubscription{bus: b, messages: make(chan []byte)}, nil
}

type closeTrackingSubscription struct {
	bus      *closeTrackingBus
	messages chan []byte
}

func (s *closeTrackingSubscription) Messages() <-chan []byte { return s.messages }
func (s *closeTrackingSubscription) Close() error {
	close(s.bus.closed)
	return nil
}

func TestAnswerSSEReleasesSubscriptionWhenClientLeaves(t *testing.T) {
	fixture := newContractFixture(t)
	bus := &closeTrackingBus{closed: make(chan struct{})}
	fixture.service.AnswerEvents = bus
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	server := httptest.NewServer(fixture.handler)
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	next := openSSE(t, ctx, server, "/api/v1/admin/answers/stream", admin)
	if event := next(); event.name != "ready" {
		t.Fatalf("first event = %#v, want ready", event)
	}
	cancel()
	select {
	case <-bus.closed:
	case <-time.After(3 * time.Second):
		t.Fatal("subscription was not closed after the client disconnected")
	}
}
