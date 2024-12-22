package game

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

const MaxMoves = 5

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
	lastMove      int
	lastMoveCount int
}

func (s *MazeSolver) NextMove() ebiten.Key {
	if s.lastMoveCount == MaxMoves {
		s.lastMoveCount = 0
		s.lastMove = (s.lastMove + 1) % len(moves)
	}
	s.lastMoveCount += 1
	logrus.Infof("maze solver move: %d, count: %d", s.lastMove, s.lastMoveCount)
	return moves[s.lastMove]
}

func (s *MazeSolver) FeedResult(arcade.Result) {
}
