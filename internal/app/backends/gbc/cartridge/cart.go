package cartridge

import (
	"fmt"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/banking"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// CART represents a gameboy game cartridge.
type CART struct {
	title        [16]byte
	manufacturer [4]byte
	mode         Mode
	licenseeCode [2]byte
	sgbFlag      bool
	cartType     CartType
	romSize      uint32
	ramSize      uint32

	BankController banking.Controller
}

// NewCart creates a valid CART struct from the given rom data
func NewCart(rom []byte) *CART {
	c := new(CART)

	// Cartridge mode
	switch rom[0x0143] {
	case 0x80:
		c.mode = DMG | CGB // Cartridge supports both DMG and CGB mode
	case 0xC0:
		c.mode = CGB // Cartridge supports only CGB mode
	default:
		c.mode = DMG // Cartridge supports only DMB mode
	}

	// Title and manufacturer code
	if c.mode&CGB == 0 {
		copy(c.title[:], rom[0x0134:0x0144])
	} else {
		copy(c.title[:], rom[0x0134:0x013F])
		copy(c.manufacturer[:], rom[0x013F:0x0143])
	}
	copy(c.licenseeCode[:], rom[0x0144:0x0146])

	// SGB capabilities
	switch rom[0x0146] {
	case 0x03:
		c.sgbFlag = true
	default:
		c.sgbFlag = false
	}

	// Cartridge ROM size
	switch rom[0x0148] {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
		c.romSize = 0x10000 << rom[0x0148]
	case 0x52:
		c.romSize = 0x10000 * 72
	case 0x53:
		c.romSize = 0x10000 * 80
	case 0x54:
		c.romSize = 0x10000 * 96
	}

	// Cartridge RAM size
	switch rom[0x0149] {
	case 0x01:
		c.ramSize = 0x800 // 2KB
	case 0x02:
		c.ramSize = 0x4000 // 8KB
	case 0x03:
		c.ramSize = 0x10000 // 32KB
	default:
		log.Warningf("Unsupported cartridge RAM size (%#02x). Defaulting to 0.", rom[0x0148])
		fallthrough
	case 0x00:
		c.ramSize = 0
	}

	// Memory bank controller
	switch rom[0x0147] {
	case 0x01:
		c.cartType = MBC1
		c.BankController = banking.NewMBC1(rom, c.romSize, 0)
	case 0x02:
		c.cartType = MBC1RAM
		c.BankController = banking.NewMBC1(rom, c.romSize, c.ramSize)
	case 0x03:
		c.cartType = MBC1RAMBAT
		c.BankController = banking.NewMBC1(rom, c.romSize, c.ramSize)
		// TODO implement BAT autosave. Every frame?
	case 0x08:
		c.cartType = ROMRAM
		c.BankController = banking.NewROMRAM(rom)
	default:
		log.Warningf("Unsupported cartridge controller type (%#02x). Defaulting to simple ROM controller.", rom[0x0147])
		fallthrough
	case 0x00:
		c.cartType = ROM
		c.BankController = banking.NewROM(rom)
	}

	return c
}

func (c *CART) Read(addr uint16) byte {
	return c.BankController.Read(addr)
}

func (c *CART) Write(addr uint16, val byte) {
	c.BankController.Write(addr, val)
}

// DebugString returns a debug string describing the CART contents
func (c *CART) DebugString() string {
	return fmt.Sprintf("title: %s\nmanufacturer: %s\ncgb mode: %s\nMBC: %s\ncartridge ROM size: %#x\ncartridge RAM size: %#x\n",
		c.title,
		c.manufacturer,
		c.mode,
		c.cartType,
		c.romSize,
		c.ramSize,
	)
}

// Title returns the cartridge title
func (c *CART) Title() string {
	return string(c.title[:])
}
