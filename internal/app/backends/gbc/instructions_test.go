package gbc

import "testing"

func TestInstructionsLD(t *testing.T) {
	instructions := []byte{
		0x3E, 0xAA, // LD A,n
		0x06, 0xDD, // LD B,n
		0x4F,             // LD C,A
		0x50,             // LD D,B
		0x21, 0x00, 0xC0, // LD HL,nn
		0x77,             // LD (HL),A
		0x5E,             // LD E,(HL)
		0x21, 0x00, 0xD0, // LD HL,nn
		0x36, 0xCC, // LD (HL),n
		0xFA, 0x00, 0xD0, // LD A,(nn)
		0x02, // LD (BC),A
		0x76, // HALT
	}

	mmu := NewMMU(instructions)

	cpu := NewCPU(mmu)

	clock := 0

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		clock += cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.AF.Hi() != 0xCC {
		t.Errorf("Expected register A to contain 0xCC, got %#2x", cpu.AF.Hi())
	}
	if cpu.BC.Hi() != 0xDD {
		t.Errorf("Expected register B to contain 0xDD, got %#2x", cpu.BC.Hi())
	}
	if cpu.BC.Lo() != 0xAA {
		t.Errorf("Expected register C to contain 0xAA, got %#2x", cpu.BC.Lo())
	}
	if cpu.DE.Hi() != 0xDD {
		t.Errorf("Expected register D to contain 0xDD, got %#2x", cpu.DE.Hi())
	}
	if cpu.DE.Lo() != 0xAA {
		t.Errorf("Expected register E to contain 0xAA, got %#2x", cpu.DE.Lo())
	}
	if cpu.HL.HiLo() != 0xD000 {
		t.Errorf("Expected register HL to contain 0xD000, got %#4x", cpu.HL.HiLo())
	}
	if mmu.Read(0xC000) != 0xAA {
		t.Errorf("Expected memory address 0xC000 to contain 0xAA, got %#2x", mmu.Read(0xC000))
	}
	if mmu.Read(0xD000) != 0xCC {
		t.Errorf("Expected memory address 0xD000 to contain 0xCC, got %#2x", mmu.Read(0xD000))
	}
	if mmu.Read(0xDDAA) != 0xCC {
		t.Errorf("Expected memory address 0xDDAA to contain 0xCC, got %#2x", mmu.Read(0xDDAA))
	}
	if clock != 25 {
		t.Errorf("Expected operation to take 25 cycles, got %d", clock)
	}
}

func TestInstructionsStack(t *testing.T) {
	instructions := []byte{
		0x01, 0xCD, 0xAB, // LD BC,nn
		0x11, 0x34, 0x12, // LD DE,nn
		0x21, 0x00, 0xD0, // LD HL,nn
		0xF9, // LD SP,HL
		0xC5, // PUSH BC
		0xD5, // PUSH DE
		0xD5, // PUSH DE
		0xD5, // PUSH DE
		0xD1, // POP DE
		0xE1, // POP HL
		0x76, // HALT
	}

	mmu := NewMMU(instructions)

	cpu := NewCPU(mmu)

	clock := 0

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		clock += cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.BC.HiLo() != 0xABCD {
		t.Errorf("Expected register BC to contain 0xABCD, got %#4x", cpu.BC.HiLo())
	}
	if cpu.DE.HiLo() != 0x1234 {
		t.Errorf("Expected register DE to contain 0x1234, got %#4x", cpu.DE.HiLo())
	}
	if cpu.HL.HiLo() != 0x1234 {
		t.Errorf("Expected register HL to contain 0x1234, got %#4x", cpu.HL.HiLo())
	}
	if cpu.SP.HiLo() != 0xCFFC {
		t.Errorf("Expected register SP to contain 0xCFFC, got %#4x", cpu.SP.HiLo())
	}
	if mmu.Read16(0xCFFE) != 0xABCD {
		t.Errorf("Expected memory address 0xCFFE to contain 0xABCD, got %#4x", mmu.Read(0xCFFE))
	}
	if mmu.Read16(0xCFFC) != 0x1234 {
		t.Errorf("Expected memory address 0xCFFC to contain 0x1234, got %#4x", mmu.Read(0xCFFC))
	}
	if mmu.Read16(0xCFFA) != 0x1234 {
		t.Errorf("Expected memory address 0xCFFA to contain 0x1234, got %#4x", mmu.Read(0xCFFA))
	}
	if mmu.Read16(0xCFF8) != 0x1234 {
		t.Errorf("Expected memory address 0xCFF8 to contain 0x1234, got %#4x", mmu.Read(0xCFF8))
	}
	if mmu.Read16(0xCFF6) != 0x0000 {
		t.Errorf("Expected memory address 0xCFF6 to contain 0x0000, got %#4x", mmu.Read(0xCFF6))
	}
	if clock != 33 {
		t.Errorf("Expected operation to take 33 cycles, got %d", clock)
	}
}

