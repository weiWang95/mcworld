package app

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type IControl interface {
	core.INode
	IRender

	GetPosition() *math32.Vector3
	SetPosition(pos *math32.Vector3)
	SetFace(face *math32.Vector3)
	GetViewport() *math32.Vector3
	GetHandLength() float32
	GetBoundBox() BoundBox
}
