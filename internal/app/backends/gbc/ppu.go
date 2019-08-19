package gbc

import (
	"image/color"
	"log"
)

// PPU represents the gameboy's graphics processing unit.
type PPU struct {
	mmu *MMU // Memory Management Unit

	memoryControl *ppuControl // Control Regsiter Device

	framebuffer []color.RGBA  // Frame Buffer
	palette     [4]color.RGBA // Color Palette

	// Control Registers
	lcdEnable    bool // Enables the entire screen
	windowMap    bool // Which window map is in use
	windowEnable bool // Enables the window display
	bgTileSelect bool // Which tileset is in use
	bgMap        bool // Which background map is in use
	spriteSize   bool // Size of the sprite
	spriteEnable bool // Enables the sprite
	bgEnable     bool // Enables rendering the background

	mode        byte // Mode Number (0: HBLANK, 1: VBLANK, 2: OAM, 3: VRAM)
	clock       int  // PPU clock
	timeInMode  int  // Number of clock cycles spent in the current mode
	line        byte // Line currently being processed
	lineCompare byte // Target line for interrupt

	interrupt0   bool // Trigger an interrupt on entering mode 0
	interrupt1   bool // Trigger an interrupt on entering mode 1
	interrupt2   bool // Trigger an interrupt on entering mode 2
	interruptLYC bool // Trigger an interrupt when line matches lineCompare

	scrollX byte // Background scroll X
	scrollY byte // Background scroll Y
}

// NewPPU constructs a valid PPU struct
func NewPPU(mmu *MMU) *PPU {
	ppu := new(PPU)

	ppu.mmu = mmu

	ppu.framebuffer = make([]color.RGBA, 160*144)
	ppu.clearScrean()

	ppu.palette = [4]color.RGBA{
		color.RGBA{255, 255, 255, 0xFF},
		color.RGBA{192, 192, 192, 0xFF},
		color.RGBA{96, 96, 96, 0xFF},
		color.RGBA{0, 0, 0, 0xFF},
	}

	ppu.mode = 2 // Start in OAM mode

	ppuControl := new(ppuControl)
	ppuControl.ppu = ppu
	ppu.memoryControl = ppuControl

	return ppu
}

