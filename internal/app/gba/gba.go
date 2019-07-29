package gba

import "time"

const (
	// ClockSpeed is the number of CPU cycles emulated per second. Increase to speed-up game time.
	ClockSpeed = 4194304

	// FPS is the target frames per second for the GPU
	FPS = 60

	// SecondsPerFrame is the number of real-time seconds per GPU frame
	SecondsPerFrame = time.Second / FPS

	// CyclesPerFrame is the number of CPU cycles per GPU frame
	CyclesPerFrame = ClockSpeed / FPS
)

// GBA is the toplevel struct containing all the gameboy systems
type GBA struct {
	mmu *MMU
	cpu *CPU
	gpu *GPU
}

// NewGBA constructs a valid GBA struct
func NewGBA() *GBA {
	gba := new(GBA)

	gba.mmu = NewMMU()
	gba.cpu = NewCPU(gba.mmu)
	gba.gpu = NewGPU(gba.mmu)

	return gba
}

// Tick runs the gameboy for a single frame-time
func (gba *GBA) Tick() {
	gba.cpu.ProcessNextInstruction()
}
