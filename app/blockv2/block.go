package blockv2

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type IDrawable interface {
	core.INode

	SetPosition(x, y, z float32)
	SetPositionVec(vpos *math32.Vector3)
	SetFaceVisible(face BlockFace, visible bool)
	SetFaceLum(face BlockFace, lum uint8)
	GetFaceLum(idx int) uint8
}

type Block struct {
	IDrawable
	*BlockAttr
}

func (b *Block) AddTo(n core.INode) {
	n.GetNode().Add(b)
}

func (b *Block) RemoveFrom(n core.INode) {
	n.GetNode().Remove(b)
	b.Dispose()
}

func (b *Block) GetPosition() math32.Vector3 {
	return b.GetNode().Position()
}

func (b *Block) Transparent() bool { return false }

type BlockAttr struct {
	BaseBlock
	Lumable
	Diggable
	Stackable
}

type BlockId uint64

type BaseBlock struct {
	Id       BlockId  `json:"id"`
	Name     string   `json:"name"`
	Textures []string `json:"textures"`
}

func (b *BaseBlock) GetId() BlockId {
	return b.Id
}

func (b *BaseBlock) GetState() uint8 {
	return 0
}
