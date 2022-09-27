package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/loader"
)

type EntityBlock struct {
	Block

	mesh *graphic.Mesh
	mats []material.IMaterial
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

func (b *EntityBlock) Visible() bool {
	if b.mesh == nil {
		return false
	}

	return b.mesh.Visible()
}

func (b *EntityBlock) AddTo(n core.INode) {
	mats := loader.LoadBlockMaterial(uint64(b.Id))
	cube := geometry.NewCube(1)

	b.mats = mats

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

func (b *EntityBlock) SetLum(lum uint8, idx int) {
	b.mats[idx].(*material.Standard).SetAmbientColor(math32.NewColor("white").MultiplyScalar(float32(lum)/15.0*0.5 + 0.5))
}

func (b *EntityBlock) RemoveFrom(n core.INode) {
	b.mesh.ClearMaterials()
	n.GetNode().Remove(b.mesh)
	b.mesh.Dispose()
}
