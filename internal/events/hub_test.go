package events

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
)

func TestHubPublishSubscribe(t *testing.T) {
	hub := NewHub()
	ch, unsubscribe := hub.Subscribe()
	defer unsubscribe()

	item := &mail.Item{
		ID:       "abc123",
		Subject:  "Hello",
		From:     "sender@example.com",
		DateSent: time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
	}
	hub.Publish(MailReceived(item))

	select {
	case evt := <-ch:
		if evt.Type != TypeMailReceived {
			t.Fatalf("type = %q, want %q", evt.Type, TypeMailReceived)
		}
		var data MailReceivedData
		if err := json.Unmarshal(evt.Data, &data); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if data.ID != item.ID || data.Subject != item.Subject {
			t.Fatalf("unexpected data: %+v", data)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestHubMultipleSubscribers(t *testing.T) {
	hub := NewHub()
	ch1, unsub1 := hub.Subscribe()
	defer unsub1()
	ch2, unsub2 := hub.Subscribe()
	defer unsub2()

	hub.Publish(MailPruned(3))

	for i, ch := range []<-chan Event{ch1, ch2} {
		select {
		case evt := <-ch:
			if evt.Type != TypeMailPruned {
				t.Fatalf("subscriber %d: type = %q", i, evt.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("subscriber %d timed out", i)
		}
	}
}

func TestHubDropsSlowSubscriber(t *testing.T) {
	hub := NewHub()
	ch, unsubscribe := hub.Subscribe()
	defer unsubscribe()

	for i := 0; i < subscriberBuffer+1; i++ {
		hub.Publish(MailPruned(int64(i)))
	}

	select {
	case _, ok := <-ch:
		if !ok {
			return
		}
	default:
	}

	// Drain buffered events; a dropped subscriber's channel closes.
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			t.Fatal("expected slow subscriber to be dropped after buffer overflow")
		}
	}
}

func TestHubConcurrentPublish(t *testing.T) {
	hub := NewHub()
	ch, unsubscribe := hub.Subscribe()
	defer unsubscribe()

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			hub.Publish(MailPruned(1))
		}()
	}
	wg.Wait()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestUnsubscribeClosesChannel(t *testing.T) {
	hub := NewHub()
	ch, unsubscribe := hub.Subscribe()
	unsubscribe()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}
