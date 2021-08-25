package player

import "github.com/g3n/engine/math32"

type Player struct {
	loc *math32.Vector3
}

func NewPlayer() *Player {
	p := new(Player)
	p.loc = &math32.Vector3{0, 0, 0}

	return p
}

func (p *Player) GetPosition() *math32.Vector3 {
	return p.loc
}
