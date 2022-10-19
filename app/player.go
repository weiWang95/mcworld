package app

import (
	"math"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"github.com/weiWang95/mcworld/app/block"
	"github.com/weiWang95/mcworld/app/blockv2"
)

const PLAYER_JUMP_SPEED = 4.85 // 1.5: 5.42 1.2: 4.85
const MaxControlDistance = 8

type OrbitEnabled int

// The possible control types.
const (
	OrbitNone OrbitEnabled = 0x00
	OrbitRot  OrbitEnabled = 0x01
	OrbitZoom OrbitEnabled = 0x02
	OrbitPan  OrbitEnabled = 0x04
	OrbitKeys OrbitEnabled = 0x08
	OrbitAll  OrbitEnabled = 0xFF
)

type orbitState int

const (
	stateNone = orbitState(iota)
	stateRotate
	stateZoom
	statePan
)

type PlayMode uint8

const (
	PlayModeCreate = iota
	PlayModeLife
)

type Player struct {
	core.Node
	core.Dispatcher

	Model     IControl
	Camera    *camera.Camera
	wreckLine *graphic.Lines
	Target    *PlayerTarget

	up     math32.Vector3
	farPos math32.Vector3

	playMode      PlayMode
	speed         float32
	vSpeed        float32
	inFall        bool
	moveDirection math32.Vector3

	enabled   OrbitEnabled
	state     orbitState
	RotSpeed  float32
	ZoomSpeed float32
	rotStart  math32.Vector2
	panStart  math32.Vector2
	zoomStart float32

	// ticker
	wreckTicker     *TickChecker
	curInventoryIdx uint8
	inventory       [10]block.BlockId
}

func NewPlayer() *Player {
	p := new(Player)
	p.Node = *core.NewNode()
	p.Dispatcher.Initialize()

	p.speed = 10

	p.enabled = OrbitAll
	p.state = stateRotate
	p.RotSpeed = 1.0
	p.ZoomSpeed = 0.1

	p.up = *math32.NewVector3(0, 1, 0)

	p.Model = NewPlayerModel()

	p.Camera = camera.New(16 / 9)
	p.ResetCamera()

	p.farPos = *p.GetViewport().Clone().Add(math32.NewVector3(p.Model.GetHandLength(), 0, 0))

	p.wreckTicker = NewTickChecker(8)

	p.playMode = PlayModeLife
	p.inventory = [10]block.BlockId{block.BlockSoil, block.BlockBrick, block.BlockLamp}

	// Subscribe to events
	gui.Manager().SetCursorFocus(p)
	gui.Manager().SubscribeID(window.OnMouseUp, &p, p.onMouse)
	gui.Manager().SubscribeID(window.OnMouseDown, &p, p.onMouse)
	// gui.Manager().SubscribeID(window.OnScroll, &p, p.onScroll)
	gui.Manager().SubscribeID(window.OnKeyDown, &p, p.onKey)
	gui.Manager().SubscribeID(window.OnKeyUp, &p, p.onKey)
	// gui.Manager().SubscribeID(window.OnKeyRepeat, &p, p.onKey)
	p.SubscribeID(window.OnCursor, &p, p.onCursor)

	return p
}

// Dispose unsubscribes from all events.
func (p *Player) Dispose() {
	gui.Manager().UnsubscribeID(window.OnMouseUp, &p)
	gui.Manager().UnsubscribeID(window.OnMouseDown, &p)
	// gui.Manager().UnsubscribeID(window.OnScroll, &p)
	gui.Manager().UnsubscribeID(window.OnKeyDown, &p)
	gui.Manager().UnsubscribeID(window.OnKeyRepeat, &p)
	p.UnsubscribeID(window.OnCursor, &p)

	gui.Manager().SetCursorFocus(nil)
}

func (p *Player) Start(a *App) {
	p.Model.Start(a)
	p.Add(p.Model)

	p.addWreckLine()
	p.Add(p.Camera)

	p.Target = NewPlayerTarget()
	p.Add(p.Target)

	a.Scene().Add(p)
}

