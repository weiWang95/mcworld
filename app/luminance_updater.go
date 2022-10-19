package app

import (
	"time"

	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/app/blockv2"
	"github.com/weiWang95/mcworld/lib/util"
)

type LuminanceUpdater struct {
	app   *App
	world *World
	log   *logger.Logger

	waitLumMap   map[string]ChunkPos
	switchLumMap map[string]ChunkPos

	lumTicker       *TickChecker
	lumSwitchTicker *TickChecker
}

func NewLuminanceUpdater(app *App) *LuminanceUpdater {
	u := new(LuminanceUpdater)
	u.app = app
	u.world = app.World()
	u.log = app.Log()

	u.lumTicker = NewTickChecker(4)
	u.lumSwitchTicker = NewTickChecker(2)

	u.waitLumMap = make(map[string]ChunkPos)

	return u
}

func (u *LuminanceUpdater) Update(app *App, t time.Duration) {
	if u.lumTicker.Next(t) {
		u.StepInitLum()
	}
	if u.lumSwitchTicker.Next(t) {
		u.StepSwitchDayNight()
	}
}

func (u *LuminanceUpdater) AddWaitLumChunk(cpos ChunkPos) {
	u.waitLumMap[cpos.Id()] = cpos
}

func (u *LuminanceUpdater) StepInitLum() {
	if len(u.waitLumMap) == 0 {
		return
	}
	// Instance().Log().Debug("wait init lum chunks:%d", len(u.waitLumMap))

	for key, cpos := range u.waitLumMap {
		// Instance().Log().Debug("init lum chunk:%v", cpos)

		u.InitChunkLum(cpos)

		delete(u.waitLumMap, key)
		break
	}
}

func (u *LuminanceUpdater) InitChunkLum(cpos ChunkPos) {
	chunk := u.world.cm.Chunk(cpos)
	if chunk == nil {
		return
	}

	updateMap := u.initChunkSunLum(chunk)
	u.updateBlocksLum(updateMap)
}

func (u *LuminanceUpdater) initChunkSunLum(chunk *Chunk) map[string]util.Pos {
	updateMap := make(map[string]util.Pos, CHUNK_WIDTH*CHUNK_WIDTH*CHUNK_HEIGHT)

	for z := int64(0); z < CHUNK_WIDTH; z++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			sunLum := u.world.sunLevel

			for y := CHUNK_HEIGHT - 1; y >= 0; y-- {
				pos := util.NewPos(x, y, z)

				cur := chunk.GetLum(pos)
				cur = cur.SetSunLum(sunLum)

				b := chunk.getBlockByPos(pos)
				if b != nil {
					cur = cur.SetSunLum(0)
				}
				chunk.SetLum(pos, cur)

				wPos := chunk.GetWorldPos(x, y, z)
				updateMap[wPos.GetId()] = wPos
			}
		}
	}

	return updateMap
}

func (u *LuminanceUpdater) UpdateSumLum() {
	u.switchLumMap = map[string]ChunkPos{}
	for _, c := range u.world.cm.loadedChunkMap {
		u.AddWaitLumChunk(*c.pos)
	}
}

func (u *LuminanceUpdater) SwitchDayNight() {
	u.switchLumMap = map[string]ChunkPos{}
	for _, c := range u.world.cm.loadedChunkMap {
		u.switchLumMap[c.pos.Id()] = *c.pos
	}
}

func (u *LuminanceUpdater) StepSwitchDayNight() {
	if len(u.switchLumMap) == 0 {
		return
	}
	// Instance().Log().Debug("wait switch lum chunks:%d", len(u.switchLumMap))

	playerPos := u.app.Player().GetPosition()
	chunkPos := make([]ChunkPos, 0, len(u.switchLumMap))
	for _, p := range u.switchLumMap {
		chunkPos = append(chunkPos, p)
	}
	nearestPos := u.world.cm.GetNearestChunk(*playerPos, chunkPos)
	for key, cpos := range u.switchLumMap {
		// Instance().Log().Debug("switch lum chunk:%v", cpos)
		if cpos.X != nearestPos.X || cpos.Z != nearestPos.Z {
			continue
		}

		chunk := u.world.cm.Chunk(cpos)
		if chunk != nil {
			u.refreshChunkLum(chunk)
		}

		delete(u.switchLumMap, key)
		break
	}
}

func (u *LuminanceUpdater) refreshChunkLum(chunk *Chunk) {
	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			for z := int64(0); z < CHUNK_WIDTH; z++ {
				u.refreshBlockLum(chunk.GetWorldPos(x, y, z))
			}
		}
	}
}

