package banking

import (
	"testing"
)

var TESTDATA = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

func bigTestData(size int, increment int) []byte {
	buf := make([]byte, size)

	for i := 0; i < size; i++ {
		buf[i] = byte((i / increment) % 0x100)
	}

	return buf
}

func TestROM(t *testing.T) {
	c1 := NewROM(nil)
	c2 := NewROM(TESTDATA)

	if c1.Read(0) != 0x00 {
		t.Errorf("Expected NewROM to initialize to zero when data is nil, got %#02X", c1.Read(0))
	}

	if c2.Read(3) != 0x04 {
		t.Errorf("Expected NewROM to accept data argument, got %#02X", c2.Read(3))
	}

	if c2.Read(8) != 0x00 {
		t.Errorf("Expected NewROM to leave remaining data as zero, got %#02X", c2.Read(8))
	}

	c2.Write(0, 0xFF)

	if c2.Read(0) != 0x01 {
		t.Errorf("Expected ROM.Write to noop.")
	}
}

func TestROMRAM(t *testing.T) {
	c1 := NewROMRAM(nil)
	c2 := NewROMRAM(TESTDATA)

	if c1.Read(0) != 0x00 {
		t.Errorf("Expected NewROMRAM to initialize to zero when data is nil, got %#02X", c1.Read(0))
	}

	if c2.Read(3) != 0x04 {
		t.Errorf("Expected NewROMRAM to accept data argument, got %#02X", c2.Read(3))
	}

	if c2.Read(8) != 0x00 {
		t.Errorf("Expected NewROMRAM to leave remaining data as zero, got %#02X", c2.Read(8))
	}

	c2.Write(0, 0xFF)
	c2.Write(0xA000, 0xAA)
	c2.Write(0xB000, 0xBB)
	c2.Write(0xBFFF, 0xBF)

	if c2.Read(0) != 0x01 {
		t.Errorf("Expected ROM.Write to noop when writing to ROM.")
	}

	if c2.Read(0xA000) != 0xAA {
		t.Errorf("Expected ROM.Write to succeed when writing to RAM.")
	}

	if c2.Read(0xB000) != 0xBB {
		t.Errorf("Expected ROM.Write to succeed when writing to RAM.")
	}

	if c2.Read(0xBFFF) != 0xBF {
		t.Errorf("Expected ROM.Write to succeed when writing to RAM.")
	}
}

func TestMBC1ROM(t *testing.T) {
	c := NewMBC1(bigTestData(0x1000000, 0x1000), 0x1000000, 0)

	if c.Read(0x1000) != 0x01 || c.Read(0x3000) != 0x03 {
		t.Errorf("Expected NewMBC1 to accept data argument.")
	}

	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC1.Read to find 0x04 at rom bank 1, got %#02X", c.Read(0x4000))
	}

	if c.Read(0x7FFF) != 0x07 {
		t.Errorf("Expected MBC1.Read to find 0x07 at rom bank 1, got %#02X", c.Read(0x7FFF))
	}

	// Select bank 4
	c.Write(0x2000, 4)
	if c.Read(0x4000) != 0x10 {
		t.Errorf("Expected MBC1 to switch to bank 4 and read 0x10, got %#02X", c.Read(0x4000))
	}

	// Select bank 0, skip to bank 1
	c.Write(0x2000, 0)
	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC1 to skip over bank 0 and read 0x04. got %#02X", c.Read(0x4000))
	}

	// Select bank 20, skip to bank 21
	c.Write(0x4000, 1)
	c.Write(0x3FFF, 0)
	if c.Read(0x4000) != 0x84 {
		t.Errorf("Expected MBC1 to skip over bank 20 and read 0x84. got %#02X", c.Read(0x4000))
	}

	// Select bank 40, skip to bank 41
	c.Write(0x4000, 2)
	c.Write(0x3FFF, 0)
	if c.Read(0x4000) != 0x04 { // Overflow from 0x104
		t.Errorf("Expected MBC1 to skip over bank 40 and read 0x04. got %#02X", c.Read(0x4000))
	}

	// Select bank 60, skip to bank 61
	c.Write(0x4000, 3)
	c.Write(0x3FFF, 0)
	if c.Read(0x4000) != 0x84 { // Overflow from 0x184
		t.Errorf("Expected MBC1 to skip over bank 60 and read 0x84. got %#02X", c.Read(0x4000))
	}

	// Select bank 62
	c.Write(0x2000, 2)
	if c.Read(0x5555) != 0x89 { // Overflow from 0x189
		t.Errorf("Expected MBC1 to switch to bank 62 and read 0x80. got %#02X", c.Read(0x5555))
	}
}

