package app

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
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
	loadTicker *TickChecker
}

var _ IRender = (*ChunkManager)(nil)

func NewChunkManager() *ChunkManager {
	cm := new(ChunkManager)

	cm.Node = *core.NewNode()

	cm.loadDistance = 3
	cm.renderDistance = 1

	cm.loadingChunkMap = make(map[ChunkPos]interface{})
	cm.loadedChunkMap = make(map[ChunkPos]*Chunk)
	cm.UnloadingChunkMap = make(map[ChunkPos]*UnloadingChunk)

	cm.loadTicker = NewTickChecker(4)

	return cm
}

func (cm *ChunkManager) Start(a *App) {
	cm.setup(a)
}

func (cm *ChunkManager) Update(a *App, t time.Duration) {
	if cm.loadTicker.Next(t) {
		cm.checkAndLoadChunks(a, a.Player().GetPosition())
	} else {
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
	if cm.centerChunk != nil && cm.centerChunk.pos.X == centerPos.X && cm.centerChunk.pos.Z == centerPos.Z {
		// 更新卸载区块信息
		for key, uc := range cm.UnloadingChunkMap {
			if uc.AliveTick >= 10 {
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

	// 更新卸载区块信息
	for key, uc := range cm.UnloadingChunkMap {
		if uc.AliveTick >= 10 {
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

	var loadedPos ChunkPos
	for pos, _ := range cm.loadingChunkMap {
		loadedPos = pos

		chunk := NewChunk(pos.X, pos.Z)
		chunk.Start(a)
		cm.Add(chunk)
		cm.loadedChunkMap[pos] = chunk
		break
	}

	delete(cm.loadingChunkMap, loadedPos)
}

func (cm *ChunkManager) GetChunk(x, y, z float32) *Chunk {
	pos := ToChunkPos(math32.NewVector3(x, 0, z))
	return cm.loadedChunkMap[pos]
}
