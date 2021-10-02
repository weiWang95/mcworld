package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
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

func (b *EntityBlock) AddTo(n core.INode, mat material.IMaterial) {
	cube := geometry.NewCube(1)

	b.mesh = graphic.NewMesh(cube, mat)
	b.mesh.SetPositionVec(&b.Pos)

	n.GetNode().Add(b.mesh)
}

func (b *EntityBlock) RemoveFrom(n core.INode) {
	b.mesh.ClearMaterials()
	n.GetNode().Remove(b.mesh)
	b.mesh.Dispose()
}
