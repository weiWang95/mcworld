package app

import (
	"math"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/app/block"
)

const DEFAULT_VIEW_DISTANCE int64 = 0
const DEFAULT_GRAVITY_SPEED float32 = -9.8
const MAX_GRAVITY_SPEED float32 = -20

type World struct {
	core.Node
	*logger.Logger

	lightCount float64
	ambLight   *light.Ambient

	wg IWorldGenerator
	cm *ChunkManager

	activePos   math32.Vector3
	activeAreas [][]*Area
}

func NewWorld() *World {
	w := new(World)
	w.Node = *core.NewNode()

	return w
}

func (w *World) Start(a *App) {
	w.setup(a)
	a.Scene().Add(w)
}

func (w *World) Update(a *App, t time.Duration) {
	w.updateLight(a, t)

	w.cm.Update(a, t)
}

func (w *World) Cleanup(a *App) {
}

func (w *World) setup(a *App) {
	w.Logger = a.Log()
	w.setupLight()

	seed := time.Now().UnixNano()
	w.setupWorldGenerator(seed)

	w.cm = NewChunkManager()
	w.cm.Start(a)
	w.Add(w.cm)
}

func (w *World) refreshAreas() {
	for i := 0; i < len(w.activeAreas); i++ {
		for j := 0; j < len(w.activeAreas[i]); j++ {
			w.activeAreas[i][j].refreshBlocks(w)
		}
	}
}

func (w *World) setupLight() {
	w.ambLight = light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 1)
	w.ambLight.SetDirection(-1, -1, -1)
	w.Add(w.ambLight)
}

func (w *World) setupWorldGenerator(seed int64) {
	w.wg = &WorldGenerator{}
	w.wg.Setup(seed)
}

func (w *World) updateLight(a *App, t time.Duration) {
	intensity := math.Min(math.Max(math.Sin(w.lightCount)+1, 0.6), 1.4)
	w.ambLight.SetIntensity(float32(intensity))
	w.lightCount += float64(t / (10 * time.Second))
}

func (w *World) getArea(x, z float32) *Area {
	areaX := int64(x) - (int64(x) % AREA_WIDTH)
	if x < 0 {
		areaX = int64(x) - AREA_WIDTH - int64(x)%AREA_WIDTH
	}
	areaZ := int64(z) - (int64(z) % AREA_WIDTH)
	if z < 0 {
		areaZ = int64(z) - AREA_WIDTH - int64(z)%AREA_WIDTH
	}

	for i := 0; i < len(w.activeAreas)-1; i++ {
		for j := 0; j < len(w.activeAreas[i])-1; j++ {
			area := w.activeAreas[i][j]
			if int64(area.Pos.X) == areaX && int64(area.Pos.Z) == areaZ {
				return area
			}
		}
	}

	return nil
}

func (w *World) GetBlockByVec(vec math32.Vector3) block.IBlock {
	return w.GetBlockByPosition(vec.X, vec.Y, vec.Z)
}

func (w *World) GetBlockByPosition(x, y, z float32) block.IBlock {
	chunk := w.cm.GetChunk(x, y, z)
	if chunk == nil {
		return nil
	}

	return chunk.GetBlock(x, y, z)
}

func (w *World) activeAreasLen() int64 {
	return 2*DEFAULT_VIEW_DISTANCE + 1
}

func (w *World) WreckBlock(pos math32.Vector3) {
	w.Debug("wreck block -> %v", pos)
	// area := w.getArea(pos.X, pos.Z)
	// area.ReplaceBlock(pos, nil)
	chunk := w.cm.GetChunk(pos.X, pos.Y, pos.Z)
	chunk.ReplaceBlock(pos, nil)
}

func (w *World) WorldGenerator() IWorldGenerator {
	return w.wg
}
