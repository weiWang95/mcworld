package block

type BlockId uint64

const (
	BlockAir BlockId = iota + 1
	BlockSoil
	BlockBrick
	BlockLamp
)

func init() {
	RegisterBlock(BlockSoil, &Soil{})
	RegisterBlock(BlockBrick, &Brick{})
	RegisterBlock(BlockLamp, &Lamp{})
}

type Soil struct {
	EntityBlock
}

type Brick struct {
	EntityBlock
}

type Lamp struct {
	EntityBlock
	BaseLumable
}

func (b *Lamp) Init(id BlockId) {
	b.EntityBlock.Init(id)
	b.lum = 15
}
