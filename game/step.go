package game

import (
	"bytes"
	"io"
	"sync"
	"time"
)

// Session wraps a Game so it can be driven one input at a time
// instead of running a blocking loop.
type Session struct {
	game   *Game
	inPipe io.WriteCloser // we write commands here → game reads them
	outBuf *safeBuffer    // game writes output here → we read it
}

// safeBuffer is a bytes.Buffer protected by a mutex (safe for concurrent use).
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) Flush() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.buf.String()
	b.buf.Reset()
	return s
}

// NewSession creates a game session driven by Step() calls.
func NewSession() *Session {
	out := &safeBuffer{}

	pr, pw := io.Pipe() // pr = game reads, pw = we write
	g := NewWithIO(pr, out)
	g.SkipArt = true

	// Run the game in the background — it will block waiting for input
	go g.Run()

	return &Session{
		game:   g,
		inPipe: pw,
		outBuf: out,
	}
}

// Step sends one command to the game and returns whatever it printed back.
// done is true when the game has ended (player died or won).
func (s *Session) Step(input string) (output string, done bool) {
	// Write the command into the pipe — the game's scanner.Scan() unblocks
	_, err := io.WriteString(s.inPipe, input+"\n")
	if err != nil {
		return "", true
	}

	// Wait a moment for the game goroutine to process and write its response
	time.Sleep(80 * time.Millisecond)

	return s.outBuf.Flush(), !s.game.running
}
