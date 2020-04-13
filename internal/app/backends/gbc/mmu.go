package gbc

const (
	romlen  = 0x8000
	vramlen = 0x2000
	eramlen = 0x2000
	wramlen = 0x2000
	goamlen = 0xA0
	zramlen = 0x80

	totalramlen = 0x10000
)

// MMU represents the memory management unit.
type MMU struct {
	bios memoryDevice
	rom  memoryDevice
	vram memoryDevice
	eram memoryDevice
	wram memoryDevice
	goam memoryDevice
	zram memoryDevice

	inputs     *inputMemoryDevice
	interrupts *interruptMemoryDevice

	zero memoryDevice

	biosEnable bool

	ppu   *PPU
	timer *Timer
}

// NewMMU constructs a valid MMU struct
func NewMMU() *MMU {
	mmu := new(MMU)

	mmu.bios = BIOS
	mmu.rom = newStandardMemoryDevice(romlen)
	mmu.vram = newStandardMemoryDevice(vramlen)
	mmu.eram = newStandardMemoryDevice(eramlen)
	mmu.wram = newStandardMemoryDevice(wramlen)
	mmu.goam = newStandardMemoryDevice(goamlen)
	mmu.zram = newStandardMemoryDevice(zramlen)

	mmu.inputs = newInputMemoryDevice()
	mmu.interrupts = newInterruptMemoryDevice()

	mmu.zero = newZeroMemoryDevice()

	mmu.biosEnable = true

	return mmu
}

// LoadROM loads a ROM into memory
func (mmu *MMU) LoadROM(buf []byte) {
	buflen := len(buf)

	if buflen > totalramlen {
		logger.Warningf("Insufficient memory capacity for ROM: %#4x", buflen)
	}

	j := 0
	for i := 0; i < romlen && i+j < buflen; i++ {
		mmu.rom.Write(uint16(i), buf[i])
	}
	j += romlen

	for i := 0; i < vramlen && i+j < buflen; i++ {
		mmu.vram.Write(uint16(i), buf[i+j])
	}
	j += vramlen

	for i := 0; i < eramlen && i+j < buflen; i++ {
		mmu.eram.Write(uint16(i), buf[i+j])
	}
	j += eramlen

	for i := 0; i < wramlen && i+j < buflen; i++ {
		mmu.eram.Write(uint16(i), buf[i+j])
	}
}

// DisableBios disables the BIOS map over the main ROM
func (mmu *MMU) DisableBios() {
	mmu.biosEnable = false
}

// Read returns the 8-bit value from the address
func (mmu *MMU) Read(addr uint16) byte {
	device, offset := mmu.mmapLocation(addr)
	result := device.Read(offset)
	return result
}

// Write writes the 8-bit value to the address
func (mmu *MMU) Write(addr uint16, val byte) {
	// Traps for MMU on-write functionality
	if addr == 0XFF46 { // DMA
		logger.Tracef("Performing DMA")
		src := uint16(val) << 8
		if src > 0xF100 {
			src = 0xF100
		}

		// TODO technically should be waiting 160us for the transfer to complete.
		// Access to memory should be restricted until it is done (except for HRAM).

		for i := uint16(0); i < 0xA0; i++ {
			mmu.goam.Write(i, mmu.Read(src+i))
		}
		return
	}
	if addr == 0xFF50 { // Disable BIOS memory overlay
		logger.Tracef("Disabling BIOS memory overlay")
		mmu.DisableBios()
		return
	}

	device, offset := mmu.mmapLocation(addr)
	device.Write(offset, val)
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

func (mmu *MMU) mmapLocation(addr uint16) (md memoryDevice, offset uint16) {
	switch addr & 0xF000 {
	// BIOS is mapped over ROM on startup
	case 0x0000:
		if mmu.biosEnable && addr < 0x0100 {
			return mmu.bios, addr
		}
		fallthrough
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000:
		return mmu.rom, addr
	// PPU VRAM
	case 0x8000, 0x9000:
		return mmu.vram, addr & 0x1FFF
	// External RAM
	case 0xA000, 0xB000:
		return mmu.eram, addr & 0x1FFF
	// Working RAM
	case 0xC000, 0xD000:
		return mmu.wram, addr & 0x1FFF
	// WRAM Shadow
	case 0xE000:
		return mmu.wram, addr & 0x1FFF
	// Shadow, IO, and ZRAM
	case 0xF000:
		switch addr & 0x0F00 {
		case 0x000, 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xA00, 0xB00, 0xC00, 0xD00:
			return mmu.wram, addr & 0x1FFF
		// PPU OAM
		case 0xE00:
			if addr < 0xFEA0 {
				return mmu.goam, addr & 0xFF
			}
			// Higher addresses always 0
			return mmu.zero, 0
		// Zpage & I/O
		case 0xF00:
			// Interrupts
			switch addr {
			case 0xFF0F, 0xFFFF:
				return mmu.interrupts, addr
			}

			switch addr & 0x00F0 {
			case 0x00:
				switch addr & 0x000F {
				case 0x0:
					return mmu.inputs, 0
				case 0x4, 0x5, 0x6, 0x7:
					return mmu.timer, addr
				default:
					return mmu.zero, 0 // Unimplemented
				}
			case 0x10, 0x20, 0x30: // Unimplemented
				return mmu.zero, 0
			case 0x40, 0x50, 0x60, 0x70: // PPU Control
				return mmu.ppu.memoryControl, addr
			case 0x80, 0x90, 0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0: // ZRAM
				return mmu.zram, addr & 0x7F
			}
		}
	}

	logger.Warningf("Encountered unmapped memory location: %#4x", addr)
	return mmu.zero, 0
}
