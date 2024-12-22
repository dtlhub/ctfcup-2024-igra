package game

import (
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var MazeSolverActive = false

func init() {
	value, ok := os.LookupEnv("MAZE_SOLVER")
	MazeSolverActive = ok && value != "0"
}

const mazeStr = `################################################################
##       #           #       #   #       #               #     #
## ####### ######### # ### # # # # # ##### ### ######### ### # #
## #       #         # #   # # #   #     #   # #   #   #   # # #
## # ####### ########### ### # ######### ### # # # ### ### ### #
## #   #   #   #         #   # #   #   #     #   #   #       # #
## ### # # ### # ### ##### ### # # # # ############# ####### # #
##     # #     # #   #       # # #   #           # #     # #   #
## ##### ####### # ########### # ##### ######### # ##### # ### #
##   # #   #     #               #     #S#   # #   #     # #   #
#### # ### ####### ##### ############# # # # # ### # ##### # ###
## # #   #       # #   # #       #   # #   # #   # # #     #   #
## # ### ####### # # # ### ##### # # ####### # # # # # # ##### #
## #   # #       # # #   # #   #   #       # # # # # # #   #   #
## ### # # # ####### ### # # ############# # # # # # ### # # ###
##   #   # # #     #   #   #           #   # # #   #   # # #   #
#### ### # ### ### ### ##### ######### # ### # ####### # # ### #
##     # #     # #   # # #   #   #     # #   # #   #   # # #   #
## ##### ####### ### # # # ##### # ### # # ### # # # ### # # ###
## #     # #       #   # #     # # #   # # #   # # # # # #   # #
## # ##### # ##### ##### ##### # # ##### # # ### ### # # ##### #
## # #       #   #       #       #     #   #     #   # #   #   #
## # # ######### ##### ### ########### ######### # ### ### # # #
## # # #   #   #   # #         #     #   #   #   # #     # # # #
## # # # # # # # # # ########### ### # # # # ##### ### ### # ###
##   #   #   #   #       #       # #   # # #   #   #   #   #   #
## ################# ##### ####### ##### # ### # ### ### ##### #
## #             #   #     #   #     #   # # #   #       #     #
## # ##### ### ### ### ####### # ### # ### # ##### ####### # ###
## # #   #   #       #   #     # #   # #   #     #   #   # #   #
## ### # ### ########### # ### # # ### # ### ### ### # ### ### #
## #   #   #   #         # #   # #   # #   #   #     #     #   #
## # ##### ##### ######### ##### ### # # # ######### ### ### # #
##   #   #     #       #   #   # #   # # # #       #   # #   # #
###### ####### ##### # ### # # # # # # ### # ##### ### ### ### #
##           # #   # #   #   #   # # #   # # # #   #       #   #
## ####### ### # # ##### # ####### # ### # # # # ####### #######
##       # #   # #     # # #       # # # #   # #       # #     #
######## # # ### ##### # # # ####### # # # ### ####### ### ### #
## #     #   #       #   # # #     #   # #           #     #   #
## # ############### ##### # # ### ### # ####### ##### ##### ###
## #   #             #   # #     #   # #     #   #   # #   #   #
## ### # ##### # ##### # ######### # # ##### ##### # # # # ### #
##     #     # # #     #   #     #E# # #   # #     # # # # #   #
## ######### # ### ####### # ### ### ### # # # ##### ### # # ###
## #     # # # #   #       #   #   #     # # #   # #   # # #   #
## # ### # # # # ### ### ##### ### ### ### # ### # ### # # # # #
##   #   #   #     # #   #       #   # #   #   # #   #   # # # #
###### ##### ####### # ### ######### ### # ### # # # ##### ### #
##   # #   # #       # #   #   #     #   #   # # # #     # #   #
## ### # # ### # ####### ### # # ##### ##### # # # # ##### # # #
##   # # #     # #       #   #   # #     #   #   # #       # # #
## # # # ######### ####### ####### # ### ######### ######### # #
## # # # #   #   # # #   # #       #   #         # #         # #
## # # # # # # # # # # # # # ### # ### ######### # # ##### ### #
## #   #   #   # #   # #   #   # # # # #       # # # #   # #   #
## ######### ### ### # ######### # # # ##### # # # # # # # # ###
##   #     #   # #   #   #       #   #     # # #   #   # # # # #
#### # ### ##### # ##### # ### ########### ### ########### # # #
## #   # #     #   #     #   # #       #   #     #     #   # # #
## ##### ##### ##### ##### # ### ##### # ### ### ### # # ### # #
##                   #     #         #       #       #   #     #
################################################################`

var path []ebiten.Key

func init() {
	lines := strings.Split(mazeStr, "\n")
	height := len(lines)
	width := len(lines[0])
	grid := make([][]rune, height)
	var start, end struct{ x, y int }

	for y := 0; y < height; y++ {
		grid[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			grid[y][x] = rune(lines[y][x])
			if grid[y][x] == 'E' {
				end.x, end.y = x, y
			}
			if y == 1 && grid[y][x] == ' ' && start.x == 0 {
				start.x, start.y = x, y
			}
		}
	}

	type pos struct{ x, y int }
	queue := []pos{{start.x, start.y}}
	visited := make(map[pos]pos)
	visited[pos{start.x, start.y}] = pos{start.x, start.y}

	dirs := []pos{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	found := false

	for len(queue) > 0 && !found {
		curr := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			next := pos{curr.x + d.x, curr.y + d.y}
			if next.x < 0 || next.x >= width || next.y < 0 || next.y >= height {
				continue
			}
			if grid[next.y][next.x] != ' ' && grid[next.y][next.x] != 'E' {
				continue
			}
			if _, ok := visited[next]; ok {
				continue
			}

			visited[next] = curr
			queue = append(queue, next)

			if next.x == end.x && next.y == end.y {
				found = true
				break
			}
		}
	}

	curr := pos{end.x, end.y}
	prev := visited[curr]
	for curr != prev {
		dx := curr.x - prev.x
		dy := curr.y - prev.y

		var key ebiten.Key
		switch {
		case dx == 1:
			key = ebiten.KeyRight
		case dx == -1:
			key = ebiten.KeyLeft
		case dy == 1:
			key = ebiten.KeyDown
		case dy == -1:
			key = ebiten.KeyUp
		}
		path = append([]ebiten.Key{key}, path...)

		curr = prev
		prev = visited[curr]
	}

	logrus.Infof("path: %v", path)
}

type MazeSolver struct {
	nextMove     int
	Active       bool
	ReadyForNext bool
}

func (m *MazeSolver) NextMove() (ebiten.Key, bool) {
	if m.nextMove >= len(path) {
		return ebiten.Key(0), false
	}
	key := path[m.nextMove]
	m.nextMove++
	return key, true
}
