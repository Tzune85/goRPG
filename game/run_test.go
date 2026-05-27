package game

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

func TestRunQuit(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("Hero\n1\nq\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.Run()
	output := buf.String()

	if !strings.Contains(output, "Farewell") {
		t.Errorf("expected quit, got ;  %s", output)
	}

}

func TestRunEOF(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("Hero\n1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.Run()

	if game.running {
		t.Errorf("expected running=false after EOF")
	}
	if strings.Contains(buf.String(), "Farewell") {
		t.Errorf("EOF path should not print quit message")
	}
}

func TestRunPlayerDead(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.rng = rand.New(rand.NewSource(42))
	reader := strings.NewReader("Hero\n2\nn\n" + strings.Repeat("2\n", 15))
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.Run()
	output := buf.String()

	if !strings.Contains(output, "GAME OVER") {
		t.Errorf("expected GAME OVER, got : %s", output)
	}

}
