package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// realtimeBuffer bounds how far one SSE connection may lag. A subscriber that
// overflows is closed instead of silently losing a change; its client then
// reconnects and refetches everything on the new ready event.
const realtimeBuffer = 64

// realtimeRetry is how long the watcher waits after a store or bus failure.
const realtimeRetry = 5 * time.Second

// rankingSettle coalesces bursts of markings into one ranking check, so every
// connected client refetches the ranking at most once per burst.
const rankingSettle = 2 * time.Second

// RealtimeResync is delivered locally when this instance may have missed
// cross-instance events, so every connection should refetch its screens.
const RealtimeResync core.RealtimeKind = "RESYNC"

type realtimeHub struct {
	mu          sync.Mutex
	nextID      int
	subscribers map[int]chan core.RealtimeEvent
}

func (h *realtimeHub) subscribe() (<-chan core.RealtimeEvent, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscribers == nil {
		h.subscribers = make(map[int]chan core.RealtimeEvent)
	}
	id := h.nextID
	h.nextID++
	channel := make(chan core.RealtimeEvent, realtimeBuffer)
	h.subscribers[id] = channel
	var once sync.Once
	return channel, func() {
		once.Do(func() {
			h.mu.Lock()
			defer h.mu.Unlock()
			if current, ok := h.subscribers[id]; ok {
				delete(h.subscribers, id)
				close(current)
			}
		})
	}
}

func (h *realtimeHub) dispatch(event core.RealtimeEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, channel := range h.subscribers {
		select {
		case channel <- event:
		default:
			delete(h.subscribers, id)
			close(channel)
		}
	}
}

// SubscribeRealtime returns this instance's realtime events. The channel is
// closed when the returned cancel function runs or when the subscriber falls
// too far behind.
func (s *Service) SubscribeRealtime() (<-chan core.RealtimeEvent, func()) {
	return s.realtime.subscribe()
}

// PublishRealtime notifies every API instance after a committed write.
// Delivery is best effort: the write already succeeded, and clients refetch
// on reconnect, so a bus failure is logged instead of failing the request.
func (s *Service) PublishRealtime(ctx context.Context, event core.RealtimeEvent) {
	if s.Realtime == nil {
		s.realtime.dispatch(event)
		return
	}
	encoded, err := json.Marshal(event)
	if err == nil {
		err = s.Realtime.Publish(ctx, encoded)
	}
	if err != nil {
		log.Printf("publish realtime %s: %v", event.Kind, err)
	}
}

// RunRealtime forwards cross-instance events to local connections and
// derives clock-based events until ctx ends.
func (s *Service) RunRealtime(ctx context.Context) {
	if s.Realtime != nil {
		go s.forwardRealtime(ctx)
	}
	s.watchRealtimeClock(ctx)
}

func (s *Service) forwardRealtime(ctx context.Context) {
	for first := true; ctx.Err() == nil; first = false {
		subscription, err := s.Realtime.Subscribe(ctx)
		if err != nil {
			log.Printf("subscribe realtime bus: %v", err)
			if !sleepContext(ctx, realtimeRetry) {
				return
			}
			continue
		}
		if !first {
			s.realtime.dispatch(core.RealtimeEvent{Kind: RealtimeResync})
		}
		for payload := range subscription.Messages() {
			var event core.RealtimeEvent
			if json.Unmarshal(payload, &event) == nil && event.Kind != "" {
				s.realtime.dispatch(event)
			}
		}
		_ = subscription.Close()
		if !sleepContext(ctx, realtimeRetry) {
			return
		}
	}
}

func (s *Service) watchRealtimeClock(ctx context.Context) {
	events, cancel := s.SubscribeRealtime()
	defer func() { cancel() }()
	cursor := s.Now()
	ranking, _ := s.publicRankingFingerprint(ctx)
	timer := time.NewTimer(0)
	defer timer.Stop()
	settle := time.NewTimer(0)
	defer settle.Stop()
	<-settle.C
	settling := false
	reschedule := func() {
		timer.Stop()
		select {
		case <-timer.C:
		default:
		}
		next, ok, err := s.nextRealtimeInstant(ctx, cursor)
		switch {
		case err != nil:
			timer.Reset(realtimeRetry)
		case ok:
			timer.Reset(max(next.Sub(s.Now()), 0))
		}
	}
	checkRankingSoon := func() {
		if !settling {
			settling = true
			settle.Reset(rankingSettle)
		}
	}
	reschedule()
	for {
		select {
		case <-ctx.Done():
			return
		case event, open := <-events:
			if !open {
				events, cancel = s.SubscribeRealtime()
				checkRankingSoon()
				reschedule()
				continue
			}
			switch event.Kind {
			case core.RealtimeAnswerSubmitted, core.RealtimeContentUpdated:
				reschedule()
			case core.RealtimeCompetitionUpdated, RealtimeResync:
				checkRankingSoon()
				reschedule()
			case core.RealtimeMarkingUpdated, core.RealtimeScoresRecalculated, core.RealtimeTeamsUpdated:
				checkRankingSoon()
			}
		case <-settle.C:
			settling = false
			ranking = s.checkRanking(ctx, ranking)
		case <-timer.C:
			now := s.Now()
			if err := s.emitClockEvents(ctx, cursor, now); err == nil {
				cursor = now
			}
			checkRankingSoon()
			reschedule()
		}
	}
}

