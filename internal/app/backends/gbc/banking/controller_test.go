package banking

import (
	"testing"
)

func TestROM(t *testing.T) {
	c1 := NewROM(nil)
	c2 := NewROM([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

	if c1.Read(0) != 0x00 {
		t.Errorf("Expected NewROM to initialize to zero when data is nil, got %#02x", c1.Read(0))
	}

	if c2.Read(3) != 0x04 {
		t.Errorf("Expected NewROM to accept data argument, got %#02x", c2.Read(3))
	}

	if c2.Read(8) != 0x00 {
		t.Errorf("Expected NewROM to leave remaining data as zero, got %#02x", c2.Read(8))
	}

	c2.Write(0, 0xFF)

	if c2.Read(0) != 0x01 {
		t.Errorf("Expected ROM.Write to noop.")
	}
}

func TestROMRAM(t *testing.T) {
	c1 := NewROMRAM(nil)
	c2 := NewROMRAM([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

	if c1.Read(0) != 0x00 {
		t.Errorf("Expected NewROMRAM to initialize to zero when data is nil, got %#02x", c1.Read(0))
	}

	if c2.Read(3) != 0x04 {
		t.Errorf("Expected NewROMRAM to accept data argument, got %#02x", c2.Read(3))
	}

	if c2.Read(8) != 0x00 {
		t.Errorf("Expected NewROMRAM to leave remaining data as zero, got %#02x", c2.Read(8))
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
