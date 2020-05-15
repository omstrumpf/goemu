package gbc

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/cartridge"
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/interrupts"
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/log"
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

// GBC is the toplevel struct containing all the gameboy systems
type GBC struct {
	mmu   *MMU
	cpu   *CPU
	cart  *cartridge.CART
	ppu   *PPU
	timer *Timer

	totalClocks uint64
}

// NewGBC constructs a valid GBC struct
func NewGBC(skiplogo bool, rom []byte) *GBC {
	gbc := new(GBC)

	gbc.cart = cartridge.NewCart(rom)
	gbc.mmu = NewMMU(gbc.cart.BankController)
	gbc.timer = NewTimer(gbc.mmu)
	gbc.cpu = NewCPU(gbc.mmu)
	gbc.ppu = NewPPU(gbc.mmu)

	gbc.mmu.ppu = gbc.ppu
	gbc.mmu.timer = gbc.timer

	if skiplogo {
		gbc.skipLogo()
	}

	return gbc
}

// Set the gameboy to the correct post-boot state
func (gbc *GBC) skipLogo() {
	log.Debugf("Skipping logo boot sequence")
	gbc.cpu.PC.Set(0x0100)
	gbc.cpu.AF.Set(0x01B0)
	gbc.cpu.BC.Set(0x0013)
	gbc.cpu.DE.Set(0x00D8)
	gbc.cpu.HL.Set(0x014D)
	gbc.cpu.SP.Set(0xFFFE)
	gbc.mmu.Write(0xFF05, 0x00) // TIMA
	gbc.mmu.Write(0xFF06, 0x00) // TMA
	gbc.mmu.Write(0xFF07, 0x00) // TAC
	gbc.mmu.Write(0xFF10, 0x80) // NR10
	gbc.mmu.Write(0xFF11, 0xBF) // NR11
	gbc.mmu.Write(0xFF12, 0xF3) // NR12
	gbc.mmu.Write(0xFF14, 0xBF) // NR14
	gbc.mmu.Write(0xFF16, 0x3F) // NR21
	gbc.mmu.Write(0xFF17, 0x00) // NR22
	gbc.mmu.Write(0xFF19, 0xBF) // NR24
	gbc.mmu.Write(0xFF1A, 0x7F) // NR30
	gbc.mmu.Write(0xFF1B, 0xFF) // NR31
	gbc.mmu.Write(0xFF1C, 0x9F) // NR32
	gbc.mmu.Write(0xFF1E, 0xBF) // NR33
	gbc.mmu.Write(0xFF20, 0xFF) // NR41
	gbc.mmu.Write(0xFF21, 0x00) // NR42
	gbc.mmu.Write(0xFF22, 0x00) // NR43
	gbc.mmu.Write(0xFF23, 0xBF) // NR30
	gbc.mmu.Write(0xFF24, 0x77) // NR50
	gbc.mmu.Write(0xFF25, 0xF3) // NR51
	gbc.mmu.Write(0xFF26, 0xF1) // NR52
	gbc.mmu.Write(0xFF40, 0x91) // LCDC
	gbc.mmu.Write(0xFF42, 0x00) // SCY
	gbc.mmu.Write(0xFF43, 0x00) // SCX
	gbc.mmu.Write(0xFF45, 0x00) // LYC
	gbc.mmu.Write(0xFF47, 0xFC) // BGP
	gbc.mmu.Write(0xFF48, 0xFF) // OBP0
	gbc.mmu.Write(0xFF49, 0xFF) // OBP1
	gbc.mmu.Write(0xFF4A, 0x00) // WY
	gbc.mmu.Write(0xFF4B, 0x00) // WX
	gbc.mmu.Write(0xFFFF, 0x00) // IE
	gbc.mmu.DisableBios()
}

// Tick runs the gameboy for a single frame-time
func (gbc *GBC) Tick() {
	clocks := 0

	for clocks < CyclesPerFrame {
		if log.Logger.IsTraceEnabled() {
			fmt.Fprintln(os.Stderr, gbc.traceString()) // Bypassing log for speed and to avoid verbose prints
		}

		c := gbc.cpu.ProcessNextInstruction()
		clocks += c
		gbc.totalClocks += uint64(c)
		gbc.ppu.RunForClocks(c)
		gbc.timer.RunForClocks(c)
	}
}

// PressButton presses the given button
func (gbc *GBC) PressButton(b console.Button) {
	gbc.mmu.inputs.PressButton(b)
	gbc.mmu.interrupts.Request(interrupts.JoypadBit)
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

// traceString produces a string of the current GBC trace, for debugging
func (gbc *GBC) traceString() string {
	pc := gbc.cpu.PC.HiLo()
	_, disassembly := gbc.cpu.Disassemble(pc)

	return fmt.Sprintf("A: %02x, F: %s, BC: %04x, DE: %04x, HL: %04x, SP: %04x, (HL): %02x, ppu: %d. %#04x: %s",
		gbc.cpu.AF.Hi(),
		gbc.cpu.flagString(),
		gbc.cpu.BC.HiLo(),
		gbc.cpu.DE.HiLo(),
		gbc.cpu.HL.HiLo(),
		gbc.cpu.SP.HiLo(),
		gbc.mmu.Read(gbc.cpu.HL.HiLo()),
		gbc.ppu.mode,
		pc,
		disassembly)
}
