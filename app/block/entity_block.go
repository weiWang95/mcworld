package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type EntityBlock struct {
	Block

	mesh *graphic.Mesh
	box  *graphic.Lines
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

func (b *EntityBlock) AddTo(n core.INode, mats ...material.IMaterial) {
	cube := geometry.NewCube(1)

	if len(mats) == 1 {
		b.mesh = graphic.NewMesh(cube, mats[0])
	} else {
		b.mesh = graphic.NewMesh(cube, nil)
		for i, _ := range mats {
			b.mesh.AddGroupMaterial(mats[i], i)
		}
	}
	pos := b.Pos.Clone().Add(math32.NewVector3(0.5, 0, 0.5))
	b.mesh.SetPositionVec(pos)

	n.GetNode().Add(b.mesh)
}

func (b *EntityBlock) RemoveFrom(n core.INode) {
	b.mesh.ClearMaterials()
	n.GetNode().Remove(b.mesh)
	b.mesh.Dispose()
}
