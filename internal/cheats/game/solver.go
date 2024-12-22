package game

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

type Solver interface {
	NextMove(state *arcade.State) (ebiten.Key, bool)
}

type SolverState struct {
	Solver
	Ready bool
}

func GetSolver() Solver {
	if MazeSolverActive {
		logrus.Error("Maze solver is not implemented")
	}
	if SnakeSolverActive {
		return NewSnakeSolver()
	}
	return nil
}