func (p *Player) Update(a *App, t time.Duration) {
	delta := float32(t) / float32(time.Second)
	vSpeed := p.vSpeed * delta

	pos := p.GetPosition()
	// if vSpeed != 0 {
	// 	Instance().Log().Debug("pos: %v fall:%v hv:%.02f hd:%.02f", pos, p.inFall, p.vSpeed, vSpeed)
	// }

	if !p.IsCreatePlayMode() {
		if p.vSpeed > 0 {
			npos := math32.NewVector3(pos.X, p.Model.GetBoundBox().BY+vSpeed, pos.Z)
			block, _ := a.World().GetBlockByVec(*npos)
			if block != nil {
				vSpeed = float32(int64(pos.Y)) - pos.Y
				p.vSpeed = 0
				p.inFall = true
			} else {
				p.vSpeed += DEFAULT_GRAVITY_SPEED * delta
				p.vSpeed = math32.Clamp(p.vSpeed, MAX_GRAVITY_SPEED, 40)
				p.inFall = true
			}
		} else {
			block, _ := a.World().GetBlockByVec(*pos.Clone().Add(math32.NewVector3(0, vSpeed-0.01, 0)))
			if block == nil {
				p.vSpeed += DEFAULT_GRAVITY_SPEED * delta
				p.vSpeed = math32.Clamp(p.vSpeed, MAX_GRAVITY_SPEED, 40)
				p.inFall = true
			} else {
				vSpeed = float32(int64(pos.Y)) - pos.Y
				p.vSpeed = 0
				p.inFall = false
			}
		}
	}

	p.Move(a, p.GetSpeed()*delta, vSpeed)

	p.Model.Update(a, t)

	if p.wreckLine != nil {
		pos := p.Camera.Position()
		p.wreckLine.SetPositionVec(&pos)
		p.wreckLine.LookAt(&p.farPos, &p.up)
	}

	if p.wreckTicker.Next(t) {
		p.Target.SetTarget(p.GetTarget())
	}
}

func (p *Player) Cleanup() {

}

func (p *Player) GetViewport() *math32.Vector3 {
	return p.Model.GetViewport()
}

func (p *Player) GetPosition() *math32.Vector3 {
	pos := p.Node.Position()
	return &pos
}

func (p *Player) SetPositionVec(pos math32.Vector3) {
	p.Node.SetPositionVec(&pos)
	p.Model.SetPosition(&pos)

	viewport := p.GetViewport()
	p.Camera.SetPositionVec(viewport.Clone().Add(math32.NewVector3(-1, 0, 0)))
	p.Camera.LookAt(p.GetViewport(), &p.up)
}

func (p *Player) ResetCamera() {
	p.Camera.SetAspect(16 / 9)
	p.Camera.UpdateSize(3)
	p.Camera.SetProjection(camera.Perspective)
	viewport := p.GetViewport()
	p.Camera.SetPositionVec(viewport.Clone().Add(math32.NewVector3(-1, 0, 0)))
	p.Camera.LookAt(viewport, math32.NewVector3(0, 1, 0))
}

func (p *Player) GetSpeed() float32 {
	return p.speed
}

func (p *Player) GetJumpPower() float32 {
	return PLAYER_JUMP_SPEED
}

func (p *Player) Rotate(thetaDelta, phiDelta float32) {
	const EPS = 0.0001

	// Compute direction vector from target to camera
	tcam := p.Camera.Position()
	viewport := p.GetViewport()
	tcam.Sub(viewport)

	// Calculate angles based on current camera position plus deltas
	radius := tcam.Length()
	theta := math32.Atan2(tcam.X, tcam.Z) + thetaDelta
	phi := math32.Acos(tcam.Y/radius) + phiDelta

	// Restrict phi and theta to be between desired limits
	phi = math32.Clamp(phi, 0, math32.Pi)
	phi = math32.Clamp(phi, EPS, math32.Pi-EPS)
	theta = math32.Clamp(theta, float32(math.Inf(-1)), float32(math.Inf(1)))

	// Calculate new cartesian coordinates
	tcam.X = radius * math32.Sin(phi) * math32.Sin(theta)
	tcam.Y = radius * math32.Cos(phi)
	tcam.Z = radius * math32.Sin(phi) * math32.Cos(theta)

	handLength := p.Model.GetHandLength()
	x := handLength * math32.Sin(phi) * math32.Sin(theta)
	y := handLength * math32.Cos(phi)
	z := handLength * math32.Sin(phi) * math32.Cos(theta)

	p.Camera.SetPositionVec(viewport.Clone().Add(&tcam))
	p.Camera.LookAt(viewport, &p.up)
	p.farPos = *viewport.Clone().Add(math32.NewVector3(-x, -y, -z))
}

