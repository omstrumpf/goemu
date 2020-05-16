package cartridge

// Mode is the mode (CGB or DMG) that the cartridge is intended to be run on.
type Mode int

// Modes
const (
	DMG Mode = 1 // DMG disables CGB functionality and acts as a DMG.
	CGB Mode = 2 // CGB enables CGB functionality
)

func (m Mode) String() string {
	switch m {
	case DMG:
		return "DMG"
	case CGB:
		return "CGB"
	default:
		return "UNKNOWN"
	}
}
