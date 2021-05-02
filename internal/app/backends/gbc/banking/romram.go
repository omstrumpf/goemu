package banking

import "github.com/omstrumpf/goemu/internal/app/log"

// ROMRAM is a basic cartridge memory controller with a fixed 32K ROM, and 8K RAM
type ROMRAM struct {
	rom [0x8000]byte
	ram [0x2000]byte
}

// NewROMRAM constructs a valid ROMRAM struct
func NewROMRAM(data []byte) *ROMRAM {
	romram := new(ROMRAM)

	if len(data) > 0x8000 {
		log.Warningf("ROMRAM controller loading oversized ROM. Data will be truncated.")
	}

	copy(romram.rom[:], data)

	return romram
}

// RunForClocks is unused on the ROMRAM controller
func (romram *ROMRAM) RunForClocks(clocks int) {}

func (romram *ROMRAM) Read(addr uint16) byte {
	if addr < 0x8000 {
		return romram.rom[addr]
	} else if addr >= 0xA000 && addr < 0xC000 {
		return romram.ram[addr-0xA000]
	} else {
		log.Errorf("ROMRAM controller encountered read out of range: %#04x", addr)
		return 0xFF
	}
}

func (romram *ROMRAM) Write(addr uint16, val byte) {
	if addr >= 0xA000 && addr < 0xC000 {
		romram.ram[addr-0xA000] = val
	} else if addr < 0x8000 {
		log.Warningf("ROMRAM controller ROM write encountered: %#04x", addr)
	} else {
		log.Errorf("ROMRAM controller write out of range: %#04x", addr)
	}
}