// Zoom moves the camera closer or farther from the target the specified amount
// and also updates the camera's orthographic size to match.
func (p *Player) Zoom(delta float32) {
	viewport := p.GetViewport()

	// Compute direction vector from target to camera
	tcam := p.Camera.Position()
	tcam.Sub(viewport)

	// Calculate new distance from target and apply limits
	dist := tcam.Length() * (1 + delta/10)
	dist = math32.Max(1.0, math32.Min(float32(math.Inf(1)), dist))
	tcam.SetLength(dist)

	// Update orthographic size and camera position with new distance
	p.Camera.UpdateSize(tcam.Length())
	p.Camera.SetPositionVec(viewport.Clone().Add(&tcam))
}

// Pan pans the camera and target the specified amount on the plane perpendicular to the viewing direction.
func (p *Player) Pan(deltaX, deltaY float32) {
	viewport := p.GetViewport()

	// Compute direction vector from camera to target
	position := p.Camera.Position()
	vdir := viewport.Clone().Sub(&position)

	// Conversion constant between an on-screen cursor delta and its projection on the target plane
	c := 2 * vdir.Length() * math32.Tan((p.Camera.Fov()/2.0)*math32.Pi/180.0) / p.winSize()

	// Calculate pan components, scale by the converted offsets and combine them
	var pan, panX, panY math32.Vector3
	panX.CrossVectors(&p.up, vdir).Normalize()
	panY.CrossVectors(vdir, &panX).Normalize()
	panY.MultiplyScalar(c * deltaY)
	panX.MultiplyScalar(c * deltaX)
	pan.AddVectors(&panX, &panY)

	// Add pan offset to camera and target
	p.Camera.SetPositionVec(position.Add(&pan))
	viewport.Add(&pan)
}

func (p *Player) Move(a *App, speed float32, vSpeed float32) {
	if p.moveDirection.X == 0 && p.moveDirection.Z == 0 && p.moveDirection.Y == 0 && vSpeed == 0 {
		return
	}
	viewport := p.GetViewport()

	// Compute direction vector from target to camera
	tcam := p.Camera.Position()
	tcam.Sub(viewport)

	theta := math32.Atan2(tcam.X, tcam.Z)
	theta = math32.Clamp(theta, float32(math.Inf(-1)), float32(math.Inf(1)))

	tcam.X = speed*math32.Sin(theta+math.Pi)*p.moveDirection.X + speed*math32.Sin(theta+0.5*math.Pi)*p.moveDirection.Z
	tcam.Z = speed*math32.Cos(theta+math.Pi)*p.moveDirection.X + speed*math32.Cos(theta+0.5*math.Pi)*p.moveDirection.Z

	pos := p.GetPosition()
	if tcam.X > 0 {
		// xBlock := a.World().GetBlockByPosition(pos.X+p.Model.GetBoundBox().X/2+tcam.X, pos.Y, pos.Z)
		xBlock, _ := a.World().GetBlockByPosition(p.Model.GetBoundBox().BX+tcam.X, pos.Y, pos.Z)
		if xBlock != nil {
			tcam.X = 0
		}
	} else if tcam.X < 0 {
		// xBlock := a.World().GetBlockByPosition(pos.X-p.Model.GetBoundBox().X/2+tcam.X, pos.Y, pos.Z)
		xBlock, _ := a.World().GetBlockByPosition(p.Model.GetBoundBox().X+tcam.X, pos.Y, pos.Z)
		if xBlock != nil {
			tcam.X = 0
		}
	}

	if tcam.Z > 0 {
		zBlock, _ := a.World().GetBlockByPosition(pos.X, pos.Y, p.Model.GetBoundBox().BZ+tcam.Z)
		if zBlock != nil {
			tcam.Z = 0
		}

	} else if tcam.Z < 0 {
		zBlock, _ := a.World().GetBlockByPosition(pos.X, pos.Y, p.Model.GetBoundBox().Z+tcam.Z)
		if zBlock != nil {
			tcam.Z = 0
		}
	}

	if viewport.Y+vSpeed < -10 {
		vSpeed = float32(CHUNK_HEIGHT) - 1 - viewport.Y
		viewport.Y = float32(CHUNK_HEIGHT) - 1
	}
	viewport.Add(math32.NewVector3(tcam.X, vSpeed, tcam.Z))

	newPos := pos.Add(math32.NewVector3(tcam.X, vSpeed, tcam.Z))
	camPos := p.Camera.Position()

	p.SetPositionVec(*newPos)

	p.Camera.SetPositionVec(camPos.Add(math32.NewVector3(tcam.X, vSpeed, tcam.Z)))
	p.Camera.LookAt(p.GetViewport(), &p.up)
}

