package game

import (
	"math/rand"
	"time"
)

type Room struct {
	Name           string
	Description    string
	Exits          map[string]string
	EnemyName      string
	Items          []int
	Cleared        bool
	IsBoss         bool
	IsShop         bool
	TranslationKey string // if set, reuses another room's name/desc translations
}

func BuildWorld() map[string]*Room {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return map[string]*Room{
		"entrance": {
			Exits: map[string]string{"north": "corridor"},
		},
		"corridor": {
			Exits:     map[string]string{"south": "entrance", "north": "armory", "west": "crypt", "east": "fungal_cavern"},
			EnemyName: RandomEnemy(1, rng).Name,
		},
		"armory": {
			Exits:     map[string]string{"south": "corridor", "west": "shop", "east": "spider_corridor"},
			EnemyName: RandomEnemy(1, rng).Name,
			Items:     []int{1},
		},
		"crypt": {
			Exits:     map[string]string{"east": "corridor", "south": "altar", "north": "shop"},
			EnemyName: RandomEnemy(2, rng).Name,
			Items:     []int{1},
		},
		"altar": {
			Exits:     map[string]string{"north": "crypt", "west": "boss_chamber"},
			EnemyName: RandomEnemy(2, rng).Name,
			Items:     []int{1},
		},
		"boss_chamber": {
			Exits:     map[string]string{"east": "altar"},
			EnemyName: "Ancient Dragon",
			IsBoss:    true,
		},
		"shop": {
			Exits:  map[string]string{"east": "armory", "south": "crypt"},
			IsShop: true,
		},
		"spider_corridor": {
			Exits: map[string]string{"west": "armory", "east": "spider_den"},
		},
		"spider_den": {
			Exits:     map[string]string{"west": "spider_corridor", "south": "spider_corridor2"},
			EnemyName: "Spider Lord",
			Items:     []int{1, 1},
		},
		"spider_corridor2": {
			Exits:          map[string]string{"west": "fungal_cavern", "north": "spider_den"},
			TranslationKey: "spider_corridor",
		},
		"fungal_cavern": {
			Exits:     map[string]string{"west": "corridor", "east": "spider_corridor2"},
			EnemyName: RandomEnemy(2, rng).Name,
			Items:     []int{1},
		},
	}
}
