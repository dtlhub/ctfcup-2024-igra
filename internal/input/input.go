package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"

	gameserverpb "github.com/c4t-but-s4d/ctfcup-2024-igra/proto/go/gameserver"
)

var interestingKeys = []ebiten.Key{
	ebiten.KeyA,
	ebiten.KeyB,
	ebiten.KeyC,
	ebiten.KeyD,
	ebiten.KeyE,
	ebiten.KeyF,
	ebiten.KeyG,
	ebiten.KeyH,
	ebiten.KeyI,
	ebiten.KeyJ,
	ebiten.KeyK,
	ebiten.KeyL,
	ebiten.KeyM,
	ebiten.KeyN,
	ebiten.KeyO,
	ebiten.KeyP,
	ebiten.KeyQ,
	ebiten.KeyR,
	ebiten.KeyS,
	ebiten.KeyT,
	ebiten.KeyU,
	ebiten.KeyV,
	ebiten.KeyW,
	ebiten.KeyX,
	ebiten.KeyY,
	ebiten.KeyZ,
	ebiten.KeyDigit0,
	ebiten.KeyDigit1,
	ebiten.KeyDigit2,
	ebiten.KeyDigit3,
	ebiten.KeyDigit4,
	ebiten.KeyDigit5,
	ebiten.KeyDigit6,
	ebiten.KeyDigit7,
	ebiten.KeyDigit8,
	ebiten.KeyDigit9,
	ebiten.KeySpace,
	ebiten.KeyComma,
	ebiten.KeyPeriod,
	ebiten.KeySlash,
	ebiten.KeyEscape,
	ebiten.KeyEnter,
	ebiten.KeyBackspace,
	ebiten.KeyShiftLeft,
	ebiten.KeyArrowUp,
	ebiten.KeyArrowDown,
	ebiten.KeyArrowLeft,
	ebiten.KeyArrowRight,
	//
	ebiten.KeyMinus,
	ebiten.KeyEqual,
	ebiten.KeyShift,
}

type Input struct {
	pressedKeys      map[ebiten.Key]struct{}
	newlyPressedKeys map[ebiten.Key]struct{}
}

func New() *Input {
	return &Input{
		pressedKeys:      make(map[ebiten.Key]struct{}),
		newlyPressedKeys: make(map[ebiten.Key]struct{}),
	}
}

func NewFromProto(p *gameserverpb.ClientEvent_KeysPressed) *Input {
	i := New()
	for _, key := range p.KeysPressed {
		i.pressedKeys[ebiten.Key(key)] = struct{}{}
	}
	for _, key := range p.NewKeysPressed {
		i.newlyPressedKeys[ebiten.Key(key)] = struct{}{}
	}
	return i
}

func (i *Input) UpdateFromRewind(move *gameserverpb.ClientEvent_KeysPressed) {
	i.pressedKeys = make(map[ebiten.Key]struct{}, len(move.KeysPressed))
	for _, key := range move.KeysPressed {
		i.pressedKeys[ebiten.Key(key)] = struct{}{}
	}

	i.newlyPressedKeys = make(map[ebiten.Key]struct{}, len(move.NewKeysPressed))
	for _, key := range move.NewKeysPressed {
		i.newlyPressedKeys[ebiten.Key(key)] = struct{}{}
	}
}

func (i *Input) Update() {
	oldPressedKeys := i.pressedKeys

	i.pressedKeys = make(map[ebiten.Key]struct{})
	i.newlyPressedKeys = make(map[ebiten.Key]struct{})
	for _, key := range interestingKeys {
		if ebiten.IsKeyPressed(key) {
			i.pressedKeys[key] = struct{}{}
			if _, ok := oldPressedKeys[key]; !ok {
				i.newlyPressedKeys[key] = struct{}{}
			}
		}
	}
}

func (i *Input) IsKeyPressed(key ebiten.Key) bool {
	_, ok := i.pressedKeys[key]
	return ok
}

func (i *Input) IsKeyNewlyPressed(key ebiten.Key) bool {
	_, ok := i.newlyPressedKeys[key]
	return ok
}

func (i *Input) JustPressedKeys() []ebiten.Key {
	return lo.Uniq(lo.Keys(i.newlyPressedKeys))
}

func (i *Input) PressedKeys() []ebiten.Key {
	return lo.Uniq(lo.Keys(i.pressedKeys))
}

func (i *Input) ToProto() *gameserverpb.ClientEvent_KeysPressed {
	return &gameserverpb.ClientEvent_KeysPressed{
		KeysPressed: lo.Map(lo.Keys(i.pressedKeys), func(key ebiten.Key, _ int) int32 {
			return int32(key)
		}),
		NewKeysPressed: lo.Map(lo.Keys(i.newlyPressedKeys), func(key ebiten.Key, _ int) int32 {
			return int32(key)
		}),
	}
}

func (i *Input) AddKeyPressed(key ebiten.Key) {
	i.pressedKeys[key] = struct{}{}
}

func (i *Input) RemoveKeyPressed(key ebiten.Key) {
	delete(i.pressedKeys, key)
}

func (i *Input) AddKeyNewlyPressed(key ebiten.Key) {
	i.newlyPressedKeys[key] = struct{}{}
}

func (i *Input) RemoveKeyNewlyPressed(key ebiten.Key) {
	delete(i.newlyPressedKeys, key)
}