func (p *Player) Jump() {
	Instance().log.Debug("pos -> %v, vSpeed -> %v, fall: %v \n", p.GetPosition(), p.vSpeed, p.inFall)

	if p.inFall {
		return
	}

	p.vSpeed = p.GetJumpPower()
}

func (p *Player) WreckBlock() {
	block, _ := p.GetTarget()
	if block == nil {
		return
	}

	Instance().curWorld.WreckBlock(block.GetPosition())
}

func (p *Player) PlaceBlock() {
	Instance().log.Debug("place block! start:%v, end:%v", p.GetViewport(), p.farPos)
	b, hitPos := p.GetTarget()
	if b == nil {
		return
	}

	face := GetBlockFace(b.GetPosition(), *hitPos)
	pos := b.GetPosition()
	switch face {
	case block.BlockFaceFront:
		pos = *(pos.Add(math32.NewVector3(0, 0, -1)))
	case block.BlockFaceBack:
		pos = *(pos.Add(math32.NewVector3(0, 0, 1)))
	case block.BlockFaceLeft:
		pos = *(pos.Add(math32.NewVector3(-1, 0, 0)))
	case block.BlockFaceRight:
		pos = *(pos.Add(math32.NewVector3(1, 0, 0)))
	case block.BlockFaceTop:
		pos = *(pos.Add(math32.NewVector3(0, 1, 0)))
	case block.BlockFaceBottom:
		pos = *(pos.Add(math32.NewVector3(0, -1, 0)))
	default:
		return
	}
	Instance().Log().Debug("place block -> b:%v hit:%v face:%v place pos: %v", b.GetPosition(), hitPos, face, pos)

	placeBlock, _ := Instance().curWorld.GetBlockByVec(pos)
	if placeBlock != nil {
		return
	}

	// nb := block.NewBlock(p.inventory[p.curInventoryIdx], pos)
	nb := Instance().bm.NewBlock(blockv2.BlockId(p.inventory[p.curInventoryIdx]))
	nb.SetPositionVec(&pos)
	Instance().curWorld.PlaceBlock(nb, pos)
}

func (p *Player) GetTarget() (*blockv2.Block, *math32.Vector3) {
	b, pos := RayTraceBlock(Instance().curWorld, *p.GetViewport(), p.farPos)
	if b == nil || pos == nil {
		return nil, nil
	}

	return b, pos
}

// onMouse is called when an OnMouseDown/OnMouseUp event is received.
func (p *Player) onMouse(evname string, ev interface{}) {

	switch evname {
	case window.OnMouseDown:
		// gui.Manager().SetCursorFocus(oc)
		mev := ev.(*window.MouseEvent)
		switch mev.Button {
		case window.MouseButtonLeft:
			p.WreckBlock()
		case window.MouseButtonMiddle:
		case window.MouseButtonRight:
			p.PlaceBlock()
		}
	case window.OnMouseUp:
		// gui.Manager().SetCursorFocus(nil)
		// oc.state = stateNone
	}
}

