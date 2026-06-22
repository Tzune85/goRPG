package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"

	"goRpg/game"
	"nhooyr.io/websocket"
)

//go:embed static
var staticFiles embed.FS

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

// wsReader bridges a WebSocket connection to io.Reader.
type wsReader struct {
	conn *websocket.Conn
	ctx  context.Context
	buf  []byte
}

func (r *wsReader) Read(p []byte) (int, error) {
	for len(r.buf) == 0 {
		_, msg, err := r.conn.Read(r.ctx)
		if err != nil {
			return 0, io.EOF
		}
		r.buf = msg
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

// wsWriter bridges io.Writer to WebSocket text frames.
type wsWriter struct {
	conn *websocket.Conn
	ctx  context.Context
}

func (w *wsWriter) Write(p []byte) (int, error) {
	clean := ansiRe.ReplaceAll(p, nil)
	if len(clean) == 0 {
		return len(p), nil
	}
	err := w.conn.Write(w.ctx, websocket.MessageText, clean)
	return len(p), err
}

// Score is a single leaderboard entry.
type Score struct {
	Name  string    `json:"name"`
	Class string    `json:"class"`
	Level int       `json:"level"`
	Gold  int       `json:"gold"`
	Score int       `json:"score"`
	At    time.Time `json:"at"`
}

var (
	scoresMu sync.Mutex
	scores   []Score
)

func addScore(s Score) {
	s.Score = s.Level * s.Gold
	scoresMu.Lock()
	defer scoresMu.Unlock()
	scores = append(scores, s)
	sort.Slice(scores, func(i, j int) bool { return scores[i].Score > scores[j].Score })
	if len(scores) > 50 {
		scores = scores[:50]
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("ws accept:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "bye")

	ctx := r.Context()

	// Keepalive ping every 30 s to survive proxy timeouts.
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_ = conn.Ping(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	reader := &wsReader{conn: conn, ctx: ctx}
	writer := &wsWriter{conn: conn, ctx: ctx}

	g := game.NewWithIO(reader, writer)
	g.Run()

	if g.Player != nil && g.Victory() {
		addScore(Score{
			Name:  g.Player.Name,
			Class: string(g.Player.Class),
			Level: g.Player.Level,
			Gold:  g.Player.Gold,
			At:    time.Now(),
		})
	}
}

func serveScoresAPI(w http.ResponseWriter, r *http.Request) {
	scoresMu.Lock()
	defer scoresMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(scores)
}

func newMux() (*http.ServeMux, error) {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	serveHTML := func(name string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f, err := staticFS.Open(name)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", serveWS)
	mux.HandleFunc("/api/scores", serveScoresAPI)
	mux.HandleFunc("/play", serveHTML("play.html"))
	mux.HandleFunc("/scores", serveHTML("scores.html"))
	mux.HandleFunc("/", serveHTML("index.html"))
	return mux, nil
}

func main() {
	mux, err := newMux()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Dungeon of Shadows server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
