package game

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"
)

var (
	alwaysAttack = func(prompt string) (string, bool) { return "1", true }
	mockPotion   = func(prompt string) (string, bool) { return "2", true }
	mockRun      = func(prompt string) (string, bool) { return "3", true }
	testTrans, _ = newTranslator("en")
)

func TestPlayerKillsEnemy(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)
	e := Enemy{Name: "TestDummy", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	rng := rand.New(rand.NewSource(42))

	result, _ := RunCombat(p, &e, alwaysAttack, io.Discard, rng, testTrans)

	if !result {
		t.Error("expected player to win against 1HP enemy")
	}
}

func TestPlayerKillsEnemyNegativeAttack(t *testing.T) {
	p := Player{Name: "TestDummy", Stats: Stats{HP: 100, MaxHP: 100}, Attack: -20}
	e := Enemy{Name: "TestDummy", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}

	result, _ := RunCombat(&p, &e, alwaysAttack, io.Discard, rand.New(rand.NewSource(0)), testTrans)

	if !result {
		t.Error("expected player to win against 1HP enemy")
	}
}

func TestEnemyKillPlayer(t *testing.T) {
	p := &Player{Name: "Dummy", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	e, _ := NewEnemy("Goblin")
	rng := rand.New(rand.NewSource(42))

	result, _ := RunCombat(p, &e, alwaysAttack, io.Discard, rng, testTrans)

	if result {
		t.Error("expected enemy to win against 1HP player")
	}
}

func TestEnemyKillPlayerNegative(t *testing.T) {
	p := &Player{Name: "Dummy", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	e := Enemy{Name: "TestDummy", Stats: Stats{HP: 100, MaxHP: 100}, Attack: -20}
	rng := rand.New(rand.NewSource(42))

	result, _ := RunCombat(p, &e, alwaysAttack, io.Discard, rng, testTrans)

	if result {
		t.Error("expected enemy to win against 1HP player")
	}
}

func TestPlayerUsePotion(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "DamagedHero", Stats: Stats{HP: 1, MaxHP: 100}, Attack: 100}
	p.Items = []int{1}

	e, _ := NewEnemy("Ancient Dragon")
	rng := rand.New(rand.NewSource(42))

	RunCombat(p, &e, mockPotion, &buf, rng, testTrans)
	output := buf.String()

	if !strings.Contains(output, "You drink a potion and recover 30 HP!") {
		t.Errorf("Expected potion drank, got : %s", output)
	}
}

func TestPlayerNoPotion(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "DamagedHero", Stats: Stats{HP: 1, MaxHP: 100}, Attack: 100}

	e, _ := NewEnemy("Ancient Dragon")
	rng := rand.New(rand.NewSource(42))

	RunCombat(p, &e, mockPotion, &buf, rng, testTrans)
	output := buf.String()

	if !strings.Contains(output, "You have no potions!") {
		t.Errorf("Expected no potion, got : %s", output)
	}
}

func TestPlayerUsePotionFullHealthText(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "TestHero", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	p.Items = []int{1}

	e, _ := NewEnemy("Ancient Dragon")
	rng := rand.New(rand.NewSource(42))

	RunCombat(p, &e, mockPotion, &buf, rng, testTrans)
	output := buf.String()

	if !strings.Contains(output, "Your health is full!") {
		t.Errorf("Expected to not use potion, got : %s", output)
	}
}

func TestPlayerUsePotionFullHealth(t *testing.T) {
	p := &Player{Name: "TestHero", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	p.Items = []int{1}

	e, _ := NewEnemy("Ancient Dragon")
	rng := rand.New(rand.NewSource(42))

	RunCombat(p, &e, mockPotion, io.Discard, rng, testTrans)
	if len(p.Items) == 0 {
		t.Error("potion should not have been consumed at full HP")
	}
}

func TestPlayerCombatOptions(t *testing.T) {
	var buf bytes.Buffer

	p := &Player{Name: "TestHero", Stats: Stats{HP: 1, MaxHP: 1}, Attack: 1}
	e, _ := NewEnemy("Ancient Dragon")
	rng := rand.New(rand.NewSource(42))

	count := 0
	mixedAction := func(prompt string) (string, bool) {
		count++
		if count == 1 {
			return "test", true
		}
		return "1", true
	}

	RunCombat(p, &e, mixedAction, &buf, rng, testTrans)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("Expected Unknown action, got: %s", output)
	}
}

func TestPlayerRun(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "Arthas", Stats: Stats{HP: 9999, MaxHP: 9999}}
	e, _ := NewEnemy("Goblin")

	// seed 0: primo Intn(2) == 0 → escape immediato
	RunCombat(p, &e, mockRun, &buf, rand.New(rand.NewSource(0)), testTrans)
	output := buf.String()

	if !strings.Contains(output, "You escaped!") {
		t.Errorf("Expected escape, got: %s", output)
	}
}

func TestPlayerCantEscape(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "Arthas", Stats: Stats{HP: 9999, MaxHP: 9999}}
	e, _ := NewEnemy("Goblin")

	called := false
	run := func(_ string) (string, bool) {
		if called {
			return "", false // dopo il primo tentativo: stop
		}
		called = true
		return "3", true
	}

	// seed 1: primo Intn(2) == 1 → escape fallisce
	RunCombat(p, &e, run, &buf, rand.New(rand.NewSource(1)), testTrans)
	output := buf.String()

	if !strings.Contains(output, "You couldn't escape!") {
		t.Errorf("Expected failed escape, got: %s", output)
	}
}

func TestPlayerRunWithShoes(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Name: "Arthas", Stats: Stats{HP: 9999, MaxHP: 9999}, hasShoes: true}
	e, _ := NewEnemy("Goblin")

	called := false
	run := func(_ string) (string, bool) {
		if called {
			return "", false // dopo il primo tentativo: stop
		}
		called = true
		return "3", true
	}

	// seed 1: primo Intn(2) == 1 → escape fallisce
	RunCombat(p, &e, run, &buf, rand.New(rand.NewSource(1)), testTrans)
	output := buf.String()

	if !strings.Contains(output, "You escaped!") {
		t.Errorf("Expected escaped, got: %s", output)
	}
}
