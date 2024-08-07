package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/weiWang95/mcworld/app/blockv2"
)

var instance *App

type App struct {
	*app.Application

	log      *logger.Logger
	dirData  string // Full path of the data directory
	scene    *core.Node
	curWorld *World
	sm       ISaveManager
	bm       *blockv2.BlockManager

	seed int64

	grid       *helper.Grid
	frameRater *util.FrameRater // Render loop frame rater

	// GUI
	mainPanel  *gui.Panel
	labelFPS   *gui.Label // header FPS label
	debugPanel *DebugPanel
	cursor     *gui.Panel
	playerGui  *PlayerGui

	// OldPlayer
	// player *OldPlayer
	player *Player

	debugMode bool
}

func Instance() *App {
	return instance
}

func Create() *App {
	if instance != nil {
		return instance
	}

	a := new(App)
	instance = a

	a.Application = app.App(800, 600, "Mc World")
	a.debugMode = true

	a.log = logger.New("main", nil)
	a.log.AddWriter(logger.NewConsole(false))
	a.log.SetFormat(logger.FTIME | logger.FMICROS)
	a.log.SetLevel(logger.DEBUG)

	a.Gls().Enable(gls.CULL_FACE)
	a.Gls().Enable(gls.DEPTH_TEST)

	// Create Scene
	a.scene = core.NewNode()

	// Creates a grid helper and saves its pointer in the test state
	a.grid = helper.NewGrid(50, 1, &math32.Color{0.4, 0.4, 0.4})
	a.scene.Add(a.grid)

	// Create camera
	// w, h := a.GetSize()
	// aspect := float32(w) / float32(h)
	// a.camera = camera.New(aspect)
	// a.scene.Add(a.camera)
	// a.orbit = camera.NewOrbitControl(a.camera)

	// Create frame rater
	a.frameRater = util.NewFrameRater(60)

	// a.player = NewOldPlayer()
	// a.player.ResetPosition(*math32.NewVector3(0, 50, 0))
	// a.scene.Add(a.player)
	// a.scene.Add(a.player.Camera)
	// a.scene.Add(a.player)
	// a.scene.Add(a.player.Camera)

	a.dirData = a.checkDirData("data")
	a.log.Info("Using data directory:%s", a.dirData)

	a.setupScene()

	a.sm = newFileSaveManager(a)
	a.initSeed()

	a.bm = blockv2.NewBlockManager(a.log, a.dirData)

	a.curWorld = NewWorld()
	a.curWorld.Start(a)

	a.player = NewPlayer()
	a.player.Start(a)
	a.player.SetPositionVec(*math32.NewVector3(0, 50, 0))

	a.buildGui()

	// Register Listen
	gui.Manager().SubscribeID(window.OnKeyDown, &a, a.OnKeyDown)
	a.Subscribe(window.OnWindowSize, a.OnWindowSize)
	a.OnWindowSize("", nil)

	return instance
}

func (a *App) setupScene() {
	if a.curWorld != nil {
		a.curWorld.Cleanup(a)
	}

	a.UnsubscribeAllID(a)

	a.DisposeAllCustomCursors()
	a.SetCursor(window.ArrowCursor)
	window.Get().(*window.GlfwWindow).SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Set default background color
	a.Gls().ClearColor(0.6, 0.6, 0.6, 1.0)

	// Reset renderer z-sorting flag
	a.Renderer().SetObjectSorting(true)

	// Create and add an axis helper to the scene
	a.scene.Add(helper.NewAxes(1))
}

func (a *App) initSeed() {
	a.seed = a.sm.LoadSeed()
	if a.seed != 0 {
		a.Log().Debug("load seed:%d", a.seed)
		return
	}

	a.seed = time.Now().Unix()
	a.Log().Debug("new seed:%d", a.seed)
	a.sm.SaveSeed(a.seed)
}

