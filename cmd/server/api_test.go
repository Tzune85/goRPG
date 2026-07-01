package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusCreated, map[string]string{"key": "abc"})

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if result["key"] != "abc" {
		t.Errorf("expected key=abc, got %s", result["key"])
	}
}

func TestAgentFromRequest(t *testing.T) {
	store := NewSessionStore()
	t.Run("missing header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/state", nil)

		result := store.agentFromRequest(w, r)

		if result != nil {
			t.Errorf("expected nil")
		}
		if w.Code != 401 {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
	t.Run("invalid key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/state", nil)
		r.Header.Set("X-API-Key", "chiave-inventata")

		result := store.agentFromRequest(w, r)

		if result != nil {
			t.Errorf("expected nil")
		}
		if w.Code != 401 {
			t.Errorf("expected 401, got %d", w.Code)
		}

	})
	t.Run("valid key", func(t *testing.T) {
		key := store.Create("TestBot")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/state", nil)
		r.Header.Set("X-API-Key", key)

		result := store.agentFromRequest(w, r)

		if result == nil {
			t.Fatalf("expected to be not nil")
		}
		if result.Name != "TestBot" {
			t.Errorf("expected TestBot, got %s", result.Name)
		}
	})
}

func TestHandleState(t *testing.T) {
	store := NewSessionStore()
	t.Run("missing header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/state", nil)

		store.handleState(w, r)

		if w.Code != 401 {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
	t.Run("valid key", func(t *testing.T) {
		key := store.Create("TestBot")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/state", nil)
		r.Header.Set("X-API-Key", key)

		store.handleState(w, r)

		if w.Code != 200 {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var result map[string]any
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("response is not valid JSON: %v", err)
		}

		if _, ok := result["output"]; !ok {
			t.Error("expected output field")
		}
		if _, ok := result["done"]; !ok {
			t.Error("expected done field")
		}
	})
}

func TestHandleAction(t *testing.T) {
	store := NewSessionStore()

	t.Run("missing header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/action", nil)

		store.handleAction(w, r)

		if w.Code != 401 {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("missing action", func(t *testing.T) {
		key := store.Create("TestBot")
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/action", nil)
		r.Header.Set("X-API-Key", key)

		store.handleAction(w, r)

		if w.Code != 400 {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
	t.Run("valid action", func(t *testing.T) {
		key := store.Create("TestBot")
		body := strings.NewReader(`{"action": "north"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/action", body)
		r.Header.Set("X-API-Key", key)

		store.handleAction(w, r)

		if w.Code != 200 {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var result map[string]any
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("response is not valid JSON: %v", err)
		}

		if _, ok := result["output"]; !ok {
			t.Error("expected output field")
		}
		if _, ok := result["done"]; !ok {
			t.Error("expected done field")
		}
	})
}

func TestHandleRegister(t *testing.T) {
	store := NewSessionStore()
	t.Run("valid name", func(t *testing.T) {
		body := strings.NewReader(`{"name": "TestBot"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/register", body)

		store.handleRegister(w, r)

		if w.Code != 201 {
			t.Errorf("expected 201, got %d", w.Code)
		}
		var result map[string]any
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("response is not valid JSON: %v", err)
		}
		// key should be present
		if _, ok := result["key"]; !ok {
			t.Error("expected key field")
		}
		// key should be populated
		key, ok := result["key"].(string)
		if !ok || key == "" {
			t.Error("expected key populated")
		}
	})
	t.Run("missing name", func(t *testing.T) {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/register", nil)

		store.handleRegister(w, r)

		if w.Code != 201 {
			t.Errorf("expected 201, got %d", w.Code)
		}
		var result map[string]any
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("response is not valid JSON: %v", err)
		}
		// key should be present
		if _, ok := result["key"]; !ok {
			t.Error("expected key field")
		}
		// key should be populated
		key, ok := result["key"].(string)
		if !ok || key == "" {
			t.Error("expected key populated")
		}
		// key should be Unknow Agent
		if store.Get(key).Name != "Unknown Agent" {
			t.Errorf("expected Unknow Agent, got %s", store.Get(key).Name)
		}
	})
}
