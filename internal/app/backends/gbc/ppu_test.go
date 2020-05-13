package gbc

import (
	"image/color"
	"testing"
)

func TestPPUInit(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	black := color.RGBA{0, 0, 0, 0xFF}
	darkGray := color.RGBA{96, 96, 96, 0xFF}
	lightGray := color.RGBA{192, 192, 192, 0xFF}
	white := color.RGBA{255, 255, 255, 0xFF}

	if ppu.framebuffer[0] != black || ppu.framebuffer[23039] != black {
		t.Error("FrameBuffer should initialize to black")
	}

	if ppu.bgPalette[0] != white || ppu.bgPalette[1] != lightGray || ppu.bgPalette[2] != darkGray || ppu.bgPalette[3] != black {
		t.Error("Palette should initialize correctly")
	}

	if ppu.mode != 2 {
		t.Errorf("Expected mode to initialize to 2, got %d", ppu.mode)
	}
}

func TestPPUTileAddress(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	mmu.Write(0x9800, 0x80) // int8(-128)
	mmu.Write(0x9A00, 0x00) // int8(0)
	mmu.Write(0x9BFF, 0x7F) // int8(127)
	mmu.Write(0x9C00, 0)
	mmu.Write(0x9E00, 127)
	mmu.Write(0x9FFF, 255)

	if ppu.tileSelect {
		t.Error("Expected tileSelect to init to false")
	}

	if ppu.getTileAddress(0x9800) != 0x8800 {
		t.Errorf("Expected tile at address 0x9800 to be 0x8800, got %#4x", ppu.getTileAddress(0x9800))
	}
	if ppu.getTileAddress(0x9A00) != 0x9000 {
		t.Errorf("Expected tile at address 0x9A00 to be 0x9000, got %#4x", ppu.getTileAddress(0x9A00))
	}
	if ppu.getTileAddress(0x9BFF) != 0x97F0 {
		t.Errorf("Expected tile at address 0x9BFF to be 0x97F0, got %#4x", ppu.getTileAddress(0x9BFF))
	}

	ppu.tileSelect = true

	if ppu.getTileAddress(0x9C00) != 0x8000 {
		t.Errorf("Expected tile at address 0x9C00 to be 0x8000, got %#4x", ppu.getTileAddress(0x9C00))
	}
	if ppu.getTileAddress(0x9E00) != 0x87F0 {
		t.Errorf("Expected tile at address 0x9E00 to be 0x87F0, got %#4x", ppu.getTileAddress(0x9E00))
	}
	if ppu.getTileAddress(0x9FFF) != 0x8FF0 {
		t.Errorf("Expected tile at address 0x9FFF to be 0x8FF0, got %#4x", ppu.getTileAddress(0x9FFF))
	}
}

func TestPPURenderPixel(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	mmu.Write16(0x8000, 0x3A5C)
	mmu.Write16(0x8002, 0x5C3A)

	if ppu.getTileVal(0x8000, 0, 0) != 0 {
		t.Errorf("Expected tile value at 0x8000,0,0 to be 0, got %d", ppu.getTileVal(0x8000, 0, 0))
	}
	if ppu.getTileVal(0x8000, 1, 0) != 1 {
		t.Errorf("Expected tile value at 0x8000,1,0 to be 1, got %d", ppu.getTileVal(0x8000, 1, 0))
	}
	if ppu.getTileVal(0x8000, 2, 0) != 2 {
		t.Errorf("Expected tile value at 0x8000,2,0 to be 2, got %d", ppu.getTileVal(0x8000, 2, 0))
	}
	if ppu.getTileVal(0x8000, 3, 0) != 3 {
		t.Errorf("Expected tile value at 0x8000,3,0 to be 3, got %d", ppu.getTileVal(0x8000, 3, 0))
	}
	if ppu.getTileVal(0x8000, 4, 0) != 3 {
		t.Errorf("Expected tile value at 0x8000,4,0 to be 3, got %d", ppu.getTileVal(0x8000, 4, 0))
	}
	if ppu.getTileVal(0x8000, 5, 0) != 1 {
		t.Errorf("Expected tile value at 0x8000,5,0 to be 1, got %d", ppu.getTileVal(0x8000, 5, 0))
	}
	if ppu.getTileVal(0x8000, 6, 0) != 2 {
		t.Errorf("Expected tile value at 0x8000,6,0 to be 2, got %d", ppu.getTileVal(0x8000, 6, 0))
	}
	if ppu.getTileVal(0x8000, 7, 0) != 0 {
		t.Errorf("Expected tile value at 0x8000,7,0 to be 0, got %d", ppu.getTileVal(0x8000, 7, 0))
	}
	if ppu.getTileVal(0x8000, 0, 1) != 0 {
		t.Errorf("Expected tile value at 0x8000,0,1 to be 0, got %d", ppu.getTileVal(0x8000, 0, 1))
	}
	if ppu.getTileVal(0x8000, 1, 1) != 2 {
		t.Errorf("Expected tile value at 0x8000,1,1 to be 2, got %d", ppu.getTileVal(0x8000, 1, 1))
	}
	if ppu.getTileVal(0x8000, 2, 1) != 1 {
		t.Errorf("Expected tile value at 0x8000,2,1 to be 1, got %d", ppu.getTileVal(0x8000, 2, 1))
	}
	if ppu.getTileVal(0x8000, 3, 1) != 3 {
		t.Errorf("Expected tile value at 0x8000,3,1 to be 3, got %d", ppu.getTileVal(0x8000, 3, 1))
	}
}

