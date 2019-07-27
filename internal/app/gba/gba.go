package gba

// GBA is the toplevel struct containing all the gameboy systems
type GBA struct {
	mmu *MMU
	cpu *CPU
}

// NewGBA constructs a valid GBA struct
func NewGBA() *GBA {
	gba := new(GBA)

	gba.mmu = NewMMU()
	gba.cpu = NewCPU(gba.mmu)

	return gba
}
