package block

import (
	"github.com/g3n/engine/core"
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
	path := ""
	b.EntityBlock.AddTo(n, path)
}
