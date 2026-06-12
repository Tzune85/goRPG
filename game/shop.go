package game

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

type Item struct {
	ID          int
	Name        string
	Description string
	Price       int
	Modifier    func(p *Player, sign int)
}

var inventory = map[int]Item{
	1: {ID: 1, Name: "Health Potion", Description: "Heal 30 HP", Price: 10},
	2: {ID: 2, Name: "Vorpal Sword +5", Description: "Add 5 to Attack", Price: 30, Modifier: func(p *Player, sign int) {
		p.Attack += 5 * sign
	}},
	3: {ID: 3, Name: "Ring of Vitality", Description: "Add 20 HP", Price: 25, Modifier: func(p *Player, sign int) {
		p.MaxHP += 20 * sign
		p.HP += 20 * sign
	}},
	4: {ID: 4, Name: "Winged Shoes", Description: "You can always Run", Price: 20, Modifier: func(p *Player, sign int) {
		p.hasShoes = sign > 0
	}},
}

func RunShop(p *Player, readLine func(string) (string, bool), out io.Writer) {
	for {
		fmt.Fprintln(out, "\nActions: [1] Buy   [2] Sell   [3] Exit")

		options, ok := readLine("> ")
		if !ok {
			return
		}

		switch options {
		case "1", "buy":
			BuyShop(p, readLine, out)
		case "2", "sell":
			SellShop(p, readLine, out)
		case "3", "exit":
			fmt.Fprintln(out, "Come back soon!")
			return
		default:
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
		}
	}
}

func BuyShop(p *Player, readline func(string) (string, bool), out io.Writer) {
	fmt.Fprintln(out)
	for {
		fmt.Fprintf(out, "You have: %d gold\n", p.Gold)
		fmt.Fprintln(out)
		ids := make([]int, 0, len(inventory))
		for id := range inventory {
			ids = append(ids, id)
		}
		slices.Sort(ids)

		for _, id := range ids {
			item := inventory[id]
			fmt.Fprintf(out, "[%d] %s = %d Gold\n", id, item.Name, item.Price)
			if item.Description != "" {
				fmt.Fprintf(out, "\t(%s)\n", item.Description)
			}
		}
		fmt.Fprintln(out, "[0] Back")
		fmt.Fprintln(out)

		choice, ok := readline(">")
		if !ok {
			return
		}

		if choice == "0" || choice == "back" {
			return
		}

		id, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}

		item, found := inventory[id]
		if !found {
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}

		if p.Gold < item.Price {
			fmt.Fprintln(out, "You haven't enough gold!")
			fmt.Fprintln(out)
			continue
		}
		p.Gold -= item.Price
		p.Items = append(p.Items, item.Name)
		if item.Modifier != nil {
			item.Modifier(p, +1)
		}
		fmt.Fprintf(out, "Here your %s!\n", item.Name)
		fmt.Fprintln(out)
	}
}

func SellShop(p *Player, readline func(string) (string, bool), out io.Writer) {
	fmt.Fprintln(out)
	for {
		ids := make([]int, 0, len(inventory))
		for id := range inventory {
			ids = append(ids, id)
		}
		slices.Sort(ids)

		for _, id := range ids {
			item := inventory[id]
			fmt.Fprintf(out, "[%d] %s = %d Gold\n", id, item.Name, item.Price/2)
		}
		fmt.Fprintln(out, "[0] Back")
		fmt.Fprintln(out)

		choice, ok := readline(">")
		if !ok {
			return
		}

		if choice == "0" || choice == "back" {
			return
		}

		id, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}

		item, found := inventory[id]
		if !found {
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}

		if len(p.Items) == 0 {
			fmt.Fprintln(out, "You don't have any items to sell!")
			fmt.Fprintln(out)
			continue
		}
		if !slices.Contains(p.Items, item.Name) {
			fmt.Fprintln(out, "You don't have this item!")
			fmt.Fprintln(out)
			continue
		}

		i := slices.Index(p.Items, item.Name)
		p.Items = append(p.Items[:i], p.Items[i+1:]...)
		p.Gold += item.Price / 2
		if item.Modifier != nil {
			item.Modifier(p, -1)
		}
		fmt.Fprintf(out, "Sold! You received %d Gold.\n", item.Price/2)
		fmt.Fprintln(out)
	}
}
