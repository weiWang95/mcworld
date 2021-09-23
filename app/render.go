package app

import "time"

type IRender interface {
	Start(a *App)
	Update(a *App, t time.Duration)
	Cleanup()
}
