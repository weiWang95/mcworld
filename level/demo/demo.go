package demo

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/weiWang95/mcworld/app"
)

func init() {
	app.RegisterLevel("demo", &Demo{})
}

type Demo struct {
	core.Node
}

func (d *Demo) setup(a *app.App) {
	d.Node = *core.NewNode()

	floorMesh := graphic.NewMesh(geometry.NewBox(100, 2, 10), material.NewStandard(&math32.Color{1, 1, 1}))
	floorMesh.SetPositionY(-1)
	d.Add(floorMesh)

	cube := graphic.NewMesh(geometry.NewCube(1), material.NewStandard(&math32.Color{1, 1, 1}))
	d.Add(cube)

	p1 := geometry.NewPlane(10, 10)
	m1 := material.NewStandard(&math32.Color{1, 1, 1})
	texfile := a.DirData() + "/images/moss.png"
	tex3, err := texture.NewTexture2DFromImage(texfile)
	if err != nil {
		a.Log().Fatal("Error:%s loading texture:%s", err, texfile)
	}
	m1.AddTexture(tex3)
	mesh1 := graphic.NewMesh(p1, m1)
	d.Add(mesh1)

	lightColor := &math32.Color{1, 1, 1}
	geom := geometry.NewSphere(0.05, 16, 8)
	mat := material.NewStandard(lightColor)
	mat.SetUseLights(0)
	mat.SetEmissiveColor(lightColor)
	lMesh := graphic.NewMesh(geom, mat)
	lMesh.SetVisible(true)

	light := light.NewPoint(lightColor, 1.0)
	light.SetPosition(0, 0, 0)
	light.SetLinearDecay(1)
	light.SetQuadraticDecay(1)
	light.SetVisible(true)
	lMesh.Add(light)

	lMesh.SetPosition(1, 4, 1)

	d.Add(lMesh)
}

func (d *Demo) Start(a *app.App) {
	d.setup(a)
	a.Scene().Add(d)
}

func (d *Demo) Update(a *app.App, t time.Duration) {

}

func (d *Demo) Cleanup(a *app.App) {

}
