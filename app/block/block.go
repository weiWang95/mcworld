package block

import "github.com/g3n/engine/math32"

type Block struct {
	Id    BlockId
	State BlockState
	Pos   math32.Vector3
}

func (b *Block) Init() {
	b.Pos = math32.Vector3{}
}

func (b *Block) GetId() BlockId {
	return b.Id
}

func (b *Block) GetState() BlockState {
	return b.State
}

func (b *Block) GetBlockLum() uint8 {
	return 0
}

func (b *Block) Lumable() bool {
	return false
}

func (b *Block) Transparent() bool {
	return false
}

func (b *Block) SetPosition(pos math32.Vector3) {
	b.Pos.X = pos.X
	b.Pos.Y = pos.Y
	b.Pos.Z = pos.Z
}

func (b *Block) GetPosition() math32.Vector3 {
	return b.Pos
}

func (b *Block) SetLum(lum uint8, idx int) {}
