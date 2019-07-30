package gbc

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

	for i, op := range instructions {
		mmu.rom.Write(uint16(i), op)
	}

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
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

func TestInstructionsStack(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	instructions := []byte{
		0x01, 0xCD, 0xAB, // LD BC,nn
		0x11, 0x34, 0x12, // LD DE,nn
		0x21, 0x00, 0x10, // LD HL,nn
		0xF9, // LD SP,HL
		0xC5, // PUSH BC
		0xD5, // PUSH DE
		0xD5, // PUSH DE
		0xD5, // PUSH DE
		0xD1, // POP DE
		0xE1, // POP HL
		0x76, // HALT
	}

	for i, op := range instructions {
		mmu.rom.Write(uint16(i), op)
	}

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
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
	if cpu.SP.HiLo() != 0x0FFC {
		t.Errorf("Expected register SP to contain 0x0FFC, got %#4x", cpu.SP.HiLo())
	}
	if mmu.Read16(0x0FFE) != 0xABCD {
		t.Errorf("Expected memory address 0x0FFE to contain 0xABCD, got %#4x", mmu.Read(0x0FFE))
	}
	if mmu.Read16(0x0FFC) != 0x1234 {
		t.Errorf("Expected memory address 0x0FFC to contain 0x1234, got %#4x", mmu.Read(0x0FFC))
	}
	if mmu.Read16(0x0FFA) != 0x1234 {
		t.Errorf("Expected memory address 0x0FFA to contain 0x1234, got %#4x", mmu.Read(0x0FFA))
	}
	if mmu.Read16(0x0FF8) != 0x1234 {
		t.Errorf("Expected memory address 0x0FF8 to contain 0x1234, got %#4x", mmu.Read(0x0FF8))
	}
	if mmu.Read16(0x0FF6) != 0x0000 {
		t.Errorf("Expected memory address 0x0FF6 to contain 0x0000, got %#4x", mmu.Read(0x0FF6))
	}
	if cpu.clock != 33 {
		t.Errorf("Expected operation to take 33 cycles, got %d", cpu.clock)
	}
}

func TestInstructionsALU(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	instructions := []byte{
		0x21, 0x00, 0x10, // LD HL,nn
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
		0x21, 0x34, 0x12, // LD HL,nn
		0x36, 0x04, // LD (HL),n
		0x86,             // ADD A,(HL)
		0xF5,             // PUSH AF
		0x11, 0x34, 0x12, // LD DE,nn
		0x21, 0x78, 0x56, // LD HL,nn
		0x19, // ADD HL,DE
		0x13, // INC DE
		0x2B, // DEC HL
		0x76, // HALT
	}

	for i, op := range instructions {
		mmu.rom.Write(uint16(i), op)
	}

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
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
	if cpu.clock != 93 {
		t.Errorf("Expected operation to take 93 cycles, got %d", cpu.clock)
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

	addr := uint16(0x0FFF)
	for i, v := range expectedValues {
		if mmu.Read(addr) != v {
			t.Errorf("Expected memory value %d to contain %#2x, got %#2x", i+1, v, mmu.Read(addr))
		}

		addr -= 2
	}
}

func TestInstructionsRot(t *testing.T) {
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	instructions := []byte{
		0x21, 0x00, 0x10, // LD HL,nn
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

	for i, op := range instructions {
		mmu.rom.Write(uint16(i), op)
	}

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.clock != 27 {
		t.Errorf("Expected operation to take 27 cycles, got %d", cpu.clock)
	}

	expectedValues := []byte{
		0x55,
		0xAB,
		0xD5,
		0xEA,
		0x00,
	}

	addr := uint16(0x0FFF)
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
	mmu := NewMMU()

	cpu := NewCPU(mmu)

	instructions := []byte{
		0x21, 0x00, 0x10, // LD HL,nn
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

	for i, op := range instructions {
		mmu.rom.Write(uint16(i), op)
	}

	mmu.DisableBios()
	cpu.PC.Set(0)

	for range instructions {
		if cpu.halt {
			break
		}
		cpu.ProcessNextInstruction()
	}

	if !cpu.IsHalted() {
		t.Errorf("Expected CPU to have halted")
	}
	if cpu.clock != 55 {
		t.Errorf("Expected operation to take 55 cycles, got %d", cpu.clock)
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

	addr := uint16(0x0FFF)
	for i, v := range expectedValues {
		if mmu.Read(addr) != v {
			t.Errorf("Expected memory value %d to contain %#2x, got %#2x", i+1, v, mmu.Read(addr))
		}

		addr -= 2
	}
}

// TODO test bit instructions