func TestPPUWritePixel(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	black := color.RGBA{0, 0, 0, 0xFF}

	for _, p := range ppu.framebuffer {
		if p != black {
			t.Error("Expected framebuffer to init to black")
		}
	}

	ppu.writePixel(color.RGBA{1, 2, 3, 0xFF}, 0, 0)
	ppu.writePixel(color.RGBA{4, 5, 6, 0xFF}, 1, 0)
	ppu.writePixel(color.RGBA{7, 8, 9, 0xFF}, 0, 1)
	ppu.writePixel(color.RGBA{10, 11, 12, 0xFF}, 1, 1)
	ppu.writePixel(color.RGBA{13, 14, 15, 0xFF}, 159, 0)
	ppu.writePixel(color.RGBA{16, 17, 18, 0xFF}, 0, 143)
	ppu.writePixel(color.RGBA{19, 20, 21, 0xFF}, 159, 143)

	expected := color.RGBA{1, 2, 3, 0xFF}
	if ppu.framebuffer[0] != expected {
		t.Error("Expected coordinate (0, 0) to be {1, 2, 3, 0xFF}")
	}
	expected = color.RGBA{4, 5, 6, 0xFF}
	if ppu.framebuffer[1] != expected {
		t.Error("Expected coordinate (1, 0) to be {4, 5, 6, 0xFF}")
	}
	expected = color.RGBA{7, 8, 9, 0xFF}
	if ppu.framebuffer[160] != expected {
		t.Error("Expected coordinate (0, 1) to be {7, 8, 9, 0xFF}")
	}
	expected = color.RGBA{10, 11, 12, 0xFF}
	if ppu.framebuffer[161] != expected {
		t.Error("Expected coordinate (1, 1) to be {10, 11, 12, 0xFF}")
	}
	expected = color.RGBA{13, 14, 15, 0xFF}
	if ppu.framebuffer[159] != expected {
		t.Error("Expected coordinate (159, 0) to be {13, 14, 15, 0xFF}")
	}
	expected = color.RGBA{16, 17, 18, 0xFF}
	if ppu.framebuffer[22880] != expected {
		t.Error("Expected coordinate (0, 143) to be {16, 17, 18, 0xFF}")
	}
	expected = color.RGBA{19, 20, 21, 0xFF}
	if ppu.framebuffer[23039] != expected {
		t.Error("Expected coordinate (159, 143) to be {19, 20, 21, 0xFF}")
	}
}

