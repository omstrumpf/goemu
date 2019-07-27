package gba

// MMU represents the memory management unit.
type MMU struct {
}

// NewMMU constructs a valid MMU struct
func NewMMU() *MMU {
	mmu := new(MMU)

	return mmu
}

// Read returns the 8-bit value from the address
func (mmu *MMU) Read(addr uint16) byte {
	return byte(0) // TODO
}

// Read16 returns the 16-bit value from the address
func (mmu *MMU) Read16(addr uint16) uint16 {
	return uint16(0) // TODO
}

// Write writes the 8-bit value to the address
func (mmu *MMU) Write(addr uint16, val byte) {
	// TODO
}

// Write16 writes the 16-bit value to the address
func (mmu *MMU) Write16(addr uint16, val uint16) {
	// TODO
}