// UpdateToClock runs the PPU until the given clock cycle
func (ppu *PPU) UpdateToClock(clock int) {
	for ppu.clock < clock {
		switch ppu.mode {
		case 0: // HBLANK
			if ppu.timeInMode == 51 {
				ppu.timeInMode = 0

				ppu.line++

				if ppu.interruptLYC && ppu.line == ppu.lineCompare {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}

				if ppu.line == 143 {
					ppu.mode = 1
					if ppu.interrupt1 {
						ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
					}
				} else {
					ppu.mode = 2
					if ppu.interrupt2 {
						ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
					}
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

				if ppu.interrupt2 {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}
			} else {
				advance := min(clock-ppu.clock, 1140-ppu.timeInMode)
				ppu.clock += advance
				ppu.timeInMode += advance
				ppu.line = 144
				if ppu.interruptLYC && ppu.lineCompare >= 143 {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}
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

				if ppu.interrupt0 {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}

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
	if !ppu.lcdEnable {
		ppu.clearScrean()
		return
	}

	// Base VRAM address for the background map
	var bgAddr uint16
	if ppu.bgMap {
		bgAddr = 0x9C00
	} else {
		bgAddr = 0x9800
	}

	// The first tile pointer to be used
	tileShiftY := uint16((ppu.line+ppu.scrollY)>>3) * 32
	tileShiftX := uint16(ppu.scrollX >> 3)
	tilePointer := bgAddr + tileShiftY + tileShiftX

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

		val := ppu.getTileVal(tileAddr, tileX, tileY)
		pixel := ppu.palette[val]
		ppu.writePixel(pixel, screenX, screenY)

		screenX++
		tileX++
	}
}

// writePixel writes the given RGBA value into the framebuffer at coordinates (x, y)
func (ppu *PPU) writePixel(val color.RGBA, x int, y int) {
	ppu.framebuffer[(y*160)+x] = val
}

// getTileVal returns the 2 bit value for a tile pixel
func (ppu *PPU) getTileVal(tileAddr uint16, tileX byte, tileY byte) byte {
	bit := byte(1 << (7 - tileX))

	tileLo := ppu.mmu.Read(tileAddr + uint16(tileY*2))
	tileHi := ppu.mmu.Read(tileAddr + uint16(tileY*2) + 1)

	val := byte(0)
	if tileLo&bit > 0 {
		val++
	}
	if tileHi&bit > 0 {
		val += 2
	}

	return val
}

// getTileAddress returns the mmu address of the tile pointed to by the mapAddr
func (ppu *PPU) getTileAddress(mapAddr uint16) uint16 {
	if ppu.bgTileSelect {
		tileNum := uint16(ppu.mmu.Read(mapAddr))
		return uint16(0x8000) + (tileNum * 16)
	}

	tileNum := int16(int8(ppu.mmu.Read(mapAddr)))
	return uint16(int32(0x9000) + int32(tileNum*16))
}

// clearScrean sets the framebuffer to all black
func (ppu *PPU) clearScrean() {
	for i := range ppu.framebuffer {
		ppu.framebuffer[i] = color.RGBA{0, 0, 0, 0xFF}
	}
}

//// Control Registers ////
type ppuControl struct {
	ppu *PPU
}

func (ppc *ppuControl) Read(addr uint16) byte {
	switch addr {
	case 0x00:
		var ret byte
		if ppc.ppu.lcdEnable {
			ret |= 0x80
		}
		if ppc.ppu.windowMap {
			ret |= 0x40
		}
		if ppc.ppu.windowEnable {
			ret |= 0x20
		}
		if ppc.ppu.bgTileSelect {
			ret |= 0x10
		}
		if ppc.ppu.bgMap {
			ret |= 0x08
		}
		if ppc.ppu.spriteSize {
			ret |= 0x04
		}
		if ppc.ppu.spriteEnable {
			ret |= 0x02
		}
		if ppc.ppu.bgEnable {
			ret |= 0x01
		}
		return ret
	case 0x01:
		var ret byte
		ret = ppc.ppu.mode
		if ppc.ppu.line == ppc.ppu.lineCompare {
			ret |= 0x04
		}
		if ppc.ppu.interrupt0 {
			ret |= 0x08
		}
		if ppc.ppu.interrupt1 {
			ret |= 0x10
		}
		if ppc.ppu.interrupt2 {
			ret |= 0x20
		}
		if ppc.ppu.interruptLYC {
			ret |= 0x40
		}
		return ret
	case 0x02:
		return ppc.ppu.scrollY
	case 0x03:
		return ppc.ppu.scrollX
	case 0x04:
		return ppc.ppu.line
	case 0x05:
		return ppc.ppu.lineCompare
	}

	log.Printf("Encountered read with unknown PPU control address: %#02x", addr)
	return 0
}

func (ppc *ppuControl) Write(addr uint16, val byte) {
	switch addr {
	case 0x00:
		ppc.ppu.lcdEnable = (val&0x80 != 0)
		ppc.ppu.windowMap = (val&0x40 != 0)
		ppc.ppu.windowEnable = (val&0x20 != 0)
		ppc.ppu.bgTileSelect = (val&0x10 != 0)
		ppc.ppu.bgMap = (val&0x08 != 0)
		ppc.ppu.spriteSize = (val&0x04 != 0)
		ppc.ppu.spriteEnable = (val&0x02 != 0)
		ppc.ppu.bgEnable = (val&0x01 != 0)
		return
	case 0x01:
		ppc.ppu.interrupt0 = (val&0x08 != 0)
		ppc.ppu.interrupt1 = (val&0x10 != 0)
		ppc.ppu.interrupt2 = (val&0x20 != 0)
		ppc.ppu.interruptLYC = (val&0x40 != 0)
		return
	case 0x02:
		ppc.ppu.scrollY = val
		return
	case 0x03:
		ppc.ppu.scrollX = val
		return
	case 0x04:
		ppc.ppu.line = val
		return
	case 0x05:
		ppc.ppu.lineCompare = val
		return
	case 0x07:
		for i := uint8(0); i < 4; i++ {
			switch (val >> (i * 2)) & 3 {
			case 0:
				ppc.ppu.palette[i] = color.RGBA{255, 255, 255, 0xFF}
			case 1:
				ppc.ppu.palette[i] = color.RGBA{192, 192, 192, 0xFF}
			case 2:
				ppc.ppu.palette[i] = color.RGBA{96, 96, 96, 0xFF}
			case 3:
				ppc.ppu.palette[i] = color.RGBA{0, 0, 0, 0xFF}
			}
		}
		return
	}

	log.Printf("Encountered write with unknown PPU control address: %#2x", addr)
}

//// Helpers ////
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
