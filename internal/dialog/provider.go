package dialog

import (
	"fmt"
)

type Provider interface {
	Get(id string) (Dialog, error)
	DisplayInput() bool
}

type StandardProvider struct {
	showInput bool
}

func NewStandardProvider(showInput bool) *StandardProvider {
	return &StandardProvider{
		showInput: showInput,
	}
}

func (sp *StandardProvider) DisplayInput() bool {
	return sp.showInput
}

func (sp *StandardProvider) Get(id string) (Dialog, error) {
	switch id {
	case "test-npc":
		return NewDummy("Hello, I'm a test NPC!\n 2 + 2 = ?", "4"), nil
	case "rop-npc":
		return NewBinary("./internal/resources/dialogs/rop", "HELLO, COMRAD.CAN YOU GIVE ME THE ADDRESS FROM WHERE THEY WANT TO ATTACK US?", "you win", true), nil
	case "crackme-npc":
		return NewBinary("./internal/resources/dialogs/crackme", "Give me the password to enter the covenant:", "gj", false), nil
	case "guess-npc":
		return NewGuess("Hello.")
	case "steve-npc":
		return NewSteve(), nil
	case "khajiit-npc":
		return NewKhajiit("CD Player", 1000), nil
	case "slon-npc":
		return NewBinary("./internal/resources/dialogs/slon/slon.js", "Hello, I'm slonik! Let's play the game!\n What's your name: ", "YOU WIN", false), nil
	default:
		return nil, fmt.Errorf("unknown dialog id: %s", id)
	}
}

type ClientProvider struct{}

func (cp *ClientProvider) Get(_ string) (Dialog, error) {
	return NewClientDialog(), nil
}

func (cp *ClientProvider) DisplayInput() bool {
	return true
}
