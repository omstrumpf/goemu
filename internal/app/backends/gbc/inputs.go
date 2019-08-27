package gbc

import (
	"log"

	"github.com/omstrumpf/goemu/internal/app/buttons"
)

type inputMemoryDevice struct {
	standardButtons    byte
	directionalButtons byte
	directional        bool
}

func newInputMemoryDevice() *inputMemoryDevice {
	imd := new(inputMemoryDevice)

	imd.directional = true
	imd.standardButtons = 0xF
	imd.directionalButtons = 0xF

	return imd
}

func (imd *inputMemoryDevice) PressButton(b buttons.Button) {
	switch b {
	case buttons.ButtonDown:
		imd.directionalButtons &= ^byte(1 << 3)
	case buttons.ButtonUp:
		imd.directionalButtons &= ^byte(1 << 2)
	case buttons.ButtonLeft:
		imd.directionalButtons &= ^byte(1 << 1)
	case buttons.ButtonRight:
		imd.directionalButtons &= ^byte(1)
	case buttons.ButtonStart:
		imd.standardButtons &= ^byte(1 << 3)
	case buttons.ButtonSelect:
		imd.standardButtons &= ^byte(1 << 2)
	case buttons.ButtonB:
		imd.standardButtons &= ^byte(1 << 1)
	case buttons.ButtonA:
		imd.standardButtons &= ^byte(1)
	default:
		log.Printf("Attempted to press unrecognized button %d", b)
	}
}

func (imd *inputMemoryDevice) ReleaseButton(b buttons.Button) {
	switch b {
	case buttons.ButtonDown:
		imd.directionalButtons |= byte(1 << 3)
	case buttons.ButtonUp:
		imd.directionalButtons |= byte(1 << 2)
	case buttons.ButtonLeft:
		imd.directionalButtons |= byte(1 << 1)
	case buttons.ButtonRight:
		imd.directionalButtons |= byte(1)
	case buttons.ButtonStart:
		imd.standardButtons |= byte(1 << 3)
	case buttons.ButtonSelect:
		imd.standardButtons |= byte(1 << 2)
	case buttons.ButtonB:
		imd.standardButtons |= byte(1 << 1)
	case buttons.ButtonA:
		imd.standardButtons |= byte(1)
	default:
		log.Printf("Attempted to release unrecognized button %d", b)
	}
}

func (imd *inputMemoryDevice) Read(addr uint16) byte {
	if imd.directional {
		return 0xE0 | imd.directionalButtons
	}

	return 0xD0 | imd.standardButtons
}

func (imd *inputMemoryDevice) Write(addr uint16, val byte) {
	if val&1<<4 == 0 {
		imd.directional = true
	}
	if val&1<<5 == 0 {
		imd.directional = false
	}
}
