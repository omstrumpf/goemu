package banking

import "github.com/omstrumpf/goemu/internal/app/log"

// MBC1 is the first banked memory controller for gameboy
type MBC1 struct {
	rom []byte
	ram []byte

	romBank uint16
	ramBank uint16

	ramEnable        bool
	romRAMModeSelect bool
}

// NewMBC1 constructs a valid MBC1 struct with the given rom amd ram sizes
func NewMBC1(data []byte, romSize uint32, ramSize uint32) *MBC1 {
	mbc1 := &MBC1{
		rom:              make([]byte, romSize),
		ram:              make([]byte, ramSize),
		romBank:          1,
		ramBank:          0,
		ramEnable:        false,
		romRAMModeSelect: false,
	}

	if len(data) > len(mbc1.rom) {
		log.Warningf("MBC1 controller loading oversized ROM. Data will be truncated.")
	}

	copy(mbc1.rom[:], data)

	return mbc1
}

// RunForClocks is unused on the MBC1
func (mbc1 *MBC1) RunForClocks(clocks int) {}

func (mbc1 *MBC1) Read(addr uint16) byte {
	if addr < 0x4000 {
		return mbc1.rom[addr]
	} else if addr < 0x8000 {
		bankOffset := 0x4000 * uint32(mbc1.romBank)
		romOffset := int(uint32(addr-0x4000) + bankOffset)
		if romOffset < len(mbc1.rom) {
			return mbc1.rom[romOffset]
		}
		log.Debugf("MBC1 encountered ROM read out of range: %#04x", addr)
		return 0xFF
	} else if addr >= 0xA000 && addr < 0xC000 {
		if mbc1.ramEnable {
			bankOffset := 0x2000 * uint32(mbc1.ramBank)
			ramOffset := int(uint32(addr-0xA000) + bankOffset)
			if ramOffset < len(mbc1.ram) {
				return mbc1.ram[ramOffset]
			}
			log.Debugf("MBC1 encountered RAM read out of range: %#04x", addr)
			return 0xFF
		}
		log.Debugf("MBC1 encountered RAM read with RAM disabled: %#04x", addr)
		return 0xFF
	} else {
		log.Errorf("MBC1 encountered read out of range: %#04x", addr)
		return 0xFF
	}
}

func (mbc1 *MBC1) Write(addr uint16, val byte) {
	if addr < 0x2000 {
		mbc1.ramEnable = (val&0x0A != 0)
		log.Tracef("MBC1 setting RAM enable: %t", mbc1.ramEnable)
	} else if addr < 0x4000 {
		mbc1.romBank = (mbc1.romBank & 0x60) | uint16(val&0x1F)
		log.Tracef("MBC1 switching ROM bank to %d", mbc1.romBank)
	} else if addr < 0x6000 {
		if mbc1.romRAMModeSelect { // RAM mode
			mbc1.ramBank = uint16(val & 3)
			log.Tracef("MBC1 switching RAM bank to %d", mbc1.ramBank)
		} else { // ROM mode
			mbc1.romBank = (uint16(val&3) << 5) | (mbc1.romBank & 0x1F)
		}
	} else if addr < 0x8000 {
		mbc1.romRAMModeSelect = (val&1 == 1)
		if !mbc1.romRAMModeSelect {
			mbc1.ramBank = 0
		}
	} else if addr >= 0xA000 && addr < 0xC000 {
		if mbc1.ramEnable {
			bankOffset := 0x2000 * uint32(mbc1.ramBank)
			ramOffset := int(uint32(addr-0xA000) + bankOffset)
			if ramOffset < len(mbc1.ram) {
				mbc1.ram[ramOffset] = val
			} else {
				log.Debugf("MBC1 encountered RAM write out of range: %#04x = %#02x", addr, val)
			}
		} else {
			log.Debugf("MBC1 encountered RAM write with RAM disabled: %#04x = %#02x", addr, val)
		}
	} else {
		log.Errorf("MBC1 encountered write out of range: %#04x = %#02x", addr, val)
	}

	// These rombank values cannot be accessed and always point to the following bank
	if mbc1.romBank == 0x00 || mbc1.romBank == 0x20 || mbc1.romBank == 0x40 || mbc1.romBank == 0x60 {
		mbc1.romBank++
	}
}
