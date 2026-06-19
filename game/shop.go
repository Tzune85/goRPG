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
	1: {ID: 1, Name: "Health Potion", Description: "item_1_desc", Price: 10},
	2: {ID: 2, Name: "Vorpal Sword +5", Description: "item_2_desc", Price: 30, Modifier: func(p *Player, sign int) {
		p.Attack += 5 * sign
	}},
	3: {ID: 3, Name: "Ring of Vitality", Description: "item_3_desc", Price: 25, Modifier: func(p *Player, sign int) {
		p.MaxHP += 20 * sign
		p.HP += 20 * sign
	}},
	4: {ID: 4, Name: "Winged Shoes", Description: "item_4_desc", Price: 20, Modifier: func(p *Player, sign int) {
		p.hasShoes = sign > 0
	}},
	//TODO
	//5: {ID: 5, Name: "Antidote", Description: "item_5_desc", Price: 15, Modifier: func(p *Player, sign int) {
	//	p.isPoisoned = sign > 0
	//}},
}

func RunShop(p *Player, readLine func(string) (string, bool), out io.Writer, t *Translator) {
	for {
		fmt.Fprintln(out, t.T("shop_actions"))

		options, ok := readLine("> ")
		if !ok {
			return
		}

		switch options {
		case "1", "buy":
			BuyShop(p, readLine, out, t)
		case "2", "sell":
			SellShop(p, readLine, out, t)
		case "3", "exit":
			fmt.Fprintln(out, t.T("shop_exit"))
			return
		default:
			fmt.Fprintln(out, t.T("shop_unknown_action"))
			fmt.Fprintln(out)
		}
	}
}

func BuyShop(p *Player, readline func(string) (string, bool), out io.Writer, t *Translator) {
	fmt.Fprintln(out)
	for {
		fmt.Fprintln(out, t.T("shop_your_gold", p.Gold))
		fmt.Fprintln(out)
		ids := make([]int, 0, len(inventory))
		for id := range inventory {
			ids = append(ids, id)
		}
		slices.Sort(ids)

		for _, id := range ids {
			item := inventory[id]
			fmt.Fprintln(out, t.T("shop_item_row", id, t.T("item_"+strconv.Itoa(id)+"_name"), item.Price))
			if item.Description != "" {
				fmt.Fprintln(out, t.T("shop_item_desc", t.T(item.Description)))
			}
		}
		fmt.Fprintln(out, t.T("shop_back"))
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
			fmt.Fprintln(out, t.T("shop_unknown_action"))
			fmt.Fprintln(out)
			continue
		}

		item, found := inventory[id]
		if !found {
			fmt.Fprintln(out, t.T("shop_unknown_action"))
			fmt.Fprintln(out)
			continue
		}

		if p.Gold < item.Price {
			fmt.Fprintln(out, t.T("shop_not_enough_gold"))
			fmt.Fprintln(out)
			continue
		}
		p.Gold -= item.Price
		p.Items = append(p.Items, item.ID)
		if item.Modifier != nil {
			item.Modifier(p, +1)
		}
		fmt.Fprintln(out, t.T("shop_buy_confirm", t.T("item_"+strconv.Itoa(item.ID)+"_name")))
		fmt.Fprintln(out)
	}
}

func SellShop(p *Player, readline func(string) (string, bool), out io.Writer, t *Translator) {
	fmt.Fprintln(out)
	for {
		ids := make([]int, 0, len(inventory))
		for id := range inventory {
			ids = append(ids, id)
		}
		slices.Sort(ids)

		for _, id := range ids {
			item := inventory[id]
			fmt.Fprintln(out, t.T("shop_item_row", id, t.T("item_"+strconv.Itoa(id)+"_name"), item.Price/2))
		}
		fmt.Fprintln(out, t.T("shop_back"))
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
			fmt.Fprintln(out, t.T("shop_unknown_action"))
			fmt.Fprintln(out)
			continue
		}

		item, found := inventory[id]
		if !found {
			fmt.Fprintln(out, t.T("shop_unknown_action"))
			fmt.Fprintln(out)
			continue
		}

		if len(p.Items) == 0 {
			fmt.Fprintln(out, t.T("shop_no_items_to_sell"))
			fmt.Fprintln(out)
			continue
		}
		if !slices.Contains(p.Items, item.ID) {
			fmt.Fprintln(out, t.T("shop_dont_have_item"))
			fmt.Fprintln(out)
			continue
		}

		i := slices.Index(p.Items, item.ID)
		p.Items = append(p.Items[:i], p.Items[i+1:]...)
		p.Gold += item.Price / 2
		if item.Modifier != nil {
			item.Modifier(p, -1)
		}
		fmt.Fprintln(out, t.T("shop_sold", item.Price/2))
		fmt.Fprintln(out)
	}
}
