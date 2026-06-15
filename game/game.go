package game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	Player  *Player
	World   map[string]*Room
	Current string
	scanner *bufio.Scanner
	running bool
	in      io.Reader
	out     io.Writer
	rng     *rand.Rand
	t       *Translator
}

var aliases = map[string]string{
	"north": "north", "n": "north", "nord": "north",
	"south": "south", "s": "south", "sud": "south",
	"east": "east", "e": "east", "est": "east",
	"west": "west", "w": "west", "ovest": "west",
}

func New() *Game {
	t, _ := newTranslator("en")
	return &Game{
		World:   BuildWorld(),
		Current: "entrance",
		scanner: bufio.NewScanner(os.Stdin),
		running: false,
		in:      os.Stdin,
		out:     os.Stdout,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		t:       t,
	}
}

func (g *Game) readLine(prompt string) (string, bool) {
	fmt.Fprint(g.out, prompt)
	ok := g.scanner.Scan()
	return strings.TrimSpace(g.scanner.Text()), ok
}

func (g *Game) Run() {
	g.setup()
	if g.Player == nil {
		return
	}

	g.running = true
	for g.running && g.Player.IsAlive() {
		input, ok := g.readLine("\n> ")
		if !ok {
			g.running = false
			return
		}
		g.handleInput(input)
	}

	if !g.Player.IsAlive() {
		fmt.Fprintln(g.out, g.t.T("game_over"))
	}
}

func (g *Game) chooseLanguage() (*Translator, bool) {
	fmt.Fprint(g.out, "\033[H\033[2J")
	fmt.Fprint(g.out, "\tSelect language / Seleziona lingua:\n\t[1] English    / [2] Italiano\n> ")
	for {
		choice, ok := g.readLine("> ")
		if !ok {
			return nil, false
		}

		switch choice {
		case "1", "english":
			t, _ := newTranslator("en")
			return t, true
		case "2", "italiano":
			t, _ := newTranslator("it")
			return t, true
		default:
			fmt.Fprintln(g.out, "Invalid selection / Selezione sbagliata")
		}

	}
}

func (g *Game) setup() {
	lang, ok := g.chooseLanguage()
	if !ok {
		return
	}
	g.t = lang
	fmt.Fprintln(g.out, g.t.T("setup_title"))

	name, ok := g.readLine(g.t.T("setup_name_prompt"))
	if !ok {
		return
	}
	if name == "" {
		name = g.t.T("setup_default_name")
	}

	class, ok := g.chooseClass()
	if !ok {
		return
	}
	g.Player = NewPlayer(name, class)

	fmt.Fprintln(g.out)
	g.describeRoom()
}

func (g *Game) chooseClass() (Class, bool) {
	fmt.Fprintln(g.out, g.t.T("setup_choose_class"))
	fmt.Fprintln(g.out, g.t.T("setup_class_warrior"))
	fmt.Fprintln(g.out, g.t.T("setup_class_mage"))
	fmt.Fprintln(g.out, g.t.T("setup_class_thief"))

	for {
		choice, ok := g.readLine("> ")
		if !ok {
			return ClassUndefinied, false
		}
		switch choice {
		case "1", "warrior", "guerriero":
			return Warrior, true
		case "2", "mage", "mago":
			return Mage, true
		case "3", "thief", "ladro":
			return Thief, true
		case "42", "god":
			return God, true
		default:
			fmt.Fprintln(g.out, g.t.T("setup_class_invalid"))
		}
	}
}

