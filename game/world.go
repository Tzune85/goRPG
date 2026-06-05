package game

import "math/rand"

type Room struct {
	Name        string
	Description string
	Exits       map[string]string
	EnemyName   string
	Items       []string
	Cleared     bool
	IsBoss      bool
	IsShop      bool
}

func BuildWorld() map[string]*Room {
	rng := rand.New(rand.NewSource(0))
	return map[string]*Room{
		"entrance": {
			Name:        "Dungeon Entrance",
			Description: "You stand at the threshold of darkness. A rusted iron gate hangs ajar, its hinges weeping with age. Cold air seeps from the corridor ahead, carrying the faint stench of rot and old blood. Whatever lies within... it is waiting.",
			Exits:       map[string]string{"north": "corridor"},
		},
		"corridor": {
			Name:        "Dark Corridor",
			Description: "Torches long extinguished line the walls, their sconces caked with black soot. Your footsteps echo unnervingly — yet something else echoes back, just a half-beat too late. The shadows here are not empty.",
			Exits:       map[string]string{"south": "entrance", "north": "armory", "east": "crypt"},
			EnemyName:   RandomEnemy(1, rng).Name,
		},
		"armory": {
			Name:        "Abandoned Armory",
			Description: "Weapon racks lie toppled and ransacked, their blades long since claimed by rust or fleeing soldiers. A tattered banner bearing a forgotten crest sags from the ceiling. Someone left in a great hurry — and did not return.",
			Exits:       map[string]string{"south": "corridor", "west": "shop"},
			EnemyName:   RandomEnemy(1, rng).Name,
			Items:       []string{"Health Potion"},
		},
		"crypt": {
			Name:        "Ancient Crypt",
			Description: "Stone sarcophagi line the walls like sleeping sentinels, their carved faces worn to featureless ovals by centuries of damp. Latin inscriptions warn of things best left undisturbed. One lid has been pushed aside — from the inside.",
			Exits:       map[string]string{"west": "corridor", "south": "altar"},
			EnemyName:   RandomEnemy(2, rng).Name,
			Items:       []string{"Health Potion"},
		},
		"altar": {
			Name:        "Ritual Altar",
			Description: "A black stone altar dominates the chamber, still stained with sacrifices older than memory. Candles of dark wax burn without flame, casting a cold violet light. The air itself feels wrong here — thick, watchful, hungry.",
			Exits:       map[string]string{"north": "crypt", "east": "boss_chamber"},
			EnemyName:   RandomEnemy(2, rng).Name,
			Items:       []string{"Health Potion"},
		},
		"boss_chamber": {
			Name:        "Dragon's Lair",
			Description: "The ceiling disappears into darkness above a cathedral of bones. Treasure glints beneath centuries of ash. Then — a sound like a furnace door opening — and two amber eyes, each the size of a shield, slide open in the dark.",
			Exits:       map[string]string{"west": "altar"},
			EnemyName:   "Ancient Dragon",
			IsBoss:      true,
		},
		"shop": {
			Name:        "Dwarf's Shop",
			Description: "Wedged improbably between crumbling dungeon walls is a cluttered merchant stall, lit by a cheerful lantern. A stout dwarf with a braided beard eyes you from behind a counter of peculiar wares. 'Coin talks, adventurer. Everything else walks.'",
			Exits:       map[string]string{"east": "armory"},
			IsShop:      true,
		},
	}
}
