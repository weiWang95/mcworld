package block

import (
	"github.com/g3n/engine/core"
	"github.com/weiWang95/mcworld/app/loader"
)

func init() {
	RegisterBlock(BlockSoil, &Soil{})
}

type Soil struct {
	EntityBlock
}

func (b *Soil) Init() {
}

func (b *Soil) AddTo(n core.INode) {
	mat := loader.LoadBlockMaterial(uint64(BlockSoil))
	b.EntityBlock.AddTo(n, mat)
}

func (b *Soil) RemoveFrom(n core.INode) {
	b.EntityBlock.RemoveFrom(n)
}
