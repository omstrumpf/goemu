package gbc

import (
	"image/color"
	"time"
)

const (
	// ClockSpeed is the number of CPU cycles emulated per second. Increase to speed-up game time.
	ClockSpeed = 4194304

	// FPS is the target frames per second for the PPU
	FPS = 60

	// FrameTime is the real-time duration of a single frame
	FrameTime = time.Second / FPS

	// CyclesPerFrame is the number of cycles in a single frame
	CyclesPerFrame = ClockSpeed / FPS

	// ScreenWidth is the width of the screen
	ScreenWidth = 160

	// ScreenHeight is the height of the screen
	ScreenHeight = 144

	// ConsoleName is the name of this console
	ConsoleName = "GameBoy Color"
)

// GBC is the toplevel struct containing all the gameboy systems
type GBC struct {
	mmu *MMU
	cpu *CPU
	ppu *PPU
}

// NewGBC constructs a valid GBC struct
func NewGBC() *GBC {
	gbc := new(GBC)

	gbc.mmu = NewMMU()
	gbc.cpu = NewCPU(gbc.mmu)
	gbc.ppu = NewPPU(gbc.mmu)

	return gbc
}

// Tick runs the gameboy for a single frame-time
func (gbc *GBC) Tick() {
	clockStart := gbc.cpu.clock

	for gbc.cpu.clock-clockStart < CyclesPerFrame {
		gbc.cpu.ProcessNextInstruction()
		gbc.ppu.UpdateToClock(gbc.cpu.clock)
	}
}

// LoadROM loads the given ROM bytes into the MMU
func (gbc *GBC) LoadROM(rom []byte) {
	gbc.mmu.LoadROM(rom)
}

// GetFrameBuffer returns the gameboy's frame buffer, a slice of RGBA values
func (gbc *GBC) GetFrameBuffer() []color.RGBA {
	return gbc.ppu.framebuffer
}

// GetFrameTime returns the real-time duration of a single frame
func (gbc *GBC) GetFrameTime() time.Duration {
	return FrameTime
}

// GetScreenWidth returns the width of the screen
func (gbc *GBC) GetScreenWidth() int {
	return ScreenWidth
}

// GetScreenHeight returns the height of the screen
func (gbc *GBC) GetScreenHeight() int {
	return ScreenHeight
}

// GetConsoleName returns the name of this console
func (gbc *GBC) GetConsoleName() string {
	return ConsoleName
}
