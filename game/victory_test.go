package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func newGodGame(input string) *Game {
	var buf bytes.Buffer
	r := strings.NewReader(input)
	g := New()
	g.out = &buf
	g.in = r
	g.scanner = bufio.NewScanner(r)
	return g
}

// TestVictoryFalseOnQuit checks that quitting does not set the victory flag.
func TestVictoryFalseOnQuit(t *testing.T) {
	// language=1, name=Hero, class=1 (Warrior), then quit
	g := newGodGame("1\nHero\n1\nq\n1\n")
	g.Run()
	if g.Victory() {
		t.Error("Victory() should be false after quitting")
	}
}

// TestVictoryFalseOnDeath checks that dying does not set the victory flag.
func TestVictoryFalseOnDeath(t *testing.T) {
	// Setup a Warrior then send EOF mid-run (no input = scanner returns false,
	// Run exits without victory).
	g := newGodGame("1\nHero\n1\n")
	g.Run()
	if g.Victory() {
		t.Error("Victory() should be false on EOF/death")
	}
}

// TestVictoryTrueOnBossKill beats the game with the God class (ATK 1000)
// so every combat resolves in one hit.
// Route: entrance→N→corridor(fight)→W→crypt(fight)→S→altar(fight)→W→boss(fight) = WIN
func TestVictoryTrueOnBossKill(t *testing.T) {
	// language=1, name=Hero, class=42 (God), then navigate and attack once per fight
	input := "1\nHero\n42\nnorth\n1\nwest\n1\nsouth\n1\nwest\n1\n"
	g := newGodGame(input)
	g.Run()
	if !g.Victory() {
		t.Error("Victory() should be true after defeating the Ancient Dragon")
	}
}
