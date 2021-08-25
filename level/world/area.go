package world

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/level/world/block"
)

const (
	AREA_WIDTH  int64 = 16
	AREA_HEIGHT int64 = 60
)

type Area struct {
	core.Node

	d [AREA_HEIGHT][AREA_WIDTH][AREA_WIDTH]block.IBlock

	Pos math32.Vector3
}

func NewArea(pos math32.Vector3) *Area {
	return &Area{
		Pos:  pos,
		Node: *core.NewNode(),
	}
}

func (a *Area) Load(wg IWorldGenerator) {
	posX := int64(a.Pos.X)
	posZ := int64(a.Pos.Z)

	for y := int64(0); y < AREA_HEIGHT; y++ {
		for x := int64(0); x < AREA_WIDTH; x++ {
			for z := int64(0); z < AREA_WIDTH; z++ {
				pos := math32.NewVector3(float32(posX+x), float32(y), float32(posZ+z))
				id := wg.GetBlock(float64(pos.X), float64(pos.Y), float64(pos.Z))
				b := block.NewBlock(id, *pos)
				b.AddTo(a)

				a.d[y][x][z] = b
			}
		}
	}

	for y := int64(0); y < AREA_HEIGHT; y++ {
		for x := int64(0); x < AREA_WIDTH; x++ {
			for z := int64(0); z < AREA_WIDTH; z++ {
				if y == AREA_HEIGHT-1 {
					a.d[y][x][z].SetVisible(true)
					continue
				}

				if (y > 0 && a.d[y-1][x][z].Transparent()) || a.d[y+1][x][z].Transparent() ||
					(x > 0 && a.d[y][x-1][z].Transparent()) || (x < AREA_WIDTH-1 && a.d[y][x+1][z].Transparent()) ||
					(z > 0 && a.d[y][x][z-1].Transparent()) || (z < AREA_WIDTH-1 && a.d[y][x][z+1].Transparent()) {
					a.d[y][x][z].SetVisible(true)
					continue
				}

				a.d[y][x][z].SetVisible(false)
			}
		}
	}
}

func (a *Area) Update() {
}
