package gbc

import (
	"image/color"
	"testing"
)

func TestPPUInit(t *testing.T) {
	mmu := NewMMU()
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	black := color.RGBA{0, 0, 0, 0xFF}
	darkGray := color.RGBA{96, 96, 96, 0xFF}
	lightGray := color.RGBA{192, 192, 192, 0xFF}
	white := color.RGBA{255, 255, 255, 0xFF}

	if ppu.framebuffer[0] != black || ppu.framebuffer[23039] != black {
		t.Error("FrameBuffer should initialize to black")
	}

	if ppu.palette[0] != white || ppu.palette[1] != lightGray || ppu.palette[2] != darkGray || ppu.palette[3] != black {
		t.Error("Palette should initialize correctly")
	}

	if !ppu.lcdEnable {
		t.Error("lcdEnable should initialize to true")
	}
	if !ppu.windowMap {
		t.Error("windowMap should initialize to true")
	}
	if !ppu.bgEnable {
		t.Error("bgEnable should initialize to true")
	}
	if ppu.mode != 2 {
		t.Errorf("Expected mode to initialize to 2, got %d", ppu.mode)
	}
}

func TestPPUTileAddress(t *testing.T) {
	mmu := NewMMU()
	ppu := NewPPU(mmu)
	mmu.ppu = ppu

	mmu.Write(0x9800, 0x80) // int8(-128)
	mmu.Write(0x9A00, 0x00) // int8(0)
	mmu.Write(0x9BFF, 0x7F) // int8(127)
	mmu.Write(0x9C00, 0)
	mmu.Write(0x9E00, 127)
	mmu.Write(0x9FFF, 255)

	if ppu.bgTileSelect {
		t.Error("Expected bgTileSelect to init to false")
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

	ppu.bgTileSelect = true

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

func TestPPUWritePixel(t *testing.T) {
	mmu := NewMMU()
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
	// TODO invoke renderLine directly
}

func TestPPUTiming(t *testing.T) {
	// TODO Test UpdateToClock timing
}

func TestPPUUpdateToClock(t *testing.T) {
	// General test of PPU
}

func TestPPUControlDevice(t *testing.T) {
	// TODO test control memory device read/write
}
