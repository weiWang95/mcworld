package app

import (
	"math"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/app/blockv2"
	"github.com/weiWang95/mcworld/lib/util"
)

const DEFAULT_VIEW_DISTANCE int64 = 0
const DEFAULT_GRAVITY_SPEED float32 = -9.8
const MAX_GRAVITY_SPEED float32 = -20

const DAY_TOTAL_TIME int64 = 12000          // 每日时长
const DAY_NIGHT_TRANSITION_TIME int64 = 600 // 昼夜交替过渡时长
const MIN_SUN_LEVEL = 0                     // 最小阳光登录

type World struct {
	core.Node
	*logger.Logger

	lightCount float64
	ambLight   *light.Ambient

	sunLevel   uint8
	curTime    int64
	timeTicker *TickChecker

	wg IWorldGenerator
	cm *ChunkManager
	bu *BlockUpdater
	lu *LuminanceUpdater
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
	// w.updateLight(a, t)

	w.cm.Update(a, t)
	w.bu.Update(a, t)
	w.lu.Update(a, t)

	if w.timeTicker.Next(t) {
		w.curTime += 1
		if w.curTime > DAY_TOTAL_TIME {
			w.curTime = 0
		}
		newSunLevel := w.CalSunLevel()
		if w.sunLevel != newSunLevel {
			a.Log().Debug("sun level update: t:%v %v -> %v", w.curTime, w.sunLevel, newSunLevel)
			w.sunLevel = newSunLevel
			w.lu.UpdateSumLum()
		}
	}
}

func (w *World) Cleanup(a *App) {
}

func (w *World) setup(a *App) {
	w.Logger = a.Log()
	// w.setupLight()

	// seed := time.Now().UnixNano()
	// seed := int64(202210080000000)
	w.setupWorldGenerator(a.seed)

	w.cm = NewChunkManager(a)
	w.cm.Start(a)
	w.Add(w.cm)

	w.bu = NewBlockUpdater(a)
	w.lu = NewLuminanceUpdater(a)

	w.timeTicker = NewTickChecker(1)
	w.sunLevel = MIN_SUN_LEVEL
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

func (w *World) GetBlockByVec(vec math32.Vector3) (block *blockv2.Block, chunkLoaded bool) {
	return w.GetBlockByPosition(vec.X, vec.Y, vec.Z)
}

func (w *World) GetBlockByPosition(x, y, z float32) (block *blockv2.Block, chunkLoaded bool) {
	chunk := w.cm.GetChunk(x, y, z)
	if chunk == nil {
		return nil, false
	}

	return chunk.GetBlock(x, y, z), true
}

func (w *World) GetLum(x, y, z float32) (lum Luminance, chunkLoaded bool) {
	chunk := w.cm.GetChunk(x, y, z)
	if chunk == nil {
		return 0, false
	}

	return chunk.GetLumByWorldPos(x, y, z), true
}

func (w *World) GetLumByVec(vec math32.Vector3) (lum Luminance, chunkLoaded bool) {
	return w.GetLum(vec.X, vec.Y, vec.Z)
}

func (w *World) WreckBlock(pos math32.Vector3) {
	w.Debug("wreck block -> %v", pos)
	// area := w.getArea(pos.X, pos.Z)
	// area.ReplaceBlock(pos, nil)
	chunk := w.cm.GetChunk(pos.X, pos.Y, pos.Z)
	if chunk.ReplaceBlock(pos, nil) {
		w.bu.TiggerUpdate(util.NewPosFromVec3(pos))
		w.lu.TiggerUpdate(util.NewPosFromVec3(pos))
	}
}

func (w *World) PlaceBlock(block *blockv2.Block, pos math32.Vector3) {
	w.Debug("place block:%T -> %v", block, pos)
	chunk := w.cm.GetChunk(pos.X, pos.Y, pos.Z)
	if chunk.ReplaceBlock(pos, block) {
		w.bu.TiggerUpdate(util.NewPosFromVec3(pos))
		w.lu.TiggerUpdate(util.NewPosFromVec3(pos))
	}
}

func (w *World) WorldGenerator() IWorldGenerator {
	return w.wg
}

func (w *World) CalSunLevel() uint8 {
	tHalf := DAY_NIGHT_TRANSITION_TIME / 2
	dawnStart := DAY_TOTAL_TIME/2 - DAY_NIGHT_TRANSITION_TIME
	dawnEnd := DAY_TOTAL_TIME/2 - tHalf
	duskStart := DAY_TOTAL_TIME - DAY_NIGHT_TRANSITION_TIME
	duskEnd := DAY_TOTAL_TIME

	speed := 2
	stepTime := DAY_NIGHT_TRANSITION_TIME / (int64(MAX_LUM - MIN_SUN_LEVEL)) * int64(speed)

	if w.curTime >= 0 && w.curTime < dawnStart {
		return MIN_SUN_LEVEL
	} else if w.curTime >= dawnStart && w.curTime < dawnEnd {
		// return MIN_SUN_LEVEL
		addSun := (w.curTime - dawnStart) / stepTime
		return MIN_SUN_LEVEL + uint8(addSun)*uint8(speed)
	} else if w.curTime >= dawnEnd && w.curTime < duskStart {
		return MAX_LUM
	} else if w.curTime >= duskStart && w.curTime <= duskEnd {
		// return MAX_LUM
		addSun := (w.curTime - duskStart) / stepTime
		return MAX_LUM - uint8(addSun)*uint8(speed)
	}

	return MAX_LUM
}

func (w *World) SunLumRate() float32 {
	return float32(w.sunLevel) / float32(MAX_LUM)
}
