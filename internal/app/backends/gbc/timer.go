package gbc

import (
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/interrupts"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// TODO write timer tests

// Timer is the gameboy's timer device
type Timer struct {
	mmu *MMU

	div  byte
	tima byte
	tma  byte
	tac  byte

	speedMod       int
	timerCounter   int
	dividerCounter int

	enable bool
}

// NewTimer constructs a valid Timer struct
func NewTimer(mmu *MMU) *Timer {
	t := new(Timer)

	t.mmu = mmu

	return t
}

// RunForClocks runs the Timer for the given number of clock cycles
func (t *Timer) RunForClocks(clocks int) {

	for c := 0; c < clocks; c++ {
		t.dividerCounter++

		if t.dividerCounter%64 == 0 { // 16384 Hz
			t.div++
		}

		if t.enable {
			t.timerCounter++

			if t.timerCounter%t.speedMod == 0 {
				t.tima++

				if t.tima == 0 {
					t.tima = t.tma
					t.mmu.interrupts.Request(interrupts.TimerBit)
				}
			}
		}
	}
}

func (t *Timer) Read(addr uint16) byte {
	switch addr {
	case 0xFF04:
		return t.div
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}

	log.Warningf("Encountered unexpected timer read: %#4x", addr)
	return 0xFF
}

func (t *Timer) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF04:
		t.div = 0
		t.dividerCounter = 0
		t.timerCounter = 0
		return
	case 0xFF05:
		t.tima = val
		t.timerCounter = 0
		return
	case 0xFF06:
		t.tma = val
		return
	case 0xFF07:
		switch val & 0x3 {
		case 0:
			t.speedMod = 256 // 4096 Hz
		case 1:
			t.speedMod = 4 // 262144 Hz
		case 2:
			t.speedMod = 16 // 65536 Hz
		case 3:
			t.speedMod = 64 // 16384 Hz
		}

		t.enable = (val&0x4 != 0)
		return
	}

	log.Warningf("Encountered unexpected timer write: %#4x", addr)
}
