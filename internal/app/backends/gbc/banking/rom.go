package banking

import "github.com/omstrumpf/goemu/internal/app/log"

// ROM is a basic cartridge memory controller with a fixed 32K ROM, and no RAM
type ROM struct {
	buf [0x8000]byte
}

// NewROM constructs a valid ROM struct
func NewROM(data []byte) *ROM {
	rom := new(ROM)

	if len(data) > 0x8000 {
		log.Warningf("ROM controller loading oversized ROM. Data will be truncated.")
	}

	copy(rom.buf[:], data)

	return rom
}

// RunForClocks is unused on the ROM controller
func (rom *ROM) RunForClocks(clocks int) {}

func (rom *ROM) Read(addr uint16) byte {
	if addr >= 0x8000 {
		log.Errorf("ROM controller encountered read out of range: %#04x", addr)
		return 0xFF
	}

	return rom.buf[addr]
}

func (rom *ROM) Write(addr uint16, val byte) {
	log.Warningf("ROM controller write encountered: %#04x", addr)
}

func (rom *ROM) GetRamSave() []byte {
	return []byte{}
}

func (rom *ROM) LoadRamSave(data []byte) {
	if len(data) > 0 {
		log.Warningf("RAM controller cannot load RAM save file.")
	}
}
