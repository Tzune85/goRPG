package game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
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
}

func New() *Game {
	return &Game{
		World:   BuildWorld(),
		Current: "entrance",
		scanner: bufio.NewScanner(os.Stdin),
		running: false,
		in:      os.Stdin,
		out:     os.Stdout,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Game) readLine(prompt string) (string, bool) {
	fmt.Fprint(g.out, prompt)
	ok := g.scanner.Scan()
	return strings.TrimSpace(g.scanner.Text()), ok
}

func (g *Game) Run() {
	g.setup()

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
		fmt.Fprintln(g.out, "\nYou have fallen in the dungeon. GAME OVER.")
	}
}

func (g *Game) setup() {
	fmt.Fprint(g.out, "\033[H\033[2J")
	fmt.Fprintln(g.out, `
╔══════════════════════════════════════╗
║       DUNGEON OF SHADOWS  v1.0       ║
║        A Text Adventure RPG          ║
╚══════════════════════════════════════╝`)

	name, ok := g.readLine("Enter your name, adventurer: ")
	if !ok {
		return
	}
	if name == "" {
		name = "Hero"
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
	fmt.Fprintln(g.out, "\nChoose your class:")
	fmt.Fprintln(g.out, "  [1] Warrior — 100 HP, 15 ATK")
	fmt.Fprintln(g.out, "  [2] Mage    —  70 HP, 25 ATK")
	fmt.Fprintln(g.out, "  [3] Thief   —  85 HP, 18 ATK")

	for {
		choice, ok := g.readLine("> ")
		if !ok {
			return ClassUndefinied, false
		}
		switch choice {
		case "1", "warrior":
			return Warrior, true
		case "2", "mage":
			return Mage, true
		case "3", "thief":
			return Thief, true
		case "god", "42":
			return God, true
		default:
			fmt.Fprintln(g.out, "Please choose 1, 2, or 3.")
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
	case "go", "move":
		g.move(arg)
	case "north", "n":
		g.move("north")
	case "south", "s":
		g.move("south")
	case "east", "e":
		g.move("east")
	case "west", "w":
		g.move("west")
	case "look", "l":
		g.describeRoom()
	case "inventory", "inv", "i":
		g.showInventory()
	case "status":
		fmt.Fprintln(g.out, g.Player.DetailedStatus())
	case "help", "?":
		g.printHelp()
	case "quit", "q":
		fmt.Fprintln(g.out, "Farewell, adventurer...")
		g.running = false
	default:
		fmt.Fprintf(g.out, "Unknown command '%s'. Type 'help' for commands.\n", raw)
	}
}

func (g *Game) move(direction string) {
	if direction == "" {
		fmt.Fprintln(g.out, "Move where? (north, south, east, west)")
		return
	}

	room := g.World[g.Current]
	nextID, exists := room.Exits[direction]
	if !exists {
		fmt.Fprintf(g.out, "You can't go %s from here.\n", direction)
		return
	}

	g.Current = nextID
	next := g.World[nextID]
	g.describeRoom()

	if next.EnemyName != "" && !next.Cleared {
		enemy, found := NewEnemy(next.EnemyName)
		if found {
			won, completed := RunCombat(g.Player, &enemy, g.readLine, g.out, g.rng)

			if !completed {
				return
			}

			if !g.Player.IsAlive() {
				return
			}

			if won {
				next.Cleared = true
				fmt.Fprintf(g.out, "You gain %d XP and %d Gold!\n", enemy.XP, enemy.Gold)
				g.Player.AddXP(enemy.XP)
				g.Player.Gold += enemy.Gold

				if next.IsBoss {
					g.victory()
					g.running = false
					return
				}
			}
		}
	}

	if len(next.Items) > 0 {
		for _, item := range next.Items {
			g.Player.Items = append(g.Player.Items, item)
			fmt.Fprintf(g.out, "You found: %s!\n", item)
		}
		next.Items = nil
	}

}

func (g *Game) describeRoom() {
	room := g.World[g.Current]

	var exits []string
	for dir := range room.Exits {
		exits = append(exits, dir)
	}

	fmt.Fprintf(g.out, "\n=== %s ===\n", strings.ToUpper(room.Name))
	fmt.Fprintln(g.out, room.Description)
	fmt.Fprintf(g.out, "\nExits: %s\n", strings.Join(exits, " | "))
}

func (g *Game) showInventory() {
	if len(g.Player.Items) == 0 {
		fmt.Fprintln(g.out, "Your inventory is empty.")
		return
	}
	fmt.Fprintln(g.out, "Inventory:")
	for _, item := range g.Player.Items {
		fmt.Fprintf(g.out, "  - %s\n", item)
	}
}

func (g *Game) victory() {
	fmt.Fprintf(g.out, "\nYOU WIN! Final score: Level %d | Gold %d\n",
		g.Player.Level, g.Player.Gold)
}

func (g *Game) printHelp() {
	fmt.Fprintln(g.out, `
Commands:
  go [north|south|east|west]   Move
  n / s / e / w                Shortcut directions
  look  (l)                    Describe current room
  inventory  (i)               Show items
  status                       Show character sheet
  help  (?)                    Show this help
  quit  (q)                    Exit`)
}
