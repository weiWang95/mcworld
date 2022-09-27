package block

func init() {
	RegisterBlock(BlockSoil, &Soil{})
}

type Soil struct {
	EntityBlock
}

func (b *Soil) Init() {
	b.EntityBlock.Init()
	b.Id = BlockSoil
}
