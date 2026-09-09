package memory

import (
	"context"
	"sync"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

type sessionEntry struct {
	data      session.Data
	expiresAt time.Time
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]sessionEntry
	now      func() time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]sessionEntry), now: time.Now}
}

func (s *SessionStore) Create(_ context.Context, data session.Data, ttl time.Duration) (string, error) {
	token, err := session.Token()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	data.CreatedAt = s.now().UTC()
	s.sessions[token] = sessionEntry{data: data, expiresAt: s.now().Add(ttl)}
	return token, nil
}

func (s *SessionStore) Get(_ context.Context, token string, expected session.Kind) (session.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.sessions[token]
	if !ok || !s.now().Before(entry.expiresAt) || entry.data.Kind != expected {
		delete(s.sessions, token)
		return session.Data{}, session.ErrNotFound
	}
	return entry.data, nil
}

func (s *SessionStore) Consume(_ context.Context, token string, expected session.Kind) (session.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.sessions[token]
	if !ok || !s.now().Before(entry.expiresAt) || entry.data.Kind != expected {
		if ok && !s.now().Before(entry.expiresAt) {
			delete(s.sessions, token)
		}
		return session.Data{}, session.ErrNotFound
	}
	delete(s.sessions, token)
	return entry.data, nil
}

func (s *SessionStore) Delete(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
	return nil
}

var _ session.Store = (*SessionStore)(nil)
