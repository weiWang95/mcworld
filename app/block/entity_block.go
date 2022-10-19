package block

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/loader"
)

var _ IBlock = (*EntityBlock)(nil)

type EntityBlock struct {
	Block
	BaseDestructible

	mesh *graphic.Mesh
	mats []material.IMaterial
}

func (b *EntityBlock) Init(id BlockId) {
	b.Block.Init(id)
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
	pos := b.Pos.Clone().Add(math32.NewVector3(0.5, 0.5, 0.5))
	b.mesh.SetPositionVec(pos)

	n.GetNode().Add(b.mesh)
}

func (b *EntityBlock) SetLum(lum uint8, idx int) {
	b.mats[idx].(*material.Standard).SetColor(math32.NewColor("white").MultiplyScalar(float32(lum)/15.0*0.8 + 0.2))
}

func (b *EntityBlock) GetFaceLum(idx int) uint8 {
	return uint8(b.mats[idx].(*material.Standard).AmbientColor().B * 15)
}

func (b *EntityBlock) RefreshLum() {
	// b.mesh.ClearMaterials()
	// for i, _ := range b.mats {
	// 	b.mesh.AddGroupMaterial(b.mats[i], i)
	// }
	// b.mesh.SetChanged(true)
}

func (b *EntityBlock) RemoveFrom(n core.INode) {
	if b.mesh == nil {
		return
	}

	b.mesh.ClearMaterials()
	n.GetNode().Remove(b.mesh)
	b.mesh.Dispose()
}
