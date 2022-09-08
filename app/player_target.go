package app

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/block"
)

type PlayerTarget struct {
	core.Node

	b   block.IBlock
	box *graphic.Lines
}

func NewPlayerTarget() *PlayerTarget {
	target := new(PlayerTarget)
	target.Node = *core.NewNode()

	target.addBox()
	target.SetTarget(nil, block.BlockFaceNone)

	return target
}

func (target *PlayerTarget) SetTarget(b block.IBlock, face block.BlockFace) {
	target.b = b
	target.SetVisible(target.b != nil)
	if target.b != nil {
		pos := b.GetPosition()
		target.SetPositionVec(&pos)
		target.box.SetPositionVec(&pos)
	}
}

func (target *PlayerTarget) addBox() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		0, 0, 0,
		1, 0, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 0, 1,
		1, 1, 1,
		0, 1, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 0,
	)
	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	target.box = graphic.NewLines(geom, mat)
	target.Add(target.box)
}
