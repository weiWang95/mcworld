package world

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app"
)

const DEFAULT_VIEW_DISTANCE int64 = 3

func init() {
	app.RegisterLevel("world", &World{})
}

type World struct {
	core.Node

	wg          IWorldGenerator
	activeAreas [][]*Area
}

func (d *World) Start(a *app.App) {
	d.Node = *core.NewNode()

	d.setup(a)
	a.Scene().Add(d)
}

func (d *World) Update(a *app.App, t time.Duration) {
	for i := 0; i < len(d.activeAreas); i++ {
		for j := 0; j < len(d.activeAreas[i]); j++ {
			d.activeAreas[i][j].Update()
		}
	}
}

func (d *World) Cleanup(a *app.App) {
}

func (d *World) setup(a *app.App) {
	seed := time.Now().UnixNano()
	d.wg = &WorldGenerator{}
	d.wg.Setup(seed)

	d.loadAreas(math32.Vector3{})
}

func (d *World) loadAreas(pos math32.Vector3) {
	curAreaX := int64(pos.X) - (int64(pos.X) % AREA_WIDTH)
	curAreaZ := int64(pos.Z) - (int64(pos.Z) % AREA_WIDTH)

	areas := make([][]*Area, 0, 2*DEFAULT_VIEW_DISTANCE+1)
	for x := -DEFAULT_VIEW_DISTANCE; x < DEFAULT_VIEW_DISTANCE+1; x++ {
		as := make([]*Area, 0, 2*DEFAULT_VIEW_DISTANCE+1)
		for z := -DEFAULT_VIEW_DISTANCE; z < DEFAULT_VIEW_DISTANCE+1; z++ {
			areaPos := math32.Vector3{
				X: float32(curAreaX + x),
				Z: float32(curAreaZ + z),
			}
			area := NewArea(areaPos)
			area.Load(d.wg)
			d.Add(area)

			as = append(as, area)
		}
		areas = append(areas, as)
	}
	d.activeAreas = areas
}
