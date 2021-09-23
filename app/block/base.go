package block

import (
	"fmt"
	"reflect"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type BlockId uint64

const (
	BlockAir BlockId = iota + 1
	BlockSoil
)

var blockMap map[BlockId]IBlock

type IBlock interface {
	Init()
	Transparent() bool
	SetVisible(state bool)
	SetPosition(pos math32.Vector3)
	GetPosition() math32.Vector3
	AddTo(n core.INode)
	RemoveFrom(n core.INode)
}

func RegisterBlock(id BlockId, b IBlock) {
	if blockMap == nil {
		blockMap = make(map[BlockId]IBlock)
	}

	if _, ok := blockMap[id]; ok {
		panic(fmt.Sprintf("dup block id: %d", id))
	}

	blockMap[id] = b
}

func NewBlock(id BlockId, pos math32.Vector3) IBlock {
	v := reflect.New(reflect.TypeOf(blockMap[id]).Elem())
	b := v.Interface().(IBlock)
	b.Init()
	b.SetPosition(pos)
	return b
}

func GetBlockTexturePath(id BlockId) string {
	return fmt.Sprintf("blocks/%d.png", id)
}
