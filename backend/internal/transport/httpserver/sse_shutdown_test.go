package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

// stalledWriter blocks every Write like a client that stopped reading, until
// its write deadline is moved.
type stalledWriter struct {
	header   http.Header
	once     sync.Once
	deadline chan struct{}
}

func (w *stalledWriter) Header() http.Header { return w.header }
func (w *stalledWriter) WriteHeader(int)     {}
func (w *stalledWriter) Write([]byte) (int, error) {
	<-w.deadline
	return 0, os.ErrDeadlineExceeded
}
func (w *stalledWriter) SetWriteDeadline(time.Time) error {
	w.once.Do(func() { close(w.deadline) })
	return nil
}

func TestShutdownUnblocksStalledStreamWrites(t *testing.T) {
	done := make(chan struct{})
	writeErr := make(chan error, 1)
	handler := unblockStreamsOnShutdown(done, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("event: ready\n\n"))
		writeErr <- err
	}))
	writer := &stalledWriter{header: http.Header{}, deadline: make(chan struct{})}
	served := make(chan struct{})
	go func() {
		handler.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/contestant/answers/stream", nil))
		close(served)
	}()
	select {
	case <-writeErr:
		t.Fatal("write returned before shutdown")
	case <-time.After(50 * time.Millisecond):
	}
	close(done)
	select {
	case err := <-writeErr:
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Fatalf("write error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("stalled write was not unblocked at shutdown")
	}
	<-served
}

func TestShutdownLeavesOrdinaryResponsesAlone(t *testing.T) {
	done := make(chan struct{})
	close(done)
	handler := unblockStreamsOnShutdown(done, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	writer := &stalledWriter{header: http.Header{}, deadline: make(chan struct{})}
	finished := make(chan struct{})
	go func() {
		handler.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/contestant/problems", nil))
		close(finished)
	}()
	select {
	case <-finished:
		t.Fatal("ordinary response had its write deadline moved")
	case <-time.After(200 * time.Millisecond):
	}
	writer.SetWriteDeadline(time.Time{})
	<-finished
}
