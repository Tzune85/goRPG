package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"goRpg/game"
)

type AgentSession struct {
	Session   *game.Session
	Name      string
	CreatedAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*AgentSession // key = API key
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*AgentSession),
	}
}

// Create adds a new session and returns the generated API key.
func (s *SessionStore) Create(name string) string {
	key := generateKey() // we'll write this next
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[key] = &AgentSession{
		Session:   game.NewSession(),
		Name:      name,
		CreatedAt: time.Now(),
	}
	return key
}

// Get retrieves a session by API key. Returns nil if not found.
func (s *SessionStore) Get(key string) *AgentSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[key]
}

// generateKey makes a random 16-byte hex string (looks like: "a3f1c2d4e5b6...")
func generateKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
