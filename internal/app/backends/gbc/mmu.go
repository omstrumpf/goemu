package gbc

import (
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/banking"
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/bios"
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/interrupts"
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/memory"
	"github.com/omstrumpf/goemu/internal/app/log"
)

const (
	vramlen = 0x2000
	wramlen = 0x2000
	zramlen = 0x80

	totalramlen = 0x10000
)

// MMU represents the memory management unit.
type MMU struct {
	bios memory.Device
	vram memory.Device
	wram memory.Device
	zram memory.Device

	bankController banking.Controller

	inputs     *inputMemoryDevice // TODO why isn't this just a memoryDevice. Move central dispatch/control elsewhere.
	interrupts *interrupts.InterruptDevice

	zero memory.Device
	high memory.Device

	biosEnable bool

	ppu   *PPU
	timer *Timer
}

// NewMMU constructs a valid MMU struct
func NewMMU(bankController banking.Controller) *MMU {
	mmu := new(MMU)

	mmu.bios = bios.BIOS
	mmu.vram = memory.NewSimple(vramlen)
	mmu.wram = memory.NewSimple(wramlen)
	mmu.zram = memory.NewSimple(zramlen)

	mmu.bankController = bankController

	mmu.inputs = newInputMemoryDevice()
	mmu.interrupts = interrupts.NewInterruptDevice()

	mmu.zero = memory.NewZero()
	mmu.high = memory.NewHigh()

	mmu.biosEnable = true

	return mmu
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
		log.Tracef("Performing DMA")
		src := uint16(val) << 8
		if src > 0xF100 {
			src = 0xF100
		}

		// TODO technically should be waiting 160us for the transfer to complete.
		// Access to memory should be restricted until it is done (except for HRAM).

		for i := uint16(0); i < 0xA0; i++ {
			mmu.ppu.oam.Write(i, mmu.Read(src+i))
		}
		return
	}
	if addr == 0xFF50 { // Disable BIOS memory overlay
		log.Tracef("Disabling BIOS memory overlay")
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

func (mmu *MMU) mmapLocation(addr uint16) (md memory.Device, offset uint16) {
	switch addr & 0xF000 {
	// BIOS is mapped over ROM on startup
	case 0x0000:
		if mmu.biosEnable && addr < 0x0100 {
			return mmu.bios, addr
		}
		fallthrough
	// Cartridge ROM
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000:
		return mmu.bankController, addr
	// PPU VRAM
	case 0x8000, 0x9000:
		return mmu.vram, addr & 0x1FFF
	// Cartridge RAM
	case 0xA000, 0xB000:
		return mmu.bankController, addr
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
				return mmu.ppu.oam, addr & 0xFF
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

	log.Warningf("Encountered unmapped memory location: %#4x", addr)
	return mmu.high, 0
}
