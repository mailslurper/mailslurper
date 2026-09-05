package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/auth"
	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/events"
	"github.com/p0vidl0/mylslurper/internal/mail"
)

func newTestAPIWithEvents(t *testing.T, cfg *config.Config) *API {
	t.Helper()
	a := newTestAPI(t, cfg)
	a.Events = events.NewHub()
	return a
}

func TestEventsRequiresAuthWhenEnabled(t *testing.T) {
	cfg := config.Default()
	cfg.AuthenticationScheme = config.AuthSchemeBasic
	cfg.AuthSecret = "test-secret"
	hash, _ := auth.HashPassword([]byte("s3cret"))
	cfg.Credentials = map[string]string{"alex": string(hash)}

	a := newTestAPIWithEvents(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/events", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without credentials, got %d", rec.Code)
	}
}

func TestEventsStreamFormat(t *testing.T) {
	cfg := config.Default()
	a := newTestAPIWithEvents(t, cfg)
	router := a.Router()

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		router.ServeHTTP(rec, req)
		close(done)
	}()

	item := &mail.Item{
		ID:       "evt-1",
		Subject:  "Test",
		From:     "a@example.com",
		DateSent: time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC),
	}

	time.Sleep(50 * time.Millisecond)
	a.Events.Publish(events.MailReceived(item))
	time.Sleep(50 * time.Millisecond)

	cancel()
	<-done

	body := rec.Body.String()
	if !strings.Contains(body, "event: mail.received") || !strings.Contains(body, `"id":"evt-1"`) {
		t.Fatalf("unexpected SSE body: %q", body)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
}

func TestPrunePublishesEvent(t *testing.T) {
	cfg := config.Default()
	a := newTestAPIWithEvents(t, cfg)
	router := a.Router()

	ch, unsub := a.Events.Subscribe()
	defer unsub()

	if _, err := a.Store.StoreMail(context.Background(), &mail.Item{
		ID: "m1", From: "a@b.com", Subject: "x", DateSent: time.Now(),
	}); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/mail/prune", jsonBody(`{"code":"all"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("prune = %d: %s", rec.Code, rec.Body.String())
	}

	select {
	case evt := <-ch:
		if evt.Type != events.TypeMailPruned {
			t.Fatalf("event type = %q, want %q", evt.Type, events.TypeMailPruned)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for prune event")
	}
}

func TestEventsUnavailableWithoutHub(t *testing.T) {
	cfg := config.Default()
	a := newTestAPI(t, cfg)
	a.Events = nil
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/events", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 without hub, got %d", rec.Code)
	}
}
