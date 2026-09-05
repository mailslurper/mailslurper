// Package events provides a lightweight pub/sub hub for server-push
// notifications to connected HTTP clients.
package events

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
)

const (
	TypeMailReceived = "mail.received"
	TypeMailPruned   = "mail.pruned"

	subscriberBuffer = 16
)

// Event is a single notification pushed to subscribers.
type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MailReceivedData is the payload for mail.received events.
type MailReceivedData struct {
	ID          string `json:"id"`
	Subject     string `json:"subject"`
	FromAddress string `json:"fromAddress"`
	DateSent    string `json:"dateSent"`
}

// MailPrunedData is the payload for mail.pruned events.
type MailPrunedData struct {
	DeletedCount int64 `json:"deletedCount"`
}

// Hub fans out events to SSE subscribers without blocking publishers.
type Hub struct {
	mu          sync.RWMutex
	subscribers map[chan Event]struct{}
}

// NewHub returns an empty event hub.
func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[chan Event]struct{}),
	}
}

// Subscribe registers a buffered channel receiver. The returned unsubscribe
// function removes the subscriber and closes its channel.
func (h *Hub) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, subscriberBuffer)

	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		if _, ok := h.subscribers[ch]; ok {
			delete(h.subscribers, ch)
			close(ch)
		}
		h.mu.Unlock()
	}

	return ch, unsubscribe
}

// Publish delivers an event to all subscribers. Slow subscribers that fill
// their buffer are dropped so SMTP ingestion never blocks.
func (h *Hub) Publish(event Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var dropped []chan Event
	for ch := range h.subscribers {
		select {
		case ch <- event:
		default:
			dropped = append(dropped, ch)
		}
	}
	for _, ch := range dropped {
		delete(h.subscribers, ch)
		close(ch)
	}
}

// MailReceived builds a mail.received event from a stored item.
func MailReceived(item *mail.Item) Event {
	data, _ := json.Marshal(MailReceivedData{
		ID:          item.ID,
		Subject:     item.Subject,
		FromAddress: item.From,
		DateSent:    item.DateSent.UTC().Format(time.RFC3339),
	})
	return Event{Type: TypeMailReceived, Data: data}
}

// MailPruned builds a mail.pruned event.
func MailPruned(deletedCount int64) Event {
	data, _ := json.Marshal(MailPrunedData{DeletedCount: deletedCount})
	return Event{Type: TypeMailPruned, Data: data}
}
