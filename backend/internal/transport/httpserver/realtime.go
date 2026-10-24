package httpserver

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

// realtimeNotice is the empty payload of events that name no target.
type realtimeNotice struct{}

type contestantRealtimeProblemEvent struct {
	ProblemCode string `json:"problem_code"`
}

type adminRealtimeAnswerEvent struct {
	Reason string                   `json:"reason"`
	Answer api.AdminAnswerReference `json:"answer"`
}

// realtimeMessage is one SSE event; an empty name drops the event.
type realtimeMessage struct {
	name string
	data any
}

func (h *Handler) StreamContestantEvents(ctx context.Context, _ api.StreamContestantEventsRequestObject) (api.StreamContestantEventsResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	team := contestant.TeamCode
	events, cancel := h.service.SubscribeRealtime()
	body := h.realtimeStream(ctx, events, cancel, func(event core.RealtimeEvent) realtimeMessage {
		return h.contestantRealtimeMessage(ctx, team, event)
	})
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamContestantEvents200TexteventStreamResponse{
		Body: body, Headers: api.StreamContestantEvents200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

func (h *Handler) StreamAdminEvents(ctx context.Context, _ api.StreamAdminEventsRequestObject) (api.StreamAdminEventsResponseObject, error) {
	events, cancel := h.service.SubscribeRealtime()
	body := h.realtimeStream(ctx, events, cancel, adminRealtimeMessage)
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamAdminEvents200TexteventStreamResponse{
		Body: body, Headers: api.StreamAdminEvents200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

// contestantRealtimeMessage keeps other teams' activity and every private
// marking out of the stream: a team hears about a marking only once
// ContestantScores would already show it.
func (h *Handler) contestantRealtimeMessage(ctx context.Context, team int64, event core.RealtimeEvent) realtimeMessage {
	switch event.Kind {
	case service.RealtimeResync:
		return realtimeMessage{"ready", realtimeNotice{}}
	case core.RealtimeAnswerSubmitted:
		if event.TeamCode == team {
			return realtimeMessage{"answers", contestantRealtimeProblemEvent{ProblemCode: event.ProblemCode}}
		}
	case core.RealtimeAnswerPublishDelayElapsed:
		if event.TeamCode == team {
			return realtimeMessage{"scores", contestantRealtimeProblemEvent{ProblemCode: event.ProblemCode}}
		}
	case core.RealtimeMarkingUpdated:
		if event.TeamCode != team {
			break
		}
		if visible, err := h.service.TeamSeesMarkingNow(ctx, event.TeamCode, event.ProblemCode, event.AnswerNumber); err == nil && visible {
			return realtimeMessage{"scores", contestantRealtimeProblemEvent{ProblemCode: event.ProblemCode}}
		}
	case core.RealtimeRankingUpdated:
		return realtimeMessage{"ranking", realtimeNotice{}}
	case core.RealtimeContentUpdated, core.RealtimeScheduleReached:
		return realtimeMessage{"content", realtimeNotice{}}
	case core.RealtimeTeamsUpdated:
		return realtimeMessage{"teams", realtimeNotice{}}
	case core.RealtimeCompetitionUpdated:
		return realtimeMessage{"competition", realtimeNotice{}}
	}
	return realtimeMessage{}
}

func adminRealtimeMessage(event core.RealtimeEvent) realtimeMessage {
	answer := func(reason string) realtimeMessage {
		return realtimeMessage{"answer", adminRealtimeAnswerEvent{Reason: reason, Answer: api.AdminAnswerReference{
			TeamCode: event.TeamCode, ProblemCode: event.ProblemCode, AnswerNumber: event.AnswerNumber,
		}}}
	}
	switch event.Kind {
	case service.RealtimeResync:
		return realtimeMessage{"ready", realtimeNotice{}}
	case core.RealtimeAnswerSubmitted:
		return answer("SUBMITTED")
	case core.RealtimeAnswerWorkflowUpdated:
		return answer("WORKFLOW_UPDATED")
	case core.RealtimeMarkingUpdated:
		return answer("MARKING_UPDATED")
	case core.RealtimeAnswerPublishDelayElapsed:
		return answer("PUBLICATION_DELAY_ELAPSED")
	case core.RealtimeScoresRecalculated:
		return realtimeMessage{"scores", realtimeNotice{}}
	case core.RealtimeRankingUpdated:
		return realtimeMessage{"ranking", realtimeNotice{}}
	case core.RealtimeContentUpdated, core.RealtimeScheduleReached:
		return realtimeMessage{"content", realtimeNotice{}}
	case core.RealtimeTeamsUpdated:
		return realtimeMessage{"teams", realtimeNotice{}}
	case core.RealtimeCompetitionUpdated:
		return realtimeMessage{"competition", realtimeNotice{}}
	}
	return realtimeMessage{}
}

// realtimeStream writes ready on every (re)connection so clients refetch
// what they missed, then maps hub events until the client leaves. A closed
// hub channel means this connection fell behind; ending the stream makes the
// browser reconnect and receive a fresh ready event.
func (h *Handler) realtimeStream(ctx context.Context, events <-chan core.RealtimeEvent, cancel func(), mapEvent func(core.RealtimeEvent) realtimeMessage) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		defer cancel()
		defer writer.Close()
		if writeSSE(writer, "ready", uuid.NewString(), realtimeNotice{}) != nil {
			return
		}
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
					return
				}
			case event, ok := <-events:
				if !ok {
					return
				}
				message := mapEvent(event)
				if message.name != "" && writeSSE(writer, message.name, uuid.NewString(), message.data) != nil {
					return
				}
			}
		}
	}()
	return reader
}
