package app

import "github.com/g3n/engine/gui"

type PlayerGui struct {
	gui.Panel

	app *App

	inventory *GuiPlayerInventory
}

func NewPlayerGui(app *App) *PlayerGui {
	g := new(PlayerGui)
	g.app = app
	g.init()
	return g
}

func (g *PlayerGui) init() {
	w, h := g.app.GetSize()
	g.Panel = *gui.NewPanel(float32(w), float32(h))

	g.inventory = NewGuiPlayerInventory(g.app)
	g.Add(g.inventory)
}
