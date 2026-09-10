package main

import (
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	"github.com/ictsc/ictsc-regalia/backend/internal/transport/httpserver"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"testing"
)

func TestDevContentSeedsOnlyEmptyStore(t *testing.T) {
	store := memory.NewCompetitionStore()
	ctx := t.Context()
	if err := seedDevContent(ctx, store); err != nil {
		t.Fatal(err)
	}
	snapshot, err := store.ActiveContent(ctx)
	if err != nil || len(snapshot.Manifest.Problems) != 3 {
		t.Fatalf("snapshot: %v", err)
	}
	replacement := snapshot
	replacement.CommitSHA = "1111111111111111111111111111111111111111"
	if _, err := store.ActivateContent(ctx, replacement, snapshot.CommitSHA); err != nil {
		t.Fatal(err)
	}
	if err := seedDevContent(ctx, store); err != nil {
		t.Fatal(err)
	}
	got, _ := store.ActiveContent(ctx)
	if got.CommitSHA != replacement.CommitSHA {
		t.Fatal("existing content was overwritten")
	}
}

func TestDevContentParticipantResponses(t *testing.T) {
	ctx := t.Context()
	store := memory.NewCompetitionStore()
	sessions := memory.NewSessionStore()
	if err := seedDevContent(ctx, store); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateTeam(ctx, core.Team{Code: 2, Name: "Test", Organization: "Test", MemberLimit: 5, Color: core.DefaultTeamColor}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RegisterContestant(ctx, core.Contestant{Name: "alice", DisplayName: "Alice", DiscordID: "123", TeamCode: 2}); err != nil {
		t.Fatal(err)
	}
	token, err := sessions.Create(ctx, session.Data{Kind: session.KindContestant, ContestantName: "alice"}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	svc := service.New(store, sessions, service.Config{})
	svc.Now = func() time.Time { return time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC) }
	handler, err := httpserver.New(svc, httpserver.Options{AllowedOrigins: []string{"https://example.test"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/v1/contestant/sections", "/api/v1/contestant/problems", "/api/v1/contestant/problems/M01", "/api/v1/contestant/problems/M02", "/api/v1/contestant/problems/M03"} {
		r := httptest.NewRequest("GET", path, nil)
		r.AddCookie(&http.Cookie{Name: "regalia-user-session", Value: token})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
	}
}
