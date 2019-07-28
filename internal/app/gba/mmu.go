package gba

import "log"

// MMU represents the memory management unit.
type MMU struct {
	bios []byte
	rom  []byte
	vram []byte
	eram []byte
	wram []byte
	goam []byte
	zram []byte

	zero []byte

	biosEnable bool
}

// NewMMU constructs a valid MMU struct
func NewMMU() *MMU {
	mmu := new(MMU)

	mmu.bios = make([]byte, 0x100)
	mmu.rom = make([]byte, 0x8000)
	mmu.vram = make([]byte, 0x8000)
	mmu.eram = make([]byte, 0x2000)
	mmu.wram = make([]byte, 0x2000)
	mmu.goam = make([]byte, 0xA0)
	mmu.zram = make([]byte, 0x80)

	mmu.zero = []byte{0}

	mmu.biosEnable = true

	return mmu
}

// DisableBios disables the BIOS map over the main ROM
func (mmu *MMU) DisableBios() {
	mmu.biosEnable = false
}

// Read returns the 8-bit value from the address
func (mmu *MMU) Read(addr uint16) byte {
	buffer, offset := mmu.mmapLocation(addr)
	return buffer[offset]
}

// Write writes the 8-bit value to the address
func (mmu *MMU) Write(addr uint16, val byte) {
	buffer, offset := mmu.mmapLocation(addr)

	// Do not write to zero buffer
	if len(buffer) == 1 && buffer[0] == 0 {
		return
	}

	buffer[offset] = val
}

// Read16 returns the 16-bit value from the address
func (mmu *MMU) Read16(addr uint16) uint16 {
	return uint16(mmu.Read(addr)) + (uint16(mmu.Read(addr+1)) << 8)
}

// Write16 writes the 16-bit value to the address
func (mmu *MMU) Write16(addr uint16, val uint16) {
	mmu.Write(addr, byte(val))
	mmu.Write(addr+1, byte(val>>8))
}

func (mmu *MMU) mmapLocation(addr uint16) (buffer []byte, offset uint16) {
	switch addr & 0xF000 {
	// BIOS is mapped over ROM on startup
	case 0x0000:
		if mmu.biosEnable && addr < 0x0100 {
			return mmu.bios, addr
		}
		fallthrough
	case 0x1000:
		fallthrough
	case 0x2000:
		fallthrough
	case 0x3000:
		fallthrough
	case 0x4000:
		fallthrough
	case 0x5000:
		fallthrough
	case 0x6000:
		fallthrough
	case 0x7000:
		return mmu.rom, addr
	case 0x8000:
		fallthrough
	case 0x9000:
		return mmu.vram, addr & 0x1FFF
	case 0xA000:
		fallthrough
	case 0xB000:
		return mmu.eram, addr & 0x1FFF
	case 0xC000:
		fallthrough
	case 0xD000:
		return mmu.wram, addr & 0x1FFF
	// WRAM Shadow
	case 0xE000:
		return mmu.wram, addr & 0x1FFF
	// Shadow, IO, and ZRAM
	case 0xF000:
		switch addr & 0x0F00 {
		case 0x000:
			fallthrough
		case 0x100:
			fallthrough
		case 0x200:
			fallthrough
		case 0x300:
			fallthrough
		case 0x400:
			fallthrough
		case 0x500:
			fallthrough
		case 0x600:
			fallthrough
		case 0x700:
			fallthrough
		case 0x800:
			fallthrough
		case 0x900:
			fallthrough
		case 0xA00:
			fallthrough
		case 0xB00:
			fallthrough
		case 0xC00:
			fallthrough
		case 0xD00:
			return mmu.wram, addr & 0x1FFF
		// GPU OAM
		case 0xE00:
			if addr < 0xFEA0 {
				return mmu.goam, addr & 0xFF
			}
			// Higher addresses always 0
			return mmu.zero, 0
		// Zpage & I/O
		case 0xF00:
			// ZPAGE
			if addr >= 0xFF80 {
				return mmu.zram, addr & 0x7F
			}
			// I/O
			return mmu.zero, 0
		}
	}

	log.Printf("Encountered unmapped memory location: %#4x", addr)
	return mmu.zero, 0
}
