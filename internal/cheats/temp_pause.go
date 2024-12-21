package cheats

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type TempPauseState int

const (
	TempPauseStateNone TempPauseState = iota
	TempPauseStateJustUnpaused
	TempPauseStatePaused
	TempPauseStateJustPaused
)

type TempPause struct {
	state TempPauseState
}

func (t *TempPause) TryPause() {
	if t.state == TempPauseStateNone {
		t.state = TempPauseStateJustPaused
	}
}

func (t *TempPause) TryUnpause() {
	if t.state == TempPauseStatePaused {
		t.state = TempPauseStateJustUnpaused
	}
}

func (t *TempPause) HandleStateUpdate(paused bool, inp *input.Input) {
	switch t.state {
	case TempPauseStateJustPaused:
		if paused {
			inp.AddKeyNewlyPressed(ebiten.KeyP)
		}
		t.state = TempPauseStatePaused

	case TempPauseStateJustUnpaused:
		if !paused {
			inp.AddKeyNewlyPressed(ebiten.KeyP)
		}
		t.state = TempPauseStateNone

	case TempPauseStateNone, TempPauseStatePaused:
		// No action needed for these states
	}
}
