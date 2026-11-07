// Package fanout shares one subscription of an event bus among many local
// subscribers, so an API instance holds one Redis Pub/Sub connection per
// channel instead of one per open SSE stream.
package fanout

import (
	"context"
	"sync"
	"time"

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
	// ready is closed once the inner subscription is set up; err tells
	// whether that failed. Both fields below are set before ready closes.
	ready        chan struct{}
	err          error
	subscription service.Subscription
	cancel       context.CancelFunc
	// waiting counts Subscribe calls between joining and receiving their
	// local subscriber, so the shared subscription is not stopped under them.
	waiting  int
	locals   map[*local]struct{}
	stopOnce sync.Once
}

// subscribeTimeout bounds setting up the inner subscription, so a stalled
// Redis cannot hold every new stream on this instance indefinitely.
const subscribeTimeout = 10 * time.Second

func New(inner service.EventBus) *Bus { return &Bus{inner: inner} }

func (b *Bus) Publish(ctx context.Context, event []byte) error {
	return b.inner.Publish(ctx, event)
}

// Subscribe returns once the shared subscription is active, so events
// published after it returns are delivered. Setting it up does not hold the
// lock, and a caller whose ctx ends stops waiting for it.
func (b *Bus) Subscribe(ctx context.Context) (service.Subscription, error) {
	b.mu.Lock()
	s := b.shared
	if s == nil {
		s = &shared{ready: make(chan struct{}), locals: make(map[*local]struct{})}
		b.shared = s
		go b.start(s)
	}
	s.waiting++
	b.mu.Unlock()

	select {
	case <-s.ready:
	case <-ctx.Done():
		b.mu.Lock()
		s.waiting--
		b.mu.Unlock()
		b.stopIfIdle(s)
		return nil, ctx.Err()
	}

	b.mu.Lock()
	s.waiting--
	if s.err != nil {
		b.mu.Unlock()
		return nil, s.err
	}
	if s.locals == nil {
		// It already ended, for example after Redis reconnected; start a
		// fresh one.
		b.mu.Unlock()
		return b.Subscribe(ctx)
	}
	l := &local{bus: b, shared: s, messages: make(chan []byte, bufferSize)}
	s.locals[l] = struct{}{}
	b.mu.Unlock()
	return l, nil
}

// start sets up the inner subscription. It outlives the request that
// started it, so it uses its own context.
func (b *Bus) start(s *shared) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(subscribeTimeout, cancel)
	subscription, err := b.inner.Subscribe(ctx)
	if !timer.Stop() && err == nil {
		// The timeout fired just as setup finished; the subscription is
		// already being torn down.
		_ = subscription.Close()
		err = context.DeadlineExceeded
	}
	b.mu.Lock()
	s.subscription, s.cancel, s.err = subscription, cancel, err
	if err != nil && b.shared == s {
		b.shared = nil
	}
	close(s.ready)
	b.mu.Unlock()
	if err != nil {
		cancel()
		return
	}
	go b.forward(s)
	// Every caller may have given up while it was being set up.
	b.stopIfIdle(s)
}

// stopIfIdle stops a set-up shared subscription that nobody uses or waits for.
func (b *Bus) stopIfIdle(s *shared) {
	b.mu.Lock()
	idle := b.shared == s && s.waiting == 0 && len(s.locals) == 0
	select {
	case <-s.ready:
		idle = idle && s.err == nil
	default:
		idle = false
	}
	if idle {
		b.shared = nil
	}
	b.mu.Unlock()
	if idle {
		s.stop()
	}
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
		b.mu.Unlock()
		b.stopIfIdle(l.shared)
	})
	return nil
}

var _ service.EventBus = (*Bus)(nil)
