package game

import (
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

type Solver interface {
	NextMove() ebiten.Key
	FeedResult(arcade.Result)
}

func ScreensEqual(a, b *arcade.State) bool {
	type color struct {
		r, g, b, a uint32
	}
	mapA := make(map[color]bool)

	for i := 0; i < len(a.Screen); i++ {
		for j := 0; j < len(a.Screen[i]); j++ {
			r, g, b, a := a.Screen[i][j].RGBA()
			mapA[color{r, g, b, a}] = true
		}
	}
	logrus.Infof("COLOR: %v", mapA)
	return true
}

func GetPlayerLocation(s *arcade.State) (int, int) {
	for i := 0; i < len(s.Screen); i++ {
		for j := 0; j < len(s.Screen[i]); j++ {
			r, g, b, a := s.Screen[i][j].RGBA()
			if r == 0 && g == 65535 && b == 0 && a == 65535 {
				return i, j
			}
		}
	}
	return -1, -1
}
