package app

import (
	"time"
)

var Levels map[string]ILevel

type ILevel interface {
	Start(a *App)
	Update(a *App, deltaTime time.Duration)
	Cleanup(a *App)
}

func RegisterLevel(name string, level ILevel) {
	if Levels == nil {
		Levels = make(map[string]ILevel)
	}

	if _, ok := Levels[name]; ok {
		return
	}

	Levels[name] = level
}
