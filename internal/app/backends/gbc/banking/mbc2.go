package banking

import "github.com/omstrumpf/goemu/internal/app/log"

// MBC2 is the second banked memory controller for gameboy
type MBC2 struct {
	// 256k ROM. 16 banks of 16k each
	rom [0x40000]byte

	// 2k total RAM, accessed 4 bits at a time cause there are only 4 data lines.
	ram [0x800]byte

	romBank uint8
	ramBank uint8

	ramEnable bool
}

// NewMBC2 constructs a valid MBC2 struct. ROM/RAM size is fixed for mbc2
func NewMBC2(data []byte) *MBC2 {
	mbc2 := &MBC2{
		romBank:   1,
		ramBank:   0,
		ramEnable: false,
	}

	if len(data) > len(mbc2.rom) {
		log.Warningf("MBC2 controller loading oversized ROM. Data will be truncated.")
	}

	copy(mbc2.rom[:], data)

	return mbc2
}

// RunForClocks is unused on the MBC2
func (mbc2 *MBC2) RunForClocks(clocks int) {}

func (mbc2 *MBC2) Read(addr uint16) byte {
	if addr < 0x4000 {
		// Fixed ROM bank 0
		return mbc2.rom[addr]
	} else if addr < 0x8000 {
		// Variable ROM bank
		bankOffset := 0x4000 * uint32(mbc2.romBank)
		romOffset := int(uint32(addr-0x4000) + bankOffset)
		if romOffset < len(mbc2.rom) {
			return mbc2.rom[romOffset]
		}
		log.Debugf("MBC2 encountered ROM read out of range: %#04x", addr)
		return 0xFF
	} else if addr >= 0xA000 && addr < 0xC000 {
		// RAM
		if mbc2.ramEnable {
			ramOffset := int(addr&0b0000_0001_1111_1111) / 2
			hilo := int(addr-0xA000) % 2

			if ramOffset < len(mbc2.ram) {
				// RAM cells are 4 bits each
				if hilo == 1 {
					return (mbc2.ram[ramOffset] & 0b1111_0000) >> 4
				}
				return mbc2.ram[ramOffset] & 0b0000_1111
			}
			log.Debugf("MBC2 encountered RAM read out of range: %#04x", addr)
			return 0xFF
		}
		return 0xFF
	} else {
		log.Errorf("MBC2 encountered read out of range: %#04x", addr)
		return 0xFF
	}
}

func (mbc2 *MBC2) Write(addr uint16, val byte) {
	if addr < 0x2000 {
		// RAM enable
		if addr&0b0000_0001_0000_0000 == 0 {
			// The least significant bit of the upper address byte must be '0' to set RAM enable
			mbc2.ramEnable = (val&0x0A != 0)
			log.Tracef("MBC2 setting RAM enable: %t", mbc2.ramEnable)
		}
	} else if addr < 0x4000 {
		// ROM bank select
		if val == 0 {
			// Bank 0 is not selectable
			val = 1
		}
		if addr&0b0000_0001_0000_0000 > 0 {
			// The least significant bit of the upper address byte must be '1' to select a ROM bank
			mbc2.romBank = uint8(val & 0b0000_1111)
		}
	} else if addr >= 0xA000 && addr < 0xC000 {
		// RAM
		if mbc2.ramEnable {
			ramOffset := int(addr&0b0000_0001_1111_1111) / 2
			hilo := int(addr-0xA000) % 2

			if ramOffset < len(mbc2.ram) {
				// RAM cells are 4 bits each
				if hilo == 1 {
					mbc2.ram[ramOffset] = (mbc2.ram[ramOffset] & 0b0000_1111) | ((val << 4) & 0b1111_0000)
				} else {
					mbc2.ram[ramOffset] = (mbc2.ram[ramOffset] & 0b1111_0000) | (val & 0b0000_1111)
				}
			}
			log.Debugf("MBC2 encountered RAM write out of range: %#04x = %#02x", addr, val)
		}
	} else {
		log.Errorf("MBC2 encountered write out of range: %#04x = %#02x", addr, val)
	}
}
