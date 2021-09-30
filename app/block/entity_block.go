package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

type EntityBlock struct {
	Block

	mesh *graphic.Mesh
}

func (b *EntityBlock) Init() {
	b.Block.Init()
}

func (b *EntityBlock) SetVisible(state bool) {
	if b.mesh == nil {
		return
	}

	b.mesh.SetVisible(state)
}

func (b *EntityBlock) AddTo(n core.INode, tex *texture.Texture2D) {
	cube := geometry.NewCube(1)

	mat := material.NewStandard(&math32.Color{1, 1, 1})
	mat.AddTexture(tex)

	b.mesh = graphic.NewMesh(cube, mat)
	b.mesh.SetPositionVec(&b.Pos)

	n.GetNode().Add(b.mesh)
}

func (b *EntityBlock) RemoveFrom(n core.INode) {
	n.GetNode().Remove(b.mesh)
	b.mesh.Dispose()
}
