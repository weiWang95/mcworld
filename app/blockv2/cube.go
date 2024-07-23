package blockv2

import (
	"math"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

type BlockFace int

const (
	BlockFaceNone   BlockFace = iota - 1
	BlockFaceBack             // 后 0
	BlockFaceFront            // 前 1
	BlockFaceTop              // 上 2
	BlockFaceBottom           // 下 3
	BlockFaceRight            // 右 4
	BlockFaceLeft             // 左 5
)

var facePosAdditionMap = map[BlockFace]*math32.Vector3{
	BlockFaceLeft:   math32.NewVector3(0, 0.5, 0.5),
	BlockFaceRight:  math32.NewVector3(1, 0.5, 0.5),
	BlockFaceFront:  math32.NewVector3(0.5, 0.5, 0),
	BlockFaceBack:   math32.NewVector3(0.5, 0.5, 1),
	BlockFaceTop:    math32.NewVector3(0.5, 1, 0.5),
	BlockFaceBottom: math32.NewVector3(0.5, 0, 0.5),
}

type Cube struct {
	core.Node

	left   *graphic.Mesh
	right  *graphic.Mesh
	front  *graphic.Mesh
	back   *graphic.Mesh
	top    *graphic.Mesh
	bottom *graphic.Mesh

	meshs []*graphic.Mesh
}

func NewCube() *Cube {
	b := new(Cube)
	b.init()
	return b
}

func (b *Cube) init() {
	n := core.NewNode()
	n.SetPosition(0, 0, 0)

	left := b.buildPlane()
	left.RotateY(-math.Pi / 2)
	left.SetPosition(0, 0.5, 0.5)
	b.left = left
	n.Add(b.left)

	right := b.buildPlane()
	right.RotateY(math.Pi / 2)
	right.SetPosition(1, 0.5, 0.5)
	b.right = right
	n.Add(b.right)

	front := b.buildPlane()
	front.RotateY(-math.Pi)
	front.SetPosition(0.5, 0.5, 0)
	b.front = front
	n.Add(b.front)

	back := b.buildPlane()
	back.SetPosition(0.5, 0.5, 1)
	b.back = back
	n.Add(b.back)

	top := b.buildPlane()
	top.RotateX(-math.Pi / 2)
	top.SetPosition(0.5, 1, 0.5)
	b.top = top
	n.Add(b.top)

	bottom := b.buildPlane()
	bottom.RotateX(math.Pi / 2)
	bottom.SetPosition(0.5, 0, 0.5)
	b.bottom = bottom
	n.Add(b.bottom)

	b.meshs = []*graphic.Mesh{
		b.back, b.front, b.top, b.bottom, b.right, b.left,
	}

	b.Node = *n
}

func (b *Cube) SetPosition(x, y, z float32) {
	b.SetPositionVec(math32.NewVector3(x, y, z))
}

func (b *Cube) SetPositionVec(vpos *math32.Vector3) {
	b.GetNode().SetPositionVec(vpos)
	b.left.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceLeft]))
	b.right.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceRight]))
	b.front.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceFront]))
	b.back.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceBack]))
	b.top.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceTop]))
	b.bottom.SetPositionVec(vpos.Clone().Add(facePosAdditionMap[BlockFaceBottom]))
}

func (b *Cube) SetFaceVisible(face BlockFace, visible bool) {
	b.meshs[face].SetVisible(visible)
}

func (b *Cube) SetFaceLum(face BlockFace, lum uint8) {
	b.setLum(b.meshs[face], lum)
}

func (b *Cube) setLum(mesh *graphic.Mesh, lum uint8) {
	ms := mesh.Materials()
	if len(ms) == 0 {
		return
	}
	ms[0].IMaterial().(*material.Standard).SetColor(math32.NewColor("white").MultiplyScalar(float32(lum)/15.0*0.8 + 0.2))
}

func (b *Cube) GetFaceLum(idx int) uint8 {
	ms := b.meshs[idx].Materials()
	if len(ms) == 0 {
		return 0
	}
	return uint8((ms[0].IMaterial().(*material.Standard).AmbientColor().B - 0.2) / 0.8 * 15.0)
}

func (b *Cube) SetTextures(textures []texture.Texture2D) {
	if len(textures) == 0 {
		return
	}

	for i, _ := range b.meshs {
		if len(textures) != 6 {
			b.setTexture(b.meshs[i], &textures[0])
		} else {
			b.setTexture(b.meshs[i], &textures[i])
		}
	}
}

func (b *Cube) setTexture(mesh *graphic.Mesh, tex *texture.Texture2D) {
	mesh.Materials()[0].IMaterial().(*material.Standard).AddTexture(tex)
}

func (b *Cube) buildPlane() *graphic.Mesh {
	p := geometry.NewPlane(1, 1)

	mat := material.NewStandard(math32.NewColor("black"))
	mat.SetSide(material.SideFront)

	mesh := graphic.NewMesh(p, mat)
	return mesh
}

func (b *Cube) Dispose() {
	for i, _ := range b.meshs {
		b.meshs[i].ClearMaterials()
	}
	b.Node.Dispose()
}
