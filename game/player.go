package game

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
	Name     string
	Class    Class
	Attack   int
	Gold     int
	Level    int
	XP       int
	Items    []int
	hasShoes bool
	//TODO
	//isPoisoned bool
}

func NewPlayer(name string, class Class) *Player {
	p := &Player{
		Name:  name,
		Class: class,
		Level: 1,
		Gold:  10,
		Items: []int{},
	}

	switch class {
	case Warrior:
		p.Stats = Stats{HP: 100, MaxHP: 100}
		p.Attack = 15
		p.Items = []int{1, 1}
		p.Gold = 1
	case Mage:
		p.Stats = Stats{HP: 70, MaxHP: 70}
		p.Attack = 25
		p.Items = []int{1}
		p.Gold = 5
	case Thief:
		p.Stats = Stats{HP: 85, MaxHP: 85}
		p.Attack = 18
		p.Items = []int{}
		p.Gold = 10
	case God:
		p.Stats = Stats{HP: 1000, MaxHP: 1000}
		p.Attack = 1000
		p.Items = []int{1, 2}
		p.Gold = 10000
	}

	return p
}

func (p *Player) AddXP(amount int, out io.Writer, t *Translator) {
	p.XP += amount
	for threshold := p.Level * 100; p.XP >= threshold; threshold = p.Level * 100 {
		p.XP -= threshold
		p.Level++
		p.MaxHP += 15
		p.HP = p.MaxHP
		p.Attack += 3
		fmt.Fprintf(out, t.T("status_level_up"), p.Level, p.MaxHP, p.Attack)
	}
}

func (p *Player) Status(t *Translator) string {
	return fmt.Sprintf(t.T("status_bar"),
		p.Name, t.T("class_"+strings.ToLower(string(p.Class))), p.HP, p.MaxHP, p.Gold, p.Level)
}

func (p *Player) DetailedStatus(t *Translator) string {
	inv := t.T("status_items_empty")
	if len(p.Items) > 0 {
		names := make([]string, len(p.Items))
		for i, id := range p.Items {
			names[i] = t.T("item_" + strconv.Itoa(id) + "_name")
		}
		inv = strings.Join(names, ", ")
	}
	return fmt.Sprintf(t.T("status_sheet"),
		p.Name, t.T("class_"+strings.ToLower(string(p.Class))),
		p.Level, p.XP, p.Level*100,
		p.HP, p.MaxHP,
		p.Attack,
		p.Gold,
		inv,
	)
}
