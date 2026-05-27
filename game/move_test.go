package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestDescribeRoom(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.describeRoom()
	output := buf.String()

	if !strings.Contains(output, "ENTRANCE") {
		t.Errorf("expected ENTRANCE in output, got: %s", output)
	}
}

func TestDescribeRoomExits(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.describeRoom()
	output := buf.String()

	if !strings.Contains(output, "north") {
		t.Errorf("expected north in output, got : %s", output)
	}
}

func TestMoveToRoom(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	locationA := game.Current
	game.Player = NewPlayer("Test", Warrior)

	reader := strings.NewReader("1\n1\n1\n1\n1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.move("north")
	locationB := game.Current

	if locationA == locationB {
		t.Errorf("expected different locations, got : %s", locationB)
	}
}

func TestMoveInvalidDirection(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.Player = NewPlayer("Test", Warrior)
	game.move("east")
	output := buf.String()

	if !strings.Contains(output, "can't go") {
		t.Errorf("expected to not be able to move in that direction, got : %s", output)
	}
}
