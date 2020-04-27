package gbc

import (
	"image/color"
	"math"
)

// PPU represents the gameboy's graphics processing unit.
type PPU struct {
	mmu *MMU // Memory Management Unit

	memoryControl *ppuControl // Control Regsiter Device
	oam           *oam        // Object Attribute Memory (sprites)

	framebuffer    []color.RGBA  // Frame Buffer
	bgPalette      [4]color.RGBA // Background Color Palette
	spritePalette0 [4]color.RGBA // Sprite Color Palette 0
	spritePalette1 [4]color.RGBA // Sprite Color Palette 1

	// Control Registers
	lcdEnable    bool // Enables the entire screen
	windowMap    bool // Which window map is in use
	windowEnable bool // Enables the window display
	tileSelect   bool // Which tileset is in use
	bgMap        bool // Which background map is in use
	spriteSize   bool // Size of all sprites
	spriteEnable bool // Enables rendering sprites
	bgEnable     bool // Enables rendering the background

	mode        byte // Mode Number (0: HBLANK, 1: VBLANK, 2: OAM, 3: VRAM)
	timeInMode  int  // Number of clock cycles spent in the current mode
	line        byte // Line currently being processed
	lineCompare byte // Target line for interrupt

	interrupt0   bool // Trigger an interrupt on entering mode 0
	interrupt1   bool // Trigger an interrupt on entering mode 1
	interrupt2   bool // Trigger an interrupt on entering mode 2
	interruptLYC bool // Trigger an interrupt when line matches lineCompare

	bgScrollX  byte // Background scroll X
	bgScrollY  byte // Background scroll Y
	wScrollXm7 byte // Window scroll X, minus 7
	wScrollY   byte // Window scroll Y
}

// NewPPU constructs a valid PPU struct
func NewPPU(mmu *MMU) *PPU {
	ppu := new(PPU)

	ppu.mmu = mmu

	ppu.framebuffer = make([]color.RGBA, ScreenHeight*ScreenWidth)
	ppu.clearScrean()

	ppu.bgPalette = [4]color.RGBA{
		{255, 255, 255, 0xFF},
		{192, 192, 192, 0xFF},
		{96, 96, 96, 0xFF},
		{0, 0, 0, 0xFF},
	}

	ppu.mode = 2 // Start in OAM mode

	ppuControl := new(ppuControl)
	ppuControl.ppu = ppu
	ppu.memoryControl = ppuControl

	ppu.oam = new(oam)

	return ppu
}

// RunForClocks runs the PPU for the given number of clock cycles
func (ppu *PPU) RunForClocks(clocks int) {
	for c := 0; c < clocks; c++ {
		ppu.timeInMode++

		switch ppu.mode {
		case 0: // HBLANK
			if ppu.timeInMode == 50 {
				ppu.timeInMode = 0

				ppu.line++

				if ppu.interruptLYC && ppu.line == ppu.lineCompare {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}

				if ppu.line == 143 {
					ppu.mode = 1
					ppu.mmu.interrupts.RequestInterrupt(interruptVBlankBit)
					if ppu.interrupt1 {
						ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
					}
				} else {
					ppu.mode = 2
					if ppu.interrupt2 {
						ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
					}
				}
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
				if ppu.interruptLYC && ppu.lineCompare == 143 {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}
				ppu.line = byte(144 + math.Round(9*(float64(ppu.timeInMode)/1140)))
			}
		case 2: // OAM
			if ppu.timeInMode == 21 {
				ppu.timeInMode = 0

				ppu.mode = 3
			}
		case 3: // VRAM
			if ppu.timeInMode == 43 {
				ppu.timeInMode = 0

				ppu.mode = 0

				if ppu.interrupt0 {
					ppu.mmu.interrupts.RequestInterrupt(interruptLCDBit)
				}

				ppu.renderLine()
			}
		}
	}
}

