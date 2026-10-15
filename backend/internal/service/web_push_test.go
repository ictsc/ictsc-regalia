package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

type recordingPushSender struct {
	payloads [][]byte
}

func (s *recordingPushSender) Send(_ context.Context, _ core.WebPushSubscription, payload []byte) (bool, error) {
	s.payloads = append(s.payloads, append([]byte(nil), payload...))
	return false, nil
}

func TestDispatchAnnouncementPushesOnlyOnceForAnnouncementsPublishedAfterSubscription(t *testing.T) {
	ctx := t.Context()
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	store := memory.NewCompetitionStore()
	_, err := store.ActivateContent(ctx, core.ContentSnapshot{
		CommitSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Manifest: core.Manifest{Announcements: []core.Announcement{
			{Slug: "old", Title: "購読前", EffectiveFrom: now.Add(-time.Hour)},
			{Slug: "current", Title: "新しい通知", EffectiveFrom: now},
			{Slug: "future", Title: "公開前", EffectiveFrom: now.Add(time.Hour)},
		}},
	}, core.NoActiveContentCommit)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.UpsertWebPushSubscription(ctx, core.WebPushSubscription{
		ContestantName: "contestant",
		Endpoint:       "https://push.example.test/subscription",
		P256DH:         "key",
		Auth:           "auth",
		CreatedAt:      now,
		UpdatedAt:      now,
	}); err != nil {
		t.Fatal(err)
	}
	sender := &recordingPushSender{}
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	svc.Now = func() time.Time { return now }
	svc.WebPush = sender
	svc.VAPIDPublicKey = "public"

	if err := svc.DispatchAnnouncementPushes(ctx); err != nil {
		t.Fatal(err)
	}
	if err := svc.DispatchAnnouncementPushes(ctx); err != nil {
		t.Fatal(err)
	}
	if len(sender.payloads) != 1 {
		t.Fatalf("sent payloads = %d, want 1", len(sender.payloads))
	}
}
