package httpserver

import (
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
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func TestSecondaryQualifierDefersGradingAndHidesResults(t *testing.T) {
	fixture := newContractFixture(t)
	fixture.service.Config.SecondaryQualifierMode = true
	now := testNow
	fixture.service.Now = func() time.Time { return now }
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	ctx := context.Background()
	for number := 1; number <= 2; number++ {
		if _, _, err := fixture.store.SubmitAnswer(ctx, core.Answer{
			TeamCode: contestant.TeamCode, ProblemCode: "P1", AuthorName: contestant.Name, Body: "answer",
			SubmittedAt: testNow.Add(-time.Duration(3-number) * 30 * time.Minute), ContentCommit: testContentCommit, MaxScore: 100,
		}, 0); err != nil {
			t.Fatal(err)
		}
	}
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	team := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	mark := func(number string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/marking-results", strings.NewReader(`{"answer":{"team_code":2,"problem_code":"P1","answer_number":`+number+`},"score":80,"rationale":"ok"}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(admin)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		return response
	}
	readAnswer := func(number string) api.AdminAnswer {
		t.Helper()
		response := fixture.request(t, http.MethodGet, "/api/v1/admin/answers/2/P1/"+number, "", admin)
		var body api.AdminAnswerResponse
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		return body.Answer
	}

	if answer := readAnswer("2"); answer.CompetitionMode != api.SECONDARYQUALIFIER || !answer.Latest || answer.MarkingBlockedReason == nil || *answer.MarkingBlockedReason != api.MarkingBlockedReasonGradingNotStarted {
		t.Fatalf("running contest answer: %#v", answer)
	}
	if response := mark("2"); response.Code != http.StatusConflict || !strings.Contains(response.Body.String(), "grading_not_started") {
		t.Fatalf("grading before end: %d %s", response.Code, response.Body.String())
	}

	now = testNow.Add(2 * time.Hour)
	if answer := readAnswer("1"); answer.Latest || answer.MarkingBlockedReason == nil || *answer.MarkingBlockedReason != api.MarkingBlockedReasonAnswerNotLatest {
		t.Fatalf("older answer: %#v", answer)
	}
	if response := mark("1"); response.Code != http.StatusConflict || !strings.Contains(response.Body.String(), "answer_not_latest") {
		t.Fatalf("grading older answer: %d %s", response.Code, response.Body.String())
	}
	if answer := readAnswer("2"); answer.MarkingBlockedReason != nil {
		t.Fatalf("latest answer after end: %#v", answer)
	}
	if response := mark("2"); response.Code != http.StatusCreated {
		t.Fatalf("grading latest answer: %d %s", response.Code, response.Body.String())
	}

	var ranking api.RankingResponse
	if err := json.Unmarshal(fixture.request(t, http.MethodGet, "/api/v1/contestant/ranking", "", team).Body.Bytes(), &ranking); err != nil {
		t.Fatal(err)
	}
	if !ranking.Hidden || len(ranking.Ranking) != 0 {
		t.Fatalf("ranking before reveal: %#v", ranking)
	}
	var answers api.AnswersResponse
	if err := json.Unmarshal(fixture.request(t, http.MethodGet, "/api/v1/contestant/problems/P1/answers", "", team).Body.Bytes(), &answers); err != nil {
		t.Fatal(err)
	}
	if len(answers.Answers) != 2 || answers.Answers[0].Score != nil {
		t.Fatalf("answers before reveal: %#v", answers)
	}

	if _, err := fixture.store.RevealFinal(ctx, now, "Judge"); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(fixture.request(t, http.MethodGet, "/api/v1/contestant/ranking", "", team).Body.Bytes(), &ranking); err != nil {
		t.Fatal(err)
	}
	if ranking.Hidden || len(ranking.Ranking) != 1 || ranking.Ranking[0].Score != 80 {
		t.Fatalf("ranking after reveal: %#v", ranking)
	}
}

func TestDevContestPhaseShiftsScheduleOnlyInDevSecondaryQualifier(t *testing.T) {
	fixture := newContractFixture(t)
	fixture.seedActiveContent(t)
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	put := func(phase string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(http.MethodPut, "/api/v1/admin/dev/contest-phase", strings.NewReader(`{"phase":"`+phase+`"}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(admin)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		return response
	}
	if response := put("AFTER"); response.Code != http.StatusNotFound {
		t.Fatalf("switch outside dev secondary qualifier: %d %s", response.Code, response.Body.String())
	}
	fixture.service.Config.DevTools = true
	fixture.service.Config.SecondaryQualifierMode = true
	for _, phase := range []api.ContestPhase{api.BEFORE, api.DURING, api.AFTER} {
		response := put(string(phase))
		var body api.DevContestPhaseResponse
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: %d %s", phase, response.Code, response.Body.String())
		}
		if response.Code != http.StatusOK || !body.ContestPhase.Enabled || body.ContestPhase.Phase != phase {
			t.Fatalf("%s: %d %#v", phase, response.Code, body)
		}
		section := body.ContestPhase.Sections[0]
		if section.Ending.Sub(section.Beginning) != 2*time.Hour {
			t.Fatalf("%s: section length changed: %#v", phase, section)
		}
	}
}

func TestSecondaryQualifierRejectsRankingFreeze(t *testing.T) {
	fixture := newContractFixture(t)
	fixture.service.Config.SecondaryQualifierMode = true
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	put := func(body string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(http.MethodPut, "/api/v1/admin/dashboard-schedule", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(admin)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		return response
	}
	if response := put(`{"ranking_freeze_at":"2026-08-30T06:00:00Z"}`); response.Code != http.StatusConflict {
		t.Fatalf("set freeze: %d %s", response.Code, response.Body.String())
	}
	response := put(`{"ranking_freeze_at":null}`)
	var body api.DashboardScheduleResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || response.Code != http.StatusOK || body.CompetitionMode != api.SECONDARYQUALIFIER {
		t.Fatalf("clear freeze: %d %s", response.Code, response.Body.String())
	}
}

func TestContentImagesAreServedFromTheActiveSnapshot(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	ctx := context.Background()
	fixture.seedActiveContent(t)
	active, err := fixture.store.ActiveContent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	next := active
	next.CommitSHA = strings.Repeat("e", 40)
	next.Manifest.Problems = append([]core.Problem(nil), active.Manifest.Problems...)
	next.Manifest.Problems[0].Body = "![diagram](/images/diagram.png)"
	next.Assets = []core.ContentAsset{
		{Name: "diagram.png", ContentType: "image/png", SHA: strings.Repeat("a", 40), Data: []byte("PNGDATA")},
		{Name: "unreferenced.png", ContentType: "image/png", SHA: strings.Repeat("b", 40), Data: []byte("HIDDEN")},
	}
	if _, err := fixture.store.ActivateContent(ctx, next, active.CommitSHA); err != nil {
		t.Fatal(err)
	}
	team := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	for _, request := range []struct {
		path   string
		cookie *http.Cookie
	}{{"/api/v1/contestant/images/diagram.png", team}, {"/api/v1/admin/images/diagram.png", admin}} {
		response := fixture.request(t, http.MethodGet, request.path, "", request.cookie)
		if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "image/png" || response.Body.String() != "PNGDATA" ||
			response.Header().Get("X-Content-Type-Options") != "nosniff" || response.Header().Get("ETag") != `"`+strings.Repeat("a", 40)+`"` {
			t.Fatalf("%s: %d %#v %q", request.path, response.Code, response.Header(), response.Body.String())
		}
		cached := httptest.NewRequest(http.MethodGet, request.path, nil)
		cached.Header.Set("If-None-Match", response.Header().Get("ETag"))
		cached.AddCookie(request.cookie)
		revalidated := httptest.NewRecorder()
		fixture.handler.ServeHTTP(revalidated, cached)
		if revalidated.Code != http.StatusNotModified || revalidated.Body.Len() != 0 {
			t.Fatalf("%s revalidation: %d %q", request.path, revalidated.Code, revalidated.Body.String())
		}
	}
	if response := fixture.request(t, http.MethodGet, "/api/v1/contestant/images/unreferenced.png", "", team); response.Code != http.StatusNotFound {
		t.Fatalf("unreferenced image for contestant: %d", response.Code)
	}
	if response := fixture.request(t, http.MethodGet, "/api/v1/admin/images/unreferenced.png", "", admin); response.Code != http.StatusOK {
		t.Fatalf("unreferenced image for admin: %d", response.Code)
	}
	if response := fixture.request(t, http.MethodGet, "/api/v1/contestant/images/missing.png", "", team); response.Code != http.StatusNotFound {
		t.Fatalf("missing image: %d", response.Code)
	}
	if response := fixture.request(t, http.MethodGet, "/api/v1/contestant/images/diagram.png", ""); response.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous image: %d", response.Code)
	}
}

func TestAnswerLengthLimitAndLogAttachments(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	ctx := context.Background()
	active, err := fixture.store.ActiveContent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	next := active
	next.CommitSHA = strings.Repeat("f", 40)
	next.Manifest.Problems = append([]core.Problem(nil), active.Manifest.Problems...)
	next.Manifest.Problems[0].AnswerMaxLength = 10
	next.Manifest.Problems[0].AllowAttachments = true
	if _, err := fixture.store.ActivateContent(ctx, next, active.CommitSHA); err != nil {
		t.Fatal(err)
	}
	team := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	submit := func(body string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/contestant/problems/P1/answers", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(team)
		response := httptest.NewRecorder()
		fixture.handler.ServeHTTP(response, request)
		return response
	}
	var problem api.ProblemResponse
	if err := json.Unmarshal(fixture.request(t, http.MethodGet, "/api/v1/contestant/problems/P1", "", team).Body.Bytes(), &problem); err != nil {
		t.Fatal(err)
	}
	if problem.Problem.AnswerMaxLength != 10 || !problem.Problem.AllowAttachments {
		t.Fatalf("problem limits: %#v", problem.Problem)
	}
	if response := submit(`{"body":"12345678901"}`); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("over-long answer: %d %s", response.Code, response.Body.String())
	}
	if response := submit(`{"body":"ok","attachments":[{"name":"../x.txt","content":"x"}]}`); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unsafe attachment name: %d %s", response.Code, response.Body.String())
	}
	if response := submit(`{"body":"ok","attachments":[{"name":"x.log","content":"a\u0000b"}]}`); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("binary attachment: %d %s", response.Code, response.Body.String())
	}
	tooLarge := httptest.NewRequest(http.MethodPost, "/api/v1/contestant/problems/P1/answers", strings.NewReader(strings.Repeat("a", 9<<20)))
	tooLarge.Header.Set("Content-Type", "application/json")
	tooLarge.Header.Set("Origin", testOrigin)
	tooLarge.AddCookie(team)
	rejected := httptest.NewRecorder()
	fixture.handler.ServeHTTP(rejected, tooLarge)
	if rejected.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized request: %d", rejected.Code)
	}
	response := submit(`{"body":"1234567890\n\n","attachments":[{"name":"syslog.txt","content":"<script>alert(1)</script>\n"}]}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("submit: %d %s", response.Code, response.Body.String())
	}
	var created api.AnswerResponse
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil || len(created.Answer.Attachments) != 1 || created.Answer.Attachments[0].Size != 26 {
		t.Fatalf("created answer: %s", response.Body.String())
	}
	var attachments api.AnswerAttachmentsResponse
	if err := json.Unmarshal(fixture.request(t, http.MethodGet, "/api/v1/contestant/problems/P1/answers/1/attachments", "", team).Body.Bytes(), &attachments); err != nil ||
		len(attachments.Attachments) != 1 || attachments.Attachments[0].Content != "<script>alert(1)</script>\n" {
		t.Fatalf("contestant attachments: %#v %v", attachments, err)
	}
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "Judge"}, service.AdminTTL)
	if response := fixture.request(t, http.MethodGet, "/api/v1/admin/answers/2/P1/1/attachments", "", admin); response.Code != http.StatusOK {
		t.Fatalf("admin attachments: %d", response.Code)
	}
}
