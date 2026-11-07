package service

import (
	"context"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type recordingBus struct{ errs []error }

func (b *recordingBus) Publish(ctx context.Context, _ []byte) error {
	b.errs = append(b.errs, ctx.Err())
	return nil
}

func (b *recordingBus) Subscribe(context.Context) (Subscription, error) { return nil, nil }

func TestPublishAnswerChangeOutlivesTheRequest(t *testing.T) {
	bus := &recordingBus{}
	svc := &Service{AnswerEvents: bus}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	svc.PublishAnswerChange(ctx, core.AnswerChange{TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1, Kind: core.AnswerSubmitted})
	if len(bus.errs) != 1 || bus.errs[0] != nil {
		t.Fatalf("publish context errors = %v, want one live context", bus.errs)
	}
}
