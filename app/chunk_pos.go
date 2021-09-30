package app

import (
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/lib/util"
)

type ChunkPos struct {
	X int64
	Z int64
}

func ToChunkPos(pos *math32.Vector3) ChunkPos {
	x := util.FloorFloat(pos.X / float32(CHUNK_WIDTH))
	z := util.FloorFloat(pos.Z / float32(CHUNK_WIDTH))

	return ChunkPos{X: x, Z: z}
}
