package game

import (
	"math/rand"
	"slices"
	"testing"
	"time"
)

func TestNewEnemy(t *testing.T) {
	e, found := NewEnemy("Goblin")

	if !found {
		t.Error("expected to find Goblin")
	}
	if e.Name != "Goblin" {
		t.Errorf("expected name Goblin, got %s", e.Name)
	}
	if e.HP != 30 {
		t.Errorf("expected HP 30, got %d", e.HP)
	}
}

func TestNewEnemyNotFound(t *testing.T) {
	_, found := NewEnemy("Unicorn")

	if found {
		t.Error("expected not to find Unicorn")
	}
}

func TestEnemyTakeDamage(t *testing.T) {
	e, _ := NewEnemy("Goblin")

	e.TakeDamage(10)

	if e.HP != 20 {
		t.Errorf("expected HP 20, got %d", e.HP)
	}
}

func TestRandomEnemy(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	e := RandomEnemy(3, rng)

	if e.Name != "Ancient Dragon" {
		t.Errorf("exoected dragon, got %s", e.Name)
	}
}

func TestRandomEnemyTier1(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	e := RandomEnemy(1, rng)

	var tier1Name []string
	for _, enemy := range catalog {
		if enemy.Tier == 1 {
			tier1Name = append(tier1Name, enemy.Name)
		}
	}
	if !slices.Contains(tier1Name, e.Name) {
		t.Errorf("expected tier 1 monster, got %s", e.Name)
	}
}
