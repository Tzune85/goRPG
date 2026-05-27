package game

import "testing"

func TestNewPlayer(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)

	if p.Name != "Arthas" {
		t.Errorf("expected name Arthas, got %s", p.Name)
	}
	if p.HP != 100 {
		t.Errorf("expected HP 100, got %d", p.HP)
	}
	if p.Attack != 15 {
		t.Errorf("expected Attack 15, got %d", p.Attack)
	}
}

func TestTakeDamage(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)

	p.TakeDamage(30)

	if p.HP != 70 {
		t.Errorf("expected HP 70, got %d", p.HP)
	}
}

func TestTakeDamageCannotGoBelowZero(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)

	p.TakeDamage(9999)

	if p.HP != 0 {
		t.Errorf("expected HP 0, got %d", p.HP)
	}
}

func TestHeal(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)
	p.TakeDamage(50)

	p.Heal(20)

	if p.HP != 70 {
		t.Errorf("expected HP 70, got %d", p.HP)
	}
}

func TestHealCannotExceedMaxHP(t *testing.T) {
	p := NewPlayer("Arthas", Warrior)

	p.Heal(9999)

	if p.HP != p.MaxHP {
		t.Errorf("expected HP %d, got %d", p.MaxHP, p.HP)
	}
}
