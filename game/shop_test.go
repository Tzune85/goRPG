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

	RunShop(p, scriptedInput("1", "1", "3", "3"), &buf)
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

	RunShop(p, scriptedInput("1", "1", "3", "3"), &buf)
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

	RunShop(p, scriptedInput("2", "1", "3", "3"), &buf)
	currentGold := p.Gold

	if currentGold != 5 {
		t.Errorf("expected 5 gold, got : %d", currentGold)
	}
	if len(p.Items) != 0 {
		t.Error("expected NO Potion in inventory")
	}
}