// nextRealtimeInstant returns the first instant after `after` at which a
// screen may change without any write: an answer's publication delay ending,
// a section opening or closing, an announcement taking effect, or the ranking
// freeze. Every instant derives from data the contestant already sees or from
// the answer submission time, so firing never reveals whether a marking exists.
func (s *Service) nextRealtimeInstant(ctx context.Context, after time.Time) (time.Time, bool, error) {
	instants, err := s.realtimeInstants(ctx)
	if err != nil {
		return time.Time{}, false, err
	}
	var next time.Time
	found := false
	for _, instant := range instants {
		if instant.at.After(after) && (!found || instant.at.Before(next)) {
			next, found = instant.at, true
		}
	}
	return next, found, nil
}

// emitClockEvents delivers the clock events in (from, to] to local
// connections only; every instance derives the same events independently.
func (s *Service) emitClockEvents(ctx context.Context, from, to time.Time) error {
	instants, err := s.realtimeInstants(ctx)
	if err != nil {
		return err
	}
	scheduled := false
	for _, instant := range instants {
		if !instant.at.After(from) || instant.at.After(to) {
			continue
		}
		if instant.event.Kind == core.RealtimeScheduleReached {
			scheduled = true
			continue
		}
		if instant.event.Kind != "" {
			s.realtime.dispatch(instant.event)
		}
	}
	if scheduled {
		s.realtime.dispatch(core.RealtimeEvent{Kind: core.RealtimeScheduleReached})
	}
	return nil
}

type realtimeInstant struct {
	at    time.Time
	event core.RealtimeEvent
}

func (s *Service) realtimeInstants(ctx context.Context) ([]realtimeInstant, error) {
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{})
	if err != nil {
		return nil, err
	}
	state, err := s.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	instants := make([]realtimeInstant, 0, len(answers)+1)
	for _, answer := range answers {
		instants = append(instants, realtimeInstant{
			at: answer.SubmittedAt.Add(core.PublishDelay),
			event: core.RealtimeEvent{
				Kind: core.RealtimeAnswerPublishDelayElapsed, TeamCode: answer.TeamCode,
				ProblemCode: answer.ProblemCode, AnswerNumber: answer.Number,
			},
		})
	}
	if state.RankingFreezeAt != nil {
		// The freeze changes only the ranking, which checkRanking detects.
		instants = append(instants, realtimeInstant{at: *state.RankingFreezeAt})
	}
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		var domainErr *core.Error
		if errors.As(err, &domainErr) && domainErr.Code == "content_not_available" {
			return instants, nil
		}
		return nil, err
	}
	schedule := core.RealtimeEvent{Kind: core.RealtimeScheduleReached}
	for _, section := range snapshot.Manifest.Sections {
		instants = append(instants, realtimeInstant{at: section.Beginning, event: schedule}, realtimeInstant{at: section.Ending, event: schedule})
	}
	for _, announcement := range snapshot.Manifest.Announcements {
		instants = append(instants, realtimeInstant{at: announcement.EffectiveFrom, event: schedule})
	}
	return instants, nil
}

// checkRanking announces a public ranking change. Comparing the rendered
// public ranking keeps frozen and hidden rankings silent while they stay
// unchanged, regardless of what happened to private scores.
func (s *Service) checkRanking(ctx context.Context, previous string) string {
	current, err := s.publicRankingFingerprint(ctx)
	if err != nil {
		return previous
	}
	if current != previous {
		s.realtime.dispatch(core.RealtimeEvent{Kind: core.RealtimeRankingUpdated})
	}
	return current
}

func (s *Service) publicRankingFingerprint(ctx context.Context) (string, error) {
	entries, frozen, frozenAt, err := s.Ranking(ctx, false)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal(struct {
		Entries  []core.RankingEntry
		Frozen   bool
		FrozenAt *time.Time
	}{entries, frozen, frozenAt})
	return string(encoded), err
}

func sleepContext(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
