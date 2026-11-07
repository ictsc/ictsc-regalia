// Package fanout shares one subscription of an event bus among many local
// subscribers, so an API instance holds one Redis Pub/Sub connection per
// channel instead of one per open SSE stream.
package fanout

import (
	"context"
	"sync"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

// bufferSize is how many events a local subscriber may fall behind before it
// is closed. Closing ends its SSE response, so the browser reconnects and
// refetches instead of silently missing events.
const bufferSize = 64

// Bus subscribes to the inner bus while it has local subscribers and closes
// that subscription when the last one leaves. When the inner subscription
// ends, for example after Redis reconnected and lost events, every local
// subscriber is closed and the next Subscribe starts a fresh one.
type Bus struct {
	inner service.EventBus

	mu     sync.Mutex
	shared *shared
}

type shared struct {
	subscription service.Subscription
	cancel       context.CancelFunc
	locals       map[*local]struct{}
	stopOnce     sync.Once
}

func New(inner service.EventBus) *Bus { return &Bus{inner: inner} }

func (b *Bus) Publish(ctx context.Context, event []byte) error {
	return b.inner.Publish(ctx, event)
}

// Subscribe returns once the shared subscription is active, so events
// published after it returns are delivered.
func (b *Bus) Subscribe(context.Context) (service.Subscription, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.shared == nil {
		// The shared subscription outlives the request that started it.
		ctx, cancel := context.WithCancel(context.Background())
		subscription, err := b.inner.Subscribe(ctx)
		if err != nil {
			cancel()
			return nil, err
		}
		b.shared = &shared{subscription: subscription, cancel: cancel, locals: make(map[*local]struct{})}
		go b.forward(b.shared)
	}
	l := &local{bus: b, shared: b.shared, messages: make(chan []byte, bufferSize)}
	b.shared.locals[l] = struct{}{}
	return l, nil
}

func (b *Bus) forward(s *shared) {
	for payload := range s.subscription.Messages() {
		b.mu.Lock()
		for l := range s.locals {
			select {
			case l.messages <- append([]byte(nil), payload...):
			default:
				// Too far behind: end it rather than drop an event.
				delete(s.locals, l)
				close(l.messages)
			}
		}
		b.mu.Unlock()
	}
	b.mu.Lock()
	if b.shared == s {
		b.shared = nil
	}
	for l := range s.locals {
		close(l.messages)
	}
	s.locals = nil
	b.mu.Unlock()
	s.stop()
}

func (s *shared) stop() {
	s.stopOnce.Do(func() {
		s.cancel()
		_ = s.subscription.Close()
	})
}

type local struct {
	bus      *Bus
	shared   *shared
	messages chan []byte
	once     sync.Once
}

func (l *local) Messages() <-chan []byte { return l.messages }

func (l *local) Close() error {
	l.once.Do(func() {
		b := l.bus
		b.mu.Lock()
		if _, ok := l.shared.locals[l]; ok {
			delete(l.shared.locals, l)
			close(l.messages)
		}
		last := len(l.shared.locals) == 0 && b.shared == l.shared
		if last {
			b.shared = nil
		}
		b.mu.Unlock()
		if last {
			l.shared.stop()
		}
	})
	return nil
}

var _ service.EventBus = (*Bus)(nil)
