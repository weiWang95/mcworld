package app

import (
	"time"

	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/app/blockv2"
	"github.com/weiWang95/mcworld/lib/util"
)

type BlockUpdater struct {
	app   *App
	world *World
	log   *logger.Logger
}

func NewBlockUpdater(app *App) *BlockUpdater {
	u := new(BlockUpdater)
	u.app = app
	u.world = app.World()
	u.log = app.Log()
	return u
}

func (u *BlockUpdater) Update(a *App, t time.Duration) {

}

func (u *BlockUpdater) RefreshChunkBlocks(chunk *Chunk) {
	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			for z := int64(0); z < CHUNK_WIDTH; z++ {
				pos := chunk.GetWorldPos(x, y, z)
				u.updateBlock(pos)

				if x == 0 {
					u.updateBlock(pos.SubX(1))
				} else if x == CHUNK_WIDTH-1 {
					u.updateBlock(pos.AddX(1))
				}

				if z == 0 {
					u.updateBlock(pos.SubZ(1))
				} else if z == CHUNK_WIDTH-1 {
					u.updateBlock(pos.AddZ(1))
				}
			}
		}
	}

	// chunk.Rendered(u.app)
}

func (u *BlockUpdater) TiggerUpdate(tiggerPos util.Pos) {
	updatePos := []util.Pos{
		tiggerPos,
		tiggerPos.SubX(1), tiggerPos.AddX(1),
		tiggerPos.SubY(1), tiggerPos.AddY(1),
		tiggerPos.SubZ(1), tiggerPos.AddZ(1),
	}

	for _, pos := range updatePos {
		u.updateBlock(pos)
	}
}

func (u *BlockUpdater) updateBlock(pos util.Pos) bool {
	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	if !loaded || b == nil {
		return false
	}

	leftVisible := !u.BlockExist(pos.SubX(1))
	b.SetFaceVisible(blockv2.BlockFaceLeft, leftVisible)

	rightVisible := !u.BlockExist(pos.AddX(1))
	b.SetFaceVisible(blockv2.BlockFaceRight, rightVisible)

	frontVisible := !u.BlockExist(pos.SubZ(1))
	b.SetFaceVisible(blockv2.BlockFaceFront, frontVisible)

	backVisible := !u.BlockExist(pos.AddZ(1))
	b.SetFaceVisible(blockv2.BlockFaceBack, backVisible)

	bottomVisible := !u.BlockExist(pos.SubY(1))
	b.SetFaceVisible(blockv2.BlockFaceBottom, bottomVisible)

	topVisible := !u.BlockExist(pos.AddY(1))
	b.SetFaceVisible(blockv2.BlockFaceTop, topVisible)

	// if !u.BlockExist(pos.SubX(1)) || !u.BlockExist(pos.AddX(1)) ||
	// 	!u.BlockExist(pos.SubY(1)) || !u.BlockExist(pos.AddY(1)) ||
	// 	!u.BlockExist(pos.SubZ(1)) || !u.BlockExist(pos.AddZ(1)) {
	// 	b.SetVisible(true)
	// 	return true
	// }

	visible := leftVisible || rightVisible || frontVisible || backVisible || bottomVisible || topVisible

	b.SetVisible(visible)
	return visible
}

func (u *BlockUpdater) BlockExist(pos util.Pos) bool {
	b, loaded := u.world.GetBlockByVec(pos.ToVec3())
	return !loaded || b != nil || pos.Y < 0
}
