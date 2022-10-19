package app

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/lib/util"
)

const MAX_ALIVE_TICK = 100

type UnloadingChunk struct {
	*Chunk

	AliveTick int64
}

type ChunkManager struct {
	core.Node
	*logger.Logger

	app *App

	loadDistance   int64
	renderDistance int64

	centerChunk       *Chunk
	loadingChunkMap   map[string]ChunkPos
	loadedChunkMap    map[string]*Chunk
	UnloadingChunkMap map[string]*UnloadingChunk
	renderedCount     int64
	unrenderedCount   int64

	// ticker
	loadTicker    *TickChecker
	loadingTicker *TickChecker
	inViewTicker  *TickChecker
}

var _ IRender = (*ChunkManager)(nil)

func NewChunkManager(app *App) *ChunkManager {
	cm := new(ChunkManager)

	cm.Node = *core.NewNode()
	cm.app = app

	cm.loadDistance = 1
	cm.renderDistance = 1

	cm.loadingChunkMap = make(map[string]ChunkPos)
	cm.loadedChunkMap = make(map[string]*Chunk)
	cm.UnloadingChunkMap = make(map[string]*UnloadingChunk)

	cm.loadTicker = NewTickChecker(4)
	cm.loadingTicker = NewTickChecker(1)
	cm.inViewTicker = NewTickChecker(4)

	return cm
}

func (cm *ChunkManager) Start(a *App) {
	cm.setup(a)
}

func (cm *ChunkManager) Update(a *App, t time.Duration) {
	if cm.loadTicker.Next(t) {
		cm.checkAndLoadChunks(a, a.Player().GetPosition())
	} else if cm.loadingTicker.Next(t) {
		cm.StepLoadChunk(a)
	}

	if cm.inViewTicker.Next(t) {
		cm.UpdateChunkInView()
	}
}

func (cm *ChunkManager) Cleanup() {
}

func (cm *ChunkManager) setup(a *App) {
	cm.Logger = a.Log()

	curPos := &math32.Vector3{}

	player := a.Player()
	if player != nil {
		curPos = player.GetPosition()
	}

	cm.checkAndLoadChunks(a, curPos)
}

func (cm *ChunkManager) checkAndLoadChunks(a *App, curPos *math32.Vector3) {
	centerPos := ToChunkPos(curPos)
	// a.Log().Debug("center pos: %v", centerPos)
	// a.Log().Debug("loaded chunks: %d", len(cm.loadedChunkMap))
	// for _, c := range cm.loadedChunkMap {
	// 	a.Log().Debug("X:%d Z:%d", c.pos.X, c.pos.Z)
	// }

	if cm.centerChunk != nil && cm.centerChunk.pos.X == centerPos.X && cm.centerChunk.pos.Z == centerPos.Z {
		// 更新卸载区块信息
		// a.Log().Debug("unloading chunks: %d", len(cm.UnloadingChunkMap))
		cm.StepUnloadChunk(a)

		cm.renderedCount, cm.unrenderedCount = 0, 0
		for x := -cm.renderDistance; x <= cm.renderDistance; x++ {
			for z := -cm.renderDistance; z <= cm.renderDistance; z++ {
				pos := ChunkPos{X: centerPos.X + x, Z: centerPos.Z + z}
				posId := pos.Id()

				chunk, ok := cm.loadedChunkMap[posId]
				if !ok {
					continue
				}

				// 渲染
				if x >= -cm.renderDistance && x <= cm.renderDistance && z >= -cm.renderDistance && z <= cm.renderDistance && chunk.inView {
					chunk.Rendered(a)
					cm.renderedCount++
				} else {
					chunk.Unrendered()
					cm.unrenderedCount++
				}
			}
		}

		return
	}

	// 将所有加载的区块标记为待卸载
	willUnload := make(map[string]bool, len(cm.loadedChunkMap))
	for _, c := range cm.loadedChunkMap {
		willUnload[c.pos.Id()] = true
	}

	for x := -cm.loadDistance; x <= cm.loadDistance; x++ {
		for z := -cm.loadDistance; z <= cm.loadDistance; z++ {
			pos := ChunkPos{X: centerPos.X + x, Z: centerPos.Z + z}
			// 无需重新载入，移除待卸载状态
			posId := pos.Id()
			if _, ok := cm.loadedChunkMap[posId]; ok {
				willUnload[pos.Id()] = false
			} else if chunk, ok := cm.UnloadingChunkMap[posId]; ok {
				// 从正在卸载的区块中恢复
				cm.loadedChunkMap[posId] = chunk.Chunk
				delete(cm.UnloadingChunkMap, posId)
			} else {
				// 加载新区块
				cm.loadingChunkMap[posId] = pos
			}

			chunk, ok := cm.loadedChunkMap[posId]
			if !ok {
				continue
			}

			// 渲染
			if x >= -cm.renderDistance && x <= cm.renderDistance && z >= -cm.renderDistance && z <= cm.renderDistance && chunk.inView {
				chunk.Rendered(a)
			} else {
				chunk.Unrendered()
			}
		}
	}

	// a.Log().Debug("unloading chunks: %d", len(cm.UnloadingChunkMap))
	// 更新卸载区块信息
	cm.StepUnloadChunk(a)

	// 处理这次将要卸载的区块
	for posId, unload := range willUnload {
		if !unload {
			continue
		}
		// a.Log().Debug("will unload chunk: %v", posId)

		cm.UnloadingChunkMap[posId] = &UnloadingChunk{Chunk: cm.loadedChunkMap[posId]}
		delete(cm.loadedChunkMap, posId)
	}

	// 更新中心区块
	cm.centerChunk = cm.loadedChunkMap[centerPos.Id()]
}

