package block

import (
	"github.com/g3n/engine/core"
)

func init() {
	RegisterBlock(BlockAir, &Air{})
}

type Air struct {
	TransparentBlock
}

func (b *Air) Init() {
}

func (a *Air) AddTo(n core.INode) {

}