// onCursor is called when an OnCursor event is received.
func (p *Player) onCursor(evname string, ev interface{}) {
	gui.Manager().SetCursorFocus(p)

	// If nothing enabled ignore event
	if p.enabled == OrbitNone || p.state == stateNone {
		return
	}

	mev := ev.(*window.CursorEvent)
	switch p.state {
	case stateRotate:
		c := -2 * math32.Pi * p.RotSpeed / p.winSize()
		p.Rotate(c*(mev.Xpos-p.rotStart.X),
			c*(mev.Ypos-p.rotStart.Y))
		p.rotStart.Set(mev.Xpos, mev.Ypos)
	case stateZoom:
		p.Zoom(p.ZoomSpeed * (mev.Ypos - p.zoomStart))
		p.zoomStart = mev.Ypos
	case statePan:
		p.Pan(mev.Xpos-p.panStart.X,
			mev.Ypos-p.panStart.Y)
		p.panStart.Set(mev.Xpos, mev.Ypos)
	}
}

// onScroll is called when an OnScroll event is received.
func (p *Player) onScroll(evname string, ev interface{}) {
	if p.enabled&OrbitZoom != 0 {
		sev := ev.(*window.ScrollEvent)
		p.Zoom(-sev.Yoffset)
	}
}

// onKey is called when an OnKeyDown/OnKeyRepeat event is received.
func (p *Player) onKey(evname string, ev interface{}) {

	// If keyboard control is disabled ignore event
	if p.enabled&OrbitKeys == 0 {
		return
	}

	kev := ev.(*window.KeyEvent)
	// kev.Mods == window.ModShift
	switch evname {
	case window.OnKeyDown:
		switch kev.Key {
		case window.KeyUp, window.KeyW:
			p.moveDirection.X += 1
		case window.KeyDown, window.KeyS:
			p.moveDirection.X -= 1
		case window.KeyLeft, window.KeyA:
			p.moveDirection.Z -= 1
		case window.KeyRight, window.KeyD:
			p.moveDirection.Z += 1
		case window.KeySpace:
			p.Jump()
		case window.KeyLeftShift:
			if p.IsCreatePlayMode() {
				p.vSpeed = -PLAYER_JUMP_SPEED
			}
		case window.Key1:
			p.curInventoryIdx = 0
		case window.Key2:
			p.curInventoryIdx = 1
		case window.Key3:
			p.curInventoryIdx = 2
		case window.Key9:
			p.playMode = PlayModeCreate
		case window.Key0:
			p.playMode = PlayModeLife
		}
	case window.OnKeyUp:
		switch kev.Key {
		case window.KeyUp, window.KeyW:
			p.moveDirection.X -= 1
		case window.KeyDown, window.KeyS:
			p.moveDirection.X += 1
		case window.KeyLeft, window.KeyA:
			p.moveDirection.Z += 1
		case window.KeyRight, window.KeyD:
			p.moveDirection.Z -= 1
		case window.KeySpace, window.KeyLeftShift:
			if p.IsCreatePlayMode() {
				p.vSpeed = 0
			}
		}
	}
}

// winSize returns the window height or width based on the camera reference axis.
func (p *Player) winSize() float32 {

	width, size := window.Get().GetSize()
	if p.Camera.Axis() == camera.Horizontal {
		size = width
	}
	return float32(size)
}

func (p *Player) addWreckLine() {
	// Creates geometry
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	vertices.Append(
		0, 0, 0,
		0, 0, MaxControlDistance+1,
	)
	colors := math32.NewArrayF32(0, 16)
	colors.Append(
		1.0, 0.0, 0.0, // red
		1.0, 0.0, 0.0,
	)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates basic material
	mat := material.NewBasic()

	// Creates lines with the specified geometry and material
	p.wreckLine = graphic.NewLines(geom, mat)
	p.wreckLine.SetPositionVec(p.GetViewport())
	p.wreckLine.SetPosition(2, 2, 2)
	p.Add(p.wreckLine)
}

func (p *Player) IsCreatePlayMode() bool {
	return p.playMode == PlayModeCreate
}
