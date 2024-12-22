package game

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
)

type Solver interface {
	NextMove() ebiten.Key
	FeedResult(arcade.Result)
}

func ScreensEqual(a, b *arcade.State) bool {
	for i := 0; i < len(a.Screen); i++ {
		for j := 0; j < len(a.Screen[i]); j++ {
			if a.Screen[i][j] != b.Screen[i][j] {
				return false
			}
		}
	}
	return true
}
