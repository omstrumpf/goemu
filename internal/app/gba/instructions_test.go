package gba

import "testing"

func TestInstructionsLD(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	instructions := []byte{
		0x3E, 0xAA, // LD A,n
		0x06, 0xBB, // LD B,n
		0x4F,             // LD C,A
		0x50,             // LD D,B
		0x21, 0x34, 0x12, // LD HL,nn
		0x77,             // LD (HL),A
		0x5E,             // LD E,(HL)
		0x21, 0x78, 0x56, // LD HL,nn
		0x36, 0xCC, // LD (HL),n
		0xFA, 0x78, 0x56, // LD A,(nn)
		0x02, // LD (BC),A
		0x76, // HALT
	}

	copy(mmu.rom, instructions)

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
	}

	if !cpu.halt {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.AF.Hi() != 0xCC {
		t.Errorf("Expected register A to contain 0xCC, got %#2x", cpu.AF.Hi())
	}
	if cpu.BC.Hi() != 0xBB {
		t.Errorf("Expected register B to contain 0xBB, got %#2x", cpu.BC.Hi())
	}
	if cpu.BC.Lo() != 0xAA {
		t.Errorf("Expected register C to contain 0xAA, got %#2x", cpu.BC.Lo())
	}
	if cpu.DE.Hi() != 0xBB {
		t.Errorf("Expected register D to contain 0xBB, got %#2x", cpu.DE.Hi())
	}
	if cpu.DE.Lo() != 0xAA {
		t.Errorf("Expected register E to contain 0xAA, got %#2x", cpu.DE.Lo())
	}
	if cpu.HL.HiLo() != 0x5678 {
		t.Errorf("Expected register HL to contain 0x5678, got %#4x", cpu.HL.HiLo())
	}
	if mmu.Read(0x1234) != 0xAA {
		t.Errorf("Expected memory address 0x1234 to contain 0xAA, got %#2x", mmu.Read(0x1234))
	}
	if mmu.Read(0x5678) != 0xCC {
		t.Errorf("Expected memory address 0x5678 to contain 0xCC, got %#2x", mmu.Read(0x5678))
	}
	if mmu.Read(0xBBAA) != 0xCC {
		t.Errorf("Expected memory address 0xBBAA to contain 0xCC, got %#2x", mmu.Read(0xBBAA))
	}
	if cpu.clock != 25 {
		t.Errorf("Expected operation to take 25 cycles, got %d", cpu.clock)
	}
}
