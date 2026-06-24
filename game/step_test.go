package game

import (
	"strings"
	"testing"
)

func TestSafeBufferWriteAndFlush(t *testing.T) {
	b := &safeBuffer{}

	b.Write([]byte("hello"))
	b.Write([]byte(" world"))

	got := b.Flush()
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}

	// Flush should reset — second call must return empty
	if second := b.Flush(); second != "" {
		t.Errorf("expected empty after flush, got %q", second)
	}
}

func TestNewSession(t *testing.T) {
	s := NewSession()
	if s == nil {
		t.Fatal("NewSession returned nil")
	}
	if s.game == nil {
		t.Error("expected game to be set")
	}
	if s.inPipe == nil {
		t.Error("expected inPipe to be set")
	}
	if s.outBuf == nil {
		t.Error("expected outBuf to be set")
	}
}

func TestStepSetup(t *testing.T) {
	s := NewSession()

	// game.running is false during setup, so done=true is expected for the first two steps
	out, _ := s.Step("1") // English
	if !strings.Contains(out, "adventurer") {
		t.Errorf("expected name prompt after language, got %q", out)
	}

	out, _ = s.Step("Hero") // name
	if !strings.Contains(out, "class") || !strings.Contains(strings.ToLower(out), "warrior") {
		t.Errorf("expected class selection prompt, got %q", out)
	}

	// After class selection setup completes, running=true — done must be false now
	out, done := s.Step("1") // Warrior
	if done {
		t.Fatal("game should not be done after setup completes")
	}
	if !strings.Contains(out, "DUNGEON ENTRANCE") {
		t.Errorf("expected first room description, got %q", out)
	}
}

func TestStepDoneOnEOF(t *testing.T) {
	s := NewSession()

	s.Step("1")     // language
	s.Step("Hero")  // name
	s.Step("1")     // class

	// Close the input pipe — simulates the agent disconnecting
	s.inPipe.Close()

	_, done := s.Step("look")
	if !done {
		t.Error("expected done=true after pipe closed")
	}
}
