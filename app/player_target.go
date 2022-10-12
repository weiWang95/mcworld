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

	b     block.IBlock
	box   *graphic.Lines
	point *graphic.Mesh
}

func NewPlayerTarget() *PlayerTarget {
	target := new(PlayerTarget)
	target.Node = *core.NewNode()

	target.addBox()
	target.SetTarget(nil, nil)

	return target
}

func (target *PlayerTarget) SetTarget(b block.IBlock, hit *math32.Vector3) {
	target.b = b
	target.SetVisible(target.b != nil)
	if target.b != nil {
		pos := b.GetPosition()
		target.SetPositionVec(&pos)
		target.box.SetPositionVec(&pos)
		target.point.SetPositionVec(hit)
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

		1, 1, 1,
		1, 1, 0,

		0, 1, 0,
		0, 1, 1,

		0, 1, 0,
		1, 1, 0,

		1, 0, 1,
		0, 0, 1,

		1, 0, 1,
		1, 0, 0,

		1, 0, 0,
		1, 1, 0,

		0, 0, 1,
		0, 1, 1,
	)

	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	target.box = graphic.NewLines(geom, mat)
	target.Add(target.box)

	s1 := geometry.NewSphere(0.05, 5, 5)
	target.point = graphic.NewMesh(s1, material.NewStandard(math32.NewColor("red")))
	target.Add(target.point)
}
