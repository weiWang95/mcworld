package block

import "github.com/g3n/engine/core"

type TransparentBlock struct {
	Block
}

func (b *TransparentBlock) Init() {
	b.Block.Init()
}

func (b *TransparentBlock) Transparent() bool {
	return true
}

func (b *TransparentBlock) SetVisible(state bool) {

}

func (b *TransparentBlock) AddTo(n core.INode, materialPath string) {

}
