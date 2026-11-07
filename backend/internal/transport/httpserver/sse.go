package httpserver

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const sseKeepaliveInterval = 15 * time.Second

func subscribe(ctx context.Context, bus service.EventBus) (service.Subscription, error) {
	if bus == nil {
		return nil, nil
	}
	return bus.Subscribe(ctx)
}

func closeSubscription(subscription service.Subscription) {
	if subscription != nil {
		_ = subscription.Close()
	}
}

// streamSubscription writes a keepalive comment every 15 seconds and one SSE
// event for each bus message that mapMessage accepts, until the client goes
// away, the subscription ends or the server shuts down. A nil subscription
// only sends keepalives. Browsers reconnect by themselves after it returns.
func (h *Handler) streamSubscription(ctx context.Context, writer io.Writer, subscription service.Subscription, mapMessage func([]byte) (event, id string, body any, ok bool)) {
	ticker := time.NewTicker(sseKeepaliveInterval)
	defer ticker.Stop()
	var messages <-chan []byte
	if subscription != nil {
		messages = subscription.Messages()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-h.options.StreamsDone:
			return
		case <-ticker.C:
			if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
				return
			}
		case payload, ok := <-messages:
			if !ok {
				return
			}
			event, id, body, include := mapMessage(payload)
			if include && writeSSE(writer, event, id, body) != nil {
				return
			}
		}
	}
}

// unblockStreamsOnShutdown lets SSE responses end promptly at shutdown even
// when a client stopped reading. Without it, a write stuck on a full TCP
// buffer would keep both the response and the stream producer from seeing
// StreamsDone, and Shutdown would wait for its whole grace period. Moving the
// write deadline to now fails that write, the response returns and closes
// the pipe, and the producer's next write fails too.
func unblockStreamsOnShutdown(done <-chan struct{}, next http.Handler) http.Handler {
	if done == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isSSEPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		finished := make(chan struct{})
		defer close(finished)
		go func() {
			select {
			case <-done:
				_ = http.NewResponseController(w).SetWriteDeadline(time.Now())
			case <-finished:
			}
		}()
		next.ServeHTTP(w, r)
	})
}
