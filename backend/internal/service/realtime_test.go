package service

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type realtimeStore struct {
	visibilityStore
	snapshot *core.ContentSnapshot
	teams    []core.Team
}

func (s *realtimeStore) ActiveContent(context.Context) (core.ContentSnapshot, error) {
	if s.snapshot == nil {
		return core.ContentSnapshot{}, core.NewError(http.StatusBadGateway, "content_not_available", "No valid content snapshot is available")
	}
	return *s.snapshot, nil
}

func (s *realtimeStore) ListTeams(context.Context) ([]core.Team, error) { return s.teams, nil }

func (s *realtimeStore) GetAnswer(_ context.Context, teamCode int64, problemCode string, number int32) (core.Answer, error) {
	for _, answer := range s.answers {
		if answer.TeamCode == teamCode && answer.ProblemCode == problemCode && answer.Number == number {
			return answer, nil
		}
	}
	return core.Answer{}, core.NewError(http.StatusNotFound, "resource_not_found", "Answer was not found")
}

func drainRealtime(events <-chan core.RealtimeEvent) []core.RealtimeEvent {
	var drained []core.RealtimeEvent
	for {
		select {
		case event := <-events:
			drained = append(drained, event)
		default:
			return drained
		}
	}
}

func TestRealtimeHubClosesLaggingSubscriberInsteadOfDroppingEvents(t *testing.T) {
	var hub realtimeHub
	lagging, cancelLagging := hub.subscribe()
	defer cancelLagging()
	for range realtimeBuffer + 1 {
		hub.dispatch(core.RealtimeEvent{Kind: core.RealtimeAnswerSubmitted})
	}
	count := 0
	for range lagging {
		count++
	}
	if count != realtimeBuffer {
		t.Fatalf("lagging subscriber received %d events before closing, want %d", count, realtimeBuffer)
	}
	fresh, cancelFresh := hub.subscribe()
	cancelFresh()
	cancelFresh()
	if _, open := <-fresh; open {
		t.Fatal("cancelled subscriber channel stayed open")
	}
}

func TestRealtimeClockEventsFollowSubmissionTimesAndSchedule(t *testing.T) {
	submitted := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	sectionStart := submitted.Add(5 * time.Minute)
	store := &realtimeStore{
		visibilityStore: visibilityStore{answers: []core.Answer{{TeamCode: 2, ProblemCode: "P1", Number: 1, SubmittedAt: submitted}}},
		snapshot: &core.ContentSnapshot{Manifest: core.Manifest{
			Sections:      []core.Section{{Slug: "day-1", Beginning: sectionStart, Ending: submitted.Add(2 * time.Hour)}},
			Announcements: []core.Announcement{{Slug: "a", EffectiveFrom: sectionStart}},
		}},
	}
	service := New(store, nil, Config{})
	events, cancel := service.SubscribeRealtime()
	defer cancel()

	next, ok, err := service.nextRealtimeInstant(context.Background(), submitted)
	if err != nil || !ok || !next.Equal(sectionStart) {
		t.Fatalf("next instant = %v %v %v, want section start", next, ok, err)
	}
	if err := service.emitClockEvents(context.Background(), submitted, sectionStart); err != nil {
		t.Fatal(err)
	}
	got := drainRealtime(events)
	if len(got) != 1 || got[0].Kind != core.RealtimeScheduleReached {
		t.Fatalf("section start events = %#v, want one SCHEDULE_REACHED", got)
	}

	next, ok, err = service.nextRealtimeInstant(context.Background(), sectionStart)
	if err != nil || !ok || !next.Equal(submitted.Add(core.PublishDelay)) {
		t.Fatalf("next instant = %v %v %v, want publication delay end", next, ok, err)
	}
	if err := service.emitClockEvents(context.Background(), sectionStart, next); err != nil {
		t.Fatal(err)
	}
	got = drainRealtime(events)
	want := core.RealtimeEvent{Kind: core.RealtimeAnswerPublishDelayElapsed, TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("delay events = %#v, want %#v", got, want)
	}
}

