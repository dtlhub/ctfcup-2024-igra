package game

import (
	"os"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var BrodilkaSolverActive = false

func init() {
	value, ok := os.LookupEnv("BRODILKA_SOLVER")
	BrodilkaSolverActive = ok && value != "0"
}

type BrodilkaSolver struct {
	moves    []ebiten.Key
	nextMove int
}

func NewBrodilkaSolver() *BrodilkaSolver {
	moves := []ebiten.Key{}
	for i := 0; i < int('W'); i++ {
		moves = append(moves, ebiten.KeyArrowDown)
		moves = append(moves, ebiten.KeyArrowUp)
	}
	moves = append(moves, ebiten.KeyArrowRight)
	for i := 1; i < int('I'); i++ {
		moves = append(moves, ebiten.KeyArrowDown)
		moves = append(moves, ebiten.KeyArrowUp)
	}
	moves = append(moves, ebiten.KeyArrowRight)
	for i := 1; i < int('N'); i++ {
		moves = append(moves, ebiten.KeyArrowDown)
		moves = append(moves, ebiten.KeyArrowUp)
	}
	return &BrodilkaSolver{
		moves: moves,
	}
}

func (m *BrodilkaSolver) NextMove(state *arcade.State) (ebiten.Key, bool) {
	if state.Result == arcade.ResultWon {
		return ebiten.Key(0), false
	}
	if state.Result == arcade.ResultLost {
		m.nextMove = 0
		return ebiten.KeyR, true
	}
	if m.nextMove >= len(m.moves) {
		logrus.Warn("out of moves")
		return ebiten.Key(0), false
	}

	move := m.moves[m.nextMove]
	m.nextMove++
	return move, true
}
