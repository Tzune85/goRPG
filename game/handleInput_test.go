package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestHandleInputUnknow(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf

	game.handleInput("test")
	output := buf.String()

	if !strings.Contains(output, "Unknown command") {
		t.Errorf("expected helper, got : %s", output)
	}
}

func TestHandleInputQuit(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.running = true
	game.handleInput("quit")
	output := buf.String()

	if !strings.Contains(output, "Farewell") {
		t.Errorf("expected farewell, got : %s", output)
	}
	if game.running {
		t.Error("expected running to be false after quit")
	}
}

func TestHandleInputLook(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf

	game.handleInput("look")
	output := buf.String()

	if !strings.Contains(output, "ENTRANCE") {
		t.Errorf("expected entrace, got : %s", output)
	}

}

func TestHandleInputEmpty(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf

	game.handleInput("\n")
	output := buf.String()

	if output != "" {
		t.Errorf("expected nothing, got : %s", output)
	}
}

func TestHandleInputMove(t *testing.T) {
	var alwaysAttack = "1\n1\n1\n1\n1\n"

	cases := []struct {
		name       string
		command    string
		startRoom  string
		wantRoom   string
		wantOutput string
	}{

		{"alias n", "n", "entrance", "corridor", ""},
		{"north complete", "north", "entrance", "corridor", ""},
		{"go north", "go north", "entrance", "corridor", ""},
		{"move north", "move north", "entrance", "corridor", ""},

		{"alias s", "s", "corridor", "entrance", ""},
		{"south complete", "south", "corridor", "entrance", ""},
		{"go south", "go south", "corridor", "entrance", ""},
		{"move south", "move south", "corridor", "entrance", ""},

		{"alias e", "e", "corridor", "crypt", ""},
		{"east complete", "east", "corridor", "crypt", ""},
		{"go east", "go east", "corridor", "crypt", ""},
		{"move east", "move east", "corridor", "crypt", ""},

		{"alias w", "w", "armory", "shop", ""},
		{"west complete", "west", "armory", "shop", ""},
		{"go west", "go west", "armory", "shop", ""},
		{"move west", "move west", "armory", "shop", ""},

		{"invalid direction from shop", "n", "shop", "shop", "can't go"},
		{"invalid direction from entrance", "w", "entrance", "entrance", "can't go"},
		{"invalid direction from boss", "n", "boss_chamber", "boss_chamber", "can't go"},

		{"empty go", "go", "entrance", "entrance", "Move where?"},
		{"empty move", "move", "entrance", "entrance", "Move where?"},

		{"garbage direction", "go banana", "entrance", "entrance", "can't go"},

		{"uppercase NORTH", "NORTH", "entrance", "corridor", ""},
		{"whitespace padded", "  north  ", "entrance", "corridor", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			game := New()
			game.out = &buf
			p := NewPlayer("Lemmy", God)
			game.Player = p
			reader := strings.NewReader(alwaysAttack)
			game.in = reader
			game.scanner = bufio.NewScanner(reader)
			game.Current = c.startRoom

			game.handleInput(c.command)
			endRoom := game.Current
			output := buf.String()

			if !strings.Contains(output, c.wantOutput) {
				t.Errorf("input %q: expected output to contain %q, got %q", c.command, c.wantOutput, output)
			}

			if endRoom != c.wantRoom {
				t.Errorf("input %q: expected %v, got %v", c.command, c.wantRoom, endRoom)
			}
		})

	}
}
