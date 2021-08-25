package app

import (
	"github.com/g3n/engine/window"
)

func (a *App) OnWindowSize(evname string, ev interface{}) {
	w, h := a.GetFramebufferSize()
	aspect := float32(w) / float32(h)

	a.log.Debug("OnWindowSize: w: %d, h: %d, camera aspect: %.6f", w, h, aspect)

	a.Gls().Viewport(0, 0, int32(w), int32(h))
	a.camera.SetAspect(aspect)

	a.mainPanel.SetSize(float32(w), float32(h))
}

func (a *App) OnKeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyEscape:
		a.Exit()
	}
}
