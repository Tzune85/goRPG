package game

import (
	"fmt"
	"io"
	"math/rand"
)

func RunCombat(p *Player, e *Enemy, readLine func(string) (string, bool), out io.Writer, rng *rand.Rand, t *Translator) (bool, bool) {
	for p.IsAlive() && e.IsAlive() {
		fmt.Fprintln(out, p.Status(t))
		fmt.Fprintln(out, e.Status())
		fmt.Fprintln(out, t.T("combat_actions"))

		choice, ok := readLine("> ")
		if !ok {
			return false, false
		}

		switch choice {
		case "1", "attack":
			playerAttack(p, e, out, rng, t)
			if e.IsAlive() {
				enemyAttack(e, p, out, rng, t)
			}
			fmt.Fprintln(out)
		case "2", "potion":
			usePotion(p, out, t)
			enemyAttack(e, p, out, rng, t)
			fmt.Fprintln(out)
		case "3", "run":
			if p.hasShoes {
				fmt.Fprintln(out, t.T("combat_escaped"))
				return false, true
			}
			if rng.Intn(2) == 0 {
				fmt.Fprintln(out, t.T("combat_escaped"))
				return false, true
			}
			fmt.Fprintln(out, t.T("combat_no_escape"))
			enemyAttack(e, p, out, rng, t)
			fmt.Fprintln(out)
		default:
			fmt.Fprintln(out, t.T("combat_unknown_action"))
			fmt.Fprintln(out)
		}
	}

	return p.IsAlive(), true
}

func playerAttack(p *Player, e *Enemy, out io.Writer, rng *rand.Rand, t *Translator) {
	dmg := p.Attack + rng.Intn(6) - 2
	if dmg < 1 {
		dmg = 1
	}
	e.TakeDamage(dmg)
	fmt.Fprintln(out, t.T("combat_player_hit", e.Name, dmg))
}

func enemyAttack(e *Enemy, p *Player, out io.Writer, rng *rand.Rand, t *Translator) {
	dmg := e.Attack + rng.Intn(4) - 1
	if dmg < 1 {
		dmg = 1
	}
	p.TakeDamage(dmg)
	fmt.Fprintln(out, t.T("combat_enemy_hit", e.Name, dmg))
}

func usePotion(p *Player, out io.Writer, t *Translator) {
	if p.HP == p.MaxHP {
		fmt.Fprintln(out, t.T("combat_hp_full", p.HP, p.MaxHP))
		return
	}
	for i, item := range p.Items {
		if item == 1 {
			p.Items = append(p.Items[:i], p.Items[i+1:]...)
			p.Heal(30)
			fmt.Fprintln(out, t.T("combat_potion_used", p.HP, p.MaxHP))
			return
		}
	}
	fmt.Fprintln(out, t.T("combat_no_potion"))
}
