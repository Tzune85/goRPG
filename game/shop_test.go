package game

import (
	"bytes"
	"strings"
	"testing"
)

func scriptedInput(inputs ...string) func(string) (string, bool) {
	i := 0
	return func(_ string) (string, bool) {
		if i >= len(inputs) {
			return "", false
		}
		s := inputs[i]
		i++
		return s, true
	}
}

func TestShopGreet(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	game.Player = NewPlayer("Test", God)
	game.Current = "armory"
	game.move("west")
	output := buf.String()

	if !strings.Contains(output, "Coin talks, adventurer") {
		t.Errorf("expected greetings, got %s", output)

	}
}

func TestHasShop(t *testing.T) {
	world := BuildWorld()

	room, _ := world["shop"]

	if !room.IsShop {
		t.Error("expected to find a shop")
	}

}

func TestShopBuyPotion(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("1", "1", "0", "3"), &buf)
	currentGold := p.Gold

	if currentGold != 10 {
		t.Errorf("expected 10 gold, got : %d", currentGold)
	}
	if len(p.Items) == 0 || p.Items[0] != "Health Potion" {
		t.Error("expected Health Potion in inventory")
	}
}

func TestShopBuyPotionNoGold(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 0, Items: []string{}}

	RunShop(p, scriptedInput("1", "1", "0", "3"), &buf)
	output := buf.String()

	if !strings.Contains(output, "enough gold") {
		t.Errorf("expected not gold, got : %s", output)
	}
	if len(p.Items) != 0 {
		t.Error("expected NO Potion in inventory")
	}
}

func TestShopSellPotion(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 0, Items: []string{"Health Potion"}}

	RunShop(p, scriptedInput("2", "1", "0", "3"), &buf)
	currentGold := p.Gold

	if currentGold != 5 {
		t.Errorf("expected 5 gold, got : %d", currentGold)
	}
	if len(p.Items) != 0 {
		t.Error("expected NO Potion in inventory")
	}
}

func TestShopSellItemNotOwned(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 0, Items: []string{"Health Potion"}}

	// player has a potion but tries to sell item ID 2 (Vorpal Sword)
	RunShop(p, scriptedInput("2", "2", "0", "3"), &buf)
	output := buf.String()

	if !strings.Contains(output, "don't have") {
		t.Errorf("expected 'don't have' message, got: %s", output)
	}
	if p.Gold != 0 {
		t.Errorf("expected 0 gold, got: %d", p.Gold)
	}
	if len(p.Items) != 1 {
		t.Errorf("expected 1 item still in inventory, got: %d", len(p.Items))
	}
}

func TestShopSellRemovesCorrectItem(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 0, Items: []string{"Health Potion", "Vorpal Sword +5"}}

	// sell the sword (ID 2), potion should remain
	RunShop(p, scriptedInput("2", "2", "0", "3"), &buf)

	if len(p.Items) != 1 {
		t.Errorf("expected 1 item remaining, got: %d", len(p.Items))
	}
	if p.Items[0] != "Health Potion" {
		t.Errorf("expected Health Potion to remain, got: %s", p.Items[0])
	}
}
