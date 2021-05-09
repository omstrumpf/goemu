package banking

import (
	"testing"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/constants"
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

func TestMBC2ROM(t *testing.T) {
	c := NewMBC2(bigTestData(0x1000000, 0x1000))

	if c.Read(0x1000) != 0x01 || c.Read(0x3000) != 0x03 {
		t.Errorf("Expected NewMBC2 to accept data argument.")
	}

	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC2.Read to find 0x04 at rom bank 1, got %#02X", c.Read(0x4000))
	}

	if c.Read(0x7FFF) != 0x07 {
		t.Errorf("Expected MBC2.Read to find 0x07 at rom bank 1, got %#02X", c.Read(0x7FFF))
	}

	// Select bank 4
	c.Write(0x2100, 4)
	if c.Read(0x4000) != 0x10 {
		t.Errorf("Expected MBC2 to switch to bank 4 and read 0x10, got %#02X", c.Read(0x4000))
	}

	// Select bank 0, skip to bank 1
	c.Write(0x2100, 0)
	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC2 to skip over bank 0 and read 0x04. got %#02X", c.Read(0x4000))
	}
}

func TestMBC2RAM(t *testing.T) {
	c := NewMBC2(nil)

	// Write to disabled RAM
	c.Write(0xA000, 0xAA)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected MBC2 to not write to RAM when disabled, got %#02X", c.Read(0xA000))
	}

	// 4 bit (lo) read/write from RAM
	c.Write(0xA000, 0xCB)
	if c.Read(0xA000) != 0x0B {
		t.Errorf("Expected MBC2 to write/read 0x0B from enabled RAM, got %#02x", c.Read(0xA000))
	}

	// 4 bit (hi) read/write from RAM
	c.Write(0xA001, 0xBC)
	if c.Read(0xA001) != 0x0C {
		t.Errorf("Expected MBC2 to write/read 0x0C from enabled RAM, got %#02X", c.Read(0xA001))
	}

	// RAM access repeats throuch to 0xC000
	if c.Read(0xA200) != 0x0B {
		t.Errorf("Expected MBC2 RAM to repeat every 0x0200 and read 0x0B, got %#02X", c.Read(0xA200))
	}
	c.Write(0xB405, 0x0E)
	if c.Read(0xA605) != 0x0E {
		t.Errorf("Expected MBC2 RAM to repeat every 0x0200 and write/read 0x0E, got %#02X", c.Read(0xA605))
	}

	// Disable RAM
	c.Write(0x0000, 0x0)

	if c.Read(0xA000) != 0xFF {
		t.Errorf("Expected MBC2 to read 0xFF from disabled RAM, got %#02X", c.Read(0xA000))
	}

	c.Write(0xA000, 0xEE)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0x0B {
		t.Errorf("Expected MBC2 to read previously written value of 0x0B after re-enabling RAM, got %#02X", c.Read(0xA000))
	}
}

func TestMBC3ROM(t *testing.T) {
	c := NewMBC3(bigTestData(0x1000000, 0x1000))

	if c.Read(0x1000) != 0x01 || c.Read(0x3000) != 0x03 {
		t.Errorf("Expected NewMBC3 to accept data argument.")
	}

	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC3.Read to find 0x04 at rom bank 0, got %#02X", c.Read(0x4000))
	}

	if c.Read(0x7FFF) != 0x07 {
		t.Errorf("Expected MBC3.Read to find 0x07 at rom bank 0, got %#02X", c.Read(0x7FFF))
	}

	// Select bank 4
	c.Write(0x2000, 4)
	if c.Read(0x4000) != 0x10 {
		t.Errorf("Expected MBC3 to switch to bank 4 and read 0x10, got %#02X", c.Read(0x4000))
	}

	// 0x0000 still shows bank 0
	if c.Read(0x1000) != 0x01 {
		t.Errorf("Expected MBC3.Read to find 0x07 at rom bank 0 while bank 4 is selected, got %#02X", c.Read(0x7FFF))
	}

	// Select bank 0
	c.Write(0x2000, 0)
	if c.Read(0x4000) != 0x04 {
		t.Errorf("Expected MBC3 to switch to bank 0 and read 0x04. got %#02X", c.Read(0x4000))
	}
}

func TestMBC3RAM(t *testing.T) {
	c := NewMBC3(nil)

	// Write to disabled RAM
	c.Write(0xA000, 0xAA)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected MBC3 to not write to RAM when disabled, got %#02X", c.Read(0xA000))
	}

	// Write to enabled RAM
	c.Write(0xA000, 0xBB)
	if c.Read(0xA000) != 0xBB {
		t.Errorf("Expected MBC3 to write/read from enabled ram successfuly, got %#02X", c.Read(0xA000))
	}

	// Switch to bank 3
	c.Write(0x2000, 0x03)

	c.Write(0xA000, 0xCC)
	c.Write(0xB000, 0xDD)

	if c.Read(0xA000) != 0xCC {
		t.Errorf("Expected MBC3 to switch to bank 3 and write/read 0xCC, got %#02X", c.Read(0xA000))
	}

	if c.Read(0xB000) != 0xDD {
		t.Errorf("Expected MBC3 to switch to bank 3 and write/read 0xDD, got %#02X", c.Read(0xB000))
	}

	// Disable RAM
	c.Write(0x0000, 0x0)

	if c.Read(0xA000) != 0xFF {
		t.Errorf("Expected MBC3 to read 0xFF from disabled RAM, got %#02X", c.Read(0xA000))
	}

	c.Write(0xA000, 0xEE)

	// Enable RAM
	c.Write(0x0000, 0x0A)

	if c.Read(0xA000) != 0xCC {
		t.Errorf("Expected MBC3 to read previously written value of 0xCC after re-enabling RAM, got %#02X", c.Read(0xA000))
	}
}

