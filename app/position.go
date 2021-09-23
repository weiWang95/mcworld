package app

import "github.com/g3n/engine/math32"

type IPosition interface {
	GetPosition() math32.Vector3
	SetPosition(vec math32.Vector3)
}
