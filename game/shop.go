package game

import (
	"fmt"
	"io"
	"slices"
)

var inventory = map[string]int{
	"potion": 10,
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
	for {
		fmt.Fprintln(out, "Items:\n [1] Health Potion = 10 Gold \n[3] Back")

		choice, ok := readline(">")
		if !ok {
			return
		}

		switch choice {
		case "1", "potion":
			if p.Gold < inventory["potion"] {
				fmt.Fprintln(out, "You haven't enough gold!")
				fmt.Fprintln(out)
				continue
			}
			p.Gold = p.Gold - inventory["potion"]
			p.Items = append(p.Items, "Health Potion")
			fmt.Fprintln(out, "Here your Health Potion!")
			fmt.Fprintln(out)
			continue
		case "3", "back":
			return
		default:
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}
	}
}

func SellShop(p *Player, readline func(string) (string, bool), out io.Writer) {
	for {
		fmt.Fprintln(out, "You need gold?\n [1] Health Potion = 5 Gold \n[3] Back")

		choice, ok := readline(">")
		if !ok {
			return
		}

		switch choice {
		case "1", "potion":
			if len(p.Items) == 0 {
				fmt.Fprintln(out, "You don't have any items to sell!")
				fmt.Fprintln(out)
				continue
			}
			if !slices.Contains(p.Items, "Health Potion") {
				fmt.Fprintln(out, "You don't have any potion to sell!")
				fmt.Fprintln(out)
				continue
			}

			continue
		case "3", "back":
			return
		default:
			fmt.Fprintln(out, "Unknown action.")
			fmt.Fprintln(out)
			continue
		}

	}
}
