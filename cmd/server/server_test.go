package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux, err := newMux()
	if err != nil {
		t.Fatalf("newMux: %v", err)
	}
	return httptest.NewServer(mux)
}

func TestHTTPRoutes(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	routes := []string{"/", "/play", "/scores"}
	for _, path := range routes {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s: want 200, got %d", path, resp.StatusCode)
		}
		ct := resp.Header.Get("Content-Type")
		if ct != "text/html; charset=utf-8" {
			t.Errorf("GET %s: Content-Type want text/html, got %q", path, ct)
		}
	}
}

func TestScoresAPIEmpty(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/scores")
	if err != nil {
		t.Fatalf("GET /api/scores: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/scores: want 200, got %d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type want application/json, got %q", ct)
	}

	var data []Score
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
}

func TestScoresAPIRecordsWin(t *testing.T) {
	// Reset global scores for this test.
	scoresMu.Lock()
	scores = nil
	scoresMu.Unlock()

	addScore(Score{Name: "Hero", Class: "Warrior", Level: 3, Gold: 50})

	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/scores")
	if err != nil {
		t.Fatalf("GET /api/scores: %v", err)
	}
	defer resp.Body.Close()

	var data []Score
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("want 1 score, got %d", len(data))
	}
	if data[0].Name != "Hero" || data[0].Score != 150 {
		t.Errorf("unexpected score entry: %+v", data[0])
	}
}
