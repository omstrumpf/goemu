package banking

import (
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/constants"
	"github.com/omstrumpf/goemu/internal/app/log"
)

type rtc struct {
	latched [5]byte
	live    [5]byte

	latchState bool

	clocksSinceLastSecond int
}

func (rtc *rtc) incDays() {
	days := int(rtc.live[3]) | (int(rtc.live[4]&0b0000_0001) << 8)

	days++

	if days == 512 {
		days = 0
		// Set the day carry bit
		rtc.live[4] |= 0b1000_0000
	}

	rtc.live[3] = byte(days & 0xFF)
	if days > 0xFF {
		rtc.live[4] |= 0b0000_0001
	} else {
		rtc.live[4] &= 0b1111_1110
	}
}

func (rtc *rtc) incHours() {
	rtc.live[2] = (rtc.live[2] + 1) & 0b0001_1111

	if rtc.live[2] == 24 || rtc.live[2] > 0b0001_1111 {
		rtc.live[2] = 0
		rtc.incDays()
	}
}

func (rtc *rtc) incMinutes() {
	rtc.live[1] = (rtc.live[1] + 1) & 0b0011_1111

	if rtc.live[1] == 60 || rtc.live[1] > 0b0011_1111 {
		rtc.live[1] = 0
		rtc.incHours()
	}
}

func (rtc *rtc) incSeconds() {
	rtc.live[0] = (rtc.live[0] + 1) & 0b0011_1111

	rtc.live[0] &= 0b0011_1111

	if rtc.live[0] == 60 {
		rtc.live[0] = 0
		rtc.incMinutes()
	}
}

func (rtc *rtc) runForClocks(clocks int) {
	if rtc.live[4]&0b0100_0000 > 0 {
		// RTC is halted
		return
	}

	rtc.clocksSinceLastSecond += clocks

	for rtc.clocksSinceLastSecond >= constants.ClockSpeed {
		rtc.clocksSinceLastSecond -= constants.ClockSpeed
		rtc.incSeconds()
	}
}

func (rtc *rtc) read(addr uint16) byte {
	if addr < 0x05 {
		return rtc.latched[addr]
	}

	log.Errorf("MBC3 RTC encountered read out of range: %#04x", addr)
	return 0xFF
}

func (rtc *rtc) write(addr uint16, val byte) {
	switch addr {
	case 0x00, 0x01:
		rtc.live[addr] = val & 0b0011_1111
	case 0x02:
		rtc.live[addr] = val & 0b0001_1111
	case 0x03:
		rtc.live[addr] = val
	case 0x04:
		rtc.live[addr] = val & 0b1100_0001
	default:
		log.Errorf("MBC3 RTC encountered write out of range: %#04x = %#02x", addr, val)
	}
}

func (rtc *rtc) latch(val byte) {
	switch val {
	case 0x00:
		rtc.latchState = false
	case 0x01:
		if !rtc.latchState {
			// If we became latched, copy live into latched
			rtc.latched = rtc.live
		}
		rtc.latchState = true
	}
}

// MBC3 is the second banked memory controller for gameboy
type MBC3 struct {
	// 2MB ROM, in 64 banks of 16KB.
	rom [0x100000]byte

	// 32KB RAM, in 4 banks of 8KB
	ram [0x8000]byte

	romBank    uint16
	ramTimBank uint16 // RAM and RTC share a single banking register on the MBC3

	rtc rtc

	ramTimEnable     bool
	romRAMModeSelect bool
}

// NewMBC3 constructs a valid MBC3 struct.
func NewMBC3(data []byte) *MBC3 {
	mbc3 := &MBC3{
		romBank:      1,
		ramTimBank:   0,
		ramTimEnable: false,
	}

	if len(data) > len(mbc3.rom) {
		log.Warningf("MBC3 controller loading oversized ROM. Data will be truncated.")
	}

	copy(mbc3.rom[:], data)

	return mbc3
}

// RunForClocks runs the MBC3's RTC for the given number of clock cycles.
func (mbc3 *MBC3) RunForClocks(clocks int) {
	mbc3.rtc.runForClocks(clocks)
}

