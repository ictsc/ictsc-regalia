package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// PublishAnswerChange tells open screens that an answer changed so they can
// refetch it. Delivery is best effort: a lost event is recovered when the
// client reconnects and refetches, so a failure is only logged.
func (s *Service) PublishAnswerChange(ctx context.Context, change core.AnswerChange) {
	if s.AnswerEvents == nil {
		return
	}
	encoded, err := json.Marshal(change)
	if err == nil {
		err = s.AnswerEvents.Publish(ctx, encoded)
	}
	if err != nil {
		log.Printf("publish answer change %d/%s/%d: %v", change.TeamCode, change.ProblemCode, change.AnswerNumber, err)
	}
}
