package game

import (
	"os"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
)

var MazeSolverActive = false

const size = 200

func init() {
	value, ok := os.LookupEnv("MAZE_SOLVER")
	MazeSolverActive = ok && value != "0"
}

type cell int

const (
	unknown cell = iota
	empty   cell = iota
	wall
	finish
)

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

type point struct {
	x, y int
}

type MazeSolver struct {
	move           int
	prohibitedMove *ebiten.Key
	lastState      *arcade.State
	playerLocation point
	movedToCell    point
	maze           [size][size]int
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
	if s.playerLocation.x != x || s.playerLocation.y != y {
		s.playerLocation = point{x, y}
	}
}

func GetPlayerLocation(s *arcade.State) (int, int) {
	for i := 0; i < len(s.Screen); i++ {
		for j := 0; j < len(s.Screen[i]); j++ {
			r, g, b, a := s.Screen[i][j].RGBA()
			if r == 65535 && g == 0 && b == 0 && a == 65535 {
				return i, j
			}
		}
	}
	return -1, -1
}
