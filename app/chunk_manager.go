package app

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/lib/util"
)

type UnloadingChunk struct {
	*Chunk

	AliveTick int64
}

type ChunkManager struct {
	core.Node
	*logger.Logger

	loadDistance   int64
	renderDistance int64

	centerChunk       *Chunk
	loadingChunkMap   map[ChunkPos]interface{}
	loadedChunkMap    map[ChunkPos]*Chunk
	UnloadingChunkMap map[ChunkPos]*UnloadingChunk

	// ticker
	loadTicker    *TickChecker
	loadingTicker *TickChecker
}

var _ IRender = (*ChunkManager)(nil)

func NewChunkManager() *ChunkManager {
	cm := new(ChunkManager)

	cm.Node = *core.NewNode()

	cm.loadDistance = 1
	cm.renderDistance = 1

	cm.loadingChunkMap = make(map[ChunkPos]interface{})
	cm.loadedChunkMap = make(map[ChunkPos]*Chunk)
	cm.UnloadingChunkMap = make(map[ChunkPos]*UnloadingChunk)

	cm.loadTicker = NewTickChecker(4)
	cm.loadingTicker = NewTickChecker(1)

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
	// for k, _ := range cm.loadedChunkMap {
	// 	a.Log().Debug("X:%d Z:%d", k.X, k.Z)
	// }

	if cm.centerChunk != nil && cm.centerChunk.pos.X == centerPos.X && cm.centerChunk.pos.Z == centerPos.Z {
		// 更新卸载区块信息
		// a.Log().Debug("unloading chunks: %d", len(cm.UnloadingChunkMap))
		for key, uc := range cm.UnloadingChunkMap {
			if uc.AliveTick >= 10 {
				a.Log().Debug("unload chunk from mem: %v", uc.pos)
				a.SaveManager().SaveChunk(uc.Chunk)
				// 卸载区块
				cm.Remove(uc)
				uc.Cleanup()
				delete(cm.UnloadingChunkMap, key)
			} else {
				uc.AliveTick++
			}
		}
		return
	}

	// 将所有加载的区块标记为待卸载
	willUnload := make(map[ChunkPos]bool, len(cm.loadedChunkMap))
	for k, _ := range cm.loadedChunkMap {
		willUnload[k] = true
	}

	for x := -cm.loadDistance; x <= cm.loadDistance; x++ {
		for z := -cm.loadDistance; z <= cm.loadDistance; z++ {
			pos := ChunkPos{X: centerPos.X + x, Z: centerPos.Z + z}
			// 无需重新载入，移除待卸载状态
			if _, ok := cm.loadedChunkMap[pos]; ok {
				willUnload[pos] = false
			} else if chunk, ok := cm.UnloadingChunkMap[pos]; ok {
				// 从正在卸载的区块中恢复
				cm.loadedChunkMap[pos] = chunk.Chunk
				delete(cm.UnloadingChunkMap, pos)
			} else {
				// 加载新区块
				cm.loadingChunkMap[pos] = nil
			}

			chunk, ok := cm.loadedChunkMap[pos]
			if !ok {
				continue
			}

			// 渲染
			if x >= -cm.renderDistance && x <= cm.renderDistance && z >= -cm.renderDistance && z <= cm.renderDistance {
				chunk.Rendered(a)
			} else {
				chunk.Unrendered()
			}
		}
	}

	// a.Log().Debug("unloading chunks: %d", len(cm.UnloadingChunkMap))
	// 更新卸载区块信息
	for key, uc := range cm.UnloadingChunkMap {
		if uc.AliveTick >= 10 {
			a.Log().Debug("unload chunk from mem: %v", uc.pos)
			a.SaveManager().SaveChunk(uc.Chunk)
			// 卸载区块
			cm.Remove(uc)
			uc.Cleanup()
			delete(cm.UnloadingChunkMap, key)
		} else {
			uc.AliveTick++
		}
	}

	// 处理这次将要卸载的区块
	for pos, unload := range willUnload {
		if !unload {
			continue
		}
		a.Log().Debug("will unload chunk: %v", pos)

		cm.UnloadingChunkMap[pos] = &UnloadingChunk{Chunk: cm.loadedChunkMap[pos]}
		delete(cm.loadedChunkMap, pos)
	}

	// 更新中心区块
	cm.centerChunk = cm.loadedChunkMap[centerPos]
}

func (cm *ChunkManager) StepLoadChunk(a *App) {
	if len(cm.loadingChunkMap) == 0 {
		return
	}
	a.Log().Debug("loading chunk count: %v", len(cm.loadingChunkMap))

	playerPos := a.Player().GetPosition()
	chunkPos := make([]ChunkPos, 0, len(cm.loadingChunkMap))
	for p, _ := range cm.loadingChunkMap {
		chunkPos = append(chunkPos, p)
	}
	nearestPos := cm.GetNearestChunk(*playerPos, chunkPos)

	a.Log().Debug("load chunk: %v", nearestPos)

	chunk := NewChunk(nearestPos.X, nearestPos.Z)
	chunk.Start(a)
	cm.Add(chunk)
	a.SaveManager().SaveChunk(chunk)
	cm.loadedChunkMap[nearestPos] = chunk
	delete(cm.loadingChunkMap, nearestPos)

	a.World().bu.RefreshChunkBlocks(chunk)
	a.World().lu.AddWaitLumChunk(nearestPos)
}

func (cm *ChunkManager) Chunk(cpos ChunkPos) *Chunk {
	return cm.loadedChunkMap[cpos]
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
