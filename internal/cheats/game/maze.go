package game

import (
	"os"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var MazeSolverActive = false

func init() {
	value, ok := os.LookupEnv("MAZE_SOLVER")
	MazeSolverActive = ok && value != "0"
}

var (
	moves = []ebiten.Key{
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
		ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
	}
	opposite = map[ebiten.Key]ebiten.Key{
		ebiten.KeyArrowUp:    ebiten.KeyArrowDown,
		ebiten.KeyArrowDown:  ebiten.KeyArrowUp,
		ebiten.KeyArrowLeft:  ebiten.KeyArrowRight,
		ebiten.KeyArrowRight: ebiten.KeyArrowLeft,
	}
)

type MazeSolver struct {
	move           int
	prohibitedMove *ebiten.Key
	lastState      *arcade.State
}

func (s *MazeSolver) NextMove() ebiten.Key {
	s.move = (s.move + 1) % len(moves)
	if s.prohibitedMove != nil && moves[s.move] == *s.prohibitedMove {
		s.move = (s.move + 1) % len(moves)
	}
	return moves[s.move]
}

func (s *MazeSolver) Reset() {
	s.move = 0
	s.prohibitedMove = nil
	s.lastState = nil
}

func (s *MazeSolver) FeedState(state *arcade.State) {
	x, y := GetPlayerLocation(state)
	logrus.Infof("Player location: %v", x, y)
	if s.lastState != nil && !ScreensEqual(s.lastState, state) {
		opposite := opposite[moves[s.move]]
		s.prohibitedMove = &opposite
	} else {
		s.prohibitedMove = nil
	}
	s.lastState = state
}