func TestPPURenderLine(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	ppu.lcdEnable = true
	ppu.bgEnable = true
	ppu.windowEnable = false

	// Create 4 tiles in each tileset
	for i := byte(0); i < 8; i++ {
		mmu.Write16(0x8000+uint16(i*2), 0xFFFF) // Black
		mmu.Write16(0x9000+uint16(i*2), 0x0000) // White
	}
	for i := byte(0); i < 8; i++ {
		mmu.Write16(0x8010+uint16(i*2), 0xFF00) // Dark Gray
		mmu.Write16(0x9010+uint16(i*2), 0x00FF) // Light Gray
	}
	for i := byte(0); i < 8; i++ {
		mmu.Write16(0x8020+uint16(i*2), 0x00FF) // Light Gray
		mmu.Write16(0x9020+uint16(i*2), 0xFF00) // Dark Gray
	}
	for i := byte(0); i < 8; i++ {
		mmu.Write16(0x8030+uint16(i*2), 0x0000) // White
		mmu.Write16(0x9030+uint16(i*2), 0xFFFF) // Black
	}

	// Create two lines of background in each map
	for i := byte(0); i < 32; i++ {
		mmu.Write(0x9800+uint16(i), i%4)     // Map 0, line 0
		mmu.Write(0x9820+uint16(i), (i+1)%4) // Map 0, line 1
		mmu.Write(0x9C00+uint16(i), (i+2)%4) // Map 1, line 0
		mmu.Write(0x9C20+uint16(i), (i+3)%4) // Map 1, line 1
	}

	// Render

	// Line 0: BG 0, TS 0
	ppu.line = 0
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.renderLine()

	// Line 1: BG 0, TS 1
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = true
	ppu.renderLine()

	// line 2: BG 1, TS 0
	ppu.line++
	ppu.bgMap = true
	ppu.tileSelect = false
	ppu.renderLine()

	// Line 3: BG 1, TS 1
	ppu.line++
	ppu.bgMap = true
	ppu.tileSelect = true
	ppu.renderLine()

	// Line 4: BG 0, TS 0, ScrollX = 8
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgScrollX = 8
	ppu.renderLine()

	// Line 5: BG 0, TS 0, ScrollX = 14
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgScrollX = 14
	ppu.renderLine()

	// Line 6: BG 0, TS 0, ScrollY = 8
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgScrollX = 0
	ppu.bgScrollY = 8
	ppu.renderLine()

	// Line 7: BG 0, TS 0, ScrollX = 9, ScrollY = 1
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgScrollX = 9
	ppu.bgScrollY = 1
	ppu.renderLine()

	// Line 8: BG 0, TS 0
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgScrollX = 0
	ppu.bgScrollY = 0
	ppu.renderLine()

	// Line 9: BG 0, TS 1
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = true
	ppu.renderLine()

	// Line 10: BG 1, TS 0
	ppu.line++
	ppu.bgMap = true
	ppu.tileSelect = false
	ppu.renderLine()

	// Line 11: BG 1, TS 1
	ppu.line++
	ppu.bgMap = true
	ppu.tileSelect = true
	ppu.renderLine()

	// Line 12: BG 0, TS 0, Palette Swapped
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.bgPalette[3] = color.RGBA{255, 255, 255, 0xFF}
	ppu.bgPalette[2] = color.RGBA{192, 192, 192, 0xFF}
	ppu.bgPalette[1] = color.RGBA{96, 96, 96, 0xFF}
	ppu.bgPalette[0] = color.RGBA{0, 0, 0, 0xFF}
	ppu.renderLine()

	// Swapping palette back
	ppu.bgPalette[0] = color.RGBA{255, 255, 255, 0xFF}
	ppu.bgPalette[1] = color.RGBA{192, 192, 192, 0xFF}
	ppu.bgPalette[2] = color.RGBA{96, 96, 96, 0xFF}
	ppu.bgPalette[3] = color.RGBA{0, 0, 0, 0xFF}

	// Line 13: Window enabled but scrolled lower
	ppu.line++
	ppu.bgMap = false
	ppu.tileSelect = false
	ppu.windowEnable = true
	ppu.windowMap = true
	ppu.wScrollXm7 = 0
	ppu.wScrollY = 14
	ppu.renderLine()

	// Line 14: Window in middle
	ppu.line++
	ppu.wScrollXm7 = 30
	ppu.renderLine()

	// Line 15: Window aligned left
	ppu.line++
	ppu.wScrollXm7 = 7
	ppu.renderLine()

	// Line 16: Window overflows left
	ppu.line++
	ppu.wScrollXm7 = 0
	ppu.renderLine()

	// Line 17: Window second tile
	ppu.line++
	ppu.wScrollXm7 = 7
	ppu.wScrollY = 9
	ppu.renderLine()

	// Line 18: Other window map
	ppu.line++
	ppu.windowMap = false
	ppu.renderLine()

	// Line 19: Window with other tileset
	ppu.line++
	ppu.tileSelect = true
	ppu.renderLine()

	expected := []byte{
		// Line 0
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
		// Line 1
		3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 2
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		// Line 3
		1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
		1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
		1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
		1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
		1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
		// Line 4
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 5
		1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1,
		// Line 6
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 7
		2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2,
		2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2,
		2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2,
		2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2,
		2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2,
		// Line 8
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 9
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		// Line 10
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		// Line 11
		0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1,
		// Line 12
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		// Line 13
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 14
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 3,
		3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3,
		3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3,
		3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3,
		3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3,
		// Line 15
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		// Line 16
		2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
		2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
		2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
		2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
		2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
		// Line 17
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2,
		// Line 18
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0,
		// Line 19
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3,
	}

	for l := 0; l < 20; l++ {
		good := true
		i := 0
		for i = 0; i < 160; i++ {
			idx := (l * 160) + i
			if ppu.framebuffer[idx] != iToRGBA(expected[idx]) {
				good = false
				break
			}
		}

		if !good {
			t.Errorf("Output mismatch on framebuffer line %d. First offending index: %d", l, i)
		}
	}

	// TODO expand this test when sprites are implemented
}

