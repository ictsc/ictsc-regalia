package session_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/sessions"
	"github.com/ictsc/ictsc-regalia/backend/pkg/redistest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/rbcervilla/redisstore/v9"
)

//nolint:cyclop
func TestSessionStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookie jar: %v", err)
	}

	serv := httptest.NewServer(setupTestHandler(t))
	t.Cleanup(serv.Close)
	client := serv.Client()
	client.Jar = jar

	servURL, err := url.Parse(serv.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}

	// 最初はセッションが存在しない
	{
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/session", nil)
		if err != nil {
			t.Fatalf("failed to create get request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send get request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
		}
		if t.Failed() {
			return
		}
	}
	// セッションを作成する
	{
		sess := &session.OAuth2Session{
			State:    "state",
			Verifier: "verifier",
		}
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(sess); err != nil {
			t.Errorf("failed to encode session: %v", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, serv.URL+"/session", &buf)
		if err != nil {
			t.Fatalf("failed to create post request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send post request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
		}
		if cookies := jar.Cookies(req.URL); len(cookies) == 0 {
			t.Error("expected to have cookies")
		}
		if t.Failed() {
			return
		}
	}
	// セッションを取得できる
	{
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/session", nil)
		if err != nil {
			t.Fatalf("failed to create get request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send get request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}

		expected, actual := session.OAuth2Session{
			State: "state", Verifier: "verifier",
		}, session.OAuth2Session{}
		if err := gob.NewDecoder(resp.Body).Decode(&actual); err != nil {
			t.Fatalf("failed to decode session: %v", err)
		}
		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Errorf("want(-), got(+):\n%s", diff)
		}

		if t.Failed() {
			return
		}
	}

	sessionCookies := jar.Cookies(servURL)
	// セッションを削除できる
	{
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serv.URL+"/session", nil)
		if err != nil {
			t.Fatalf("failed to create delete request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send delete request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
		if cookies := jar.Cookies(req.URL); len(cookies) != 0 {
			t.Errorf("expected to have no cookies, but got: %v", cookies)
		}
		if t.Failed() {
			return
		}
	}

	// セッションリプレイができない
	{
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, serv.URL+"/session", nil)
		if err != nil {
			t.Fatalf("failed to create get request: %v", err)
		}

		client.Jar.SetCookies(servURL, sessionCookies)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send get request: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
		}
		if t.Failed() {
			return
		}
	}
}

func setupTestHandler(t *testing.T) http.Handler {
	t.Helper()

	mux := http.NewServeMux()

	// Read
	mux.HandleFunc("GET /session", func(w http.ResponseWriter, r *http.Request) {
		val, err := session.OAuth2SessionStore.Get(r.Context())
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			t.Errorf("failed to get session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := gob.NewEncoder(w).Encode(val); err != nil {
			t.Errorf("failed to encode session: %v", err)
		}
	})

	// Create
	mux.HandleFunc("POST /session", func(w http.ResponseWriter, r *http.Request) {
		var sess session.OAuth2Session
		if err := gob.NewDecoder(r.Body).Decode(&sess); err != nil {
			t.Errorf("failed to decode session: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := session.OAuth2SessionStore.Write(r, w, &sess, &sessions.Options{MaxAge: 60}); err != nil {
			t.Errorf("failed to write session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	// Delete
	mux.HandleFunc("DELETE /session", func(w http.ResponseWriter, r *http.Request) {
		if err := session.OAuth2SessionStore.Write(r, w, nil, nil); err != nil {
			t.Errorf("failed to delete session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	rdb := redistest.SetupRedis(t)

	store, err := redisstore.NewRedisStore(context.Background(), rdb)
	if err != nil {
		t.Fatalf("failed to create session store: %v", err)
	}

	return session.NewHandler(store)(mux)
}
