package memory

import (
	"context"
	"sync"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

type EventBus struct {
	mu          sync.Mutex
	nextID      int
	subscribers map[int]chan []byte
}

func NewEventBus() *EventBus { return &EventBus{subscribers: make(map[int]chan []byte)} }

func (b *EventBus) Publish(_ context.Context, event []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, channel := range b.subscribers {
		payload := append([]byte(nil), event...)
		select {
		case channel <- payload:
		default:
		}
	}
	return nil
}

func (b *EventBus) Subscribe(context.Context) (service.Subscription, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.nextID
	b.nextID++
	channel := make(chan []byte, 64)
	b.subscribers[id] = channel
	return &memorySubscription{bus: b, id: id, messages: channel}, nil
}

type memorySubscription struct {
	bus      *EventBus
	id       int
	messages chan []byte
	once     sync.Once
}

func (s *memorySubscription) Messages() <-chan []byte { return s.messages }

func (s *memorySubscription) Close() error {
	s.once.Do(func() {
		s.bus.mu.Lock()
		delete(s.bus.subscribers, s.id)
		close(s.messages)
		s.bus.mu.Unlock()
	})
	return nil
}

var _ service.EventBus = (*EventBus)(nil)