func TestMBC1RAM(t *testing.T) {
	c := NewMBC1(nil, 0x10000, 0x10000)

	// Write to disabled RAM
	c.Write(0xA000, 0xAA)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected MBC1 to not write to RAM when disabled, got %#02X", c.Read(0xA000))
	}

	// Write to enabled RAM
	c.Write(0xA000, 0xBB)
	if c.Read(0xA000) != 0xBB {
		t.Errorf("Expected MBC1 to write/read from enabled ram successfuly, got %#02X", c.Read(0xA000))
	}

	// Switch to RAM banking mode
	c.Write(0x6000, 0x01)

	// Switch to bank 3
	c.Write(0x4000, 0x03)

	c.Write(0xA000, 0xCC)
	c.Write(0xB000, 0xDD)

	if c.Read(0xA000) != 0xCC {
		t.Errorf("Expected MBC1 to switch to bank 3 and write/read 0xCC, got %#02X", c.Read(0xA000))
	}

	if c.Read(0xB000) != 0xDD {
		t.Errorf("Expected MBC1 to switch to bank 3 and write/read 0xDD, got %#02X", c.Read(0xB000))
	}

	// Switch back to ROM banking mode
	c.Write(0x6000, 0x00)

	if c.Read(0xA000) != 0xBB {
		t.Errorf("Expected MBC1 to switch back to bank 0 when switching to ROM banking mode, got %#02X", c.Read(0xA000))
	}

	// Disable RAM
	c.Write(0x0000, 0x0)

	if c.Read(0xA000) != 0xFF {
		t.Errorf("Expected MBC1 to read 0xFF from disabled RAM, got %#02X", c.Read(0xA000))
	}

	c.Write(0xA000, 0xEE)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0xBB {
		t.Errorf("Expected MBC1 to read previously written value of 0xBB after re-enabling RAM, got %#02X", c.Read(0xA000))
	}
}

func TestMBC1Overflow(t *testing.T) {
	c := NewMBC1(bigTestData(0x100000, 0x1000), 0x10000, 0x800)

	// Select bank 3
	c.Write(0x2000, 3)

	if c.Read(0x7FFF) != 0x0F {
		t.Errorf("Expected MBC1 to switch to bank 3 and read 0x0F, got %#02X", c.Read(0x7FFF))
	}

	// Select bank 4
	c.Write(0x2000, 4)

	if c.Read(0x4000) != 0xFF {
		t.Errorf("Expected MBC1 to switch to bank 4 and read overflow value of 0xFF, got %#02X", c.Read(0x4000))
	}

	// Read RAM when disabled
	if c.Read(0xA000) != 0xFF {
		t.Errorf("Expected MBC1 to return 0xFF when RAM is disabled, got %#02X", c.Read(0xA000))
	}

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected MBC1 to enable RAM and return initial value of 0x00")
	}

	if c.Read(0xB000) != 0xFF {
		t.Errorf("Expected MBC1 to return overflow value of 0xFF, got %#02X", c.Read(0xB000))
	}

	// Switch to RAM banking mode and select RAM bank 2
	c.Write(0x6000, 0x01)
	c.Write(0x4000, 0x02)

	c.Write(0xA000, 0xAA)

	if c.Read(0xA000) != 0xFF {
		t.Errorf("Expected MBC1 to select RAM bank 2 and return overflow value of 0xFF, got %#02X", c.Read(0xA000))
	}
}
