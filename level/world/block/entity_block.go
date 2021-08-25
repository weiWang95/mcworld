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

func (b *EntityBlock) AddTo(n core.INode, materialPath string) {
	cube := geometry.NewCube(1)
	mat := material.NewStandard(&math32.Color{1, 1, 1})

	// tex, err := texture.NewTexture2DFromImage(materialPath)
	// if err != nil {
	// 	panic(fmt.Sprintf("Error:%s loading texture:%s \n", err, materialPath))
	// }
	// mat.AddTexture(tex)

	b.mesh = graphic.NewMesh(cube, mat)
	b.mesh.SetPositionVec(&b.Pos)

	n.GetNode().Add(b.mesh)
}
