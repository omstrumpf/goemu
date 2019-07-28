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

func TestCPUFlags(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	if cpu.z() {
		t.Error("Z flag should init to false")
	}
	if cpu.n() {
		t.Error("N flag should init to false")
	}
	if cpu.h() {
		t.Error("H flag should init to false")
	}
	if cpu.c() {
		t.Error("C flag should init to false")
	}

	cpu.setZ(true)
	if !cpu.z() {
		t.Error("Expected Z flag to be set")
	}
	cpu.setN(true)
	if !cpu.n() {
		t.Error("Expected N flag to be set")
	}
	cpu.setH(true)
	if !cpu.h() {
		t.Error("Expected H flag to be set")
	}
	cpu.setC(true)
	if !cpu.c() {
		t.Error("Expected C flag to be set")
	}

	if cpu.AF.Lo() != 0xF0 {
		t.Errorf("Expected flag register to match 0xF0, got %#2x", cpu.AF.Lo())
	}

	cpu.setZ(false)
	if cpu.z() {
		t.Error("Expected Z flag to be unset")
	}
	cpu.setN(false)
	if cpu.n() {
		t.Error("Expected N flag to be unset")
	}
	cpu.setH(false)
	if cpu.h() {
		t.Error("Expected H flag to be unset")
	}
	cpu.setC(false)
	if cpu.c() {
		t.Error("Expected C flag to be unset")
	}

	if cpu.AF.Lo() != 0x00 {
		t.Errorf("Expected flag register to match 0x00, got %#2x", cpu.AF.Lo())
	}

}
