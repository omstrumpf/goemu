package cartridge

import (
	"github.com/omstrumpf/goemu/internal/app/backends/gbc/banking"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Mode is the mode (CGB or DMG) that the cartridge is intended to be run on.
type Mode int

// Modes
const (
	DMG Mode = 1 // DMG disables CGB functionality and acts as a DMG.
	CGB Mode = 2 // CGB enables CGB functionality
)

// CartType is the type of the cartridge, specifying which memory bank controller and features are present.
type CartType byte

// CartTypes
const (
	ROM CartType = iota
	MBC1
	MBC1RAM
	MBC1RAMBAT
	MBC2
	MBC2BAT
	ROMRAM
	ROMRAMBAT
	MMM01
	MM01RAM
	MM01RAMBAT
	MBC3TIMBAT
	MBC3TIMRAMBAT
	MBC3
	MBC3RAM
	MBC3RAMBAT
	MBC4
	MBC4RAM
	MBC4RAMBAT
	MBC5
	MBC5RAM
	MBC5RAMBAT
	MBC5RUMBLE
	MBC5RUMBLERAM
	MBC5RUMBLERAMBAT
	POCKETCAM
	BANDAITAMA5
	HUC3
	HUC1RAMBAT
)

// CART represents a gameboy game cartridge.
type CART struct {
	title        [16]byte
	manufacturer [4]byte
	mode         Mode
	licenseeCode [2]byte
	sgbFlag      bool
	cartType     CartType

	BankController banking.Controller
}

// NewCart creates a valid CART struct from the given rom data
func NewCart(rom []byte) *CART {
	c := new(CART)

	switch rom[0x0143] {
	case 0x80:
		c.mode = DMG | CGB // Cartridge supports both DMG and CGB mode
	case 0xC0:
		c.mode = CGB // Cartridge supports only CGB mode
	default:
		c.mode = DMG // Cartridge supports only DMB mode
	}

	if c.mode&CGB == 0 {
		copy(c.title[:], rom[0x0134:0x0144])
	} else {
		copy(c.title[:], rom[0x0134:0x013F])
		copy(c.manufacturer[:], rom[0x013F:0x0143])
	}

	copy(c.licenseeCode[:], rom[0x0144:0x0146])

	switch rom[0x0146] {
	case 0x03:
		c.sgbFlag = true
	default:
		c.sgbFlag = false
	}

	switch rom[0x0147] {
	// case 0x01:
	// 	c.cartType = MBC1
	// 	c.BankController = banking.NewMBC1(rom)
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
