package app

import (
	"fmt"
	"math"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/weiWang95/mcworld/app/block"
)

const DEFAULT_VIEW_DISTANCE int64 = 2
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

	// playerAreaPos := GetAreaPosByPos(*a.Player().GetPosition())
	// if playerAreaPos.X != w.activePos.X || playerAreaPos.Z != w.activePos.Z {
	// 	w.loadAreas(playerAreaPos)
	// 	w.refreshAreas()
	// 	w.activePos = playerAreaPos
	// }

	// for i := 0; i < len(w.activeAreas); i++ {
	// 	for j := 0; j < len(w.activeAreas[i]); j++ {
	// 		w.activeAreas[i][j].Update()
	// 	}
	// }
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

	// w.activePos = math32.Vector3{}
	// w.loadAreas(w.activePos)
	// w.refreshAreas()
}

func (w *World) loadAreas(pos math32.Vector3) {
	curAreaX := int64(pos.X) //int64(pos.X) - (int64(pos.X) % AREA_WIDTH)
	curAreaZ := int64(pos.Z) //int64(pos.Z) - (int64(pos.Z) % AREA_WIDTH)

	if len(w.activeAreas) != 0 {
		shiftX := (curAreaX - int64(w.activePos.X)) / AREA_WIDTH
		shiftZ := (curAreaZ - int64(w.activePos.Z)) / AREA_WIDTH
		if shiftX != 0 || shiftZ != 0 {
			needLoads := make(map[string]interface{})

			w.Debug("area center pos: %v, cur pos: %v", w.activePos, pos)

			// 解除关联、移动区块位置、注销移除区块
			for i := int64(0); i < w.activeAreasLen(); i++ {
				for j := int64(0); j < w.activeAreasLen(); j++ {
					w.activeAreas[i][j].ClearRelations()

					if (i-shiftX < 0 || i-shiftX >= w.activeAreasLen()) || (j-shiftZ < 0 || j-shiftZ >= w.activeAreasLen()) {
						w.Debug("area disposed: %v", w.activeAreas[i][j].Pos)
						w.Remove(w.activeAreas[i][j])
						w.activeAreas[i][j].Dispose()
					}

					// 需要载入区块
					if i+shiftX < 0 || i+shiftX >= w.activeAreasLen() || j+shiftZ < 0 || j+shiftZ >= w.activeAreasLen() {
						needLoads[fmt.Sprintf("%d-%d", i, j)] = nil
						continue
					}

					// 无需重新载入区块
					w.activeAreas[i][j] = w.activeAreas[i+shiftX][j+shiftZ]
				}
			}

			// 载入新区块
			for i := int64(0); i < w.activeAreasLen(); i++ {
				for j := int64(0); j < w.activeAreasLen(); j++ {
					if _, ok := needLoads[fmt.Sprintf("%d-%d", i, j)]; ok {
						areaPos := math32.Vector3{
							X: float32(curAreaX + (i-DEFAULT_VIEW_DISTANCE)*AREA_WIDTH),
							Z: float32(curAreaZ + (j-DEFAULT_VIEW_DISTANCE)*AREA_WIDTH),
						}

						w.Debug("load area: i:%d,j:%d -> %v", i, j, areaPos)

						area := NewArea(areaPos)
						area.Load(w.wg)

						w.Add(area)
						w.activeAreas[i][j] = area
					}
				}
			}

			// 创建关联
			for i := int64(0); i < w.activeAreasLen(); i++ {
				for j := int64(0); j < w.activeAreasLen(); j++ {
					if i > 0 {
						w.activeAreas[i][j].NorthArea = w.activeAreas[i-1][j]
						w.activeAreas[i-1][j].SouthArea = w.activeAreas[i][j]
					}

					if j > 0 {
						w.activeAreas[i][j].WestArea = w.activeAreas[i][j-1]
						w.activeAreas[i][j-1].EastArea = w.activeAreas[i][j]
					}
				}
			}
		}

		return
	}

	areas := make([][]*Area, 0, w.activeAreasLen())
	for x := -DEFAULT_VIEW_DISTANCE; x < DEFAULT_VIEW_DISTANCE+1; x++ {
		as := make([]*Area, 0, w.activeAreasLen())
		for z := -DEFAULT_VIEW_DISTANCE; z < DEFAULT_VIEW_DISTANCE+1; z++ {
			areaPos := math32.Vector3{
				X: float32(curAreaX + x*AREA_WIDTH),
				Z: float32(curAreaZ + z*AREA_WIDTH),
			}
			area := NewArea(areaPos)
			area.Load(w.wg)

			if len(as) > 0 {
				area.NorthArea = as[len(as)-1]
				as[len(as)-1].SouthArea = area
			}
			if len(areas) > 0 {
				area.WestArea = areas[len(areas)-1][len(as)]
				areas[len(areas)-1][len(as)].EastArea = area
			}

			w.Add(area)

			as = append(as, area)
		}
		areas = append(areas, as)
	}
	w.activeAreas = areas
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
