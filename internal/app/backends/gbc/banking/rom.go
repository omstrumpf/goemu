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

func (rom *ROM) Read(addr uint16) byte {
	if addr >= 0x8000 {
		log.Errorf("ROM controller encountered read out of range: %#04x", addr)
		return 0
	}

	return rom.buf[addr]
}

func (rom *ROM) Write(addr uint16, val byte) {
	log.Warningf("ROM controller write encountered: %#04x", addr)
}
