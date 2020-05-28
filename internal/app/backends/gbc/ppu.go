package gbc

import (
	"image/color"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/interrupts"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// PPU represents the gameboy's graphics processing unit.
type PPU struct {
	mmu *MMU // Memory Management Unit

	oam *oam // Object Attribute Memory (sprites)

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
					ppu.mmu.interrupts.Request(interrupts.LCDBit)
				}

				if ppu.line == 144 {
					ppu.mode = 1
					ppu.mmu.interrupts.Request(interrupts.VBlankBit)
					if ppu.interrupt1 {
						ppu.mmu.interrupts.Request(interrupts.LCDBit)
					}
				} else {
					ppu.mode = 2
					if ppu.interrupt2 {
						ppu.mmu.interrupts.Request(interrupts.LCDBit)
					}
				}
			}
		case 1: // VBLANK
			if ppu.timeInMode == 1140 {
				ppu.timeInMode = 0

				ppu.mode = 2
				ppu.line = 0

				if ppu.interrupt2 {
					ppu.mmu.interrupts.Request(interrupts.LCDBit)
				}
			} else {
				if ppu.interruptLYC && ppu.lineCompare == 143 {
					ppu.mmu.interrupts.Request(interrupts.LCDBit)
				}
				ppu.line = byte(144 + ppu.timeInMode/144)
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
					ppu.mmu.interrupts.Request(interrupts.LCDBit)
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

	// BG/window color values on the current line, for sprite transparency
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
		mapY := uint16(ppu.line+ppu.bgScrollY) >> 3
		mapX := uint16(ppu.bgScrollX >> 3)
		tileAddr := ppu.getTileAddress(bgAddr, mapX, mapY)

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
				mapX = (mapX + 1) % 32
				tileAddr = ppu.getTileAddress(bgAddr, mapX, mapY)
			}

			val := ppu.getTileVal(tileAddr, tileX, tileY)
			pixel := ppu.bgPalette[val]

			lineColors[screenX] = val
			ppu.writePixel(pixel, screenX, screenY)

			screenX++
			tileX++
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
		mapY := uint16((ppu.line - ppu.wScrollY) >> 3)
		mapX := uint16(0) // Window always starts from the left
		tileAddr := ppu.getTileAddress(wAddr, mapX, mapY)

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
				mapX = (mapX + 1) % 32
				tileAddr = ppu.getTileAddress(wAddr, mapX, mapY)
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
			for x := byte(0); x < 8; x++ {
				if screenX >= 0 && screenX < 160 {

					tileX := x
					if sprite.xFlip {
						tileX = 7 - tileX
					}

					val := ppu.getTileVal(tileAddr, tileX, tileY)

					if val != 0 && (!sprite.priority || (lineColors[screenX] == 0)) {
						pixel := palette[val]
						ppu.writePixel(pixel, screenX, screenY)
					}

				}
				screenX++
			}

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

// getTileAddress returns the mmu address of the tile at the given 32x32 BG coordinates, in the given map.
func (ppu *PPU) getTileAddress(baseAddr uint16, x uint16, y uint16) uint16 {
	if x >= 32 || y >= 32 {
		log.Errorf("PPU Tile Address out of bounds! (%d, %d)", x, y)
	}

	mapAddr := baseAddr + (y << 5) + x

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

func (ppu *PPU) Read(addr uint16) byte {
	switch addr {
	case 0xFF40:
		var ret byte
		if ppu.lcdEnable {
			ret |= 0x80
		}
		if ppu.windowMap {
			ret |= 0x40
		}
		if ppu.windowEnable {
			ret |= 0x20
		}
		if ppu.tileSelect {
			ret |= 0x10
		}
		if ppu.bgMap {
			ret |= 0x08
		}
		if ppu.spriteSize {
			ret |= 0x04
		}
		if ppu.spriteEnable {
			ret |= 0x02
		}
		if ppu.bgEnable {
			ret |= 0x01
		}
		return ret
	case 0xFF41:
		var ret byte
		ret = ppu.mode
		if ppu.line == ppu.lineCompare {
			ret |= 0x04
		}
		if ppu.interrupt0 {
			ret |= 0x08
		}
		if ppu.interrupt1 {
			ret |= 0x10
		}
		if ppu.interrupt2 {
			ret |= 0x20
		}
		if ppu.interruptLYC {
			ret |= 0x40
		}
		return ret
	case 0xFF42:
		return ppu.bgScrollY
	case 0xFF43:
		return ppu.bgScrollX
	case 0xFF44:
		return ppu.line
	case 0xFF45:
		return ppu.lineCompare
	case 0xFF47:
		var ret byte
		ret |= (rgbaToI(ppu.bgPalette[0])) << 6
		ret |= (rgbaToI(ppu.bgPalette[1])) << 4
		ret |= (rgbaToI(ppu.bgPalette[2])) << 2
		ret |= (rgbaToI(ppu.bgPalette[3]))
		return ret
	case 0xFF48:
		var ret byte
		ret |= (rgbaToI(ppu.spritePalette0[0])) << 6
		ret |= (rgbaToI(ppu.spritePalette0[1])) << 4
		ret |= (rgbaToI(ppu.spritePalette0[2])) << 2
		ret |= (rgbaToI(ppu.spritePalette0[3]))
		return ret
	case 0xFF49:
		var ret byte
		ret |= (rgbaToI(ppu.spritePalette1[0])) << 6
		ret |= (rgbaToI(ppu.spritePalette1[1])) << 4
		ret |= (rgbaToI(ppu.spritePalette1[2])) << 2
		ret |= (rgbaToI(ppu.spritePalette1[3]))
		return ret
	case 0xFF4A:
		return ppu.wScrollY
	case 0xFF4B:
		return ppu.wScrollXm7
	}

	log.Warningf("Encountered read with unknown PPU control address: %#04x", addr)
	return 0xFF
}

func (ppu *PPU) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF40:
		ppu.lcdEnable = (val&0x80 != 0)
		ppu.windowMap = (val&0x40 != 0)
		ppu.windowEnable = (val&0x20 != 0)
		ppu.tileSelect = (val&0x10 != 0)
		ppu.bgMap = (val&0x08 != 0)
		ppu.spriteSize = (val&0x04 != 0)
		ppu.spriteEnable = (val&0x02 != 0)
		ppu.bgEnable = (val&0x01 != 0)
		return
	case 0xFF41:
		ppu.interrupt0 = (val&0x08 != 0)
		ppu.interrupt1 = (val&0x10 != 0)
		ppu.interrupt2 = (val&0x20 != 0)
		ppu.interruptLYC = (val&0x40 != 0)
		return
	case 0xFF42:
		ppu.bgScrollY = val
		return
	case 0xFF43:
		ppu.bgScrollX = val
		return
	case 0xFF44:
		ppu.line = val
		return
	case 0xFF45:
		ppu.lineCompare = val
		return
	case 0xFF47:
		for i := uint8(0); i < 4; i++ {
			ppu.bgPalette[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF48:
		for i := uint8(0); i < 4; i++ {
			ppu.spritePalette0[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF49:
		for i := uint8(0); i < 4; i++ {
			ppu.spritePalette1[i] = iToRGBA((val >> (i * 2)) & 3)
		}
		return
	case 0xFF4A:
		ppu.wScrollY = val
		return
	case 0xFF4B:
		ppu.wScrollXm7 = val
		return
	}

	log.Warningf("Encountered write with unknown PPU control address: %#4x", addr)
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
