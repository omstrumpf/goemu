package gba

import "testing"

func TestCPUInit(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	if cpu.AF.HiLo() != 0 {
		t.Error("AF register should init to 0")
	}
	if cpu.BC.HiLo() != 0 {
		t.Error("BC register should init to 0")
	}
	if cpu.DE.HiLo() != 0 {
		t.Error("BC register should init to 0")
	}
	if cpu.HL.HiLo() != 0 {
		t.Error("BC register should init to 0")
	}
}
