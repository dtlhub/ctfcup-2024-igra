package engine

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/cheats/tps"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/geometry"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	keyaliases = map[ebiten.Key]ebiten.Key{
		ebiten.KeyZ: ebiten.KeyDown,
	}

	keymap = map[ebiten.Key]ebiten.Key{
		ebiten.KeyUp:    ebiten.KeyW,
		ebiten.KeyLeft:  ebiten.KeyA,
		ebiten.KeyRight: ebiten.KeyD,
	}
)

func (e *Engine) HandleCustomKeys(inp *input.Input) {
	e.HandleFreeCamKeys(inp)
}

func (e *Engine) HandleFreeCamKeys(inp *input.Input) {
	if inp.IsKeyNewlyPressed(ebiten.KeyN) {
		e.ToggleFreeCam()
	}

	if e.FreeCam.Enabled {
		if inp.IsKeyPressed(ebiten.KeyI) {
			e.FreeCam.Move(&geometry.Vector{X: 0, Y: -e.FreeCam.Speed})
		}

		if inp.IsKeyPressed(ebiten.KeyJ) {
			e.FreeCam.Move(&geometry.Vector{X: -e.FreeCam.Speed, Y: 0})
		}

		if inp.IsKeyPressed(ebiten.KeyK) {
			e.FreeCam.Move(&geometry.Vector{X: 0, Y: e.FreeCam.Speed})
		}

		if inp.IsKeyPressed(ebiten.KeyL) {
			e.FreeCam.Move(&geometry.Vector{X: e.FreeCam.Speed, Y: 0})
		}

		if inp.IsKeyPressed(ebiten.KeyO) {
			e.FreeCam.SpeedUp()
		}

		if inp.IsKeyPressed(ebiten.KeyU) {
			e.FreeCam.SpeedDown()
		}
	}
}

func (e *Engine) PreprocessKeys(inp *input.Input) {
	e.MapKeys(inp, keyaliases)

	e.HandlePauseJump(inp)
	e.HandlePauseMove(inp)
	e.HandlePauseSkip(inp)
	e.HandleTPS(inp)

	e.MapKeys(inp, keymap)
}

func (e *Engine) HandlePauseJump(inp *input.Input) {
	if inp.IsKeyPressed(ebiten.KeyUp) {
		e.JumpTempPause.TryPause()
	} else {
		e.JumpTempPause.TryUnpause()
	}

	e.JumpTempPause.HandleStateUpdate(e.Paused, inp)
}

func (e *Engine) HandlePauseMove(inp *input.Input) {
	if inp.IsKeyPressed(ebiten.KeyLeft) {
		e.FrameMovement = FrameMovementStateStarted
	}

	if inp.IsKeyPressed(ebiten.KeyRight) {
		e.FrameMovement = FrameMovementStateStarted
	}

	switch e.FrameMovement {
	case FrameMovementStateStarted:
		if e.Paused {
			inp.AddKeyNewlyPressed(ebiten.KeyP)
		}
		e.FrameMovement = FrameMovementStateAwaitingPause

	case FrameMovementStateAwaitingPause:
		if !e.Paused {
			inp.AddKeyNewlyPressed(ebiten.KeyP)
		}
		e.FrameMovement = FrameMovementStateNone

	case FrameMovementStateNone:
		// do nothing
	}
}

func (e *Engine) HandlePauseSkip(inp *input.Input) {
	if inp.IsKeyPressed(ebiten.KeyDown) {
		e.TempPause.TryPause()
	} else if !inp.IsKeyPressed(ebiten.KeyDown) {
		e.TempPause.TryUnpause()
	}

	e.TempPause.HandleStateUpdate(e.Paused, inp)
}

func (e *Engine) HandleTPS(inp *input.Input) {
	if inp.IsKeyNewlyPressed(ebiten.KeyEqual) {
		tps.Increment()
	} else if inp.IsKeyNewlyPressed(ebiten.KeyMinus) {
		tps.Decrement()
	}
}

func (e *Engine) MapKeys(inp *input.Input, mapping map[ebiten.Key]ebiten.Key) {
	for k, m := range mapping {
		if inp.IsKeyPressed(k) {
			inp.RemoveKeyPressed(k)
			inp.AddKeyPressed(m)
		}
		if inp.IsKeyNewlyPressed(k) {
			inp.RemoveKeyNewlyPressed(k)
			inp.AddKeyNewlyPressed(m)
		}
	}
}
