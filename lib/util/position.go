package util

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

type Pos struct {
	X int64
	Y int64
	Z int64
}

func NewPos(x, y, z int64) Pos {
	return Pos{x, y, z}
}

func NewPosByFloat32(x, y, z float32) Pos {
	return Pos{int64(x), int64(y), int64(z)}
}

func (p Pos) GetId() string {
	return fmt.Sprintf("%d-%d-%d", p.X, p.Y, p.Z)
}

func (p Pos) Add(pos Pos) Pos {
	return Pos{
		X: p.X + pos.X,
		Y: p.Y + pos.Y,
		Z: p.Z + pos.Z,
	}
}

func (p Pos) Sub(pos Pos) Pos {
	return Pos{
		X: p.X - pos.X,
		Y: p.Y - pos.Y,
		Z: p.Z - pos.Z,
	}
}

func (p Pos) AddX(x int64) Pos {
	return p.Add(Pos{X: x})
}

func (p Pos) AddY(y int64) Pos {
	return p.Add(Pos{Y: y})
}

func (p Pos) AddZ(z int64) Pos {
	return p.Add(Pos{Z: z})
}

func (p Pos) SubX(x int64) Pos {
	return p.Add(Pos{X: -x})
}

func (p Pos) SubY(y int64) Pos {
	return p.Add(Pos{Y: -y})
}

func (p Pos) SubZ(z int64) Pos {
	return p.Add(Pos{Z: -z})
}

func (p Pos) ToVec3() math32.Vector3 {
	return *math32.NewVector3(float32(p.X), float32(p.Y), float32(p.Z))
}
