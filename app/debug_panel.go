package app

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/app/blockv2"
)

var transparentColor = math32.Color4{0, 0, 0, 0.1}
var headerColor = math32.Color4{0, 0, 0, 0.1}
var lightTextColor = math32.Color4{0.8, 0.8, 0.8, 1}

const fontSize = 14

type DebugPanel struct {
	*gui.Panel

	app *App

	chunk    *gui.Label
	pos      *gui.Label
	viewPort *gui.Label
	camera   *gui.Label
	farPos   *gui.Label
	target   *gui.Label
}

func NewDebugPanel(a *App) *DebugPanel {
	p := new(DebugPanel)
	p.app = a

	p.init()

	return p
}

func (p *DebugPanel) GetPanel() gui.IPanel {
	return p.Panel
}

func (p *DebugPanel) init() {
	panel := gui.NewPanel(100, 200)
	panel.SetPaddings(0, 0, 0, 0)
	panel.SetColor4(&transparentColor)
	panel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})
	panel.SetLayout(gui.NewDockLayout())

	// Position
	p0 := newDefaultPanel()
	p0.Add(newDefaultLabel("Chunk:"))
	p.chunk = newDefaultLabel(" ")
	p0.Add(p.chunk)
	panel.Add(p0)

	// Position
	p1 := newDefaultPanel()
	p1.Add(newDefaultLabel("Pos:"))
	p.pos = newDefaultLabel(" ")
	p1.Add(p.pos)
	panel.Add(p1)

	// Viewport
	p2 := newDefaultPanel()
	p2.Add(newDefaultLabel("Viewport:"))
	p.viewPort = newDefaultLabel(" ")
	p2.Add(p.viewPort)
	panel.Add(p2)

	// Viewport
	p3 := newDefaultPanel()
	p3.Add(newDefaultLabel("Camera:"))
	p.camera = newDefaultLabel(" ")
	p3.Add(p.camera)
	panel.Add(p3)

	// FarPos
	p4 := newDefaultPanel()
	p4.Add(newDefaultLabel("FarPos:"))
	p.farPos = newDefaultLabel(" ")
	p4.Add(p.farPos)
	panel.Add(p4)

	// Target
	p5 := newDefaultPanel()
	p5.Add(newDefaultLabel("Target:"))
	p.target = newDefaultLabel(" ")
	p5.Add(p.target)
	panel.Add(p5)

	p.Panel = panel
}

func (p *DebugPanel) update() {
	player := p.app.Player()
	pos := player.GetPosition()
	lum, _ := p.app.World().GetLum(pos.X, pos.Y, pos.Z)

	p.chunk.SetText(fmt.Sprintf("R:%d U:%d", p.app.curWorld.cm.renderedCount, p.app.curWorld.cm.unrenderedCount))
	p.pos.SetText(fmt.Sprintf("%s %s", p.formatPos(*pos), p.formatLum(lum)))
	p.viewPort.SetText(p.formatPos(*player.GetViewport()))
	p.camera.SetText(p.formatPos(player.Camera.Position()))
	p.farPos.SetText(p.formatPos(player.farPos))

	targetPos := player.Target.Position()
	targetLum, _ := p.app.World().GetLum(targetPos.X, targetPos.Y, targetPos.Z)
	targetTopLum, _ := p.app.World().GetLum(targetPos.X, targetPos.Y+1, targetPos.Z)
	p.target.SetText(fmt.Sprintf("T: %T V:%v P:%s %s %s, PT: %s", player.Target.b, player.Target.Visible(), p.formatPos(targetPos), p.formatLum(targetLum), p.formatFaceLum(player.Target.b), p.formatLum(targetTopLum)))
}

func (p *DebugPanel) formatPos(pos math32.Vector3) string {
	return fmt.Sprintf("X: %.1f, Y: %.1f, Z: %.1f", pos.X, pos.Y, pos.Z)
}

func (p *DebugPanel) formatLum(lum Luminance) string {
	return fmt.Sprintf("Lum:[S:%v, B:%v]", lum.SunLum(), lum.BlockLum())
}

func (p *DebugPanel) formatFaceLum(b *blockv2.Block) string {
	if b == nil {
		return ""
	}

	return fmt.Sprintf(
		"Face Lum:[F:%d B:%d L:%d R:%d T:%d B:%d]",
		b.GetFaceLum(int(block.BlockFaceFront)),
		b.GetFaceLum(int(block.BlockFaceBack)),
		b.GetFaceLum(int(block.BlockFaceLeft)),
		b.GetFaceLum(int(block.BlockFaceRight)),
		b.GetFaceLum(int(block.BlockFaceTop)),
		b.GetFaceLum(int(block.BlockFaceBottom)),
	)
}

func newDefaultPanel() *gui.Panel {
	panel := gui.NewPanel(200, 20)
	panel.SetPaddings(4, 4, 4, 4)
	panel.SetColor4(&headerColor)
	panel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})
	panel.SetLayout(gui.NewHBoxLayout())
	return panel
}

func newDefaultLabel(v interface{}) *gui.Label {
	label := gui.NewLabel(" ")
	label.SetFontSize(fontSize)
	label.SetLayoutParams(&gui.HBoxLayoutParams{AlignV: gui.AlignCenter})
	label.SetText(fmt.Sprint(v))
	label.SetColor4(&lightTextColor)
	return label
}
