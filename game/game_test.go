package game

import (
	"bytes"
	"strings"
	"testing"
)

func TestShowInventory(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.Player = NewPlayer("TestInventory", Mage)
	game.Player.Items = []int{1}

	game.showInventory()
	output := buf.String()

	if !strings.Contains(output, "Potion") {
		t.Errorf("expected to have potion, got : %s", output)
	}
}

func TestShowInventoryEmpty(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.Player = NewPlayer("TestInventory", Mage)
	game.Player.Items = []int{}

	game.showInventory()
	output := buf.String()

	if !strings.Contains(output, "Your inventory is empty.") {
		t.Errorf("expected to be empty, got : %s", output)
	}
}

func TestShowHelp(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf

	game.printHelp()
	output := buf.String()

	if !strings.Contains(output, "Commands:") {
		t.Errorf("expected to have commands, got : %s", output)
	}
}

func TestVictory(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf

	p := NewPlayer("TestVictory", God)

	game.Player = p
	game.victory()
	output := buf.String()

	if !strings.Contains(output, "YOU WIN") {
		t.Errorf("expected win, got : %s", output)
	}
}