func (cm *ChunkManager) StepUnloadChunk(a *App) {
	var unloaded bool
	for key, uc := range cm.UnloadingChunkMap {
		if uc.AliveTick >= MAX_ALIVE_TICK && (!unloaded || uc.AliveTick >= 2*MAX_ALIVE_TICK) {
			// a.Log().Debug("unload chunk from mem: %v", uc.pos)
			a.SaveManager().SaveChunk(uc.Chunk)
			// 卸载区块
			cm.Remove(uc)
			uc.Cleanup()
			delete(cm.UnloadingChunkMap, key)

			unloaded = true
		} else {
			uc.AliveTick++
		}
	}
}

func (cm *ChunkManager) StepLoadChunk(a *App) {
	if len(cm.loadingChunkMap) == 0 {
		return
	}
	// a.Log().Debug("loading chunk count: %v", len(cm.loadingChunkMap))

	playerPos := a.Player().GetPosition()
	chunkPos := make([]ChunkPos, 0, len(cm.loadingChunkMap))
	for _, p := range cm.loadingChunkMap {
		chunkPos = append(chunkPos, p)
	}
	nearestPos := cm.GetNearestChunk(*playerPos, chunkPos)

	// a.Log().Debug("load chunk: %v", nearestPos)

	chunk := NewChunk(nearestPos.X, nearestPos.Z)
	chunk.Start(a)
	cm.Add(chunk)
	a.SaveManager().SaveChunk(chunk)
	cm.loadedChunkMap[nearestPos.Id()] = chunk
	delete(cm.loadingChunkMap, nearestPos.Id())

	a.World().bu.RefreshChunkBlocks(chunk)
	a.World().lu.AddWaitLumChunk(nearestPos)
}

func (cm *ChunkManager) Chunk(cpos ChunkPos) *Chunk {
	return cm.loadedChunkMap[cpos.Id()]
}

func (cm *ChunkManager) GetChunk(x, y, z float32) *Chunk {
	pos := ToChunkPos(math32.NewVector3(x, 0, z))
	return cm.Chunk(pos)
}

func (cm *ChunkManager) GetChunkByPos(pos util.Pos) *Chunk {
	p := pos.ToVec3()
	return cm.GetChunk(p.X, p.X, p.Z)
}

func (cm *ChunkManager) SaveAll() {
	for _, item := range cm.loadedChunkMap {
		Instance().SaveManager().SaveChunk(item)
	}

	for _, item := range cm.UnloadingChunkMap {
		Instance().SaveManager().SaveChunk(item.Chunk)
	}
}

