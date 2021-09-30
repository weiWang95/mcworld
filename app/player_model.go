package app

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type PlayerModel struct {
	core.Node

	target    core.INode
	axis      *graphic.Lines
	wreckLine *graphic.Lines

	w        float32
	h        float32
	viewRate float32
}

func NewPlayerModel() *PlayerModel {
	p := new(PlayerModel)

	p.Node = *core.NewNode()
	p.w = 1
	p.h = 2
	p.viewRate = 0.9

	return p
}

func (p *PlayerModel) Start(a *App) {
	box := geometry.NewBox(1, 2, 1)
	mat := material.NewStandard(math32.NewColor("red"))
	p.target = graphic.NewMesh(box, mat)
	p.Add(p.target)

	p.addAxis()
}

func (p *PlayerModel) Update(a *App, t time.Duration) {
}

func (p *PlayerModel) Cleanup() {

}

func (p *PlayerModel) addAxis() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		0, 0, 0,
		0, 2, 0,
	)
	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 0.0, 0.0, // red
		1.0, 0.0, 0.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	p.axis = graphic.NewLines(geom, mat)
	p.Add(p.axis)
}

func (p *PlayerModel) GetPosition() *math32.Vector3 {
	pos := p.Node.Position()
	return &pos
}

func (p *PlayerModel) SetPosition(pos *math32.Vector3) {
	p.Node.SetPosition(pos.X, pos.Y, pos.Z)
}

func (p *PlayerModel) SetFace(face *math32.Vector3) {
}

func (p *PlayerModel) GetViewport() *math32.Vector3 {
	return p.GetPosition().Clone().Add(math32.NewVector3(0, p.h*p.viewRate, 0))
}

func (p *PlayerModel) GetHandLength() float32 {
	return 6
}

func (p *PlayerModel) GetBoundBox() BoundBox {
	pos := p.GetPosition()

	return BoundBox{
		X:  pos.X - p.w*0.5,
		Y:  pos.Y,
		Z:  pos.Z - p.w*0.5,
		BX: pos.X + p.w*0.5,
		BY: pos.Y + p.h,
		BZ: pos.Z + p.w*0.5,
	}
}
