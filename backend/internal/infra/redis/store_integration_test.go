package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const redisTestImage = "redis:7.4-alpine"

func TestRedisStoreIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Redis integration tests are disabled by -short")
	}

	rawURL := startRedisContainer(t)

	t.Run("stores opaque sessions with a real TTL", func(t *testing.T) {
		store, err := Open(t.Context(), rawURL)
		if err != nil {
			t.Fatalf("open Redis session store: %v", err)
		}
		t.Cleanup(func() {
			if err := store.Close(); err != nil {
				t.Errorf("close Redis session store: %v", err)
			}
		})

		const ttl = 750 * time.Millisecond
		data := session.Data{
			Kind:           session.KindContestant,
			ContestantName: "integration-user",
			ImpersonatedBy: "admin-user",
		}
		token, err := store.Create(t.Context(), data, ttl)
		if err != nil {
			t.Fatalf("create session: %v", err)
		}
		if token == "" || token == data.ContestantName {
			t.Fatalf("session token is not opaque: %q", token)
		}

		remaining, err := store.client.PTTL(t.Context(), sessionPrefix+token).Result()
		if err != nil {
			t.Fatalf("read session TTL: %v", err)
		}
		if remaining <= 0 || remaining > ttl {
			t.Fatalf("session TTL = %s, want >0 and <=%s", remaining, ttl)
		}

		got, err := store.Get(t.Context(), token, session.KindContestant)
		if err != nil {
			t.Fatalf("get session: %v", err)
		}
		if got.ContestantName != data.ContestantName || got.ImpersonatedBy != data.ImpersonatedBy || got.CreatedAt.IsZero() {
			t.Fatalf("stored session = %+v, want contestant and impersonation metadata", got)
		}
		if _, err := store.Get(t.Context(), token, session.KindAdmin); !errors.Is(err, session.ErrNotFound) {
			t.Fatalf("get session with wrong kind error = %v, want ErrNotFound", err)
		}

		deadline := time.Now().Add(4 * time.Second)
		for {
			_, err := store.Get(t.Context(), token, session.KindContestant)
			if errors.Is(err, session.ErrNotFound) {
				break
			}
			if err != nil {
				t.Fatalf("wait for session expiry: %v", err)
			}
			if time.Now().After(deadline) {
				t.Fatal("session did not expire by its Redis TTL")
			}
			time.Sleep(25 * time.Millisecond)
		}

		if err := store.Delete(t.Context(), token); err != nil {
			t.Fatalf("idempotently delete expired session: %v", err)
		}
	})

	t.Run("delivers deployment events across store instances", func(t *testing.T) {
		publisherStore, err := Open(t.Context(), rawURL)
		if err != nil {
			t.Fatalf("open publisher Redis store: %v", err)
		}
		t.Cleanup(func() {
			if err := publisherStore.Close(); err != nil {
				t.Errorf("close publisher Redis store: %v", err)
			}
		})
		subscriberStore, err := Open(t.Context(), rawURL)
		if err != nil {
			t.Fatalf("open subscriber Redis store: %v", err)
		}
		t.Cleanup(func() {
			if err := subscriberStore.Close(); err != nil {
				t.Errorf("close subscriber Redis store: %v", err)
			}
		})

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		subscription, err := NewDeploymentBus(subscriberStore).Subscribe(ctx)
		if err != nil {
			t.Fatalf("subscribe to deployment bus: %v", err)
		}
		defer func() {
			if err := subscription.Close(); err != nil {
				t.Errorf("close deployment subscription: %v", err)
			}
		}()

		payload := []byte(`{"event_id":"e0189717-65ce-47f2-94b5-d07a03584a62","status":"DEPLOYING"}`)
		if err := NewDeploymentBus(publisherStore).Publish(t.Context(), payload); err != nil {
			t.Fatalf("publish deployment event: %v", err)
		}

		select {
		case got, ok := <-subscription.Messages():
			if !ok {
				t.Fatal("deployment subscription closed before delivery")
			}
			if string(got) != string(payload) {
				t.Fatalf("deployment event = %s, want %s", got, payload)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for cross-instance Redis Pub/Sub delivery")
		}
	})
}

func startRedisContainer(t *testing.T) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        redisTestImage,
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor: wait.ForLog("Ready to accept connections").
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("start %s (use -short to explicitly disable integration tests): %v", redisTestImage, err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("resolve Redis container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("resolve Redis container port: %v", err)
	}
	return fmt.Sprintf("redis://%s/0", net.JoinHostPort(host, port.Port()))
}
