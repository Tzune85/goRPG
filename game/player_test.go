package game

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	cases := []struct {
		name     string
		class    Class
		level    int
		hp       int
		attack   int
		maxHP    int
		gold     int
		lenItems int
	}{
		{"Arthas", Warrior, 1, 100, 15, 100, 1, 2},
		{"Mordred", Mage, 1, 70, 25, 70, 5, 1},
		{"Lesto", Thief, 1, 85, 18, 85, 10, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewPlayer(c.name, c.class)

			if p.Name != c.name {
				t.Errorf("expected name %s, got %s", c.name, p.Name)
			}
			if p.Class != c.class {
				t.Errorf("expected class %q, got %s", c.class, p.Class)
			}
			if p.HP != c.hp {
				t.Errorf("expected HP %d, got %d", c.hp, p.HP)
			}
			if p.Attack != c.attack {
				t.Errorf("expected Attack %d, got %d", c.attack, p.Attack)
			}
			if p.MaxHP != c.maxHP {
				t.Errorf("expected MAXHP %d, got %d", c.maxHP, p.MaxHP)
			}
			if p.Gold != c.gold {
				t.Errorf("expected Gold %d, got %d", c.gold, p.Gold)
			}
			if len(p.Items) != c.lenItems {
				t.Errorf("expected %d items, got %d items", c.lenItems, len(p.Items))
			}
		})
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
