package engine

import (
	"strings"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/cheats/game"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/cheats/tps"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/geometry"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

var (
	keyaliases = map[ebiten.Key]ebiten.Key{
		ebiten.KeyZ: ebiten.KeyDown,
	}

	keymap = map[ebiten.Key]ebiten.Key{
		// ebiten.KeyUp:    ebiten.KeyW,
		// ebiten.KeyLeft:  ebiten.KeyA,
		// ebiten.KeyRight: ebiten.KeyD,
	}
)

func (e *Engine) PreprocessKeys(inp *input.Input) {
	e.MapKeys(inp, keyaliases)

	// Customly mapped keys
	e.HandlePauseJump(inp)
	e.HandlePauseMove(inp)
	e.HandlePauseSkip(inp)
	e.HandleClipboardFeed(inp)

	e.MapKeys(inp, keymap)

	if game.MazeSolverActive {
		inp.AddKeyPressed(e.MazeSolver.NextMove())
	}

	// Keys for custom client-side actions
	e.HandleFreeCamKeys(inp)
	e.HandleTPS(inp)
}

func (e *Engine) HandleFreeCamKeys(inp *input.Input) {
	if !e.FreeRoamMode() {
		return
	}

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
	if !e.FreeRoamMode() {
		return
	}

	if inp.IsKeyNewlyPressed(ebiten.KeyEqual) {
		tps.Increment()
	} else if inp.IsKeyNewlyPressed(ebiten.KeyMinus) {
		tps.Decrement()
	}
}

func (e *Engine) MapKeys(inp *input.Input, mapping map[ebiten.Key]ebiten.Key) {
	if !e.FreeRoamMode() {
		return
	}

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

type VState = int

const (
	VStateNotPressed VState = iota
	VStateWaitingForRelease
)

type ClipboardFeed struct {
	Keys   []ebiten.Key
	VState VState
}

func (e *Engine) HandleClipboardFeed(inp *input.Input) {
	if e.activeNPC == nil {
		return
	}

	if e.ClipboardFeed.VState == VStateNotPressed && inp.IsKeyNewlyPressed(ebiten.KeyV) && inp.IsKeyPressed(ebiten.KeyControl) {
		c := clipboard.New()
		text, err := c.PasteText()
		if err != nil {
			return
		}

		if strings.ContainsRune(text, '\n') {
			logrus.Error("can't contain newline symbols in clipboard feed")
			return
		}

		keys := make([]ebiten.Key, 0, len(text))
		for _, b := range []byte(text) {
			var k ebiten.Key
			switch b {
			case ' ':
				k = ebiten.KeySpace
			case '.':
				k = ebiten.KeyPeriod
			case ',':
				k = ebiten.KeyComma
			case '/':
				k = ebiten.KeySlash
			default:
				k = ebiten.Key(0)
				if err := k.UnmarshalText([]byte{b}); err != nil {
					logrus.Errorf("can't unmarshal key: %s", err.Error())
					return
				}
			}
			keys = append(keys, k)
		}
		e.ClipboardFeed.Keys = append(keys, e.ClipboardFeed.Keys...)

		e.ClipboardFeed.VState = VStateWaitingForRelease
		inp.RemoveKeyNewlyPressed(ebiten.KeyV)
	}
	if !inp.IsKeyPressed(ebiten.KeyV) {
		e.ClipboardFeed.VState = VStateNotPressed
	}

	if len(e.ClipboardFeed.Keys) > 0 {
		inp.AddKeyNewlyPressed(e.ClipboardFeed.Keys[0])
		e.ClipboardFeed.Keys = e.ClipboardFeed.Keys[1:]
	}
}
