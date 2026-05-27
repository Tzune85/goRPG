package game

import "testing"

func TestEntranceExists(t *testing.T) {
	world := BuildWorld()

	room, exists := world["entrance"]

	if !exists {
		t.Error("expected entrance room to exist")
	}
	if room.Name != "Dungeon Entrance" {
		t.Errorf("expected name 'Dungeon Entrance', got %s", room.Name)
	}
}

func TestEntranceHasNorth(t *testing.T) {
	world := BuildWorld()

	room, _ := world["entrance"]

	destination, hasNorth := room.Exits["north"]

	if !hasNorth {
		t.Error("expected entrance to have north exit")
	}

	if destination != "corridor" {
		t.Errorf("expected north to lead to corridor, got %s", destination)
	}

}

func TestHasBoss(t *testing.T) {
	world := BuildWorld()

	room, _ := world["boss_chamber"]

	if !room.IsBoss {
		t.Error("expected to find a boss")
	}

}
