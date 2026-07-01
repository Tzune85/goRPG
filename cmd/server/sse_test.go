package main

import "testing"

func TestPublish(t *testing.T) {
	b := NewBroadcaster()
	ch := b.subscribe()

	event := AgentEvent{Agent: "TestBot", Action: "north", Output: "you moved", Done: false}
	b.Publish(event)

	received := <-ch
	if received.Agent != event.Agent {
		t.Errorf("expected agent %s, got %s", event.Agent, received.Agent)
	}
	if received.Action != event.Action {
		t.Errorf("expected action %s, got %s", event.Action, received.Action)
	}
}

func TestSubcribe(t *testing.T) {
	b := NewBroadcaster()
	ch := b.subscribe()

	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
	if _, ok := b.subscribers[ch]; !ok {
		t.Error("expected ch to be in subscribers")
	}
}

func TestUnsubcribe(t *testing.T) {
	b := NewBroadcaster()
	ch := b.subscribe()

	b.unsubscribe(ch)

	if _, ok := b.subscribers[ch]; ok {
		t.Error("expected ch to be removed from subscribers")
	}

	if _, ok := <-ch; ok {
		t.Error("expected channel to be closed")
	}
}
