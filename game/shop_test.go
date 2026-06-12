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

func TestShopOptionsUnknow(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("test"), &buf)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("expected unknown, got : %s", output)
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

func TestShopBuyUnknownString(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("1", "test"), &buf)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("expected unknown, got : %s", output)
	}
}

func TestShopBuyUnknownNumber(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("1", "99"), &buf)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("expected unknown, got : %s", output)
	}
}

func TestShopBuyVorpalSword(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 50, Attack: 5, Items: []string{}}

	RunShop(p, scriptedInput("1", "2", "0", "3"), &buf)

	if p.Attack != 10 {
		t.Errorf("expected Attack 10, got %d", p.Attack)
	}
	if len(p.Items) == 0 || p.Items[0] != "Vorpal Sword +5" {
		t.Error("expected Vorpal Sword in inventory")
	}
}

func TestShopBuyRing(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Stats: Stats{HP: 5, MaxHP: 25}, Gold: 50, Items: []string{}}

	RunShop(p, scriptedInput("1", "3", "0", "3"), &buf)

	if p.HP != 25 {
		t.Errorf("expected HP 25, got %d", p.HP)
	}
	if p.MaxHP != 45 {
		t.Errorf("expected MaxHP 45, got %d", p.MaxHP)
	}
	if len(p.Items) == 0 || p.Items[0] != "Ring of Vitality" {
		t.Error("expected Ring of Vitality in inventory")
	}
}

func TestShopBuyShoes(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 500, hasShoes: false, Items: []string{}}

	RunShop(p, scriptedInput("1", "4", "0", "3"), &buf)

	if !p.hasShoes {
		t.Errorf("expected hasShoes true, got %t", p.hasShoes)
	}
	if len(p.Items) == 0 || p.Items[0] != "Winged Shoes" {
		t.Error("expected Winged Shoes in inventory")
	}
}

//////////////////////////////////////////////////////////////
///////////////// SELL

func TestShopSellUnknownNumber(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("2", "99"), &buf)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("expected unknown, got : %s", output)
	}
}

func TestShopSellUnknownString(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("2", "test"), &buf)
	output := buf.String()

	if !strings.Contains(output, "Unknown") {
		t.Errorf("expected unknown, got : %s", output)
	}
}

func TestShopSellNothing(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 20, Items: []string{}}

	RunShop(p, scriptedInput("2", "1", "0", "3"), &buf)
	output := buf.String()

	if !strings.Contains(output, "You don't have any items to sell!") {
		t.Errorf("expected don't have any item, got : %s", output)
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

func TestShopSellShoes(t *testing.T) {
	var buf bytes.Buffer
	p := &Player{Gold: 0, hasShoes: true, Items: []string{"Winged Shoes"}}

	RunShop(p, scriptedInput("2", "4", "0", "3"), &buf)

	if p.hasShoes {
		t.Errorf("expected hasShoes false, got %t", p.hasShoes)
	}
	if len(p.Items) != 0 {
		t.Error("expected NO Winged Shoes in inventory")
	}
}
