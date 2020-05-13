package gbc

import (
	"testing"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/bios"
)

func TestMMUInit(t *testing.T) {
	mmu := NewMMU(nil)

	if mmu.biosEnable != true {
		t.Error("MMU biosEnable flag should init to true")
	}
}

func TestMMUReadWrite(t *testing.T) {
	mmu := NewMMU(nil)

	got := mmu.Read(0)
	if got != bios.BIOS.Read(0) {
		t.Errorf("Expected initial value to be %#2x, got %#2x", bios.BIOS.Read(0), got)
	}

	mmu.Write(0, 0x12)
	mmu.Write(1, 0x34)
	mmu.Write16(2, 0x7856)

	got = mmu.Read(0)
	if got != 0x12 {
		t.Errorf("Expected to read written value of 0x12, got %#2x", got)
	}
	got16 := mmu.Read16(1)
	if got16 != 0x5634 {
		t.Errorf("Expected to read written value of 0x5634, got %#4x", got16)
	}
	got = mmu.Read(3)
	if got != 0x78 {
		t.Errorf("Expected to read writte nvalue of 0x78, got %#2x", got)
	}
}
