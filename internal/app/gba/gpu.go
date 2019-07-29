package gba

// GPU represents the gameboy's graphics processing unit.
type GPU struct {
	mmu *MMU // Memory Management Unit
}

// NewGPU constructs a valid GPU struct
func NewGPU(mmu *MMU) *GPU {
	gpu := new(GPU)

	gpu.mmu = mmu

	return gpu
}
