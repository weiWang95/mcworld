package app

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type GuiPlayerInventory struct {
	gui.Panel

	app *App

	activeIndex int
}

func NewGuiPlayerInventory(app *App) *GuiPlayerInventory {
	g := new(GuiPlayerInventory)
	g.app = app
	g.init()
	return g
}

func (g *GuiPlayerInventory) init() {
	w, h := g.app.GetSize()

	g.Panel = *gui.NewPanel(400, 40)
	g.SetLayout(gui.NewGridLayout(10))
	g.SetBordersColor(math32.NewColor("grey"))
	g.SetBorders(4, 4, 4, 4)
	g.SetPosition(float32(w)/2-g.Width()/2, float32(h)-g.Height())

	l1 := gui.NewImageLabel("")
	l1.SetImageFromFile(g.app.DirData() + "/images/blocks/2_0.jpg")
	g.Add(l1)
}

func (g *GuiPlayerInventory) Switch(idx int) {

}
