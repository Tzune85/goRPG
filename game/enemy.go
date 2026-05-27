package game

import (
	"fmt"
	"math/rand"
)

type Enemy struct {
	Stats
	Name   string
	Attack int
	XP     int
	Gold   int
	Tier   int
}

var catalog = []Enemy{
	{Name: "Goblin", Stats: Stats{HP: 30, MaxHP: 30}, Attack: 8, XP: 20, Gold: 5, Tier: 1},
	{Name: "Skeleton", Stats: Stats{HP: 45, MaxHP: 45}, Attack: 11, XP: 35, Gold: 8, Tier: 1},
	{Name: "Zombie", Stats: Stats{HP: 35, MaxHP: 35}, Attack: 9, XP: 25, Gold: 6, Tier: 1},
	{Name: "Orc Warrior", Stats: Stats{HP: 60, MaxHP: 60}, Attack: 14, XP: 50, Gold: 12, Tier: 2},
	{Name: "Dark Knight", Stats: Stats{HP: 80, MaxHP: 80}, Attack: 18, XP: 80, Gold: 20, Tier: 2},
	{Name: "Ancient Dragon", Stats: Stats{HP: 100, MaxHP: 100}, Attack: 25, XP: 200, Gold: 100, Tier: 3},
}

func NewEnemy(name string) (Enemy, bool) {
	for _, e := range catalog {
		if e.Name == name {
			return e, true
		}
	}
	return Enemy{}, false
}

func RandomEnemy(tier int, rng *rand.Rand) Enemy {
	var pool []Enemy
	for _, e := range catalog {
		if e.Tier == tier {
			pool = append(pool, e)
		}
	}
	return pool[rng.Intn(len(pool))]
}

func (e *Enemy) Status() string {
	return fmt.Sprintf("[ %s | HP %d ]", e.Name, e.HP)
}