func TestInstructionsALU(t *testing.T) {
	instructions := []byte{
		0x21, 0x00, 0xD0, // LD HL,nn
		0xF9,       // LD SP,HL
		0xC6, 0x08, // ADD A,n
		0xF5,       // PUSH AF
		0x47,       // LD B,A
		0xC6, 0x08, // ADD A,n
		0xF5,       // PUSH AF
		0x80,       // ADD A,B
		0xF5,       // PUSH AF
		0xC6, 0xFF, // ADD A,n
		0xF5,       // PUSH AF
		0x88,       // ADC A,B
		0xF5,       // PUSH AF
		0x90,       // SUB A,B
		0xF5,       // PUSH AF
		0xDE, 0x01, // SBC A,n
		0xF5,       // PUSH AF
		0xF6, 0x01, // OR n
		0xF5,       // PUSH AF
		0xE6, 0x01, // AND n
		0xF5,       // PUSH AF
		0xEE, 0x01, // XOR n
		0xF5,             // PUSH AF
		0x3C,             // INC A
		0xF5,             // PUSH AF
		0x05,             // DEC B
		0x21, 0xCC, 0xCC, // LD HL,nn
		0x36, 0x04, // LD (HL),n
		0x86,             // ADD A,(HL)
		0xF5,             // PUSH AF
		0x11, 0xCC, 0xCC, // LD DE,nn
		0x21, 0xDD, 0xDD, // LD HL,nn
		0x19, // ADD HL,DE
		0x13, // INC DE
		0x2B, // DEC HL
		0x76, // HALT
	}

	mmu := NewMMU(instructions)

	cpu := NewCPU(mmu)

	clock := 0

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		clock += cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.AF.Hi() != 0x05 {
		t.Errorf("Expected register A to contain 0x05, got %#2x", cpu.AF.Hi())
	}
	if cpu.BC.Hi() != 0x07 {
		t.Errorf("Expected register B to contain 0x07, got %#2x", cpu.BC.Hi())
	}
	if clock != 93 {
		t.Errorf("Expected operation to take 93 cycles, got %d", clock)
	}

	expectedValues := []byte{
		0x08,
		0x10,
		0x18,
		0x17,
		0x20,
		0x18,
		0x17,
		0x17,
		0x01,
		0x00,
		0x01,
		0x05,
		0x00,
	}

	addr := uint16(0xCFFF)
	for i, v := range expectedValues {
		if mmu.Read(addr) != v {
			t.Errorf("Expected memory value %d to contain %#2x, got %#2x", i+1, v, mmu.Read(addr))
		}

		addr -= 2
	}
}

func TestInstructionsRot(t *testing.T) {
	instructions := []byte{
		0x21, 0x00, 0xD0, // LD HL,nn
		0xF9,       // LD SP,HL
		0x3E, 0xAA, // LD A,n
		0x07, // RLCA
		0xF5, // PUSH AF
		0x17, // RLA
		0xF5, // PUSH AF
		0x0F, // RRCA
		0xF5, // PUSH AF
		0x1F, // RRCA
		0xF5, // PUSH AF
		0x76, // HALT
	}

	mmu := NewMMU(instructions)

	cpu := NewCPU(mmu)

	clock := 0

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		clock += cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if clock != 27 {
		t.Errorf("Expected operation to take 27 cycles, got %d", clock)
	}

	expectedValues := []byte{
		0x55,
		0xAB,
		0xD5,
		0xEA,
		0x00,
	}

	addr := uint16(0xCFFF)
	for i, v := range expectedValues {
		if mmu.Read(addr) != v {
			t.Errorf("Expected memory value %d to contain %#2x, got %#2x", i+1, v, mmu.Read(addr))
		}

		addr -= 2
	}
}

// TODO test jumps

// TODO test control instructions

func TestInstructionsCBRot(t *testing.T) {
	instructions := []byte{
		0x21, 0x00, 0xD0, // LD HL,nn
		0xF9,       // LD SP,HL
		0x06, 0xBB, // LD B,n
		0xCB, 0x00, // RLC B
		0xC5,       // PUSH BC
		0xCB, 0x10, // RL B
		0xC5,       // PUSH BC
		0xCB, 0x08, // RRC B
		0xC5,       // PUSH BC
		0xCB, 0x18, // RR B
		0xC5,       // PUSH BC
		0xCB, 0x20, // SLA B
		0xC5,       // PUSH BC
		0xCB, 0x30, // SWAP B
		0xC5,       // PUSH BC
		0xCB, 0x28, // SRA B
		0xC5,       // PUSH BC
		0xCB, 0x38, // SRL B
		0xC5, // PUSH BC
		0x76, // HALT
	}

	mmu := NewMMU(instructions)

	cpu := NewCPU(mmu)

	clock := 0

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		clock += cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if clock != 55 {
		t.Errorf("Expected operation to take 55 cycles, got %d", clock)
	}

	expectedValues := []byte{
		0x77,
		0xEF,
		0xF7,
		0xFB,
		0xF6,
		0x6F,
		0x37,
		0x1B,
		0x00,
	}

	addr := uint16(0xCFFF)
	for i, v := range expectedValues {
		if mmu.Read(addr) != v {
			t.Errorf("Expected memory value %d to contain %#2x, got %#2x", i+1, v, mmu.Read(addr))
		}

		addr -= 2
	}
}

// TODO test bit instructions
