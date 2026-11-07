package httpserver

import (
	"context"
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

// The answer streams carry only the changed answer's identity. Clients
// refetch what they may see, so hidden scores never travel over the stream.

type answerStreamReady struct{}

type contestantAnswerChange struct {
	ProblemCode  string `json:"problem_code"`
	AnswerNumber int32  `json:"answer_number"`
}

type adminAnswerChange struct {
	TeamCode     int64                 `json:"team_code"`
	ProblemCode  string                `json:"problem_code"`
	AnswerNumber int32                 `json:"answer_number"`
	Kind         core.AnswerChangeKind `json:"kind"`
}

func (h *Handler) StreamContestantAnswers(ctx context.Context, _ api.StreamContestantAnswersRequestObject) (api.StreamContestantAnswersResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	team := contestant.TeamCode
	subscription, err := subscribe(ctx, h.service.AnswerEvents)
	if err != nil {
		closeSubscription(subscription)
		return nil, err
	}
	body := h.answerStream(ctx, subscription, func(change core.AnswerChange) (any, bool) {
		if change.TeamCode != team {
			return nil, false
		}
		return contestantAnswerChange{ProblemCode: change.ProblemCode, AnswerNumber: change.AnswerNumber}, true
	})
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamContestantAnswers200TexteventStreamResponse{
		Body: body, Headers: api.StreamContestantAnswers200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

func (h *Handler) StreamAdminAnswers(ctx context.Context, _ api.StreamAdminAnswersRequestObject) (api.StreamAdminAnswersResponseObject, error) {
	subscription, err := subscribe(ctx, h.service.AnswerEvents)
	if err != nil {
		closeSubscription(subscription)
		return nil, err
	}
	body := h.answerStream(ctx, subscription, func(change core.AnswerChange) (any, bool) {
		return adminAnswerChange(change), true
	})
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamAdminAnswers200TexteventStreamResponse{
		Body: body, Headers: api.StreamAdminAnswers200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

// answerStream sends a ready event first so that a reconnecting client
// refetches whatever it missed, then one answer event per accepted change.
func (h *Handler) answerStream(ctx context.Context, subscription service.Subscription, mapChange func(core.AnswerChange) (any, bool)) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		defer closeSubscription(subscription)
		defer writer.Close()
		if err := writeSSE(writer, "ready", uuid.NewString(), answerStreamReady{}); err != nil {
			return
		}
		h.streamSubscription(ctx, writer, subscription, func(payload []byte) (string, string, any, bool) {
			var change core.AnswerChange
			if json.Unmarshal(payload, &change) != nil {
				return "", "", nil, false
			}
			body, include := mapChange(change)
			return "answer", uuid.NewString(), body, include
		})
	}()
	return reader
}
