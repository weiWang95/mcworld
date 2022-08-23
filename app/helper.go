package app

import (
	"math"

	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/lib/util"
)

// 获取射线与X平面焦点
func GetIntermediateWithX(start, end math32.Vector3, x float32) *math32.Vector3 {
	return GetIntermediate(start, end, &x, nil, nil)
}

// 获取射线与Y平面焦点
func GetIntermediateWithY(start, end math32.Vector3, y float32) *math32.Vector3 {
	return GetIntermediate(start, end, nil, &y, nil)
}

// 获取射线与Z平面焦点
func GetIntermediateWithZ(start, end math32.Vector3, z float32) *math32.Vector3 {
	return GetIntermediate(start, end, nil, nil, &z)
}

// 获取射线与平面焦点
func GetIntermediate(start, end math32.Vector3, x, y, z *float32) *math32.Vector3 {
	dx := end.X - start.X
	dy := end.Y - start.Y
	dz := end.Z - start.Z

	var scale float32

	if x != nil {
		if dx*dx < 1 {
			return nil
		}

		scale = (*x - start.X) / dx
	} else if y != nil {
		if dy*dy < 1 {
			return nil
		}

		scale = (*y - start.Y) / dy
	} else if z != nil {
		if dz*dz < 1 {
			return nil
		}

		scale = (*z - start.Z) / dz
	} else {
		return nil
	}

	if scale < 0 || scale > 1 {
		return nil
	}

	return math32.NewVector3(
		start.X+dx*scale,
		start.Y+dy*scale,
		start.Z+dz*scale,
	)
}

type BoundBox struct {
	X  float32
	Y  float32
	Z  float32
	BX float32
	BY float32
	BZ float32
}

