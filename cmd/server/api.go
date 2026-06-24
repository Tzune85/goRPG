package main

import (
	"encoding/json"
	"net/http"
)

// POST /api/register
// Body:    {"name": "MyBot"}
// Returns: {"key": "a3f1c2d4..."}
func (s *SessionStore) handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct{ Name string }
	json.NewDecoder(r.Body).Decode(&body)
	if body.Name == "" {
		body.Name = "Unknown Agent"
	}

	key := s.Create(body.Name)

	// Auto-setup: send language + name + class so the agent starts in the dungeon
	agent := s.Get(key)
	agent.Session.Step("1")       // English
	agent.Session.Step(body.Name) // character name
	output, _ := agent.Session.Step("1") // Warrior class — captures first room description

	broadcaster.Publish(AgentEvent{
		Agent:  body.Name,
		Action: "joined the dungeon",
		Output: output,
	})

	writeJSON(w, http.StatusCreated, map[string]string{"key": key})
}

// POST /api/action
// Header: X-API-Key: <key>
// Body:    {"action": "north"}
// Returns: {"output": "...", "done": false}
func (s *SessionStore) handleAction(w http.ResponseWriter, r *http.Request) {
	agent := s.agentFromRequest(w, r)
	if agent == nil {
		return
	}

	var body struct{ Action string }
	json.NewDecoder(r.Body).Decode(&body)
	if body.Action == "" {
		http.Error(w, "missing action", http.StatusBadRequest)
		return
	}

	output, done := agent.Session.Step(body.Action)

	broadcaster.Publish(AgentEvent{
		Agent:  agent.Name,
		Action: body.Action,
		Output: output,
		Done:   done,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"output": output,
		"done":   done,
	})
}

// GET /api/state
// Header: X-API-Key: <key>
// Returns current room description (same as typing "look")
func (s *SessionStore) handleState(w http.ResponseWriter, r *http.Request) {
	agent := s.agentFromRequest(w, r)
	if agent == nil {
		return
	}

	output, done := agent.Session.Step("look")

	writeJSON(w, http.StatusOK, map[string]any{
		"output": output,
		"done":   done,
	})
}

// agentFromRequest extracts and validates the API key from the X-API-Key header.
func (s *SessionStore) agentFromRequest(w http.ResponseWriter, r *http.Request) *AgentSession {
	key := r.Header.Get("X-API-Key")
	if key == "" {
		http.Error(w, "missing X-API-Key header", http.StatusUnauthorized)
		return nil
	}
	agent := s.Get(key)
	if agent == nil {
		http.Error(w, "invalid API key", http.StatusUnauthorized)
		return nil
	}
	return agent
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