func (ppu *PPU) renderLine() {
	if !ppu.lcdEnable {
		ppu.clearScrean()
		return
	}

	// BG color values on the current line, for sprite transparency
	var lineColors [ScreenWidth]uint8

	// Draw the background if enabled
	if ppu.bgEnable {
		// Base VRAM address for the background map
		var bgAddr uint16
		if ppu.bgMap {
			bgAddr = 0x9C00
		} else {
			bgAddr = 0x9800
		}

		// First tile to be drawn
		tileShiftY := uint16((ppu.line+ppu.bgScrollY)>>3) << 5
		tileShiftX := uint16(ppu.bgScrollX >> 3)
		tilePointer := bgAddr + tileShiftY + tileShiftX
		tileAddr := ppu.getTileAddress(tilePointer)

		// Coordinate in the tile to start drawing
		tileX := ppu.bgScrollX & 0x7
		tileY := (ppu.line + ppu.bgScrollY) & 0x7

		// Coordinates in the framebuffer to draw
		screenX := 0
		screenY := int(ppu.line)

		for screenX < ScreenWidth {
			// Move to next tile if needed
			if tileX == 8 {
				tileX = 0
				tilePointer++
				tileAddr = ppu.getTileAddress(tilePointer)
			}

			val := ppu.getTileVal(tileAddr, tileX, tileY)
			pixel := ppu.bgPalette[val]

			lineColors[screenX] = val
			ppu.writePixel(pixel, screenX, screenY)

			screenX++
			tileX++
		}
	}

	// Draw sprites if enabled
	if ppu.spriteEnable {
		visibleSprites := ppu.oam.VisibleSpritesOnLine(ppu.line, ppu.spriteSize)
		for _, sprite := range visibleSprites {
			tileAddr := 0x8000 + (uint16(sprite.tileNum) << 4)

			palette := ppu.spritePalette0
			if sprite.paletteFlag {
				palette = ppu.spritePalette1
			}

			tileY := ppu.line - (sprite.yPos - 16)
			if sprite.yFlip {
				if ppu.spriteSize {
					tileY = 15 - tileY
				} else {
					tileY = 7 - tileY
				}
			}

			screenX := int(sprite.xPos) - 8
			screenY := int(ppu.line)
			for tileX := byte(0); tileX < 8; tileX++ {
				if screenX >= 0 && screenX < 160 {

					if sprite.xFlip {
						tileX = 7 - tileX
					}

					val := ppu.getTileVal(tileAddr, tileX, tileY)

					if val != 0 && (sprite.priority || lineColors[screenX] == 0) {
						pixel := palette[val]
						ppu.writePixel(pixel, screenX, screenY)
					}

				}
				screenX++
			}

		}
	}

	// Draw the window if enabled
	if ppu.windowEnable && ppu.line >= ppu.wScrollY {

		// Base VRAM address for the window map
		var wAddr uint16
		if ppu.windowMap {
			wAddr = 0x9C00
		} else {
			wAddr = 0x9800
		}

		// First tile to be drawn
		tileShiftY := uint16((ppu.line-ppu.wScrollY)>>3) << 5
		tileShiftX := uint16(0) // Window always starts from the left
		tilePointer := wAddr + tileShiftY + tileShiftX
		tileAddr := ppu.getTileAddress(tilePointer)

		// Coordinates in the tile to start drawing
		var tileX, tileY byte
		// Coordinates in the framebuffer to draw
		var screenX, screenY int
		if ppu.wScrollXm7 < 7 {
			tileX = 7 - ppu.wScrollXm7
			screenX = 0
		} else {
			tileX = 0
			screenX = int(ppu.wScrollXm7) - 7
		}
		tileY = (ppu.line - ppu.wScrollY) & 0x7
		screenY = int(ppu.line)

		for screenX < ScreenWidth {
			// Move to next tile if needed
			if tileX == 8 {
				tileX = 0
				tilePointer++
				tileAddr = ppu.getTileAddress(tilePointer)
			}

			val := ppu.getTileVal(tileAddr, tileX, tileY)
			pixel := ppu.bgPalette[val]

			ppu.writePixel(pixel, screenX, screenY)

			screenX++
			tileX++
		}
	}
}

