package app

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/lib/util"
)

type ChunkState uint8

const (
	Unloaded  ChunkState = 0
	Unloading ChunkState = 1
	Loaded    ChunkState = 2
	Rendered  ChunkState = 3
)

const (
	CHUNK_WIDTH        int64 = 16
	CHUNK_HEIGHT       int64 = 40
	CHUNK_UPDATE_RANGE int64 = 1
)

type Chunk struct {
	core.Node

	State ChunkState

	pos    *ChunkPos
	actPos *math32.Vector3

	blocks [CHUNK_HEIGHT][CHUNK_WIDTH][CHUNK_WIDTH]block.IBlock
	axis   core.INode
}

func NewChunk(x, z int64) *Chunk {
	c := new(Chunk)

	c.Node = *core.NewNode()

	c.pos = &ChunkPos{X: x, Z: z}
	c.actPos = math32.NewVector3(float32(x*CHUNK_WIDTH), 0, float32(z*CHUNK_WIDTH))
	c.SetPositionVec(c.actPos)

	return c
}

func (c *Chunk) Start(a *App) {
	c.setup(a)
	c.SetVisible(false)
}

func (c *Chunk) Update(a *App, t time.Duration) {
	c.axis.SetVisible(a.IsDebugMode())
}

func (c *Chunk) Cleanup() {
	if c.axis != nil {
		c.axis.Dispose()
	}

	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			for z := int64(0); z < CHUNK_WIDTH; z++ {
				if c.blocks[y][x][z] != nil {
					c.blocks[y][x][z].RemoveFrom(c)
				}
			}
		}
	}
}

func (c *Chunk) setup(a *App) {
	c.Load(a)

	c.RefreshBlocks(a)

	c.addAxis()
}

func (c *Chunk) Load(a *App) {
	wg := a.World().WorldGenerator()

	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			for z := int64(0); z < CHUNK_WIDTH; z++ {
				pos := math32.NewVector3(c.actPos.X+float32(x), float32(y), c.actPos.Z+float32(z))
				id := wg.GetBlock(float64(pos.X), float64(pos.Y), float64(pos.Z))
				if id != 0 {
					b := block.NewBlock(id, *pos)
					c.blocks[y][x][z] = b
				}
			}
		}
	}

	c.State = Loaded
}

func (c *Chunk) RefreshBlocks(a *App) {
	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			for z := int64(0); z < CHUNK_WIDTH; z++ {
				c.RefreshBlock(x, y, z)
			}
		}
	}
}

func (c *Chunk) RefreshNearbyBlocks(bx, by, bz int64) {
	for y := util.MaxInt64(0, by-CHUNK_UPDATE_RANGE); y <= util.MinInt64(by+CHUNK_UPDATE_RANGE, CHUNK_HEIGHT-1); y++ {
		for x := util.MaxInt64(0, bx-CHUNK_UPDATE_RANGE); x <= util.MinInt64(bx+CHUNK_UPDATE_RANGE, CHUNK_WIDTH-1); x++ {
			for z := util.MaxInt64(0, bz-CHUNK_UPDATE_RANGE); z <= util.MinInt64(bz+CHUNK_UPDATE_RANGE, CHUNK_WIDTH-1); z++ {
				c.RefreshBlock(x, y, z)
			}
		}
	}
}

func (c *Chunk) RefreshBlock(x, y, z int64) {
	world := Instance().World()
	if c.blocks[y][x][z] == nil {
		return
	}

	c.blocks[y][x][z].AddTo(c)

	if y == 0 && c.IsTransparent(x, y+1, z) {
		c.blocks[y][x][z].SetVisible(true)
		return
	} else if y == CHUNK_HEIGHT-1 && c.IsTransparent(x, y-1, z) {
		c.blocks[y][x][z].SetVisible(true)
		return
	}

	if (x == 0 && BlockIsTransparent(world.GetBlockByVec(c.BlockPos(x-1, y, z)))) ||
		(x == CHUNK_WIDTH-1 && BlockIsTransparent(world.GetBlockByVec(c.BlockPos(x+1, y, z)))) ||
		(z == 0 && BlockIsTransparent(world.GetBlockByVec(c.BlockPos(x, y, z-1)))) ||
		(z == CHUNK_WIDTH-1 && BlockIsTransparent(world.GetBlockByVec(c.BlockPos(x, y, z+1)))) {
		c.blocks[y][x][z].SetVisible(true)
		return
	} else if (y > 0 && c.IsTransparent(x, y-1, z)) || c.IsTransparent(x, y+1, z) ||
		(x > 0 && c.IsTransparent(x-1, y, z)) || (x < CHUNK_WIDTH-1 && c.IsTransparent(x+1, y, z)) ||
		(z > 0 && c.IsTransparent(x, y, z-1)) || (z < CHUNK_WIDTH-1 && c.IsTransparent(x, y, z+1)) {

		c.blocks[y][x][z].SetVisible(true)
		return
	}

	c.blocks[y][x][z].SetVisible(false)
}

func (c *Chunk) Rendered(a *App) {
	c.SetVisible(true)

	c.State = Rendered
}

func (c *Chunk) Unrendered() {
	c.SetVisible(false)

	c.State = Loaded
}

func (c *Chunk) addAxis() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		c.actPos.X, 0, c.actPos.Z,
		c.actPos.X+float32(CHUNK_WIDTH), 0, c.actPos.Z,
		c.actPos.X, 0, float32(c.actPos.Z),
		c.actPos.X, float32(CHUNK_HEIGHT), c.actPos.Z,
		c.actPos.X, 0, c.actPos.Z,
		c.actPos.X, 0, c.actPos.Z+float32(CHUNK_WIDTH),
	)
	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 0.0, 0.0, // red
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0, // green
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0, // blue
		0.0, 0.0, 1.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	c.axis = graphic.NewLines(geom, mat)
	c.Add(c.axis)
}

func (c *Chunk) IsTransparent(x, y, z int64) bool {
	return c.blocks[y][x][z] == nil || c.blocks[y][x][z].Transparent()
}

func (c *Chunk) GetBlock(x, y, z float32) block.IBlock {
	if y < 0 || y >= float32(CHUNK_HEIGHT) {
		return nil
	}

	bx := util.FloorFloat(x) - int64(c.actPos.X)
	bz := util.FloorFloat(z) - int64(c.actPos.Z)

	return c.blocks[int64(y)][bx][bz]
}

func (c *Chunk) BlockPos(x, y, z int64) math32.Vector3 {
	return *c.actPos.Clone().Add(math32.NewVector3(float32(x), float32(y), float32(z)))
}

func (c *Chunk) ReplaceBlock(pos math32.Vector3, block block.IBlock) bool {
	Instance().Log().Debug("replace block: pos -> %v, to -> %v", pos, block)

	if pos.Y < 0 || pos.Y >= float32(CHUNK_HEIGHT) {
		return false
	}

	bx := util.FloorFloat(pos.X) - int64(c.actPos.X)
	bz := util.FloorFloat(pos.Z) - int64(c.actPos.Z)

	if bx < 0 || bx >= CHUNK_WIDTH || bz < 0 || bz >= CHUNK_WIDTH {
		return false
	}

	b := c.blocks[int64(pos.Y)][bx][bz]
	if b != nil {
		b.RemoveFrom(c)
	}

	c.blocks[int64(pos.Y)][bx][bz] = block
	c.RefreshNearbyBlocks(bx, int64(pos.Y), bz)
	return true
}
