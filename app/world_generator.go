package app

import (
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/lib/perlin"
)

const MAX_GROUND_HEIGHT = CHUNK_HEIGHT / 2
const MIN_GROUND_HEIGHT = CHUNK_HEIGHT / 6

type IWorldGenerator interface {
	Setup(seed int64)
	GetBlock(x, y, z float64) block.BlockId
}

type WorldGenerator struct {
	seed int64
	p    *perlin.Perlin
}

func (wg *WorldGenerator) Setup(seed int64) {
	wg.seed = seed
	wg.p = perlin.NewPerlin(2, 2, 5, wg.seed)
}

func (wg *WorldGenerator) GetBlock(x, y, z float64) block.BlockId {
	h := (wg.p.Noise2D(0.015*x, 0.015*z)+1)*float64(MAX_GROUND_HEIGHT-MIN_GROUND_HEIGHT) + float64(MIN_GROUND_HEIGHT)
	if int64(y) > int64(h) {
		return 0
	}

	return block.BlockSoil
}
