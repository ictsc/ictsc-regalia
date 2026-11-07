package fanout

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

// countingBus counts subscriptions of the inner bus and lets a test end them.
type countingBus struct {
	inner *memory.EventBus

	mu            sync.Mutex
	subscriptions []service.Subscription
	closed        int
}

func newCountingBus() *countingBus { return &countingBus{inner: memory.NewEventBus()} }

func (b *countingBus) Publish(ctx context.Context, event []byte) error {
	return b.inner.Publish(ctx, event)
}

func (b *countingBus) Subscribe(ctx context.Context) (service.Subscription, error) {
	subscription, err := b.inner.Subscribe(ctx)
	if err != nil {
		return nil, err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscriptions = append(b.subscriptions, subscription)
	return &countingSubscription{Subscription: subscription, bus: b}, nil
}

func (b *countingBus) counts() (subscribed, closed int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.subscriptions), b.closed
}

type countingSubscription struct {
	service.Subscription
	bus  *countingBus
	once sync.Once
}

func (s *countingSubscription) Close() error {
	s.once.Do(func() {
		s.bus.mu.Lock()
		s.bus.closed++
		s.bus.mu.Unlock()
	})
	return s.Subscription.Close()
}

func receive(t *testing.T, subscription service.Subscription) ([]byte, bool) {
	t.Helper()
	select {
	case payload, ok := <-subscription.Messages():
		return payload, ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for the subscription")
		return nil, false
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("condition not met in time")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestBusSharesOneInnerSubscription(t *testing.T) {
	inner := newCountingBus()
	bus := New(inner)
	first, err := bus.Subscribe(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	second, err := bus.Subscribe(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if subscribed, _ := inner.counts(); subscribed != 1 {
		t.Fatalf("inner subscriptions = %d, want 1", subscribed)
	}
	if err := bus.Publish(t.Context(), []byte("event")); err != nil {
		t.Fatal(err)
	}
	for _, subscription := range []service.Subscription{first, second} {
		if payload, ok := receive(t, subscription); !ok || string(payload) != "event" {
			t.Fatalf("payload = %q %v, want event", payload, ok)
		}
	}

	_ = first.Close()
	if _, closed := inner.counts(); closed != 0 {
		t.Fatal("inner subscription closed while a subscriber remains")
	}
	_ = second.Close()
	waitFor(t, func() bool { _, closed := inner.counts(); return closed == 1 })

	// A later subscriber starts a fresh inner subscription.
	third, err := bus.Subscribe(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer third.Close()
	if subscribed, _ := inner.counts(); subscribed != 2 {
		t.Fatalf("inner subscriptions = %d, want 2", subscribed)
	}
}

func TestBusClosesSubscribersWhenInnerSubscriptionEnds(t *testing.T) {
	inner := newCountingBus()
	bus := New(inner)
	first, _ := bus.Subscribe(t.Context())
	second, _ := bus.Subscribe(t.Context())
	inner.mu.Lock()
	ended := inner.subscriptions[0]
	inner.mu.Unlock()
	_ = ended.Close() // like the Redis bus after a reconnect

	for _, subscription := range []service.Subscription{first, second} {
		if _, ok := receive(t, subscription); ok {
			t.Fatal("subscriber stayed open after the inner subscription ended")
		}
		_ = subscription.Close()
	}
	next, err := bus.Subscribe(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer next.Close()
	if subscribed, _ := inner.counts(); subscribed != 2 {
		t.Fatalf("inner subscriptions = %d, want a fresh one", subscribed)
	}
}

func TestBusClosesOnlyASubscriberThatFallsBehind(t *testing.T) {
	bus := New(newCountingBus())
	slow, _ := bus.Subscribe(t.Context())
	fast, _ := bus.Subscribe(t.Context())
	defer fast.Close()
	for range bufferSize + 1 {
		if err := bus.Publish(t.Context(), []byte("event")); err != nil {
			t.Fatal(err)
		}
		if _, ok := receive(t, fast); !ok {
			t.Fatal("a subscriber that keeps up was closed")
		}
	}
	received := 0
	for range slow.Messages() {
		received++
	}
	if received != bufferSize {
		t.Fatalf("slow subscriber received %d events before closing, want %d", received, bufferSize)
	}
	_ = slow.Close()
}

// gatedBus holds every inner Subscribe until the test opens the gate, like a
// Redis that does not acknowledge the subscription.
type gatedBus struct {
	*countingBus
	gate chan struct{}
}

func (b *gatedBus) Subscribe(ctx context.Context) (service.Subscription, error) {
	select {
	case <-b.gate:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return b.countingBus.Subscribe(ctx)
}

func TestSubscribeStopsWaitingWhenCallerLeaves(t *testing.T) {
	inner := &gatedBus{countingBus: newCountingBus(), gate: make(chan struct{})}
	bus := New(inner)

	ctx, cancel := context.WithCancel(context.Background())
	returned := make(chan error, 1)
	go func() {
		_, err := bus.Subscribe(ctx)
		returned <- err
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-returned:
		if err == nil {
			t.Fatal("canceled Subscribe returned a subscription")
		}
	case <-time.After(time.Second):
		t.Fatal("Subscribe kept waiting after its caller left")
	}

	// Nobody wants it any more, so it is stopped once it is set up.
	close(inner.gate)
	waitFor(t, func() bool {
		subscribed, closed := inner.counts()
		return subscribed == 1 && closed == 1
	})

	// A later caller starts a fresh one and receives events.
	subscription, err := bus.Subscribe(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Close()
	if err := bus.Publish(context.Background(), []byte("after")); err != nil {
		t.Fatal(err)
	}
	if payload, ok := receive(t, subscription); !ok || string(payload) != "after" {
		t.Fatalf("received %q, %v", payload, ok)
	}
}

func TestSubscribeWaitersShareOneSetup(t *testing.T) {
	inner := &gatedBus{countingBus: newCountingBus(), gate: make(chan struct{})}
	bus := New(inner)

	results := make(chan service.Subscription, 2)
	for range 2 {
		go func() {
			subscription, err := bus.Subscribe(context.Background())
			if err != nil {
				t.Error(err)
			}
			results <- subscription
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(inner.gate)
	first, second := <-results, <-results
	if subscribed, _ := inner.counts(); subscribed != 1 {
		t.Fatalf("inner subscriptions = %d, want 1", subscribed)
	}
	if err := bus.Publish(context.Background(), []byte("both")); err != nil {
		t.Fatal(err)
	}
	for _, subscription := range []service.Subscription{first, second} {
		if payload, ok := receive(t, subscription); !ok || string(payload) != "both" {
			t.Fatalf("received %q, %v", payload, ok)
		}
		_ = subscription.Close()
	}
	waitFor(t, func() bool { _, closed := inner.counts(); return closed == 1 })
}
