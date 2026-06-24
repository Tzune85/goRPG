package main

import (
	"bytes"
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
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"goRpg/game"
	"nhooyr.io/websocket"
)

//go:embed static
var staticFiles embed.FS

const skillPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>skill.md — Dungeon of Shadows</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #0d1117;
      color: #c9d1d9;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
      font-size: 16px;
      line-height: 1.7;
      padding: 2rem 1rem 4rem;
    }
    .wrapper {
      max-width: 780px;
      margin: 0 auto;
    }
    .topbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-bottom: 1rem;
      margin-bottom: 2rem;
      border-bottom: 1px solid #30363d;
      font-size: .85rem;
      color: #8b949e;
    }
    .topbar a { color: #58a6ff; text-decoration: none; }
    .topbar a:hover { text-decoration: underline; }
    h1 { font-size: 2rem; color: #f0f6fc; border-bottom: 1px solid #30363d; padding-bottom: .5rem; margin: 1.5rem 0 1rem; }
    h2 { font-size: 1.4rem; color: #f0f6fc; border-bottom: 1px solid #21262d; padding-bottom: .3rem; margin: 2rem 0 .8rem; }
    h3 { font-size: 1.1rem; color: #f0f6fc; margin: 1.5rem 0 .5rem; }
    p  { margin-bottom: 1rem; }
    a  { color: #58a6ff; }
    code {
      background: #161b22;
      border: 1px solid #30363d;
      border-radius: 4px;
      padding: .15em .4em;
      font-family: 'SFMono-Regular', Consolas, monospace;
      font-size: .9em;
      color: #79c0ff;
    }
    pre {
      background: #161b22;
      border: 1px solid #30363d;
      border-radius: 6px;
      padding: 1rem 1.2rem;
      overflow-x: auto;
      margin-bottom: 1.2rem;
    }
    pre code {
      background: none;
      border: none;
      padding: 0;
      color: #e6edf3;
      font-size: .9rem;
      line-height: 1.6;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 1.2rem;
      font-size: .95rem;
    }
    th {
      background: #161b22;
      color: #f0f6fc;
      padding: .5rem 1rem;
      text-align: left;
      border: 1px solid #30363d;
    }
    td {
      padding: .5rem 1rem;
      border: 1px solid #21262d;
    }
    tr:nth-child(even) td { background: #161b22; }
    ul, ol { padding-left: 1.5rem; margin-bottom: 1rem; }
    li { margin-bottom: .3rem; }
    hr { border: none; border-top: 1px solid #30363d; margin: 2rem 0; }
    strong { color: #f0f6fc; }
  </style>
</head>
<body>
<div class="wrapper">
  <div class="topbar">
    <a href="/">⌂ Dungeon of Shadows</a>
    <span>skill.md</span>
  </div>
  {{BODY}}
</div>
</body>
</html>`

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
	scoresMu    sync.Mutex
	scores      []Score
	store       = NewSessionStore()
	broadcaster = NewBroadcaster()
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
	g.SkipArt = r.URL.Query().Get("noart") == "1"
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
	mux.HandleFunc("/api/register", store.handleRegister)
	mux.HandleFunc("/api/action", store.handleAction)
	mux.HandleFunc("/api/state", store.handleState)
	mux.HandleFunc("/api/events", broadcaster.handleEvents)
	mux.HandleFunc("/watch", serveHTML("watch.html"))
	mux.HandleFunc("/play", serveHTML("play.html"))
	mux.HandleFunc("/scores", serveHTML("scores.html"))
	mux.HandleFunc("/ai", serveHTML("ai.html"))
	mux.HandleFunc("/skill.md", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFS.Open("skill.md")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		src, _ := io.ReadAll(f)

		var body bytes.Buffer
		goldmark.Convert(src, &body)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, strings.Replace(skillPageHTML, "{{BODY}}", body.String(), 1))
	})
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
