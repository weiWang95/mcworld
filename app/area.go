package app

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/lib/util"
)

const (
	AREA_WIDTH  int64 = 16
	AREA_HEIGHT int64 = 60
)

type Area struct {
	core.Node

	d [AREA_HEIGHT][AREA_WIDTH][AREA_WIDTH]block.IBlock

	Pos math32.Vector3

	SouthArea *Area
	NorthArea *Area
	WestArea  *Area
	EastArea  *Area
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

	a.addAxis()

	for y := int64(0); y < AREA_HEIGHT; y++ {
		for x := int64(0); x < AREA_WIDTH; x++ {
			for z := int64(0); z < AREA_WIDTH; z++ {
				pos := math32.NewVector3(float32(posX+x), float32(y), float32(posZ+z))
				id := wg.GetBlock(float64(pos.X), float64(pos.Y), float64(pos.Z))
				if id != 0 {
					b := block.NewBlock(id, *pos)
					b.AddTo(a)

					a.d[y][x][z] = b
				}
			}
		}
	}
}

func (a *Area) refreshBlocks(w *World) {
	for y := int64(0); y < AREA_HEIGHT; y++ {
		for x := int64(0); x < AREA_WIDTH; x++ {
			for z := int64(0); z < AREA_WIDTH; z++ {
				if a.d[y][x][z] == nil {
					continue
				}

				if y == AREA_HEIGHT-1 {
					a.d[y][x][z].SetVisible(true)
					continue
				}

				if (y > 0 && a.BlockTransparent(x, y-1, z)) || a.BlockTransparent(x, y+1, z) ||
					(x > 0 && a.BlockTransparent(x-1, y, z)) || (x < AREA_WIDTH-1 && a.BlockTransparent(x+1, y, z)) ||
					(z > 0 && a.BlockTransparent(x, y, z-1)) || (z < AREA_WIDTH-1 && a.BlockTransparent(x, y, z+1)) {
					a.d[y][x][z].SetVisible(true)
					continue
				}

				if (x == 0 && a.WestArea != nil && a.WestArea.BlockTransparent(AREA_WIDTH-1, y, z)) ||
					(x == AREA_WIDTH-1 && a.EastArea != nil && a.EastArea.BlockTransparent(0, y, z)) ||
					(z == 0 && a.SouthArea != nil && a.SouthArea.BlockTransparent(AREA_WIDTH-1, y, z)) ||
					(z == AREA_WIDTH-1 && a.NorthArea != nil && a.NorthArea.BlockTransparent(0, y, z)) {
					a.d[y][x][z].SetVisible(true)
					continue
				}

				a.d[y][x][z].SetVisible(false)
			}
		}
	}
}

func (a *Area) BlockTransparent(x, y, z int64) bool {
	return a.d[y][x][z] == nil || a.d[y][x][z].Transparent()
}

func (a *Area) Update() {
}

func (a *Area) GetBlock(x, y, z float32) block.IBlock {
	if y < 0 || y >= float32(AREA_HEIGHT) {
		return nil
	}

	bx := util.FloorFloat(x) - util.FloorFloat(a.Pos.X)
	bz := util.FloorFloat(z) - util.FloorFloat(a.Pos.Z)

	return a.d[int64(y)][bx][bz]
}

func (a *Area) addAxis() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		a.Pos.X, a.Pos.Y, a.Pos.Z,
		a.Pos.X+float32(AREA_WIDTH), a.Pos.Y, a.Pos.Z,
		a.Pos.X, a.Pos.Y, a.Pos.Z,
		a.Pos.X, a.Pos.Y+float32(AREA_HEIGHT), a.Pos.Z,
		a.Pos.X, a.Pos.Y, a.Pos.Z,
		a.Pos.X, a.Pos.Y, a.Pos.Z+float32(AREA_WIDTH),
	)
	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 0.0, 0.0, // red
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0, // green
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0, // blue
		0.0, 0.0, 1.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	lines := graphic.NewLines(geom, mat)
	a.Node.Add(lines)
}

func (a *Area) ClearRelations() {
	a.SouthArea = nil
	a.NorthArea = nil
	a.WestArea = nil
	a.EastArea = nil
}

func (a *Area) ReplaceBlock(pos math32.Vector3, block block.IBlock) bool {
	if pos.Y < 0 || pos.Y >= float32(AREA_HEIGHT) {
		return false
	}

	bx := util.FloorFloat(pos.X) - util.FloorFloat(a.Pos.X)
	bz := util.FloorFloat(pos.Z) - util.FloorFloat(a.Pos.Z)

	if bx < 0 || bx >= AREA_WIDTH || bz < 0 || bz >= AREA_WIDTH {
		return false
	}

	b := a.d[int64(pos.Y)][bx][bz]
	if b != nil {
		b.RemoveFrom(a)
	}

	a.d[int64(pos.Y)][bx][bz] = block
	return true
}

func GetAreaPosByPos(pos math32.Vector3) math32.Vector3 {
	x := pos.X
	if x < 0 {
		x -= 1
	}
	z := pos.Z
	if z < 0 {
		z -= 1
	}

	p := math32.Vector3{}
	p.X = float32(int64(x) - int64(x)%AREA_WIDTH)
	p.Z = float32(int64(z) - int64(z)%AREA_WIDTH)

	return p
}