func (cm *ChunkManager) GetNearestChunk(pos math32.Vector3, data []ChunkPos) ChunkPos {
	var nearestPos ChunkPos
	minDistance := float32(1000)
	for _, p := range data {
		d := SquareDistance(*math32.NewVector3(pos.X, 0, pos.Z), *math32.NewVector3(float32(p.X), 0, float32(p.Z)))
		if d < minDistance {
			minDistance = d
			nearestPos = p
		}
	}

	return nearestPos
}

func (cm *ChunkManager) UpdateChunkInView() {
	handler := cm.GetPlayerViewHanlder()
	for _, chunk := range cm.loadedChunkMap {
		chunk.inView = cm.ChunkInView(chunk, handler)
	}
}

func (cm *ChunkManager) ChunkInView(chunk *Chunk, handler func(p math32.Vector3) bool) bool {
	var inView bool
	chunk.RangePos(func(pos math32.Vector3) bool {
		if handler(pos) {
			inView = true
			return true
		}

		return false
	})

	return inView
}

func (cm *ChunkManager) GetPlayerViewHanlder() func(p math32.Vector3) bool {
	viewport := cm.app.player.GetViewport()
	lookAt := cm.app.player.farPos

	x := lookAt.X - viewport.X
	y := lookAt.Y - viewport.Y
	z := lookAt.Z - viewport.Z
	if x == 0 && z == 0 {
		return func(p math32.Vector3) bool {
			return true
		}
	}

	z1 := x * x / z
	y1 := x * x / y
	y2 := z * z / y

	var xzHander func(p math32.Vector3) bool
	if lookAt.Z == viewport.Z {
		xzHander = func(p math32.Vector3) bool {
			if x > 0 {
				return p.X >= viewport.X
			} else {
				return p.X < viewport.X
			}
		}
	} else if lookAt.X == viewport.X {
		xzHander = func(p math32.Vector3) bool {
			if z > 0 {
				return p.Z >= viewport.Z
			} else {
				return p.Z < viewport.Z
			}
		}
	} else {
		pos := math32.NewVector3(lookAt.X, 0, lookAt.Z-z-z1)
		a := (pos.Z - viewport.Z) / (pos.X - viewport.X)
		b := viewport.Z - a*viewport.X

		xzHander = func(p math32.Vector3) bool {
			pz := a*p.X + b
			if z > 0 {
				return p.Z > pz
			} else {
				return p.Z < pz
			}
		}
	}

	var yxHander func(p math32.Vector3) bool
	if lookAt.Y == viewport.Y {
		yxHander = func(p math32.Vector3) bool {
			if x > 0 {
				return p.X >= viewport.X
			} else {
				return p.X < viewport.X
			}
		}
	} else if lookAt.X == viewport.X {
		yxHander = func(p math32.Vector3) bool {
			if y > 0 {
				return p.Y >= viewport.Y
			} else {
				return p.Y < viewport.Y
			}
		}
	} else {
		pos := math32.NewVector3(lookAt.X, lookAt.Y-y-y1, 0)
		a := (pos.Y - viewport.Y) / (pos.X - viewport.X)
		b := viewport.Y - a*viewport.X

		yxHander = func(p math32.Vector3) bool {
			py := a*p.X + b
			if y > 0 {
				return p.Y > py
			} else {
				return p.Y < py
			}
		}
	}

	var yzHander func(p math32.Vector3) bool
	if lookAt.Y == viewport.Y {
		yzHander = func(p math32.Vector3) bool {
			if z > 0 {
				return p.Z >= viewport.Z
			} else {
				return p.Z < viewport.Z
			}
		}
	} else if lookAt.Z == viewport.Z {
		yzHander = func(p math32.Vector3) bool {
			if y > 0 {
				return p.Y >= viewport.Y
			} else {
				return p.Y < viewport.Y
			}
		}
	} else {
		pos := math32.NewVector3(0, lookAt.Y-y-y2, lookAt.Z)
		a := (pos.Y - viewport.Y) / (pos.Z - viewport.Z)
		b := viewport.Y - a*viewport.Z

		yzHander = func(p math32.Vector3) bool {
			py := a*p.Z + b
			if y > 0 {
				return p.Y > py
			} else {
				return p.Y < py
			}
		}
	}

	return func(p math32.Vector3) bool {
		return xzHander(p) || yxHander(p) || yzHander(p)
	}
}
