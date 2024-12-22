package boss

import (
	"math/rand/v2"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/cheats"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/damage"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/geometry"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/object"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ BOSS = &V1{}

type V1Stage string

const (
	V1StageInitial V1Stage = "initial"
	V1StageHealing V1Stage = "healing"
	V1StageDeath   V1Stage = "death"
)

type V1 struct {
	*object.Rendered `msgpack:"object"`

	BulletImg *ebiten.Image

	outside int `msgpack:"-"`
	inside  int `msgpack:"-"`

	stage           V1Stage        `msgpack:"stage"`
	health          int            `msgpack:"health"`
	startTick       int            `msgpack:"start_tick"`
	initialLocation geometry.Point `msgpack:"initial_location"`
}

func NewV1(obj *object.Rendered, bulletImg *ebiten.Image) *V1 {
	v1 := &V1{
		Rendered:        obj,
		BulletImg:       bulletImg,
		initialLocation: obj.Origin,
	}
	v1.Reset()
	return v1
}

func (b *V1) Object() *object.Rendered {
	return b.Rendered
}

func (b *V1) Tick(s *TickState) *TickResult {
	if b.startTick == 0 {
		b.startTick = s.CurrentTick
	}

	if b.health <= 0 {
		return &TickResult{Dead: true}
	}

	if b.stage == V1StageInitial && b.health < 100 {
		b.stage = V1StageHealing
	}

	if b.stage == V1StageHealing {
		if (s.CurrentTick-b.startTick)%20 == 0 {
			b.health += 50
		}
		if b.health >= 300 {
			b.health = 300
			b.stage = V1StageDeath
		}
	}

	const (
		bulletCount           = 1
		bulletSpawnSquareSize = 50
		bulletSpeed           = 2
	)
	var bulletDamage = 1
	if cheats.Enabled {
		bulletDamage = 0
	}

	res := &TickResult{}
	switch b.stage {
	case V1StageInitial, V1StageHealing:
		if (s.CurrentTick-b.startTick)%30 == 0 {
			b.health -= 4
		}

		rnd := rand.New(rand.NewPCG(0, uint64(s.CurrentTick)))
		for range bulletCount {
			dx := float64(randInt(rnd, -bulletSpawnSquareSize, bulletSpawnSquareSize))
			dy := float64(randInt(rnd, -bulletSpawnSquareSize, bulletSpawnSquareSize))
			res.Bullets = append(res.Bullets, damage.NewBullet(
				b.Origin.Add(geometry.Vector{X: dx, Y: dy}),
				b.BulletImg,
				bulletDamage,
				geometry.Vector{X: dx, Y: dy},
				bulletSpeed,
			))
		}

	case V1StageDeath:
		if (s.CurrentTick-b.startTick)%30 == 0 {
			// logrus.WithFields(logrus.Fields{
			// 	"outside":       b.outside,
			// 	"inside":        b.inside,
			// 	"outside_ratio": float64(b.outside) / (float64(b.outside) + float64(b.inside)),
			// }).Infof("health=%d", b.health)
			b.health -= 8
		}

		const scalingFactor = 10

		rnd := rand.New(rand.NewPCG(0, uint64(s.CurrentTick)))
		for range bulletCount * scalingFactor {
			dx := float64(randInt(rnd, -bulletSpawnSquareSize*scalingFactor, bulletSpawnSquareSize*scalingFactor))
			dy := float64(randInt(rnd, -bulletSpawnSquareSize*scalingFactor, bulletSpawnSquareSize*scalingFactor))

			// origin := b.Origin.Add(geometry.Vector{X: dx, Y: dy})
			// if origin.X < 1664 || origin.X > 2208 || origin.Y < 2176 || origin.Y > 2752 {
			// 	b.outside++
			// } else {
			// 	b.inside++
			// }

			res.Bullets = append(res.Bullets, damage.NewBullet(
				b.Origin.Add(geometry.Vector{X: dx, Y: dy}),
				b.BulletImg,
				bulletDamage,
				geometry.Vector{X: dx, Y: dy},
				bulletSpeed*scalingFactor,
			))
		}
	}

	res.Dead = b.health <= 0

	if !res.Dead {
		rnd := rand.New(rand.NewPCG(0, uint64(s.CurrentTick)))
		dx := float64(randInt(rnd, -10, 10))
		dy := float64(randInt(rnd, -10, 10))

		newLocation := b.Rendered.Origin.Add(geometry.Vector{X: dx, Y: dy})
		if newLocation.SubPoint(b.initialLocation).Length() > 100 {
			dx = -dx
			dy = -dy
		}

		b.Move(geometry.Vector{X: dx, Y: dy})
	}

	return res
}

func (b *V1) Health() *HealthState {
	return &HealthState{
		Health:    b.health,
		MaxHealth: 300,
	}
}

func (b *V1) Reset() {
	b.health = 300
	b.stage = V1StageInitial
	b.startTick = 0
	b.Rendered.MoveTo(b.initialLocation)
}

func randInt(rnd *rand.Rand, low, high int) int {
	return low + rnd.IntN(high-low)
}
