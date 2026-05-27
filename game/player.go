package game

import "fmt"

type Class string

const (
	Warrior         Class = "Warrior"
	Mage            Class = "Mage"
	Thief           Class = "Thief"
	God             Class = "God"
	ClassUndefinied Class = ""
)

type Player struct {
	Stats
	Name   string
	Class  Class
	Attack int
	Gold   int
	Level  int
	XP     int
	Items  []string
}

func NewPlayer(name string, class Class) *Player {
	p := &Player{
		Name:  name,
		Class: class,
		Level: 1,
		Gold:  10,
		Items: []string{},
	}

	switch class {
	case Warrior:
		p.Stats = Stats{HP: 100, MaxHP: 100}
		p.Attack = 15
		p.Items = []string{"Health Potion", "Health Potion"}
		p.Gold = 1
	case Mage:
		p.Stats = Stats{HP: 70, MaxHP: 70}
		p.Attack = 25
		p.Items = []string{"Health Potion"}
		p.Gold = 5
	case Thief:
		p.Stats = Stats{HP: 85, MaxHP: 85}
		p.Attack = 18
		p.Items = []string{"Health Potion"}
		p.Gold = 10
	case God:
		p.Stats = Stats{HP: 1000, MaxHP: 1000}
		p.Attack = 1000
		p.Items = []string{"Health Potion", "Health Potion"}
		p.Gold = 10000
	}

	return p
}

func (p *Player) AddXP(amount int) {
	p.XP += amount
	if p.XP >= p.Level*100 {
		p.Level++
		p.XP = 0
		p.MaxHP += 15
		p.HP = p.MaxHP
		p.Attack += 3
		fmt.Printf("\n⚡ LEVEL UP! Now level %d! (HP: %d, ATK: %d)\n",
			p.Level, p.MaxHP, p.Attack)
	}
}

func (p *Player) Status() string {
	return fmt.Sprintf("[ %s the %s | HP %d/%d | Gold %d | Lv.%d ]",
		p.Name, p.Class, p.HP, p.MaxHP, p.Gold, p.Level)
}

func (p *Player) DetailedStatus() string {
	inv := "empty"
	if len(p.Items) > 0 {
		inv = fmt.Sprintf("%v", p.Items)
	}
	return fmt.Sprintf(`
=== CHARACTER SHEET ===
Name   : %s
Class  : %s
Level  : %d (XP: %d/%d)
HP     : %d / %d
Attack : %d
Gold   : %d
Items  : %s
=======================`,
		p.Name, p.Class,
		p.Level, p.XP, p.Level*100,
		p.HP, p.MaxHP,
		p.Attack,
		p.Gold,
		inv,
	)
}
