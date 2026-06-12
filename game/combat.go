package game

import (
	"fmt"
	"io"
	"math/rand"
)

func RunCombat(p *Player, e *Enemy, readLine func(string) (string, bool), out io.Writer, rng *rand.Rand) (bool, bool) {
	for p.IsAlive() && e.IsAlive() {
		fmt.Fprintln(out, p.Status())
		fmt.Fprintln(out, e.Status())
		fmt.Fprintln(out, "\nActions: [1] attack   [2] potion   [3] run")

		choice, ok := readLine("> ")
		if !ok {
			return false, false
		}

		switch choice {
		case "1", "attack":
			playerAttack(p, e, out, rng)
			if e.IsAlive() {
				enemyAttack(e, p, out, rng)
			}
			fmt.Fprintln(out)
		case "2", "potion":
			usePotion(p, out)
			enemyAttack(e, p, out, rng)
			fmt.Fprintln(out)
		case "3", "run":
			if p.hasShoes {
				fmt.Fprintln(out, "You escaped!")
				return false, true
			}
			if rng.Intn(2) == 0 {
				fmt.Fprintln(out, "You escaped!")
				return false, true
			}
			fmt.Fprintln(out, "You couldn't escape!")
			enemyAttack(e, p, out, rng)
			fmt.Fprintln(out)
		default:
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
		}
	}

	return p.IsAlive(), true
}

func playerAttack(p *Player, e *Enemy, out io.Writer, rng *rand.Rand) {
	dmg := p.Attack + rng.Intn(6) - 2
	if dmg < 1 {
		dmg = 1
	}
	e.TakeDamage(dmg)
	fmt.Fprintf(out, "You hit %s for %d damage!\n", e.Name, dmg)
}

func enemyAttack(e *Enemy, p *Player, out io.Writer, rng *rand.Rand) {
	dmg := e.Attack + rng.Intn(4) - 1
	if dmg < 1 {
		dmg = 1
	}
	p.TakeDamage(dmg)
	fmt.Fprintf(out, "%s hits you for %d damage!\n", e.Name, dmg)
}

func usePotion(p *Player, out io.Writer) {
	if p.HP == p.MaxHP {
		fmt.Fprintf(out, "Your health is full! (HP: %d/%d)\n", p.HP, p.MaxHP)
		return

	}
	for i, item := range p.Items {
		if item == "Health Potion" {
			p.Items = append(p.Items[:i], p.Items[i+1:]...)
			p.Heal(30)
			fmt.Fprintf(out, "You drink a potion and recover 30 HP! (HP: %d/%d)\n", p.HP, p.MaxHP)
			return
		}
	}
	fmt.Fprintln(out, "You have no potions!")
}