func TestMBC3RTC(t *testing.T) {
	// NOTE: These are basic tests, consider running the rtc test rom from: https://github.com/aaaaaa123456789/rtc3test
	c := NewMBC3(nil)

	// Enable RAM / RTC registers
	c.Write(0x0000, 0x0A)

	// Latch
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	// Select RTC register RCHS
	c.Write(0x4000, 0x08)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected RTC to initialize to 0 seconds, got %#02x", c.Read(0xA000))
	}

	// Wait for 2 minutes 10 seconds
	c.RunForClocks(constants.BaseClockSpeed * 130)

	if c.Read(0xA000) != 0x00 {
		t.Errorf("Expected RTC registers to persist until latched, got %#02x", c.Read(0xA000))
	}

	// Latch
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	if c.Read(0xA000) != 0x0A {
		t.Errorf("Expected RTCS to count 10 seconds, got %#02x", c.Read(0xA000))
	}

	// Select RTCM
	c.Write(0x4000, 0x09)
	if c.Read(0xA000) != 0x02 {
		t.Errorf("Expected RTCMto count 2 minutes, got %#02x", c.Read(0xA000))
	}

	// Wait for 26 hours
	c.RunForClocks(constants.BaseClockSpeed * 60 * 60 * 26)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	c.Write(0x4000, 0x08)
	if c.Read(0xA000) != 10 {
		t.Errorf("Expected RTCS to count 10 seconds, got %d", c.Read(0xA000))
	}

	c.Write(0x4000, 0x09)
	if c.Read(0xA000) != 2 {
		t.Errorf("Expected RTCM to count 2 minutes, got %d", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0A)
	if c.Read(0xA000) != 2 {
		t.Errorf("Expected RTCH to count 2 hours, got %d", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0B)
	if c.Read(0xA000) != 1 {
		t.Errorf("Expected RTCDL to count 1 day, got %d", c.Read(0xA000))
	}

	// Wait for 300 days
	c.RunForClocks(constants.BaseClockSpeed * 60 * 60 * 24 * 300)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	if c.Read(0xA000) != 0x2D {
		t.Errorf("Expected RTCDL to count 0x2D, got %#02X", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0C)
	if c.Read(0xA000) != 0x01 {
		t.Errorf("Expected RTCH to read 0x01, got %#02X", c.Read(0xA000))
	}

	// Halt the RTC
	c.Write(0xA000, 0b0100_0001)

	// Wait 23 seconds while halted
	c.RunForClocks(constants.BaseClockSpeed * 23)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	c.Write(0x4000, 0x08)
	if c.Read(0xA000) != 10 {
		t.Errorf("Expected RTC not to count while halted, read %d", c.Read(0xA000))
	}

	// Resume RTC
	c.Write(0x4000, 0x0C)
	c.Write(0xA000, 0b0000_0001)

	// Wait for 300 more days
	c.RunForClocks(constants.BaseClockSpeed * 60 * 60 * 24 * 300)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	c.Write(0x4000, 0x0B)
	if c.Read(0xA000) != 89 {
		t.Errorf("Expected RTCDL to count 0x2D, got %d", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0C)
	if c.Read(0xA000) != 0b1000_0000 {
		t.Errorf("Expected RTCDH to read 0b10000000, got %b", c.Read(0xA000))
	}

	c.Write(0x4000, 0x08)
	c.Write(0xA000, 0xFF)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)
	if c.Read(0xA000) != 0b0011_1111 {
		t.Errorf("Expected RTCS to be masked to 5 bits, got %b", c.Read(0xA000))
	}

	c.Write(0x4000, 0x09)
	c.Write(0xA000, 0xFF)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)
	if c.Read(0xA000) != 0b0011_1111 {
		t.Errorf("Expected RTCM to be masked to 5 bits, got %b", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0A)
	c.Write(0xA000, 0xFF)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)
	if c.Read(0xA000) != 0b0001_1111 {
		t.Errorf("Expected RTCM to be masked to 4 bits, got %b", c.Read(0xA000))
	}

	c.Write(0x4000, 0x0C)
	c.Write(0xA000, 0xFF)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)
	if c.Read(0xA000) != 0b1100_0001 {
		t.Errorf("Expected RTCDL to be masked to 0b11000001, got %b", c.Read(0xA000))
	}
	c.Write(0xA000, 0x0)

	// Write 61 seconds to RTCS
	c.Write(0x4000, 0x08)
	c.Write(0xA000, 61)

	// Write 0 minutes to RTCM
	c.Write(0x4000, 0x09)
	c.Write(0xA000, 0)

	// Wait 2 seconds
	c.RunForClocks(constants.BaseClockSpeed * 2)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	c.Write(0x4000, 0x09)
	if c.Read(0xA000) != 0 {
		t.Errorf("Expected RTCM to remain 0, got %d", c.Read(0xA000))
	}

	c.Write(0x4000, 0x08)
	if c.Read(0xA000) != 0x3F {
		t.Errorf("Expected RTCS to read 0x3F, got %#02X", c.Read(0xA000))
	}

	// Wait 1 more second
	c.RunForClocks(constants.BaseClockSpeed)
	c.Write(0x6000, 0x00)
	c.Write(0x6000, 0x01)

	c.Write(0x4000, 0x08)
	if c.Read(0xA000) != 0 {
		t.Errorf("Expected RTCS to read 0, got %#02X", c.Read(0xA000))
	}

	c.Write(0x4000, 0x09)
	if c.Read(0xA000) != 0 {
		t.Errorf("Expected RTCM to remain 0, got %d", c.Read(0xA000))
	}
}