func (mbc3 *MBC3) Read(addr uint16) byte {
	if addr < 0x4000 {
		// Fixed ROM bank 0
		return mbc3.rom[addr]
	} else if addr < 0x8000 {
		// Variable ROM bank
		bankOffset := 0x4000 * uint32(mbc3.romBank)
		romOffset := int(uint32(addr-0x4000) + bankOffset)
		if romOffset < len(mbc3.rom) {
			return mbc3.rom[romOffset]
		}
		log.Debugf("MBC3 encountered ROM read out of range: %#04x", addr)
		return 0xFF
	} else if addr >= 0xA000 && addr < 0xC000 {
		if mbc3.ramTimEnable {
			if mbc3.ramTimBank < 0x04 {
				// Variable RAM bank
				bankOffset := 0x2000 * uint32(mbc3.ramTimBank)
				ramOffset := int(uint32(addr-0xA000) + bankOffset)
				if ramOffset < len(mbc3.ram) {
					return mbc3.ram[ramOffset]
				}
				log.Debugf("MBC3 encountered RAM read out of range: %#04x", addr)
				return 0xFF
			} else if mbc3.ramTimBank >= 0x08 && mbc3.ramTimBank < 0x0D {
				// RTC registers
				return mbc3.rtc.read(mbc3.ramTimBank - 0x08)
			}

			log.Debugf("MBC3 encountered RAM/RTC read with invalid bank: %#04x (bank %#02x)", addr, mbc3.ramTimBank)
			return 0xFF
		} else {
			log.Debugf("MBC3 encountered read from disabled RAM: %#04x", addr)
			return 0xFF
		}
	} else {
		log.Errorf("MBC3 encountered read out of range: %#04x", addr)
		return 0xFF
	}
}

func (mbc3 *MBC3) Write(addr uint16, val byte) {
	if addr < 0x2000 {
		// RAM / RTC enable
		mbc3.ramTimEnable = (val&0x0A != 0)
	} else if addr < 0x4000 {
		// ROM bank select
		if val == 0 {
			// Bank 0 is not selectable
			val = 1
		}

		mbc3.romBank = uint16(val & 0b0111_1111)
	} else if addr < 0x6000 {
		// RAM / RTC bank select

		if val > 0x0C {
			// Maximum bank
			val = 0x0C
		}

		mbc3.ramTimBank = uint16(val)
	} else if addr < 0x8000 {
		// RTC latch clock
		mbc3.rtc.latch(val)
	} else if addr >= 0xA000 && addr < 0xC000 {
		if mbc3.ramTimEnable {
			if mbc3.ramTimBank < 0x04 {
				// Variable RAM bank
				bankOffset := 0x2000 * uint32(mbc3.ramTimBank)
				ramOffset := int(uint32(addr-0xA000) + bankOffset)
				if ramOffset < len(mbc3.ram) {
					mbc3.ram[ramOffset] = val
				} else {
					log.Debugf("MBC3 encountered RAM write out of range: %#04x", addr)
				}
			} else if mbc3.ramTimBank >= 0x08 && mbc3.ramTimBank < 0x0D {
				// RTC registers
				mbc3.rtc.write(mbc3.ramTimBank-0x08, val)
			} else {
				log.Debugf("MBC3 encountered RAM/RTC write with invalid bank: %#04x = %#02x (bank %#02x)", addr, val, mbc3.ramTimBank)
			}
		} else {
			log.Debugf("MBC3 encountered write to disabled RAM: %#04x = %#02x", addr, val)
		}
	} else {
		log.Errorf("MBC3 encountered write out of range: %#04x = %#02x", addr, val)
	}
}

func (mbc3 *MBC3) GetRamSave() []byte {
	return append(mbc3.rtc.live[:], mbc3.ram[:]...)
}

func (mbc3 *MBC3) LoadRamSave(data []byte) {
	copy(mbc3.rtc.live[:], data[:5])

	if len(data)-5 > len(mbc3.ram) {
		log.Warningf("MBC3 controller loading oversized RAM save. Data will be truncated.")
	}

	copy(mbc3.ram[:], data[5:])
}
