package game

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
)

type Solver interface {
	NextMove() ebiten.Key
	FeedResult(arcade.Result)
}
