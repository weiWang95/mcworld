package app

import (
	"fmt"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
)

type GameMode int8

const (
	LifeMode GameMode = iota
	GodMode
)

const (
	PLAYER_VIEWPORT_HEIGHT float32 = 1.5
	PLAYER_JUMP_SPEED      float32 = 4.9
)

type OldPlayer struct {
	core.Node
	core.Dispatcher

	mode   GameMode
	speed  float32
	vSpeed float32
	inFall bool

	Camera    *camera.Camera
	CamOrb    *OldPlayerControl
	Model     *graphic.Mesh
	axis      *graphic.Lines
	wreckLine *graphic.Lines

	Pos    *math32.Vector3
	LookAt *math32.Vector3
}

func NewOldPlayer() *OldPlayer {
	p := new(OldPlayer)
	p.Node = *core.NewNode()

	p.Init()

	return p
}

func (p *OldPlayer) Init() {
	p.mode = LifeMode
	p.speed = 10.0

	p.Pos = math32.NewVector3(0, 0, 0)
	p.LookAt = math32.NewVector3(0, 1.5, 0)

	p.Camera = camera.New(16 / 9)
	p.Camera.UpdateSize(3)
	p.Camera.LookAt(p.LookAt, &math32.Vector3{0, 1, 0})
	p.Camera.SetProjection(camera.Perspective)
	p.Camera.SetPositionVec(p.LookAt.Clone().Add(math32.NewVector3(-2, 1.5, 0)))
	p.CamOrb = NewOldPlayerControl(p, p.Camera)

	p.Add(helper.NewAxes(1))

	p.initBody()
	p.addWreckLine()
	p.addAxis()
}

func (p *OldPlayer) Dispose() {
}

func (p *OldPlayer) Update(a *App, t time.Duration) {
	p.CamOrb.Update(a, t)
}

func (p *OldPlayer) initBody() {
	box := geometry.NewBox(1, 1, 1)
	mat := material.NewStandard(math32.NewColor("red"))

	p.Model = graphic.NewMesh(box, mat)
	p.Add(p.Model)
}

func (p *OldPlayer) addWreckLine() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		0, 0, 0,
		MaxControlDistance, 0, 0,
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
	p.wreckLine = graphic.NewLines(geom, mat)
	p.wreckLine.SetPosition(0, 1.5, 0)
	p.Add(p.wreckLine)
}

func (p *OldPlayer) addAxis() {
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

func (p *OldPlayer) GetPosition() math32.Vector3 {
	return *p.Pos
}

func (p *OldPlayer) SetPosition(vec math32.Vector3) {
	p.Pos.Set(vec.X, vec.Y, vec.Z)
	p.Model.SetPositionVec(p.Pos)
	p.axis.SetPositionVec(p.Pos)
	p.LookAt.Set(vec.X, vec.Y+1.5, vec.Z)
	p.CamOrb.SetTarget(*p.LookAt)
}

func (p *OldPlayer) ResetPosition(vec math32.Vector3) {
	p.SetPosition(vec)

	p.Camera.SetPositionVec(p.Pos.Clone().Add(math32.NewVector3(-2, 1.5, 0)))
	p.Camera.LookAt(p.LookAt, &math32.Vector3{0, 1, 0})
}

func (p *OldPlayer) SetDirection(d math32.Vector3) {

}

func (p *OldPlayer) GetSpeed() float32 {
	return p.speed
}

func (p *OldPlayer) GetJumpVt() float32 {
	return PLAYER_JUMP_SPEED
}

func (p *OldPlayer) Jump() {
	fmt.Printf("pos -> %v, vSpeed -> %v, fall: %v \n", p.Pos, p.vSpeed, p.inFall)
	if p.inFall {
		return
	}

	p.vSpeed = p.GetJumpVt()
}
