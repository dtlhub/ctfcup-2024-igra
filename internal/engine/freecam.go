package engine

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/camera"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/geometry"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/object"
)

const defaultFreeCamSpeed = 20.0

type FreeCam struct {
	*object.Base
	Enabled bool
	Speed   float64
}

func NewFreeCam() *FreeCam {
	return &FreeCam{
		Enabled: false,
		Speed:   defaultFreeCamSpeed,
	}
}

func (f *FreeCam) Reset(camera *camera.Camera) {
	f.Base = &object.Base{
		Origin: geometry.Point{
			X: camera.Origin.X,
			Y: camera.Origin.Y,
		},
		Width:  camera.Width,
		Height: camera.Height,
	}
	f.Speed = defaultFreeCamSpeed
}

func (f *FreeCam) Move(v *geometry.Vector) *FreeCam {
	f.Base.Origin.X += v.X
	f.Base.Origin.Y += v.Y
	return f
}

func (f *FreeCam) SpeedUp() {
	f.Speed *= 1.0 + 1.0/64.0
}

func (f *FreeCam) SpeedDown() {
	f.Speed /= 1.0 + 1.0/64.0
}

func (e *Engine) ToggleFreeCam() {
	e.FreeCam.Enabled = !e.FreeCam.Enabled
	if e.FreeCam.Enabled {
		e.FreeCam.Reset(e.Camera)
	}
}

func (e *Engine) CameraObject() *object.Base {
	if e.FreeCam.Enabled {
		return e.FreeCam.Base
	}
	return e.Camera.Base
}
