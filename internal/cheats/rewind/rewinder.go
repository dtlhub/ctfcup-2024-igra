package rewind

import (
	"fmt"
	"os"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/cheats/tps"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/engine"
	gameserverpb "github.com/c4t-but-s4d/ctfcup-2024-igra/proto/go/gameserver"
)

type Rewinder struct {
	loadedRewind *Rewind
	lastFrame    int
	skipEmpty    bool
	finished     bool
	engine       *engine.Engine
}

func (r *Rewinder) Initialize(engine *engine.Engine) error {
	r.engine = engine
	r.loadedRewind = &Rewind{}
	if rewindFile := os.Getenv("REWIND_FILE"); rewindFile != "" {
		if err := r.loadedRewind.Load(rewindFile); err != nil {
			return fmt.Errorf("loading rewind: %w", err)
		}
	}

	skipEmpty, ok := os.LookupEnv("REWIND_SKIP_EMPTY")
	r.skipEmpty = ok && skipEmpty != "0"

	tps.SetFromEnv("REWIND_TPS", 6_000)

	return nil
}

func (r *Rewinder) CurrentFrame() int {
	return r.lastFrame
}

func (r *Rewinder) TotalFrames() int {
	return len(r.loadedRewind.Moves)
}

func (r *Rewinder) NextFrame() (keys *gameserverpb.ClientEvent_KeysPressed, justPaused, done bool) {
	r.skipFrames()
	if r.done() {
		justStopped := !r.finished
		r.finished = true
		return nil, justStopped, true
	}

	moves := r.loadedRewind.Moves[r.lastFrame]
	r.lastFrame++
	return moves, false, false
}

func (r *Rewinder) skipFrames() {
	if r.lastFrame >= len(r.loadedRewind.Moves) {
		return
	}

	moves := r.loadedRewind.Moves[r.lastFrame]
	for (r.engine.Paused || r.skipEmpty) && r.lastFrame < len(r.loadedRewind.Moves)-1 &&
		(r.engine.IsStill() && moves.KeysPressed == nil && moves.NewKeysPressed == nil) {
		r.lastFrame++
		moves = r.loadedRewind.Moves[r.lastFrame]
	}
}

func (r *Rewinder) done() bool {
	return r.lastFrame >= len(r.loadedRewind.Moves)
}
