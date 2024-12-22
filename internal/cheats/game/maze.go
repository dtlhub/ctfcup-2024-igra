package game

import (
	"os"

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
	// lines := strings.Split(mazeStr, "\n")
	// height := len(lines)
	// width := len(lines[0])
	// grid := make([][]rune, height)
	// var start, end struct{ y, x int }

	// for y := 0; y < height; y++ {
	// 	grid[y] = make([]rune, width)
	// 	for x := 0; x < width; x++ {
	// 		grid[y][x] = rune(lines[y][x])
	// 		if grid[y][x] == 'E' {
	// 			end.y, end.x = y, x
	// 		}
	// 		if grid[y][x] == 'S' {
	// 			start.y, start.x = y, x
	// 			grid[y][x] = ' ' // Mark as empty space after finding start
	// 		}
	// 	}
	// }

	// if start.x == 0 && start.y == 0 {
	// 	logrus.Fatal("Start position not found in maze")
	// }

	// type pos struct{ y, x int }
	// queue := []pos{{start.y, start.x}}
	// visited := make(map[pos]pos)
	// visited[pos{start.y, start.x}] = pos{start.y, start.x}

	// dirs := []pos{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	// found := false

	// for len(queue) > 0 && !found {
	// 	curr := queue[0]
	// 	queue = queue[1:]

	// 	for _, d := range dirs {
	// 		next := pos{curr.y + d.y, curr.x + d.x}
	// 		if next.y < 0 || next.y >= height || next.x < 0 || next.x >= width {
	// 			continue
	// 		}
	// 		if grid[next.y][next.x] != ' ' && grid[next.y][next.x] != 'E' {
	// 			continue
	// 		}
	// 		if _, ok := visited[next]; ok {
	// 			continue
	// 		}

	// 		visited[next] = curr
	// 		queue = append(queue, next)

	// 		if next.y == end.y && next.x == end.x {
	// 			found = true
	// 			break
	// 		}
	// 	}
	// }

	// curr := pos{end.y, end.x}
	// prev := visited[curr]
	// for curr != prev {
	// 	dx := curr.x - prev.x
	// 	dy := curr.y - prev.y

	// 	var key ebiten.Key
	// 	switch {
	// 	case dx == 1:
	// 		key = ebiten.KeyRight
	// 	case dx == -1:
	// 		key = ebiten.KeyLeft
	// 	case dy == 1:
	// 		key = ebiten.KeyDown
	// 	case dy == -1:
	// 		key = ebiten.KeyUp
	// 	}
	// 	path = append([]ebiten.Key{key}, path...)

	// 	curr = prev
	// 	prev = visited[curr]
	// }
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
	logrus.Infof("next move: %v", key)
	return key, true
}
