package gbc

import (
	"image/color"
	"log"
)

// PPU represents the gameboy's graphics processing unit.
type PPU struct {
	mmu *MMU // Memory Management Unit

	memoryControl *ppuControl // Control Registers

	framebuffer []color.RGBA // Frame Buffer

	registers []byte // Control Registers

	mode       uint8  // Mode Number (0: HBLANK, 1: VBLANK, 2: OAM, 3: VRAM)
	clock      int    // PPU clock
	timeInMode int    // Number of clock cycles spent in the current mode
	line       uint16 // Line currently being processed

	bgMap   uint8 // Which background map is in use
	tileSet uint8 // Which tileset is in use

	scrollX uint16 // Background scroll X
	scrollY uint16 // Background scroll Y
}

// NewPPU constructs a valid PPU struct
func NewPPU(mmu *MMU) *PPU {
	ppu := new(PPU)

	ppu.mmu = mmu

	ppu.framebuffer = make([]color.RGBA, 160*144)
	ppu.registers = make([]byte, 64)

	ppu.mode = 2 // Start in OAM mode

	ppuControl := new(ppuControl)
	ppuControl.ppu = ppu
	ppu.memoryControl = ppuControl

	return ppu
}

// UpdateToClock runs the PPU until the given clock cycle
func (ppu *PPU) UpdateToClock(clock int) {
	for ppu.clock <= clock {
		switch ppu.mode {
		case 0: // HBLANK
			if ppu.timeInMode == 51 {
				ppu.timeInMode = 0

				ppu.line++

				if ppu.line == 143 {

					ppu.mode = 1
				} else {
					ppu.mode = 2
				}
			} else {
				advance := min(clock-ppu.clock, 51-ppu.timeInMode)
				ppu.clock += advance
				ppu.timeInMode += advance
			}
		case 1: // VBLANK
			if ppu.timeInMode == 1140 {
				ppu.timeInMode = 0

				ppu.mode = 2
				ppu.line = 0
			} else {
				advance := min(clock-ppu.clock, 1140-ppu.timeInMode)
				ppu.clock += advance
				ppu.timeInMode += advance
			}
		case 2: // OAM
			if ppu.timeInMode == 20 {
				ppu.timeInMode = 0

				ppu.mode = 3
			} else {
				advance := min(clock-ppu.clock, 20-ppu.timeInMode)
				ppu.clock += advance
				ppu.timeInMode += advance
			}
		case 3: // VRAM
			if ppu.timeInMode == 43 {
				ppu.timeInMode = 0

				ppu.mode = 0

				ppu.renderLine()
			} else {
				advance := min(clock-ppu.clock, 43-ppu.timeInMode)
				ppu.clock += advance
				ppu.timeInMode += advance
			}
		}
	}
}

func (ppu *PPU) renderLine() {
	var bgAddr uint16

	// Base VRAM address for the background map
	if ppu.bgMap == 0 {
		bgAddr = 0x9800
	} else {
		bgAddr = 0x9C00
	}

	// The first tile pointer in this line
	topLeft := bgAddr + (((ppu.line + ppu.scrollY) & 0x00FF) >> 3)

	// The first tile pointer to be used
	tilePointer := topLeft + (ppu.scrollX >> 3)

	// Address of the current tile data
	tileAddr := ppu.getTileAddress(tilePointer)

	tileX := ppu.scrollX & 0x7              // Which column within the tile to start at
	tileY := (ppu.scrollY + ppu.line) & 0x7 // Which row within the tile to use

	screenX := 0             // Which column in the framebuffer to start at
	screenY := int(ppu.line) // Which row in the framebuffer to use

	for screenX < 160 {
		// Move to next tile if needed
		if tileX == 8 {
			tileX = 0
			tilePointer++
			tileAddr = ppu.getTileAddress(tilePointer)
		}

		tileRow := ppu.mmu.Read16(tileAddr + (tileY * 2))

		val := byte(tileRow >> (14 - (2 * tileX)))

		pixel := gbToRGBA(val)
		ppu.writePixel(pixel, screenX, screenY)
		screenX++
	}
}

func (ppu *PPU) writePixel(val color.RGBA, x int, y int) {
	ppu.framebuffer[(y*160)+x] = val
}

func (ppu *PPU) getTileAddress(mapAddr uint16) uint16 {
	switch ppu.getTileSet() {
	case 0:
		tileNum := uint16(ppu.mmu.Read(mapAddr))
		return uint16(0x9000) + (tileNum * 16)
	case 1:
		tileNum := int16(int8(ppu.mmu.Read(mapAddr)))
		return uint16(int32(0x8000) + int32(tileNum*16))
	}

	log.Printf("Unrecognized tileset %d", ppu.tileSet)
	return 0x9000
}

func (ppu *PPU) getWindowEnable() bool {
	return ppu.registers[0]&0x20 != 0
}

func (ppu *PPU) getTileSet() int {
	if ppu.registers[0]&0x10 == 0 {
		return 0
	}
	return 1
}

//// Control Registers ////
type ppuControl struct {
	ppu *PPU
}

func (ppc *ppuControl) Read(addr uint16) byte {
	return 0 // TODO
}

func (ppc *ppuControl) Write(addr uint16, val byte) {
	// TODO
}

//// Helpers ////
func gbToRGBA(val byte) color.RGBA {
	switch val {
	case 0:
		return color.RGBA{0, 0, 0, 0xFF}
		// return color.RGBA{255, 255, 255, 0xFF}
	case 1:
		return color.RGBA{192, 192, 192, 0xFF}
	case 2:
		return color.RGBA{96, 96, 96, 0xFF}
	case 3:
		return color.RGBA{0, 0, 0, 0xFF}
	}

	log.Printf("Unrecognized gb color value: %d", val)
	return color.RGBA{255, 255, 255, 0xFF}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
