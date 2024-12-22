package game

import (
	"os"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
)

var SnakeSolverActive = false

func init() {
	value, ok := os.LookupEnv("SNAKE_SOLVER")
	SnakeSolverActive = ok && value != "0"
}

type SnakeSolver struct {
	sequence [][]ebiten.Key
	x        int
	y        int
	dir      ebiten.Key
}

const (
	S = ebiten.Key0
	R = ebiten.KeyD
	D = ebiten.KeyS
	L = ebiten.KeyA
	U = ebiten.KeyW
)

func NewSnakeSolver() *SnakeSolver {
	return &SnakeSolver{
		sequence: [][]ebiten.Key{
			{R, S, S, S, S, S, S, D},
			{U, S, S, S, S, S, L, S},
			{R, S, S, S, S, S, U, S},
			{U, S, S, S, S, S, L, S},
			{R, S, S, S, S, S, U, S},
			{U, S, S, S, S, S, L, S},
			{R, S, S, S, S, S, U, S},
			{U, S, S, S, S, S, S, L},
		},
		x:   3,
		y:   0,
		dir: R,
	}
}

func (m *SnakeSolver) NextMove(state *arcade.State) (ebiten.Key, bool) {
	if state.Result != arcade.ResultUnknown {
		return ebiten.Key(0), false
	}

	next := m.sequence[m.y][m.x]

	dir := next
	if next == S {
		dir = m.dir
	}
	m.dir = dir

	switch dir {
	case R:
		m.x++
	case D:
		m.y++
	case L:
		m.x--
	case U:
		m.y--
	}
	return next, next != S
}