func TestRealtimeRankingAnnouncesOnlyPublicChanges(t *testing.T) {
	submitted := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	store := &realtimeStore{
		visibilityStore: visibilityStore{answers: []core.Answer{{TeamCode: 2, ProblemCode: "P1", Number: 1, SubmittedAt: submitted, MaxScore: 100}}},
		snapshot:        &core.ContentSnapshot{Manifest: core.Manifest{Problems: []core.Problem{{Code: "P1", MaxScore: 100}}}},
		teams:           []core.Team{{Code: 2, Name: "Two"}, {Code: 3, Name: "Three"}},
	}
	service := New(store, nil, Config{})
	now := submitted.Add(time.Minute)
	service.Now = func() time.Time { return now }
	events, cancel := service.SubscribeRealtime()
	defer cancel()
	ranking, err := service.publicRankingFingerprint(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	store.markings = []core.MarkingResult{{ID: "m1", TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1, MarkedScore: 80, CreatedAt: now}}
	ranking = service.checkRanking(context.Background(), ranking)
	if got := drainRealtime(events); len(got) != 0 {
		t.Fatalf("marking during publication delay announced %#v", got)
	}

	now = submitted.Add(core.PublishDelay)
	ranking = service.checkRanking(context.Background(), ranking)
	if got := drainRealtime(events); len(got) != 1 || got[0].Kind != core.RealtimeRankingUpdated {
		t.Fatalf("published marking events = %#v, want RANKING_UPDATED", got)
	}
	_ = service.checkRanking(context.Background(), ranking)
	if got := drainRealtime(events); len(got) != 0 {
		t.Fatalf("unchanged ranking announced %#v", got)
	}
}

func TestTeamSeesMarkingNowMirrorsContestantScores(t *testing.T) {
	submitted := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	store := &realtimeStore{visibilityStore: visibilityStore{answers: []core.Answer{{TeamCode: 2, ProblemCode: "P1", Number: 1, SubmittedAt: submitted}}}}
	service := New(store, nil, Config{})
	now := submitted.Add(core.PublishDelay - time.Nanosecond)
	service.Now = func() time.Time { return now }
	if visible, err := service.TeamSeesMarkingNow(context.Background(), 2, "P1", 1); err != nil || visible {
		t.Fatalf("visible during delay = %v, %v", visible, err)
	}
	revealed := submitted
	store.state.FinalRevealedAt = &revealed
	if visible, err := service.TeamSeesMarkingNow(context.Background(), 2, "P1", 1); err != nil || !visible {
		t.Fatalf("visible after final reveal = %v, %v", visible, err)
	}
	store.state.FinalRevealedAt = nil
	now = submitted.Add(core.PublishDelay)
	if visible, err := service.TeamSeesMarkingNow(context.Background(), 2, "P1", 1); err != nil || !visible {
		t.Fatalf("visible after delay = %v, %v", visible, err)
	}
}

type channelBus struct {
	mu          sync.Mutex
	subscribers []chan []byte
}

func (b *channelBus) Publish(_ context.Context, event []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, subscriber := range b.subscribers {
		subscriber <- append([]byte(nil), event...)
	}
	return nil
}

func (b *channelBus) Subscribe(context.Context) (Subscription, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	channel := make(chan []byte, 8)
	b.subscribers = append(b.subscribers, channel)
	return channelSubscription(channel), nil
}

type channelSubscription chan []byte

func (s channelSubscription) Messages() <-chan []byte { return s }
func (s channelSubscription) Close() error            { return nil }

func TestPublishRealtimeTravelsThroughBusToLocalSubscribers(t *testing.T) {
	store := &realtimeStore{}
	service := New(store, nil, Config{})
	bus := &channelBus{}
	service.Realtime = bus
	events, cancel := service.SubscribeRealtime()
	defer cancel()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	go service.forwardRealtime(ctx)
	deadline := time.Now().Add(2 * time.Second)
	for {
		bus.mu.Lock()
		ready := len(bus.subscribers) > 0
		bus.mu.Unlock()
		if ready {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("forwarder did not subscribe to the bus")
		}
		time.Sleep(time.Millisecond)
	}
	want := core.RealtimeEvent{Kind: core.RealtimeAnswerSubmitted, TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1}
	service.PublishRealtime(context.Background(), want)
	select {
	case got := <-events:
		if got != want {
			t.Fatalf("forwarded event = %#v, want %#v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("event published to the bus never reached local subscribers")
	}
}

func TestRunRealtimeCoalescesMarkingsIntoOneRankingUpdate(t *testing.T) {
	submitted := time.Now().UTC().Add(-time.Hour)
	store := &realtimeStore{
		visibilityStore: visibilityStore{answers: []core.Answer{{TeamCode: 2, ProblemCode: "P1", Number: 1, SubmittedAt: submitted, MaxScore: 100}}},
		snapshot:        &core.ContentSnapshot{Manifest: core.Manifest{Problems: []core.Problem{{Code: "P1", MaxScore: 100}}}},
		teams:           []core.Team{{Code: 2, Name: "Two"}},
	}
	var mu sync.Mutex
	service := New(&lockedRealtimeStore{realtimeStore: store, mu: &mu}, nil, Config{})
	events, cancel := service.SubscribeRealtime()
	defer cancel()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	go service.RunRealtime(ctx)
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	store.markings = []core.MarkingResult{{ID: "m1", TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1, MarkedScore: 80, CreatedAt: submitted}}
	mu.Unlock()
	for range 3 {
		service.PublishRealtime(context.Background(), core.RealtimeEvent{Kind: core.RealtimeMarkingUpdated, TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1})
	}
	deadline := time.After(rankingSettle + 2*time.Second)
	rankingUpdates := 0
	for rankingUpdates == 0 {
		select {
		case event := <-events:
			if event.Kind == core.RealtimeRankingUpdated {
				rankingUpdates++
			}
		case <-deadline:
			t.Fatal("ranking update was never announced")
		}
	}
	time.Sleep(rankingSettle + 200*time.Millisecond)
	for _, event := range drainRealtime(events) {
		if event.Kind == core.RealtimeRankingUpdated {
			t.Fatal("a burst of markings announced more than one ranking update")
		}
	}
}

// lockedRealtimeStore lets the test mutate markings while the watcher reads.
type lockedRealtimeStore struct {
	*realtimeStore
	mu *sync.Mutex
}

func (s *lockedRealtimeStore) ListMarkingResults(ctx context.Context) ([]core.MarkingResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.realtimeStore.ListMarkingResults(ctx)
}
