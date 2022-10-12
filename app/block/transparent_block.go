package block

import "github.com/g3n/engine/core"

func init() {
	RegisterBlock(BlockAir, &TransparentBlock{})
}

var _ IBlock = (*TransparentBlock)(nil)

type TransparentBlock struct {
	Block
}

func (b *TransparentBlock) Init(id BlockId) {
	b.Block.Init(id)
}

func (b *TransparentBlock) Transparent() bool {
	return true
}

func (b *TransparentBlock) SetVisible(state bool) {

}

func (b *TransparentBlock) Visible() bool {
	return false
}

func (b *TransparentBlock) AddTo(n core.INode) {

}

func (b *TransparentBlock) RemoveFrom(n core.INode) {

}