// writePixel writes the given RGBA value into the framebuffer at coordinates (x, y)
func (ppu *PPU) writePixel(val color.RGBA, x int, y int) {
	ppu.framebuffer[(y*ScreenWidth)+x] = val
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
	if ppu.tileSelect {
		tileNum := uint16(ppu.mmu.Read(mapAddr))
		return uint16(0x8000) + (tileNum << 4)
	}

	tileNum := int16(int8(ppu.mmu.Read(mapAddr)))
	return uint16(int32(0x9000) + int32(tileNum<<4))
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
	case 0xFF40:
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
		if ppc.ppu.tileSelect {
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
	case 0xFF41:
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
	case 0xFF42:
		return ppc.ppu.bgScrollY
	case 0xFF43:
		return ppc.ppu.bgScrollX
	case 0xFF44:
		return ppc.ppu.line
	case 0xFF45:
		return ppc.ppu.lineCompare
	case 0xFF47:
		var ret byte
		ret |= (rgbaToI(ppc.ppu.bgPalette[0])) << 6
		ret |= (rgbaToI(ppc.ppu.bgPalette[1])) << 4
		ret |= (rgbaToI(ppc.ppu.bgPalette[2])) << 2
		ret |= (rgbaToI(ppc.ppu.bgPalette[3]))
		return ret
	case 0xFF48:
		var ret byte
		ret |= (rgbaToI(ppc.ppu.spritePalette0[0])) << 6
		ret |= (rgbaToI(ppc.ppu.spritePalette0[1])) << 4
		ret |= (rgbaToI(ppc.ppu.spritePalette0[2])) << 2
		ret |= (rgbaToI(ppc.ppu.spritePalette0[3]))
	case 0xFF49:
		var ret byte
		ret |= (rgbaToI(ppc.ppu.spritePalette1[0])) << 6
		ret |= (rgbaToI(ppc.ppu.spritePalette1[1])) << 4
		ret |= (rgbaToI(ppc.ppu.spritePalette1[2])) << 2
		ret |= (rgbaToI(ppc.ppu.spritePalette1[3]))
	case 0xFF4A:
		return ppc.ppu.wScrollY
	case 0xFF4B:
		return ppc.ppu.wScrollXm7
	}

	logger.Warningf("Encountered read with unknown PPU control address: %#04x", addr)
	return 0
}

func (ppc *ppuControl) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF40:
		ppc.ppu.lcdEnable = (val&0x80 != 0)
		ppc.ppu.windowMap = (val&0x40 != 0)
		ppc.ppu.windowEnable = (val&0x20 != 0)
		ppc.ppu.tileSelect = (val&0x10 != 0)
		ppc.ppu.bgMap = (val&0x08 != 0)
		ppc.ppu.spriteSize = (val&0x04 != 0)
		ppc.ppu.spriteEnable = (val&0x02 != 0)
		ppc.ppu.bgEnable = (val&0x01 != 0)
		return
	case 0xFF41:
		ppc.ppu.interrupt0 = (val&0x08 != 0)
		ppc.ppu.interrupt1 = (val&0x10 != 0)
		ppc.ppu.interrupt2 = (val&0x20 != 0)
		ppc.ppu.interruptLYC = (val&0x40 != 0)
		return
	case 0xFF42:
		ppc.ppu.bgScrollY = val
		return
	case 0xFF43:
		ppc.ppu.bgScrollX = val
		return
	case 0xFF44:
		ppc.ppu.line = val
		return
	case 0xFF45:
		ppc.ppu.lineCompare = val
		return
	case 0xFF47:
		for i := uint8(0); i < 4; i++ {
			ppc.ppu.bgPalette[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF48:
		for i := uint8(0); i < 4; i++ {
			ppc.ppu.spritePalette0[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF49:
		for i := uint8(0); i < 4; i++ {
			ppc.ppu.spritePalette1[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF4A:
		ppc.ppu.wScrollY = val
		return
	case 0xFF4B:
		ppc.ppu.wScrollXm7 = val
		return
	}

	logger.Warningf("Encountered write with unknown PPU control address: %#4x", addr)
}

//// Helpers ////

// Returns the min of two ints
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Maps a palette byte value (0-3) to RGBA
func iToRGBA(val byte) color.RGBA {
	switch val {
	case 0:
		return color.RGBA{255, 255, 255, 0xFF}
	case 1:
		return color.RGBA{192, 192, 192, 0xFF}
	case 2:
		return color.RGBA{96, 96, 96, 0xFF}
	case 3:
		return color.RGBA{0, 0, 0, 0xFF}
	}

	return color.RGBA{255, 255, 255, 0xFF}
}

// Maps an RGBA color to palette byte value (0-3)
func rgbaToI(c color.RGBA) byte {
	switch c {
	case color.RGBA{255, 255, 255, 0xFF}:
		return 0
	case color.RGBA{192, 192, 192, 0xFF}:
		return 1
	case color.RGBA{96, 96, 96, 0xFF}:
		return 2
	case color.RGBA{0, 0, 0, 0xFF}:
		return 3
	}

	return 0
}