func (g *Game) handleInput(raw string) {
	parts := strings.Fields(strings.ToLower(raw))
	if len(parts) == 0 {
		return
	}

	verb := parts[0]
	var arg string
	if len(parts) > 1 {
		arg = strings.Join(parts[1:], " ")
	}

	switch verb {
	case "go", "move", "vai":
		g.move(arg)
	case "north", "n":
		g.move("north")
	case "south", "s":
		g.move("south")
	case "east", "e":
		g.move("east")
	case "west", "w", "o":
		g.move("west")
	case "look", "l", "guarda", "g":
		g.describeRoom()
	case "inventory", "inv", "i", "inventario":
		g.showInventory()
	case "potion", "cure", "pozione", "p":
		usePotion(g.Player, g.out, g.t)
	case "status", "scheda":
		fmt.Fprintln(g.out, g.Player.DetailedStatus(g.t))
	case "help", "?", "aiuto":
		g.printHelp()
	case "quit", "q":
		g.quit()
	default:
		fmt.Fprintln(g.out, g.t.T("cmd_unknown", raw))
	}
}

func (g *Game) quit() {
	fmt.Fprintln(g.out, g.t.T("cmd_quit_guard"))
	choice, ok := g.readLine("> ")
	if !ok || choice != "1" {
		return
	}
	fmt.Fprintln(g.out, g.t.T("cmd_farewell"))
	g.running = false
}

func (g *Game) move(direction string) {
	if direction == "" {
		fmt.Fprintln(g.out, g.t.T("nav_move_where"))
		return
	}

	full, ok := aliases[direction]
	if !ok {
		fmt.Fprintln(g.out, g.t.T("nav_bad_direction", direction))
		return
	}

	direction = full

	room := g.World[g.Current]
	nextID, exists := room.Exits[direction]
	if !exists {
		fmt.Fprintln(g.out, g.t.T("nav_no_exit", direction))
		return
	}

	g.Current = nextID
	next := g.World[nextID]
	g.describeRoom()

	if next.IsShop {
		RunShop(g.Player, g.readLine, g.out, g.t)
	}

	if next.EnemyName != "" && !next.Cleared {
		enemy, found := NewEnemy(next.EnemyName)
		if found {
			won, completed := RunCombat(g.Player, &enemy, g.readLine, g.out, g.rng, g.t)

			if !completed {
				return
			}

			if !g.Player.IsAlive() {
				return
			}

			if won {
				next.Cleared = true
				fmt.Fprintln(g.out, g.t.T("game_xp_gold", enemy.XP, enemy.Gold))
				g.Player.AddXP(enemy.XP, g.out, g.t)
				g.Player.Gold += enemy.Gold

				if next.IsBoss {
					g.victory()
					g.running = false
					return
				}
			}
		}
	}

	if len(next.Items) > 0 && (next.EnemyName == "" || next.Cleared) {
		for _, id := range next.Items {
			g.Player.Items = append(g.Player.Items, id)
			fmt.Fprintln(g.out, g.t.T("nav_found_item", g.t.T("item_"+strconv.Itoa(id)+"_name")))
		}
		next.Items = nil
	}
}

func (g *Game) describeRoom() {
	var exits []string
	for dir := range g.World[g.Current].Exits {
		exits = append(exits, g.t.T("dir_"+dir))
	}

	name := g.t.T("room_" + g.Current + "_name")
	desc := g.t.T("room_" + g.Current + "_desc")

	fmt.Fprintf(g.out, g.t.T("room_header"), strings.ToUpper(name))
	fmt.Fprintln(g.out, desc)
	fmt.Fprintf(g.out, g.t.T("room_exits"), strings.Join(exits, " | "))
}

func (g *Game) showInventory() {
	if len(g.Player.Items) == 0 {
		fmt.Fprintln(g.out, g.t.T("inv_empty"))
		return
	}
	fmt.Fprintln(g.out, g.t.T("inv_header"))
	for _, id := range g.Player.Items {
		fmt.Fprintln(g.out, g.t.T("inv_item", g.t.T("item_"+strconv.Itoa(id)+"_name")))
	}
}

func (g *Game) victory() {
	fmt.Fprintf(g.out, g.t.T("game_win"), g.Player.Level, g.Player.Gold)
}

func (g *Game) printHelp() {
	fmt.Fprintln(g.out, g.t.T("cmd_help"))
}
