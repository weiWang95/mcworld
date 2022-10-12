package block

import (
	"fmt"
	"reflect"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type BlockState uint8

type BlockFace int

const (
	BlockFaceNone   BlockFace = iota - 1
	BlockFaceBack             // 后 0
	BlockFaceFront            // 前 1
	BlockFaceTop              // 上 2
	BlockFaceBottom           // 下 3
	BlockFaceRight            // 右 4
	BlockFaceLeft             // 左 5
)

var blockMap map[BlockId]IBlock

type IBlock interface {
	Init(id BlockId)
	GetId() BlockId
	GetState() BlockState
	GetBlockLum() uint8
	Lumable() bool
	Transparent() bool
	SetVisible(state bool)
	Visible() bool
	SetPosition(pos math32.Vector3)
	GetPosition() math32.Vector3
	SetLum(lum uint8, idx int)
	RefreshLum()
	AddTo(n core.INode)
	RemoveFrom(n core.INode)
	GetFaceLum(idx int) uint8
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
	b.Init(id)
	b.SetPosition(pos)
	return b
}
