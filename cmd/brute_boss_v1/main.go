package main

import (
	"math/rand/v2"

	"github.com/sirupsen/logrus"
)

const windowSize = 1125

const topLeftX = 1664
const topLeftY = 2208
const bottomRightX = 2752
const bottomRightY = 2176
const bossX = 2024
const bossY = 1746
const maxAllowedDistance = 100

func randInt(rnd *rand.Rand, low, high int) int {
	return low + rnd.IntN(high-low)
}

func isGood(seed int) bool {
	rnd := rand.New(rand.NewPCG(0, uint64(seed)))
	for range 10 {
		dx := float64(randInt(rnd, -500, 500))
		dy := float64(randInt(rnd, -500, 500))

		x := bossX + dx
		y := bossY + dy

		if topLeftX < x && x < bottomRightX && topLeftY < y && y < bottomRightY {
			return false
		}
	}
	return true
}

type Sequence struct {
	start int
	seq   [windowSize * 2]bool
}

func NewSequence() *Sequence {
	s := &Sequence{
		start: 0,
		seq:   [windowSize * 2]bool{},
	}
	for i := range s.seq {
		s.seq[i] = isGood(i)
	}
	return s
}

func (s *Sequence) MoveWindow() {
	copy(s.seq[:windowSize], s.seq[windowSize:])
	for i := windowSize; i < len(s.seq); i++ {
		s.seq[i] = isGood(i)
	}
	s.start += windowSize
}

func (s *Sequence) CheckWindows() {
	logrus.Infof("start=%d", s.start)
	collected := 0
	for i := 0; i < 2*windowSize; i++ {
		if s.seq[i] {
			collected++
		} else {
			logrus.Infof("collected=%d, start_seed=%d", collected, i-collected)
			collected = 0
		}
	}
}

func main() {
	s := NewSequence()
	for {
		s.CheckWindows()
		s.MoveWindow()
	}
}