func TestPPUTiming(t *testing.T) {
	mmu := NewMMU(nil)
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	ppu.RunForClocks(0)
	if ppu.mode != 2 {
		t.Errorf("Expected PPU to init in mode 2, got %d", ppu.mode)
	}

	ppu.RunForClocks(20)
	if ppu.mode != 2 {
		t.Errorf("Expected PPU to remain in mode 2, got %d", ppu.mode)
	}
	if ppu.timeInMode != 20 {
		t.Errorf("Expected PPU to track time in mode correctly as 20, got %d", ppu.timeInMode)
	}

	ppu.RunForClocks(1)
	if ppu.mode != 3 {
		t.Errorf("Expected PPU to advance to mode 3, got %d", ppu.mode)
	}
	if ppu.timeInMode != 0 {
		t.Errorf("Expected PPU to reset time in mode, got %d", ppu.timeInMode)
	}

	ppu.RunForClocks(42)
	if ppu.mode != 3 {
		t.Errorf("Expected PPU to remain in mode 3, got %d", ppu.mode)
	}
	if ppu.timeInMode != 42 {
		t.Errorf("Expected PPU to track time in mode correctly as 43, got %d", ppu.timeInMode)
	}

	ppu.RunForClocks(1)
	if ppu.mode != 0 {
		t.Errorf("Expected PPU to advance to mode 0, got %d", ppu.mode)
	}
	if ppu.timeInMode != 0 {
		t.Errorf("Expected PPU to reset time in mode, got %d", ppu.timeInMode)
	}

	ppu.RunForClocks(49)
	if ppu.mode != 0 {
		t.Errorf("Expected PPU to remain in mode 0, got %d", ppu.mode)
	}
	if ppu.timeInMode != 49 {
		t.Errorf("Expected PPU to track time in mode correctly as 49, got %d", ppu.timeInMode)
	}

	if ppu.line != 0 {
		t.Errorf("Expected PPU to be rendering line 0, got %d", ppu.line)
	}

	ppu.RunForClocks(1)
	if ppu.mode != 2 {
		t.Errorf("Expected PPU to advance to mode 2, got %d", ppu.mode)
	}
	if ppu.timeInMode != 0 {
		t.Errorf("Expected PPU to reset time in mode, got %d", ppu.timeInMode)
	}

	if ppu.line != 1 {
		t.Errorf("Expected PPU to advance line to 1, got %d", ppu.line)
	}

	ppu.RunForClocks(16187)
	if ppu.mode != 0 {
		t.Errorf("Expected PPU to be in mode 0, got %d", ppu.mode)
	}
	if ppu.timeInMode != 49 {
		t.Errorf("Expected PPU to track time in mode correctly as 49, got %d", ppu.timeInMode)
	}
	if ppu.line != 142 {
		t.Errorf("Expected PPU line to be 142, got %d", ppu.line)
	}

	ppu.RunForClocks(2)
	if ppu.mode != 1 {
		t.Errorf("Expected PPU to be in mode 1, got %d", ppu.mode)
	}
	if ppu.timeInMode != 1 {
		t.Errorf("Expected PPU to reset time in mode, got %d", ppu.timeInMode)
	}
	if ppu.line != 144 {
		t.Errorf("Expected PPU line to be 144, got %d", ppu.line)
	}

	ppu.RunForClocks(1138)
	if ppu.mode != 1 {
		t.Errorf("Expected PPU to be in mode 1, got %d", ppu.mode)
	}
	if ppu.timeInMode != 1139 {
		t.Errorf("Expected PPU to track time in mode correctly as 1139, got %d", ppu.timeInMode)
	}

	ppu.RunForClocks(1)
	if ppu.mode != 2 {
		t.Errorf("Expected PPU to advance to mode 2, got %d", ppu.mode)
	}
	if ppu.timeInMode != 0 {
		t.Errorf("Expected PPU to reset time in mode, got %d", ppu.timeInMode)
	}
	if ppu.line != 0 {
		t.Errorf("Expected PPU line to be reset, got %d", ppu.line)
	}
}

func TestPPUUpdateToClock(t *testing.T) {
	// General test of PPU
}

func TestPPUControlDevice(t *testing.T) {
	// TODO test control memory device read/write
}

func TestPPUInterrupts(t *testing.T) {
	// TODO test get/set of interrupts, and correct firing
}