func SquareDistance(from, to math32.Vector3) float32 {
	dx := to.X - from.X
	dy := to.Y - from.Y
	dz := to.Z - from.Z

	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// 坐标是否位于碰撞盒内
func (b *BoundBox) Inside(pos math32.Vector3) bool {
	if pos.X >= b.X && pos.X <= b.X+b.BX &&
		pos.Y >= b.Y && pos.Y <= b.Y+b.BY &&
		pos.Z >= b.Z && pos.Z <= b.Z+b.BZ {
		return true
	}

	return false
}

func NewBlockBoundBox(x, y, z int64) *BoundBox {
	return &BoundBox{
		X:  float32(x),
		Y:  float32(y),
		Z:  float32(z),
		BX: 1, BY: 1, BZ: 1,
	}
}

func CollisionRayTrace(box *BoundBox, start, end math32.Vector3) *math32.Vector3 {
	// Instance().log.Debug("start CollisionRayTrace -> %v, %v, %v", box, start, end)
	// // 以碰撞盒坐标为原点
	// newStart := start.Add(math32.NewVector3(-box.X, -box.Y, -box.Z))
	// newEnd := end.Add(math32.NewVector3(-box.X, -box.Y, -box.Z))

	// // 计算start到end与碰撞盒各平面交点
	// yz1 := GetIntermediateWithX(*newStart, *newEnd, box.X)
	// yz2 := GetIntermediateWithX(*newStart, *newEnd, box.X+box.BX)

	// xz1 := GetIntermediateWithY(*newStart, *newEnd, box.Y)
	// xz2 := GetIntermediateWithY(*newStart, *newEnd, box.Y+box.BY)

	// xy1 := GetIntermediateWithZ(*newStart, *newEnd, box.Z)
	// xy2 := GetIntermediateWithZ(*newStart, *newEnd, box.Z+box.BZ)

	yz1 := GetIntermediateWithX(start, end, box.X)
	yz2 := GetIntermediateWithX(start, end, box.X+box.BX)

	xz1 := GetIntermediateWithY(start, end, box.Y)
	xz2 := GetIntermediateWithY(start, end, box.Y+box.BY)

	xy1 := GetIntermediateWithZ(start, end, box.Z)
	xy2 := GetIntermediateWithZ(start, end, box.Z+box.BZ)

	// 交点不在碰撞盒内
	if yz1 == nil || !box.Inside(*yz1) {
		yz1 = nil
	}
	if yz2 == nil || !box.Inside(*yz2) {
		yz2 = nil
	}
	if xz1 == nil || !box.Inside(*xz1) {
		xz1 = nil
	}
	if xz2 == nil || !box.Inside(*xz2) {
		xz2 = nil
	}
	if xy1 == nil || !box.Inside(*xy1) {
		xy1 = nil
	}
	if xy2 == nil || !box.Inside(*xy2) {
		xy2 = nil
	}

	var hitPos *math32.Vector3

	if yz1 != nil && (hitPos == nil || SquareDistance(start, *yz1) < SquareDistance(start, *hitPos)) {
		hitPos = yz1
	}

	if yz2 != nil && (hitPos == nil || SquareDistance(start, *yz2) < SquareDistance(start, *hitPos)) {
		hitPos = yz2
	}

	if xz1 != nil && (hitPos == nil || SquareDistance(start, *xz1) < SquareDistance(start, *hitPos)) {
		hitPos = xz1
	}

	if xz2 != nil && (hitPos == nil || SquareDistance(start, *xz2) < SquareDistance(start, *hitPos)) {
		hitPos = xz2
	}

	if xy1 != nil && (hitPos == nil || SquareDistance(start, *xy1) < SquareDistance(start, *hitPos)) {
		hitPos = xy1
	}

	if xy2 != nil && (hitPos == nil || SquareDistance(start, *xy2) < SquareDistance(start, *hitPos)) {
		hitPos = xy2
	}

	Instance().log.Debug("hitPos -> %v", hitPos)

	return hitPos
}

func RayTraceBlock(world *World, start, end math32.Vector3) block.IBlock {
	// Instance().log.Debug("start ray trace block! start:%v, end:%v", start, end)

	startX, startY, startZ := util.FloorFloat(start.X), util.FloorFloat(start.Y), util.FloorFloat(start.Z)
	endX, endY, endZ := util.FloorFloat(end.X), util.FloorFloat(end.Y), util.FloorFloat(end.Z)

	for i := 200; i >= 0; i-- {
		// Instance().log.Debug("start check block -> %v, %v, %v", startX, startY, startZ)
		// 检测到终点方块
		if startX == endX && startY == endY && startZ == endZ {
			return world.GetBlockByPosition(end.X, end.Y, end.Z)
		}

		xChanged, yChanged, zChanged := true, true, true
		var newX, newY, newZ float32

		if endX > startX {
			newX = float32(startX) + 1
		} else if endX < startX {
			newX = float32(startX)
		} else {
			xChanged = false
		}

		if endY > startY {
			newY = float32(startY) + 1
		} else if endY < startY {
			newY = float32(startY)
		} else {
			yChanged = false
		}

		if endZ > startZ {
			newZ = float32(startZ) + 1
		} else if endZ < startZ {
			newZ = float32(startZ)
		} else {
			zChanged = false
		}

		xt, yt, zt := float32(999.0), float32(999.0), float32(999.0)
		dx := end.X - start.X
		dy := end.Y - start.Y
		dz := end.Z - start.Z

		if xChanged {
			xt = (newX - start.X) / dx
		}

		if yChanged {
			yt = (newY - start.Y) / dy
		}

		if zChanged {
			zt = (newZ - start.Z) / dz
		}

		d := 0

		if xt < yt && xt < zt {
			if endX > startX {
				d = 4
			} else {
				d = 5
			}

			start.X = newX
			start.Y += dy * xt
			start.Z += dz * xt
		} else if yt < zt {
			if endY > startY {
				d = 0
			} else {
				d = 1
			}
			start.X += dx * yt
			start.Y = newY
			start.Z += dz * yt
		} else {
			if endZ > startZ {
				d = 2
			} else {
				d = 3
			}
			start.X += dx * zt
			start.Y += dy * zt
			start.Z = newZ
		}

		startX = util.FloorFloat(start.X)
		startY = util.FloorFloat(start.Y)
		startZ = util.FloorFloat(start.Z)
		switch d {
		case 5: // X- 方向
			startX -= 1
		case 1: // Y- 方向
			startY -= 1
		case 3: // Z- 方向
			startZ -= 1
		}

		block := world.GetBlockByPosition(float32(startX), float32(startY), float32(startZ))
		if block == nil {
			continue
		}

		box := NewBlockBoundBox(startX, startY, startZ)
		pos := CollisionRayTrace(box, start, end)
		if pos != nil {
			return block
		}
	}

	return nil
}
