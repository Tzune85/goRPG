package game

import "testing"

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
