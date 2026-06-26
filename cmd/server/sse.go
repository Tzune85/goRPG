package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// AgentEvent is one thing an agent did, broadcast to all spectators.
type AgentEvent struct {
	Agent  string `json:"agent"`
	Action string `json:"action"`
	Output string `json:"output"`
	Done   bool   `json:"done"`
}

// Broadcaster fans out events to all connected spectator clients.
type Broadcaster struct {
	mu          sync.RWMutex
	subscribers map[chan AgentEvent]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[chan AgentEvent]struct{}),
	}
}

func (b *Broadcaster) subscribe() chan AgentEvent {
	ch := make(chan AgentEvent, 16)
	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) unsubscribe(ch chan AgentEvent) {
	b.mu.Lock()
	delete(b.subscribers, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Broadcaster) Publish(e AgentEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subscribers {
		select {
		case ch <- e:
		default: // drop if the subscriber is too slow
		}
	}
}

// GET /api/events — SSE stream consumed by watch.html
func (b *Broadcaster) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	ch := b.subscribe()
	defer b.unsubscribe(ch)

	fmt.Fprintf(w, "data: {\"connected\":true}\n\n")
	flusher.Flush()

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
