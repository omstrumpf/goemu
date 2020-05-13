package gbc

import (
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/log"
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

func (imd *inputMemoryDevice) PressButton(b console.Button) {
	switch b {
	case console.ButtonDown:
		imd.directionalButtons &= ^byte(1 << 3)
	case console.ButtonUp:
		imd.directionalButtons &= ^byte(1 << 2)
	case console.ButtonLeft:
		imd.directionalButtons &= ^byte(1 << 1)
	case console.ButtonRight:
		imd.directionalButtons &= ^byte(1)
	case console.ButtonStart:
		imd.standardButtons &= ^byte(1 << 3)
	case console.ButtonSelect:
		imd.standardButtons &= ^byte(1 << 2)
	case console.ButtonB:
		imd.standardButtons &= ^byte(1 << 1)
	case console.ButtonA:
		imd.standardButtons &= ^byte(1)
	default:
		log.Warningf("Attempted to press unrecognized button %d", b)
	}
}

func (imd *inputMemoryDevice) ReleaseButton(b console.Button) {
	switch b {
	case console.ButtonDown:
		imd.directionalButtons |= byte(1 << 3)
	case console.ButtonUp:
		imd.directionalButtons |= byte(1 << 2)
	case console.ButtonLeft:
		imd.directionalButtons |= byte(1 << 1)
	case console.ButtonRight:
		imd.directionalButtons |= byte(1)
	case console.ButtonStart:
		imd.standardButtons |= byte(1 << 3)
	case console.ButtonSelect:
		imd.standardButtons |= byte(1 << 2)
	case console.ButtonB:
		imd.standardButtons |= byte(1 << 1)
	case console.ButtonA:
		imd.standardButtons |= byte(1)
	default:
		log.Warningf("Attempted to release unrecognized button %d", b)
	}
}

func (imd *inputMemoryDevice) Read(addr uint16) byte {
	if imd.directional {
		return 0xE0 | imd.directionalButtons
	}

	return 0xD0 | imd.standardButtons
}

func (imd *inputMemoryDevice) Write(addr uint16, val byte) {
	if val&(1<<4) == 0 {
		imd.directional = true
	}
	if val&(1<<5) == 0 {
		imd.directional = false
	}
}