func (u *LuminanceUpdater) TiggerUpdate(pos util.Pos) {
	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if !loaded {
		return
	}
	hasBlock := b != nil
	if hasBlock {
		u.setLum(pos, u.getLum(pos).SetSunLum(0))
	}

	// 阳光直射
	isBeat := true
	if pos.Y+1 < CHUNK_HEIGHT {
		topLum, _ := u.world.GetLumByVec(pos.AddY(1).ToVec3())
		isBeat = topLum.SunLum() == MAX_LUM
	}

	updates := make(map[string]util.Pos)
	updates[pos.GetId()] = pos
	if pos.Y+1 < CHUNK_HEIGHT {
		topPos := pos.AddY(1)
		updates[topPos.GetId()] = topPos
	}
	if pos.Y-1 >= 0 {
		bottomPos := pos.SubY(1)
		updates[bottomPos.GetId()] = bottomPos
	}

	for {
		pos = pos.SubY(1)

		updates[pos.GetId()] = pos
		b, _ := u.world.GetBlockByVec(pos.ToVec3())
		if b != nil {
			break
		}

		// 更新垂直方向方块阳光强度
		if isBeat && hasBlock {
			u.setLum(pos, u.getLum(pos).SetSunLum(0))
		} else if isBeat && !hasBlock {
			u.setLum(pos, u.getLum(pos).SetSunLum(MAX_LUM))
		}

		if pos.Y-1 <= 0 {
			break
		}
	}
	u.updateBlocksLum(updates)
}

func (u *LuminanceUpdater) updateBlocksLum(updateMap map[string]util.Pos) {
	refreshBlockMap := make(map[string]util.Pos)

	for i := 0; i < 16; i++ {
		newMap := make(map[string]util.Pos)
		for _, p := range updateMap {
			for _, item := range u.affectedBlocks(p) {
				refreshBlockMap[item.GetId()] = item
			}

			for _, item := range u.updateLum(p, i) {
				newMap[item.GetId()] = item
			}
		}

		if len(newMap) == 0 {
			break
		}

		updateMap = newMap
	}

	for _, item := range refreshBlockMap {
		u.refreshBlockLum(item)
	}
}

func (u *LuminanceUpdater) updateLum(pos util.Pos, times int) []util.Pos {
	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if !loaded {
		return nil
	}
	cur := u.getLum(pos)
	oldLum := cur

	if b != nil {
		cur = NewLuminance(0, b.GetBlockLum())
		u.setLum(pos, cur)
	} else {
		max := u.getNearbyMaxLum(pos)
		maxSunLum, maxBlockLum := max.SunLum(), max.BlockLum()
		if cur.SunLum() != MAX_LUM {
			if maxSunLum > 0 {
				cur = cur.SetSunLum(maxSunLum - 1)
			} else {
				cur = cur.SetSunLum(0)
			}
			u.setLum(pos, cur)
		}

		if maxBlockLum > 0 {
			cur = cur.SetBlockLum(maxBlockLum - 1)
		} else {
			cur = cur.SetBlockLum(0)
		}
		u.setLum(pos, cur)
	}

	if (b == nil || !b.GetLumable()) && oldLum == cur {
		return nil
	}

	updates := make([]util.Pos, 0)
	pos.RangeAdjoin(func(p util.Pos, face block.BlockFace) {
		if u.needUpdateLum(p, cur) {
			updates = append(updates, p)
		}
	})

	return updates
}

func (u *LuminanceUpdater) setLum(pos util.Pos, lum Luminance) {
	chunk := u.world.cm.GetChunkByPos(pos)
	if chunk == nil {
		return
	}

	chunk.SetLum(chunk.ConvertChunkPos(pos), lum)
}

func (u *LuminanceUpdater) getLum(pos util.Pos) Luminance {
	lum, loaded := u.world.GetLumByVec(pos.ToVec3())
	if !loaded {
		return NewLuminance(0, 0)
	}

	return lum
}

func (u *LuminanceUpdater) getNearbyMaxLum(pos util.Pos) Luminance {
	var max Luminance

	pos.RangeAdjoin(func(p util.Pos, face block.BlockFace) {
		max = MaxLum(u.getLum(p), max)
	})

	return max
}

func (u *LuminanceUpdater) needUpdateLum(pos util.Pos, l Luminance) bool {
	if PosOverRange(pos) {
		return false
	}

	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if !loaded || b != nil {
		return false
	}

	lum, loaded := u.world.GetLumByVec(pos.ToVec3())
	if !loaded {
		return false
	}

	return (l.SunLum() != lum.SunLum() || l.SunLum() != 0) || (l.BlockLum() != lum.BlockLum())
}

func (u *LuminanceUpdater) affectedBlocks(pos util.Pos) []util.Pos {
	arr := make([]util.Pos, 0)
	_, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if loaded {
		arr = append(arr, pos)
	}

	pos.RangeAdjoin(func(p util.Pos, face block.BlockFace) {
		_, loaded := u.world.GetBlockByVec(p.ToVec3())
		if loaded && !PosOverRange(p) {
			arr = append(arr, p)
		}
	})

	return arr
}

func (u *LuminanceUpdater) refreshBlockLum(pos util.Pos) {
	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if !loaded || b == nil || !b.Visible() {
		return
	}

	pos.RangeAdjoin(func(p util.Pos, face block.BlockFace) {
		faceBlock, loaded := u.world.GetBlockByVec(p.ToVec3())
		if loaded && faceBlock == nil {
			// b.SetLum(u.CurLum(u.getLum(p)), int(face))
			b.SetFaceLum(blockv2.BlockFace(face), u.CurLum(u.getLum(p)))
		}
	})
	// b.RefreshLum()
}

func (u *LuminanceUpdater) CurLum(l Luminance) uint8 {
	sun := uint8(float32(l.SunLum()) * u.world.SunLumRate())
	if sun > l.BlockLum() {
		return sun
	}

	return l.BlockLum()
}
