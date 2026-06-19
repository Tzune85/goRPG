package game

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

func TestHandleInputUnknown(t *testing.T) {
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

	if !strings.Contains(output, "Quitting") {
		t.Errorf("expected guard, got : %s", output)
	}
	if !game.running {
		t.Error("expected running to be true after guard")
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

		{"move n", "move n", "entrance", "corridor", ""},
		{"go n", "go n", "entrance", "corridor", ""},

		{"alias s", "s", "corridor", "entrance", ""},
		{"south complete", "south", "corridor", "entrance", ""},
		{"go south", "go south", "corridor", "entrance", ""},
		{"move south", "move south", "corridor", "entrance", ""},

		{"alias e", "e", "corridor", "fungal_cavern", ""},
		{"east complete", "east", "corridor", "fungal_cavern", ""},
		{"go east", "go east", "corridor", "fungal_cavern", ""},
		{"move east", "move east", "corridor", "fungal_cavern", ""},

		{"alias w", "w", "armory", "shop", ""},
		{"west complete", "west", "armory", "shop", ""},
		{"go west", "go west", "armory", "shop", ""},
		{"move west", "move west", "armory", "shop", ""},

		{"invalid direction from shop", "n", "shop", "shop", "can't go"},
		{"invalid direction from entrance", "w", "entrance", "entrance", "can't go"},
		{"invalid direction from boss", "n", "boss_chamber", "boss_chamber", "can't go"},

		{"empty go", "go", "entrance", "entrance", "Move where?"},
		{"empty move", "move", "entrance", "entrance", "Move where?"},

		{"garbage direction", "go banana", "entrance", "entrance", "not a direction"},

		{"uppercase NORTH", "NORTH", "entrance", "corridor", ""},
		{"whitespace padded", "  north  ", "entrance", "corridor", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			game := New()
			game.out = &buf
			game.rng = rand.New(rand.NewSource(42))
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

func TestHandleInputInventory(t *testing.T) {
	cases := []struct {
		name    string
		command string
		items   []int
	}{
		{"digit i", "i", []int{1}},
		{"word inventory", "inventory", []int{1}},
		{"empty inventory", "i", []int{}},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		game := New()
		game.Player = NewPlayer("test", Warrior)
		game.out = &buf
		game.running = true
		game.Player.Items = c.items

		game.handleInput(c.command)
		output := buf.String()
		if len(c.items) == 0 {
			if !strings.Contains(output, "Your inventory is empty.") {
				t.Errorf("input %q: expected empty inventory, got : %s", c.command, output)
			}
		} else {
			if !strings.Contains(output, "Inventory:") {
				t.Errorf("input %q: expected inventory, got : %s", c.command, output)
			}
		}
	}
}

func TestHandleInputStatus(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.Player = NewPlayer("test", Warrior)
	game.out = &buf
	game.running = true

	game.handleInput("status")
	output := buf.String()

	if !strings.Contains(output, "=== CHARACTER SHEET ===") {
		t.Errorf("expected stastus, got : %s", output)
	}
}

func TestHandleInputHelp(t *testing.T) {
	cases := []struct {
		name    string
		command string
	}{
		{"digit ?", "?"},
		{"word help", "help"},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		game := New()

		game.out = &buf
		game.running = true

		game.handleInput(c.command)
		output := buf.String()

		if !strings.Contains(output, "Commands:") {
			t.Errorf("input %q: expected help, got : %s", c.name, output)
		}
	}
}

func TestHandleInputPotion(t *testing.T) {
	cases := []struct {
		name    string
		command string
	}{
		{"word potion", "potion"},
		{"digit p", "p"},
		{"word pozione", "pozione"},
		{"word cure", "cure"},
	}

	for _, c := range cases {

		var buf bytes.Buffer
		game := New()
		game.Player = &Player{Name: "Test", Stats: Stats{HP: 10, MaxHP: 50}, Items: []int{1}}
		game.out = &buf
		game.running = true

		game.handleInput(c.command)
		output := buf.String()
		currentHP := game.Player.HP

		if !strings.Contains(output, "You drink a potion and recover 30 HP!") {
			t.Errorf("expected stastus, got : %s", output)
		}
		if game.Player.HP != 40 {
			t.Errorf("expected 40 HP, got : %d", currentHP)
		}
	}
}
