package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/p0vidl0/mylslurper/internal/events"
)

const sseKeepAliveInterval = 30 * time.Second

func (a *API) handleEvents(w http.ResponseWriter, r *http.Request) {
	if a.Events == nil {
		writeError(w, http.StatusServiceUnavailable, "events not configured")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	eventCh, unsubscribe := a.Events.Subscribe()
	defer unsubscribe()

	keepAlive := time.NewTicker(sseKeepAliveInterval)
	defer keepAlive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-keepAlive.C:
			if _, err := fmt.Fprintf(w, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case evt, ok := <-eventCh:
			if !ok {
				return
			}
			if err := writeSSE(w, evt); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, evt events.Event) error {
	if _, err := fmt.Fprintf(w, "event: %s\n", evt.Type); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", evt.Data); err != nil {
		return err
	}
	return nil
}
