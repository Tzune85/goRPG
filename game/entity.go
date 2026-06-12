package game

type Entity interface {
	IsAlive() bool
	TakeDamage(amount int)
}

type Stats struct {
	HP    int
	MaxHP int
}

func (s *Stats) TakeDamage(amount int) {
	s.HP -= amount
	if s.HP < 0 {
		s.HP = 0
	}
}

func (s *Stats) Heal(amount int) {
	s.HP += amount
	if s.HP > s.MaxHP {
		s.HP = s.MaxHP
	}
}

func (s *Stats) IsAlive() bool {
	return s.HP > 0
}
