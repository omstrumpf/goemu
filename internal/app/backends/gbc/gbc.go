package gbc

import (
	"image/color"
	"time"

	"github.com/juju/loggo"
	"github.com/omstrumpf/goemu/internal/app/console"
)

const (
	// ClockSpeed is the number of CPU cycles emulated per second. Increase to speed-up game time.
	ClockSpeed = 1048576

	// FPS is the target frames per second for the display
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

var logger = loggo.GetLogger("goemu.gbc")

// GBC is the toplevel struct containing all the gameboy systems
type GBC struct {
	mmu   *MMU
	cpu   *CPU
	ppu   *PPU
	timer *Timer

	totalClocks uint64
}

// NewGBC constructs a valid GBC struct
func NewGBC() *GBC {
	gbc := new(GBC)

	gbc.mmu = NewMMU()
	gbc.timer = NewTimer(gbc.mmu)
	gbc.cpu = NewCPU(gbc.mmu)
	gbc.ppu = NewPPU(gbc.mmu)

	gbc.mmu.ppu = gbc.ppu
	gbc.mmu.timer = gbc.timer

	return gbc
}

// Tick runs the gameboy for a single frame-time
func (gbc *GBC) Tick() {
	clocks := 0

	for clocks < CyclesPerFrame {
		if logger.IsTraceEnabled() {
			pc := gbc.cpu.PC.HiLo()
			opcode := gbc.mmu.Read(pc)
			logger.Tracef("%016d: %04x=%02x", gbc.totalClocks, pc, opcode)
			if pc == 0xc302 {
				break
			}
		}

		c := gbc.cpu.ProcessNextInstruction()
		clocks += c
		gbc.totalClocks += uint64(c)
		gbc.ppu.RunForClocks(c)
		gbc.timer.RunForClocks(c)
	}
}

// LoadROM loads the given ROM bytes into the MMU
func (gbc *GBC) LoadROM(rom []byte) {
	gbc.mmu.LoadROM(rom)
}

// PressButton presses the given button
func (gbc *GBC) PressButton(b console.Button) {
	gbc.mmu.inputs.PressButton(b)
	gbc.mmu.interrupts.RequestInterrupt(interruptJoypadBit)
}

// ReleaseButton releases the given button
func (gbc *GBC) ReleaseButton(b console.Button) {
	gbc.mmu.inputs.ReleaseButton(b)
}

// IsStopped returns true if the gameboy is not running
func (gbc *GBC) IsStopped() bool {
	return gbc.cpu.IsStopped()
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
