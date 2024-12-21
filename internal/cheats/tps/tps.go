package tps

import (
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var breakpoints = []int{
	1,
	2,
	3,
	5,
	10,
	20,
	30,
	60,
	120,
	180,
	300,
	600,
	1500,
	6000,
}

func Increment() {
	current := ebiten.TPS()
	for i := 0; i < len(breakpoints); i++ {
		if breakpoints[i] > current {
			ebiten.SetTPS(breakpoints[i])
			return
		}
	}
}

func Decrement() {
	current := ebiten.TPS()
	for i := len(breakpoints) - 1; i >= 0; i-- {
		if breakpoints[i] < current {
			ebiten.SetTPS(breakpoints[i])
			return
		}
	}
}

func SetFromEnv(varname string, defaultValue int) {
	tps := defaultValue
	if value := os.Getenv(varname); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			logrus.Errorf("invalid TPS: %v", err)
		} else {
			tps = parsed
		}
	}
	ebiten.SetTPS(tps)
}
