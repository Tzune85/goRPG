package game

import (
	"math/rand"
	"time"
)

type Room struct {
	Name        string
	Description string
	Exits       map[string]string
	EnemyName   string
	Items       []int
	Cleared     bool
	IsBoss      bool
	IsShop      bool
}

func BuildWorld() map[string]*Room {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return map[string]*Room{
		"entrance": {
			Exits: map[string]string{"north": "corridor"},
		},
		"corridor": {
			Exits:     map[string]string{"south": "entrance", "north": "armory", "east": "crypt"},
			EnemyName: RandomEnemy(1, rng).Name,
		},
		"armory": {
			Exits:     map[string]string{"south": "corridor", "west": "shop"},
			EnemyName: RandomEnemy(1, rng).Name,
			Items:     []int{1},
		},
		"crypt": {
			Exits:     map[string]string{"west": "corridor", "south": "altar"},
			EnemyName: RandomEnemy(2, rng).Name,
			Items:     []int{1},
		},
		"altar": {
			Exits:     map[string]string{"north": "crypt", "east": "boss_chamber"},
			EnemyName: RandomEnemy(2, rng).Name,
			Items:     []int{1},
		},
		"boss_chamber": {
			Exits:     map[string]string{"west": "altar"},
			EnemyName: "Ancient Dragon",
			IsBoss:    true,
		},
		"shop": {
			Exits:  map[string]string{"east": "armory"},
			IsShop: true,
		},
	}
}
