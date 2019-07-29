package gbc

import "time"

const (
	// ClockSpeed is the number of CPU cycles emulated per second. Increase to speed-up game time.
	ClockSpeed = 4194304

	// FPS is the target frames per second for the GPU
	FPS = 60
)

// GBC is the toplevel struct containing all the gameboy systems
type GBC struct {
	mmu *MMU
	cpu *CPU
	gpu *GPU
}

// NewGBC constructs a valid GBC struct
func NewGBC() *GBC {
	gbc := new(GBC)

	gbc.mmu = NewMMU()
	gbc.cpu = NewCPU(gbc.mmu)
	gbc.gpu = NewGPU(gbc.mmu)

	return gbc
}

// GetFrameTime returns the real-time duration of a single frame
func (gbc *GBC) GetFrameTime() time.Duration {
	return time.Second / FPS
}

// Tick runs the gameboy for a single frame-time
func (gbc *GBC) Tick() {
	gbc.cpu.ProcessNextInstruction()
}
