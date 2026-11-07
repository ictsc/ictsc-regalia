package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// answerPublishTimeout bounds a publish that no longer follows the request.
const answerPublishTimeout = 5 * time.Second

// PublishAnswerChange tells open screens that an answer changed so they can
// refetch it. The change is already stored, so the publish must not fail
// just because the client went away; it ignores the request's cancellation.
// A failure is only logged: if Redis itself is unavailable, the subscriber
// connections drop too, and reconnecting browsers refetch from the database.
func (s *Service) PublishAnswerChange(ctx context.Context, change core.AnswerChange) {
	if s.AnswerEvents == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), answerPublishTimeout)
	defer cancel()
	encoded, err := json.Marshal(change)
	if err == nil {
		err = s.AnswerEvents.Publish(ctx, encoded)
	}
	if err != nil {
		log.Printf("publish answer change %d/%s/%d: %v", change.TeamCode, change.ProblemCode, change.AnswerNumber, err)
	}
}