func (a *App) buildGui() {
	width, height := a.GetSize()
	a.mainPanel = gui.NewPanel(float32(width), float32(height))
	a.mainPanel.SetRenderable(false)
	a.mainPanel.SetEnabled(false)
	a.mainPanel.SetLayout(gui.NewDockLayout())
	a.scene.Add(a.mainPanel)
	gui.Manager().Set(a.mainPanel)

	headerColor := math32.Color4{0, 0, 0, 0.1}
	lightTextColor := math32.Color4{0.8, 0.8, 0.8, 1}
	header := gui.NewPanel(100, 20)
	header.SetPaddings(4, 4, 4, 4)
	header.SetColor4(&headerColor)
	header.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})

	// Horizontal box layout for the header
	hbox := gui.NewHBoxLayout()
	header.SetLayout(hbox)
	a.mainPanel.Add(header)

	const fontSize = 14
	// FPS
	l1 := gui.NewLabel(" ")
	l1.SetFontSize(fontSize)
	l1.SetLayoutParams(&gui.HBoxLayoutParams{AlignV: gui.AlignCenter})
	l1.SetText("  FPS: ")
	l1.SetColor4(&lightTextColor)
	header.Add(l1)
	// FPS value
	a.labelFPS = gui.NewLabel(" ")
	a.labelFPS.SetFontSize(fontSize)
	a.labelFPS.SetLayoutParams(&gui.HBoxLayoutParams{AlignV: gui.AlignCenter})
	a.labelFPS.SetColor4(&lightTextColor)
	header.Add(a.labelFPS)

	// debug panel
	a.debugPanel = NewDebugPanel(a)
	a.mainPanel.Add(a.debugPanel.GetPanel())

	cursor := gui.NewPanel(20, 20)
	cursor.SetRenderable(false)
	cursor.SetEnabled(false)
	cursor.SetLayout(gui.NewDockLayout())

	label := gui.NewLabel(" ")
	label.SetFontSize(20)
	label.SetLayoutParams(&gui.DockLayoutParams{})
	label.SetText("+")
	label.SetColor4(&lightTextColor)
	cursor.Add(label)
	a.cursor = cursor
	a.cursor.SetPosition(float32(width)/2-10, float32(height)/2-10)

	a.mainPanel.Add(a.cursor)

	// Player gui
	a.playerGui = NewPlayerGui(a)
	a.scene.Add(a.playerGui)
}

func (a *App) Run() {
	a.Application.Run(a.Update)
}

func (a *App) Update(rend *renderer.Renderer, deltaTime time.Duration) {
	// Start measuring this frame
	a.frameRater.Start()

	// Clear the color, depth, and stencil buffers
	a.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT) // TODO maybe do inside renderer, and allow customization

	// Update the current running demo if any
	if a.curWorld != nil {
		a.curWorld.Update(a, deltaTime)
	}

	if a.player != nil {
		a.player.Update(a, deltaTime)
	}

	// Render scene
	err := rend.Render(a.scene, a.player.Camera)
	if err != nil {
		panic(err)
	}

	// Control and update FPS
	a.frameRater.Wait()
	a.updateFPS()
	a.updateDebug()
}

// UpdateFPS updates the fps value in the window title or header label
func (a *App) updateFPS() {
	// Get the FPS and potential FPS from the frameRater
	fps, pfps, ok := a.frameRater.FPS(time.Duration(1000) * time.Millisecond)
	if !ok {
		return
	}

	// Show the FPS in the header label
	a.labelFPS.SetText(fmt.Sprintf("%3.1f / %3.1f", fps, pfps))
}

func (a *App) updateDebug() {
	if a.debugPanel != nil {
		a.debugPanel.update()
	}
}

func (a *App) checkDirData(dirDataName string) string {
	// Check first if data directory is in the current directory
	if _, err := os.Stat(dirDataName); err != nil {
		panic(err)
	}
	dirData, err := filepath.Abs(dirDataName)
	if err != nil {
		panic(err)
	}
	return dirData
}

func (a *App) Scene() *core.Node {
	return a.scene
}

func (a *App) Log() *logger.Logger {
	return a.log
}

func (a *App) Player() *Player {
	return a.player
}

func (a *App) World() *World {
	return a.curWorld
}

func (a *App) SaveManager() ISaveManager {
	return a.sm
}

func (a *App) IsDebugMode() bool {
	return a.debugMode
}

// DirData returns the base directory for data
func (a *App) DirData() string {

	return a.dirData
}

func (a *App) OnWindowSize(evname string, ev interface{}) {
	w, h := a.GetFramebufferSize()
	aspect := float32(w) / float32(h)

	a.log.Debug("OnWindowSize: w: %d, h: %d, camera aspect: %.6f", w, h, aspect)

	a.Gls().Viewport(0, 0, int32(w), int32(h))
	a.player.Camera.SetAspect(aspect)

	a.mainPanel.SetSize(float32(w), float32(h))
}

func (a *App) OnKeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyEscape:
		a.World().cm.SaveAll()
		a.Exit()
	}
}
