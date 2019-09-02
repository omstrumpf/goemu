package gbc

import "log"

// Add adds the operand into the accumulator (A = A + op). If useCarry is set, it also adds the carry bit.
func (cpu *CPU) add(operand byte, useCarry bool) {
	original := cpu.AF.Hi()
	carry := 0
	if useCarry && cpu.c() {
		carry = 1
	}
	result := int16(original) + int16(operand) + int16(carry)

	cpu.AF.SetHi(byte(result))

	cpu.setZ(byte(result) == 0)
	cpu.setN(false)
	cpu.setH((operand&0xF)+(original*0xF)+byte(carry) > 0xF)
	cpu.setC(result > 0xFF)
}

// Sub subtracts the operand from the accumulator (A = A - op). If useCarry is set, it also subtracts the carry bit.
func (cpu *CPU) sub(operand byte, useCarry bool) {
	original := cpu.AF.Hi()
	carry := 0
	if useCarry && cpu.c() {
		carry = 1
	}
	result := int16(original) - int16(operand) - int16(carry)

	cpu.AF.SetHi(byte(result))

	cpu.setZ(byte(result) == 0)
	cpu.setN(true)
	cpu.setH(int16(original&0xF)-int16(operand&0xF)-int16(carry) < 0)
	cpu.setC(result < 0)
}

// And performs a bitwise AND between the accumulator and the operand (A = A AND op).
func (cpu *CPU) and(operand byte) {
	result := cpu.AF.Hi() & operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(true)
	cpu.setC(false)
}

// Xor performs a bitwise XOR between the accumulator and the operand (A = A XOR op).
func (cpu *CPU) xor(operand byte) {
	result := cpu.AF.Hi() ^ operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

// Or performs a bitwise OR between the accumulator and the operand (A = A OR op)
func (cpu *CPU) or(operand byte) {
	result := cpu.AF.Hi() | operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

// Cp compares the accumulator with the operand.
func (cpu *CPU) cp(operand byte) {
	original := cpu.AF.Hi()
	result := original - operand

	cpu.setZ(result == 0)
	cpu.setN(true)
	cpu.setH((original & 0x0F) > (operand & 0x0F))
	cpu.setC(original > operand)
}

// Inc increments the value into the given setter function.
func (cpu *CPU) inc(val byte, setter func(byte)) {
	result := val + 1

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(((val & 0xF) + 1) > 0xF)
}

// Dec decrements the value into the given setter function
func (cpu *CPU) dec(val byte, setter func(byte)) {
	result := val - 1

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(true)
	cpu.setH(val&0xF == 0)
}

// Add16 performs a 16 bit AND into the HL register (HL = HL + val)
func (cpu *CPU) add16(val uint16) {
	result := uint32(cpu.HL.HiLo()) + uint32(val)

	cpu.HL.Set(uint16(result))
	cpu.setN(false)
	cpu.setH(uint32(val&0xFFF) > (result & 0xFFF))
	cpu.setC(result > 0xFFFF)
}

// Add16Signed performs a 16 bit AND into the given setter function
func (cpu *CPU) add16Signed(original uint16, operand int8, setter func(uint16)) {
	result := int32(original) + int32(operand)

	setter(uint16(result))

	// xor operands and result to determine carries
	xor := original ^ uint16(operand) ^ uint16(result)

	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH((xor & 0x10) == 0x10)
	cpu.setC((xor & 0x100) == 0x100)
}

// RotLeft rotates the value into the given setter function. Rotates through the carry bit if set (9 bit rotate)
func (cpu *CPU) rotLeft(original byte, useCarry bool, setter func(byte)) {
	var result byte

	if useCarry {
		carry := 0
		if cpu.c() {
			carry = 1
		}
		result = byte(original<<1) + byte(carry)
	} else {
		result = byte(original<<1) | (original >> 7)
	}

	setter(result)

	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(original > 0x7F)
}

// RotRight rotates the value into the given setter function. Rotates through the carry bit if set (9 bit rotate)
func (cpu *CPU) rotRight(original byte, useCarry bool, setter func(byte)) {
	var result byte

	if useCarry {
		carry := 0
		if cpu.c() {
			carry = 1 << 7
		}
		result = byte(original>>1) | byte(carry)
	} else {
		result = byte(original>>1) | (original << 7)
	}

	setter(result)

	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(1&original == 1)
}

// ShiftLeft shifts the value into the given setter function.
func (cpu *CPU) shiftLeft(original byte, setter func(byte)) {
	result := original << 1

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(original > 0x7F)
}

// ShiftRight shifts the value into the given setter function.
func (cpu *CPU) shiftRight(original byte, logical bool, setter func(byte)) {
	result := original >> 1

	if !logical {
		result |= (original & 0x80)
	}

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(1&original == 1)
}

// Swap swaps the values in the bottom and top halfs of the value into the given setter function.
func (cpu *CPU) swap(original byte, setter func(byte)) {
	result := (original >> 4) | (original << 4)

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

// TestBit tests a bit of the value.
func (cpu *CPU) testBit(val byte, bit uint8) {
	cpu.setZ(val>>bit&1 != 1)
	cpu.setN(false)
	cpu.setH(true)
}

// SetBit sets a bit of the value into the given setter function.
func (cpu *CPU) setBit(original byte, bit uint8, setter func(byte)) {
	setter(original | 1<<bit)
}

// ResetBit resets a bit of the value into the given setter function.
func (cpu *CPU) resetBit(original byte, bit uint8, setter func(byte)) {
	setter(original & ^(1 << bit))
}

// Push pushes the value to the stack
func (cpu *CPU) push(val uint16) {
	cpu.SP.Dec2()
	cpu.mmu.Write16(cpu.SP.HiLo(), val)
}

// Pop pops the stack into the given setter function
func (cpu *CPU) pop(setter func(uint16)) {
	setter(cpu.mmu.Read16(cpu.SP.Inc2()))
}

// Call pushes PC to the stack, then jumps to the target.
func (cpu *CPU) call(target uint16) {
	// Push PC to stack
	cpu.push(cpu.PC.HiLo())

	// Set PC to new target
	cpu.PC.Set(target)
}

// Ret pops the return address off the stack, and jumps to it.
func (cpu *CPU) ret() {
	// Restore PC
	cpu.pop(cpu.PC.Set)
}

// PopulateInstructions populates the cpu opcode maps for instructions and cycle costs
func (cpu *CPU) PopulateInstructions() {
	cpu.instructions = [0x100]func(){
		//// 8-bit loads ////
		0x7F: func() { // LD A,A
			cpu.AF.SetHi(cpu.AF.Hi())
		},
		0x78: func() { // LD A,B
			cpu.AF.SetHi(cpu.BC.Hi())
		},
		0x79: func() { // LD A,C
			cpu.AF.SetHi(cpu.BC.Lo())
		},
		0x7A: func() { // LD A,D
			cpu.AF.SetHi(cpu.DE.Hi())
		},
		0x7B: func() { // LD A,E
			cpu.AF.SetHi(cpu.DE.Lo())
		},
		0x7C: func() { // LD A,H
			cpu.AF.SetHi(cpu.HL.Hi())
		},
		0x7D: func() { // LD A,L
			cpu.AF.SetHi(cpu.HL.Lo())
		},
		0x47: func() { // LD B,A
			cpu.BC.SetHi(cpu.AF.Hi())
		},
		0x40: func() { // LD B,B
			cpu.BC.SetHi(cpu.BC.Hi())
		},
		0x41: func() { // LD B,C
			cpu.BC.SetHi(cpu.BC.Lo())
		},
		0x42: func() { // LD B,D
			cpu.BC.SetHi(cpu.DE.Hi())
		},
		0x43: func() { // LD B,E
			cpu.BC.SetHi(cpu.DE.Lo())
		},
		0x44: func() { // LD B,H
			cpu.BC.SetHi(cpu.HL.Hi())
		},
		0x45: func() { // LD B,L
			cpu.BC.SetHi(cpu.HL.Lo())
		},
		0x4F: func() { // LD C,A
			cpu.BC.SetLo(cpu.AF.Hi())
		},
		0x48: func() { // LD C,B
			cpu.BC.SetLo(cpu.BC.Hi())
		},
		0x49: func() { // LD C,C
			cpu.BC.SetLo(cpu.BC.Lo())
		},
		0x4A: func() { // LD C,D
			cpu.BC.SetLo(cpu.DE.Hi())
		},
		0x4B: func() { // LD C,E
			cpu.BC.SetLo(cpu.DE.Lo())
		},
		0x4C: func() { // LD C,H
			cpu.BC.SetLo(cpu.HL.Hi())
		},
		0x4D: func() { // LD C,L
			cpu.BC.SetLo(cpu.HL.Lo())
		},
		0x57: func() { // LD D,A
			cpu.DE.SetHi(cpu.AF.Hi())
		},
		0x50: func() { // LD D,B
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x51: func() { // LD D,C
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x52: func() { // LD D,D
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x53: func() { // LD D,E
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x54: func() { // LD D,H
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x55: func() { // LD D,L
			cpu.DE.SetHi(cpu.BC.Hi())
		},
		0x5F: func() { // LD E,A
			cpu.DE.SetLo(cpu.AF.Hi())
		},
		0x58: func() { // LD E,B
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x59: func() { // LD E,C
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x5A: func() { // LD E,D
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x5B: func() { // LD E,E
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x5C: func() { // LD E,H
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x5D: func() { // LD E,L
			cpu.DE.SetLo(cpu.BC.Hi())
		},
		0x67: func() { // LD H,A
			cpu.HL.SetHi(cpu.AF.Hi())
		},
		0x60: func() { // LD H,B
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x61: func() { // LD H,C
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x62: func() { // LD H,D
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x63: func() { // LD H,E
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x64: func() { // LD H,H
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x65: func() { // LD H,L
			cpu.HL.SetHi(cpu.BC.Hi())
		},
		0x6F: func() { // LD L,A
			cpu.HL.SetLo(cpu.AF.Hi())
		},
		0x68: func() { // LD L,B
			cpu.HL.SetLo(cpu.BC.Hi())
		},
		0x69: func() { // LD L,C
			cpu.HL.SetLo(cpu.BC.Hi())
		},
		0x6A: func() { // LD L,D
			cpu.HL.SetLo(cpu.BC.Hi())
		},
		0x6B: func() { // LD L,E
			cpu.HL.SetLo(cpu.BC.Hi())
		},
		0x6C: func() { // LD L,H
			cpu.HL.SetLo(cpu.BC.Hi())
		},
		0x6D: func() { // LD L,L
			cpu.HL.SetLo(cpu.BC.Hi())
		},

		0x3E: func() { // LD A,n
			cpu.AF.SetHi(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x06: func() { // LD B,n
			cpu.BC.SetHi(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x0E: func() { // LD C,n
			cpu.BC.SetLo(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x16: func() { // LD D,n
			cpu.DE.SetHi(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x1E: func() { // LD E,n
			cpu.DE.SetLo(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x26: func() { // LD H,n
			cpu.HL.SetHi(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0x2E: func() { // LD L,n
			cpu.HL.SetLo(cpu.mmu.Read(cpu.PC.Inc()))
		},

		0x7E: func() { // LD A,(HL)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x46: func() { // LD B,(HL)
			cpu.BC.SetHi(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x4E: func() { // LD C,(HL)
			cpu.BC.SetLo(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x56: func() { // LD D,(HL)
			cpu.DE.SetHi(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x5E: func() { // LD E,(HL)
			cpu.DE.SetLo(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x66: func() { // LD H,(HL)
			cpu.HL.SetHi(cpu.mmu.Read(cpu.HL.HiLo()))
		},
		0x6E: func() { // LD L,(HL)
			cpu.HL.SetLo(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0x77: func() { // LD (HL),A
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.AF.Hi())
		},
		0x70: func() { // LD (HL),B
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.BC.Hi())
		},
		0x71: func() { // LD (HL),C
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.BC.Lo())
		},
		0x72: func() { // LD (HL),D
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.DE.Hi())
		},
		0x73: func() { // LD (HL),E
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.DE.Lo())
		},
		0x74: func() { // LD (HL),H
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.HL.Hi())
		},
		0x75: func() { // LD (HL),L
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.HL.Lo())
		},

		0x36: func() { // LD (HL),n
			cpu.mmu.Write(cpu.HL.HiLo(), cpu.mmu.Read(cpu.PC.Inc()))
		},

		0x0A: func() { // LD A,(BC)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.BC.HiLo()))
		},
		0x1A: func() { // LD A,(DE)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.DE.HiLo()))
		},
		0xFA: func() { // LD A,(nn)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.mmu.Read16(cpu.PC.Inc2())))
		},

		0x02: func() { // LD (BC),A
			cpu.mmu.Write(cpu.BC.HiLo(), cpu.AF.Hi())
		},
		0x12: func() { // LD (DE),A
			cpu.mmu.Write(cpu.DE.HiLo(), cpu.AF.Hi())
		},
		0xEA: func() { // LD (nn),A
			cpu.mmu.Write(cpu.mmu.Read16(cpu.PC.Inc2()), cpu.AF.Hi())
		},
		0x08: func() { // LD (nn),SP
			cpu.mmu.Write16(cpu.mmu.Read16(cpu.PC.Inc2()), cpu.SP.HiLo())
		},

		0xF2: func() { // LD A,(FF00+C)
			cpu.AF.SetHi(cpu.mmu.Read(0xFF00 + uint16(cpu.BC.Lo())))
		},
		0xE2: func() { // LD (FF00+C),A
			cpu.mmu.Write(0xFF00+uint16(cpu.BC.Lo()), cpu.AF.Hi())
		},
		0xF0: func() { // LD A,(FF00+n)
			cpu.AF.SetHi(cpu.mmu.Read(0xFF00 + uint16(cpu.mmu.Read(cpu.PC.Inc()))))
		},
		0xE0: func() { // LD (FF00+n),A
			cpu.mmu.Write(0xFF00+uint16(cpu.mmu.Read(cpu.PC.Inc())), cpu.AF.Hi())
		},

		0x22: func() { // LDI (HL),A
			cpu.mmu.Write(cpu.HL.Inc(), cpu.AF.Hi())
		},
		0x2A: func() { // LDI A,(HL)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.HL.Inc()))
		},
		0x32: func() { // LDD (HL),A
			cpu.mmu.Write(cpu.HL.Dec(), cpu.AF.Hi())
		},
		0x3A: func() { // LDD A,(HL)
			cpu.AF.SetHi(cpu.mmu.Read(cpu.HL.Dec()))
		},

		//// 16-bit loads ////
		0x01: func() { // LD BC,nn
			cpu.BC.Set(cpu.mmu.Read16(cpu.PC.Inc2()))
		},
		0x11: func() { // LD DE,nn
			cpu.DE.Set(cpu.mmu.Read16(cpu.PC.Inc2()))
		},
		0x21: func() { // LD HL,nn
			cpu.HL.Set(cpu.mmu.Read16(cpu.PC.Inc2()))
		},
		0x31: func() { // LD SP,nn
			cpu.SP.Set(cpu.mmu.Read16(cpu.PC.Inc2()))
		},

		0xF9: func() { // LD SP,HL
			cpu.SP.Set(cpu.HL.HiLo())
		},

		0xC5: func() { // PUSH BC
			cpu.push(cpu.BC.HiLo())
		},
		0xD5: func() { // PUSH DE
			cpu.push(cpu.DE.HiLo())
		},
		0xE5: func() { // PUSH HL
			cpu.push(cpu.HL.HiLo())
		},
		0xF5: func() { // PUSH AF
			cpu.push(cpu.AF.HiLo())
		},

		0xC1: func() { // POP BC
			cpu.pop(cpu.BC.Set)
		},
		0xD1: func() { // POP DE
			cpu.pop(cpu.DE.Set)
		},
		0xE1: func() { // POP HL
			cpu.pop(cpu.HL.Set)
		},
		0xF1: func() { // POP AF
			cpu.pop(cpu.AF.Set)
		},

		//// 8-bit ALU ////
		0x87: func() { // ADD A,A
			cpu.add(cpu.AF.Hi(), false)
		},
		0x80: func() { // ADD A,B
			cpu.add(cpu.BC.Hi(), false)
		},
		0x81: func() { // ADD A,C
			cpu.add(cpu.BC.Lo(), false)
		},
		0x82: func() { // ADD A,D
			cpu.add(cpu.DE.Hi(), false)
		},
		0x83: func() { // ADD A,E
			cpu.add(cpu.DE.Lo(), false)
		},
		0x84: func() { // ADD A,H
			cpu.add(cpu.HL.Hi(), false)
		},
		0x85: func() { // ADD A,L
			cpu.add(cpu.HL.Lo(), false)
		},
		0xC6: func() { // ADD A,n
			cpu.add(cpu.mmu.Read(cpu.PC.Inc()), false)
		},
		0x86: func() { // ADD A,(HL)
			cpu.add(cpu.mmu.Read(cpu.HL.HiLo()), false)
		},

		0x8F: func() { // ADC A,A
			cpu.add(cpu.AF.Hi(), true)
		},
		0x88: func() { // ADC A,B
			cpu.add(cpu.BC.Hi(), true)
		},
		0x89: func() { // ADC A,C
			cpu.add(cpu.BC.Lo(), true)
		},
		0x8A: func() { // ADC A,D
			cpu.add(cpu.DE.Hi(), true)
		},
		0x8B: func() { // ADC A,E
			cpu.add(cpu.DE.Lo(), true)
		},
		0x8C: func() { // ADC A,H
			cpu.add(cpu.HL.Hi(), true)
		},
		0x8D: func() { // ADC A,L
			cpu.add(cpu.HL.Lo(), true)
		},
		0xCE: func() { // ADC A,n
			cpu.add(cpu.mmu.Read(cpu.PC.Inc()), true)
		},
		0x8E: func() { // ADC A,(HL)
			cpu.add(cpu.mmu.Read(cpu.HL.HiLo()), true)
		},

		0x97: func() { // SUB A,A
			cpu.sub(cpu.AF.Hi(), false)
		},
		0x90: func() { // SUB A,B
			cpu.sub(cpu.BC.Hi(), false)
		},
		0x91: func() { // SUB A,C
			cpu.sub(cpu.BC.Lo(), false)
		},
		0x92: func() { // SUB A,D
			cpu.sub(cpu.DE.Hi(), false)
		},
		0x93: func() { // SUB A,E
			cpu.sub(cpu.DE.Lo(), false)
		},
		0x94: func() { // SUB A,H
			cpu.sub(cpu.HL.Hi(), false)
		},
		0x95: func() { // SUB A,L
			cpu.sub(cpu.HL.Lo(), false)
		},
		0xD6: func() { // SUB A,n
			cpu.sub(cpu.mmu.Read(cpu.PC.Inc()), false)
		},
		0x96: func() { // SUB A,(HL)
			cpu.sub(cpu.mmu.Read(cpu.HL.HiLo()), false)
		},

		0x9F: func() { // SBC A,A
			cpu.sub(cpu.AF.Hi(), true)
		},
		0x98: func() { // SBC A,B
			cpu.sub(cpu.BC.Hi(), true)
		},
		0x99: func() { // SBC A,C
			cpu.sub(cpu.BC.Lo(), true)
		},
		0x9A: func() { // SBC A,D
			cpu.sub(cpu.DE.Hi(), true)
		},
		0x9B: func() { // SBC A,E
			cpu.sub(cpu.DE.Lo(), true)
		},
		0x9C: func() { // SBC A,H
			cpu.sub(cpu.HL.Hi(), true)
		},
		0x9D: func() { // SBC A,L
			cpu.sub(cpu.HL.Lo(), true)
		},
		0xDE: func() { // SBC A,n
			cpu.sub(cpu.mmu.Read(cpu.PC.Inc()), true)
		},
		0x9E: func() { // SBC A,(HL)
			cpu.sub(cpu.mmu.Read(cpu.HL.HiLo()), true)
		},

		0xA7: func() { // AND A
			cpu.and(cpu.AF.Hi())
		},
		0xA0: func() { // AND B
			cpu.and(cpu.BC.Hi())
		},
		0xA1: func() { // AND C
			cpu.and(cpu.BC.Lo())
		},
		0xA2: func() { // AND D
			cpu.and(cpu.DE.Hi())
		},
		0xA3: func() { // AND E
			cpu.and(cpu.DE.Lo())
		},
		0xA4: func() { // AND H
			cpu.and(cpu.HL.Hi())
		},
		0xA5: func() { // AND L
			cpu.and(cpu.HL.Lo())
		},
		0xE6: func() { // AND n
			cpu.and(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xA6: func() { // AND (HL)
			cpu.and(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xAF: func() { // XOR A
			cpu.xor(cpu.AF.Hi())
		},
		0xA8: func() { // XOR B
			cpu.xor(cpu.BC.Hi())
		},
		0xA9: func() { // XOR C
			cpu.xor(cpu.BC.Lo())
		},
		0xAA: func() { // XOR D
			cpu.xor(cpu.DE.Hi())
		},
		0xAB: func() { // XOR E
			cpu.xor(cpu.DE.Lo())
		},
		0xAC: func() { // XOR H
			cpu.xor(cpu.HL.Hi())
		},
		0xAD: func() { // XOR L
			cpu.xor(cpu.HL.Lo())
		},
		0xEE: func() { // XOR n
			cpu.xor(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xAE: func() { // XOR (HL)
			cpu.xor(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xB7: func() { // OR A
			cpu.or(cpu.AF.Hi())
		},
		0xB0: func() { // OR B
			cpu.or(cpu.BC.Hi())
		},
		0xB1: func() { // OR C
			cpu.or(cpu.BC.Lo())
		},
		0xB2: func() { // OR D
			cpu.or(cpu.DE.Hi())
		},
		0xB3: func() { // OR E
			cpu.or(cpu.DE.Lo())
		},
		0xB4: func() { // OR H
			cpu.or(cpu.HL.Hi())
		},
		0xB5: func() { // OR L
			cpu.or(cpu.HL.Lo())
		},
		0xF6: func() { // OR n
			cpu.or(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xB6: func() { // OR (HL)
			cpu.or(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xBF: func() { // CP A
			cpu.cp(cpu.AF.Hi())
		},
		0xB8: func() { // CP B
			cpu.cp(cpu.BC.Hi())
		},
		0xB9: func() { // CP C
			cpu.cp(cpu.BC.Lo())
		},
		0xBA: func() { // CP D
			cpu.cp(cpu.DE.Hi())
		},
		0xBB: func() { // CP E
			cpu.cp(cpu.DE.Lo())
		},
		0xBC: func() { // CP H
			cpu.cp(cpu.HL.Hi())
		},
		0xBD: func() { // CP L
			cpu.cp(cpu.HL.Lo())
		},
		0xFE: func() { // CP n
			cpu.cp(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xBE: func() { // CP (HL)
			cpu.cp(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0x3C: func() { // INC A
			cpu.inc(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x04: func() { // INC B
			cpu.inc(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x0C: func() { // INC C
			cpu.inc(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x14: func() { // INC D
			cpu.inc(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x1C: func() { // INC E
			cpu.inc(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x24: func() { // INC H
			cpu.inc(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x2C: func() { // INC L
			cpu.inc(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x34: func() { // INC (HL)
			addr := cpu.HL.HiLo()
			cpu.inc(cpu.mmu.Read(addr), func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x3D: func() { // DEC A
			cpu.dec(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x05: func() { // DEC B
			cpu.dec(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x0D: func() { // DEC C
			cpu.dec(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x15: func() { // DEC D
			cpu.dec(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x1D: func() { // DEC E
			cpu.dec(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x25: func() { // DEC H
			cpu.dec(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x2D: func() { // DEC L
			cpu.dec(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x35: func() { // DEC (HL)
			addr := cpu.HL.HiLo()
			cpu.dec(cpu.mmu.Read(addr), func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x27: func() { // DAA
			// BCD correction on register A

			if !cpu.n() { // If last instruction was an addition
				// Add 0x60 if large or carry
				if cpu.c() || cpu.AF.Hi() > 0x99 {
					cpu.AF.SetHi(cpu.AF.Hi() + 0x60)
					cpu.setC(true)
				}

				// Add 0x06 if half-large or half-carry
				if cpu.h() || cpu.AF.Hi()&0xF > 0x9 {
					cpu.AF.SetHi(cpu.AF.Hi() + 0x06)
					cpu.setH(false)
				}
			} else if cpu.c() && cpu.h() { // If subtraction and both carries are set
				cpu.AF.SetHi(cpu.AF.Hi() + 0x9A)
				cpu.setH(false)
			} else if cpu.c() { // If subtraction and only high carry is set
				cpu.AF.SetHi(cpu.AF.Hi() + 0xA0)
			} else if cpu.h() { // If subtraction and only half carry is set
				cpu.AF.SetHi(cpu.AF.Hi() + 0xFA)
				cpu.setH(false)
			}
			cpu.setZ(cpu.AF.Hi() == 0)
		},
		0x2F: func() { // CPL
			cpu.AF.SetHi(cpu.AF.Hi() ^ 0xFF)
			cpu.setN(true)
			cpu.setH(true)
		},

		//// 16-bit ALU ////
		0x09: func() { // ADD HL,BC
			cpu.add16(cpu.BC.HiLo())
		},
		0x19: func() { // ADD HL,DE
			cpu.add16(cpu.DE.HiLo())
		},
		0x29: func() { // ADD HL,HL
			cpu.add16(cpu.HL.HiLo())
		},
		0x39: func() { // ADD HL,SP
			cpu.add16(cpu.SP.HiLo())
		},

		0x03: func() { // INC BC
			cpu.BC.Set(cpu.BC.HiLo() + 1)
		},
		0x13: func() { // INC DE
			cpu.DE.Set(cpu.DE.HiLo() + 1)
		},
		0x23: func() { // INC HL
			cpu.HL.Set(cpu.HL.HiLo() + 1)
		},
		0x33: func() { // INC SP
			cpu.SP.Set(cpu.SP.HiLo() + 1)
		},

		0x0B: func() { // DEC BC
			cpu.BC.Set(cpu.BC.HiLo() - 1)
		},
		0x1B: func() { // DEC DE
			cpu.DE.Set(cpu.DE.HiLo() - 1)
		},
		0x2B: func() { // DEC HL
			cpu.HL.Set(cpu.HL.HiLo() - 1)
		},
		0x3B: func() { // DEC SP
			cpu.SP.Set(cpu.SP.HiLo() - 1)
		},

		0xE8: func() { // ADD SP,d
			cpu.add16Signed(cpu.SP.HiLo(), int8(cpu.PC.Inc()), cpu.SP.Set)
		},
		0xF8: func() { // LD HL,SP,d
			cpu.add16Signed(cpu.SP.HiLo(), int8(cpu.PC.Inc()), cpu.HL.Set)
		},

		//// Rotate / Shift ////
		0x07: func() { // RLCA
			cpu.rotLeft(cpu.AF.Hi(), false, cpu.AF.SetHi)
		},
		0x17: func() { // RLA
			cpu.rotLeft(cpu.AF.Hi(), true, cpu.AF.SetHi)
		},
		0x0F: func() { // RRCA
			cpu.rotRight(cpu.AF.Hi(), false, cpu.AF.SetHi)
		},
		0x1F: func() { // RRA
			cpu.rotRight(cpu.AF.Hi(), true, cpu.AF.SetHi)
		},

		//// Control ////
		0x3F: func() { // CCF
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(!cpu.c())
		},
		0x37: func() { // SCF
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(true)
		},
		0x00: func() { // NOP
		},
		0x76: func() { // HALT
			cpu.halt = true
		},
		0x10: func() { // STOP
			cpu.stop = true
			cpu.PC.Inc()
		},
		0xF3: func() { // DI
			cpu.ime = false
		},
		0xFB: func() { // EI
			cpu.ime = true
		},

		//// Jump /////
		0xC3: func() { // JP nn
			cpu.PC.Set(cpu.mmu.Read16(cpu.PC.Inc2()))
		},
		0xE9: func() { // JP HL
			cpu.PC.Set(cpu.HL.HiLo())
		},

		0xC2: func() { // JP NZ,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())
			if !cpu.z() {
				cpu.PC.Set(target)
				cpu.clock++
			}
		},
		0xCA: func() { // JP Z,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())
			if cpu.z() {
				cpu.PC.Set(target)
				cpu.clock++
			}
		},
		0xD2: func() { // JP NC,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())
			if !cpu.c() {
				cpu.PC.Set(target)
				cpu.clock++
			}
		},
		0xDA: func() { // JP C,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())
			if cpu.c() {
				cpu.PC.Set(target)
				cpu.clock++
			}
		},

		0x18: func() { // JR n
			offset := int8(cpu.mmu.Read(cpu.PC.Inc()))
			cpu.PC.Set(uint16(int32(cpu.PC.HiLo()) + int32(offset)))
		},

		0x20: func() { // JR NZ,n
			offset := int8(cpu.mmu.Read(cpu.PC.Inc()))
			if !cpu.z() {
				cpu.PC.Set(uint16(int32(cpu.PC.HiLo()) + int32(offset)))
				cpu.clock++
			}
		},
		0x28: func() { // JR Z,n
			offset := int8(cpu.mmu.Read(cpu.PC.Inc()))
			if cpu.z() {
				cpu.PC.Set(uint16(int32(cpu.PC.HiLo()) + int32(offset)))
				cpu.clock++
			}
		},
		0x30: func() { // JR NC,n
			offset := int8(cpu.mmu.Read(cpu.PC.Inc()))
			if !cpu.c() {
				cpu.PC.Set(uint16(int32(cpu.PC.HiLo()) + int32(offset)))
				cpu.clock++
			}
		},
		0x38: func() { // JR C,n
			offset := int8(cpu.mmu.Read(cpu.PC.Inc()))
			if cpu.c() {
				cpu.PC.Set(uint16(int32(cpu.PC.HiLo()) + int32(offset)))
				cpu.clock++
			}
		},

		0xCD: func() { // CALL nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())

			cpu.call(target)
		},
		0xC4: func() { // CALL NZ,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())

			if !cpu.z() {
				cpu.call(target)
				cpu.clock += 3
			}
		},
		0xCC: func() { // CALL Z,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())

			if cpu.z() {
				cpu.call(target)
				cpu.clock += 3
			}
		},
		0xD4: func() { // CALL NC,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())

			if !cpu.c() {
				cpu.call(target)
				cpu.clock += 3
			}
		},
		0xDC: func() { // CALL C,nn
			target := cpu.mmu.Read16(cpu.PC.Inc2())

			if cpu.c() {
				cpu.call(target)
				cpu.clock += 3
			}
		},

		0xC9: func() { // RET
			cpu.ret()
		},
		0xC0: func() { // RET NZ
			if !cpu.z() {
				cpu.ret()
				cpu.clock += 3
			}
		},
		0xC8: func() { // RET Z
			if cpu.z() {
				cpu.ret()
				cpu.clock += 3
			}
		},
		0xD0: func() { // RET NC
			if !cpu.c() {
				cpu.ret()
				cpu.clock += 3
			}
		},
		0xD8: func() { // RET C
			if cpu.c() {
				cpu.ret()
				cpu.clock += 3
			}
		},
		0xD9: func() { // RETI
			cpu.ret()
			cpu.ime = true
		},

		0xC7: func() { // RST 0x00
			cpu.call(0x0000)
		},
		0xCF: func() { // RST 0x08
			cpu.call(0x0008)
		},
		0xD7: func() { // RST 0x10
			cpu.call(0x0010)
		},
		0xDF: func() { // RST 0x18
			cpu.call(0x0018)
		},
		0xE7: func() { // RST 0x20
			cpu.call(0x0020)
		},
		0xEF: func() { // RST 0x28
			cpu.call(0x0028)
		},
		0xF7: func() { // RST 0x30
			cpu.call(0x0030)
		},
		0xFF: func() { // RST 0x38
			cpu.call(0x0038)
		},

		//// CB Prefix ////
		0xCB: func() {
			opcode := cpu.mmu.Read(cpu.PC.Inc())

			cpu.instructionsCB[opcode]()

			cpu.clock += cpu.cyclesCB[opcode]
		},
	}

	for k, v := range cpu.instructions {
		if v == nil {
			opcode := k
			cpu.instructions[k] = func() {
				log.Printf("Encountered unknown instruction: %#2x", opcode)
				cpu.stop = true
			}
		}
	}

	cpu.instructionsCB = [0x100]func(){
		0x07: func() { // RLC A
			cpu.rotLeft(cpu.AF.Hi(), false, cpu.AF.SetHi)
		},
		0x00: func() { // RLC B
			cpu.rotLeft(cpu.BC.Hi(), false, cpu.BC.SetHi)
		},
		0x01: func() { // RLC C
			cpu.rotLeft(cpu.BC.Lo(), false, cpu.BC.SetLo)
		},
		0x02: func() { // RLC D
			cpu.rotLeft(cpu.DE.Hi(), false, cpu.DE.SetHi)
		},
		0x03: func() { // RLC E
			cpu.rotLeft(cpu.DE.Lo(), false, cpu.DE.SetLo)
		},
		0x04: func() { // RLC H
			cpu.rotLeft(cpu.HL.Hi(), false, cpu.HL.SetHi)
		},
		0x05: func() { // RLC L
			cpu.rotLeft(cpu.HL.Lo(), false, cpu.HL.SetLo)
		},
		0x06: func() { // RLC (HL)
			addr := cpu.HL.HiLo()
			cpu.rotLeft(cpu.mmu.Read(addr), false, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x17: func() { // RL A
			cpu.rotLeft(cpu.AF.Hi(), true, cpu.AF.SetHi)
		},
		0x10: func() { // RL B
			cpu.rotLeft(cpu.BC.Hi(), true, cpu.BC.SetHi)
		},
		0x11: func() { // RL C
			cpu.rotLeft(cpu.BC.Lo(), true, cpu.BC.SetLo)
		},
		0x12: func() { // RL D
			cpu.rotLeft(cpu.DE.Hi(), true, cpu.DE.SetHi)
		},
		0x13: func() { // RL E
			cpu.rotLeft(cpu.DE.Lo(), true, cpu.DE.SetLo)
		},
		0x14: func() { // RL H
			cpu.rotLeft(cpu.HL.Hi(), true, cpu.HL.SetHi)
		},
		0x15: func() { // RL L
			cpu.rotLeft(cpu.HL.Lo(), true, cpu.HL.SetLo)
		},
		0x16: func() { // RL (HL)
			addr := cpu.HL.HiLo()
			cpu.rotLeft(cpu.mmu.Read(addr), true, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x0F: func() { // RRC A
			cpu.rotRight(cpu.AF.Hi(), false, cpu.AF.SetHi)
		},
		0x08: func() { // RRC B
			cpu.rotRight(cpu.BC.Hi(), false, cpu.BC.SetHi)
		},
		0x09: func() { // RRC C
			cpu.rotRight(cpu.BC.Lo(), false, cpu.BC.SetLo)
		},
		0x0A: func() { // RRC D
			cpu.rotRight(cpu.DE.Hi(), false, cpu.DE.SetHi)
		},
		0x0B: func() { // RRC E
			cpu.rotRight(cpu.DE.Lo(), false, cpu.DE.SetLo)
		},
		0x0C: func() { // RRC H
			cpu.rotRight(cpu.HL.Hi(), false, cpu.HL.SetHi)
		},
		0x0D: func() { // RRC L
			cpu.rotRight(cpu.HL.Lo(), false, cpu.HL.SetLo)
		},
		0x0E: func() { // RRC (HL)
			addr := cpu.HL.HiLo()
			cpu.rotRight(cpu.mmu.Read(addr), false, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x1F: func() { // RR A
			cpu.rotRight(cpu.AF.Hi(), true, cpu.AF.SetHi)
		},
		0x18: func() { // RR B
			cpu.rotRight(cpu.BC.Hi(), true, cpu.BC.SetHi)
		},
		0x19: func() { // RR C
			cpu.rotRight(cpu.BC.Lo(), true, cpu.BC.SetLo)
		},
		0x1A: func() { // RR D
			cpu.rotRight(cpu.DE.Hi(), true, cpu.DE.SetHi)
		},
		0x1B: func() { // RR E
			cpu.rotRight(cpu.DE.Lo(), true, cpu.DE.SetLo)
		},
		0x1C: func() { // RR H
			cpu.rotRight(cpu.HL.Hi(), true, cpu.HL.SetHi)
		},
		0x1D: func() { // RR L
			cpu.rotRight(cpu.HL.Lo(), true, cpu.HL.SetLo)
		},
		0x1E: func() { // RR (HL)
			addr := cpu.HL.HiLo()
			cpu.rotRight(cpu.mmu.Read(addr), true, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x27: func() { // SLA A
			cpu.shiftLeft(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x20: func() { // SLA B
			cpu.shiftLeft(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x21: func() { // SLA C
			cpu.shiftLeft(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x22: func() { // SLA D
			cpu.shiftLeft(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x23: func() { // SLA E
			cpu.shiftLeft(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x24: func() { // SLA H
			cpu.shiftLeft(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x25: func() { // SLA L
			cpu.shiftLeft(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x26: func() { // SLA (HL)
			addr := cpu.HL.HiLo()
			cpu.shiftLeft(cpu.mmu.Read(addr), func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x37: func() { // SWAP A
			cpu.swap(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x30: func() { // SWAP B
			cpu.swap(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x31: func() { // SWAP C
			cpu.swap(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x32: func() { // SWAP D
			cpu.swap(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x33: func() { // SWAP E
			cpu.swap(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x34: func() { // SWAP H
			cpu.swap(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x35: func() { // SWAP L
			cpu.swap(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x36: func() { // SWAP (HL)
			addr := cpu.HL.HiLo()
			cpu.swap(cpu.mmu.Read(addr), func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x2F: func() { // SRA A
			cpu.shiftRight(cpu.AF.Hi(), false, cpu.AF.SetHi)
		},
		0x28: func() { // SRA B
			cpu.shiftRight(cpu.BC.Hi(), false, cpu.BC.SetHi)
		},
		0x29: func() { // SRA C
			cpu.shiftRight(cpu.BC.Lo(), false, cpu.BC.SetLo)
		},
		0x2A: func() { // SRA D
			cpu.shiftRight(cpu.DE.Hi(), false, cpu.DE.SetHi)
		},
		0x2B: func() { // SRA E
			cpu.shiftRight(cpu.DE.Lo(), false, cpu.DE.SetLo)
		},
		0x2C: func() { // SRA H
			cpu.shiftRight(cpu.HL.Hi(), false, cpu.HL.SetHi)
		},
		0x2D: func() { // SRA L
			cpu.shiftRight(cpu.HL.Lo(), false, cpu.HL.SetLo)
		},
		0x2E: func() { // SRA (HL)
			addr := cpu.HL.HiLo()
			cpu.shiftRight(cpu.mmu.Read(addr), false, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x3F: func() { // SRL A
			cpu.shiftRight(cpu.AF.Hi(), true, cpu.AF.SetHi)
		},
		0x38: func() { // SRL B
			cpu.shiftRight(cpu.BC.Hi(), true, cpu.BC.SetHi)
		},
		0x39: func() { // SRL C
			cpu.shiftRight(cpu.BC.Lo(), true, cpu.BC.SetLo)
		},
		0x3A: func() { // SRL D
			cpu.shiftRight(cpu.DE.Hi(), true, cpu.DE.SetHi)
		},
		0x3B: func() { // SRL E
			cpu.shiftRight(cpu.DE.Lo(), true, cpu.DE.SetLo)
		},
		0x3C: func() { // SRL H
			cpu.shiftRight(cpu.HL.Hi(), true, cpu.HL.SetHi)
		},
		0x3D: func() { // SRL L
			cpu.shiftRight(cpu.HL.Lo(), true, cpu.HL.SetLo)
		},
		0x3E: func() { // SRL (HL)
			addr := cpu.HL.HiLo()
			cpu.shiftRight(cpu.mmu.Read(addr), true, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x47: func() { // BIT 0,A
			cpu.testBit(cpu.AF.Hi(), 0)
		},
		0x40: func() { // BIT 0,B
			cpu.testBit(cpu.BC.Hi(), 0)
		},
		0x41: func() { // BIT 0,C
			cpu.testBit(cpu.BC.Lo(), 0)
		},
		0x42: func() { // BIT 0,D
			cpu.testBit(cpu.DE.Hi(), 0)
		},
		0x43: func() { // BIT 0,E
			cpu.testBit(cpu.DE.Lo(), 0)
		},
		0x44: func() { // BIT 0,H
			cpu.testBit(cpu.HL.Hi(), 0)
		},
		0x45: func() { // BIT 0,L
			cpu.testBit(cpu.HL.Lo(), 0)
		},
		0x46: func() { // BIT 0,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 0)
		},

		0x4F: func() { // BIT 1,A
			cpu.testBit(cpu.AF.Hi(), 1)
		},
		0x48: func() { // BIT 1,B
			cpu.testBit(cpu.BC.Hi(), 1)
		},
		0x49: func() { // BIT 1,C
			cpu.testBit(cpu.BC.Lo(), 1)
		},
		0x4A: func() { // BIT 1,D
			cpu.testBit(cpu.DE.Hi(), 1)
		},
		0x4B: func() { // BIT 1,E
			cpu.testBit(cpu.DE.Lo(), 1)
		},
		0x4C: func() { // BIT 1,H
			cpu.testBit(cpu.HL.Hi(), 1)
		},
		0x4D: func() { // BIT 1,L
			cpu.testBit(cpu.HL.Lo(), 1)
		},
		0x4E: func() { // BIT 1,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 1)
		},

		0x57: func() { // BIT 2,A
			cpu.testBit(cpu.AF.Hi(), 2)
		},
		0x50: func() { // BIT 2,B
			cpu.testBit(cpu.BC.Hi(), 2)
		},
		0x51: func() { // BIT 2,C
			cpu.testBit(cpu.BC.Lo(), 2)
		},
		0x52: func() { // BIT 2,D
			cpu.testBit(cpu.DE.Hi(), 2)
		},
		0x53: func() { // BIT 2,E
			cpu.testBit(cpu.DE.Lo(), 2)
		},
		0x54: func() { // BIT 2,H
			cpu.testBit(cpu.HL.Hi(), 2)
		},
		0x55: func() { // BIT 2,L
			cpu.testBit(cpu.HL.Lo(), 2)
		},
		0x56: func() { // BIT 2,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 2)
		},

		0x5F: func() { // BIT 3,A
			cpu.testBit(cpu.AF.Hi(), 3)
		},
		0x58: func() { // BIT 3,B
			cpu.testBit(cpu.BC.Hi(), 3)
		},
		0x59: func() { // BIT 3,C
			cpu.testBit(cpu.BC.Lo(), 3)
		},
		0x5A: func() { // BIT 3,D
			cpu.testBit(cpu.DE.Hi(), 3)
		},
		0x5B: func() { // BIT 3,E
			cpu.testBit(cpu.DE.Lo(), 3)
		},
		0x5C: func() { // BIT 3,H
			cpu.testBit(cpu.HL.Hi(), 3)
		},
		0x5D: func() { // BIT 3,L
			cpu.testBit(cpu.HL.Lo(), 3)
		},
		0x5E: func() { // BIT 3,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 3)
		},

		0x67: func() { // BIT 4,A
			cpu.testBit(cpu.AF.Hi(), 4)
		},
		0x60: func() { // BIT 4,B
			cpu.testBit(cpu.BC.Hi(), 4)
		},
		0x61: func() { // BIT 4,C
			cpu.testBit(cpu.BC.Lo(), 4)
		},
		0x62: func() { // BIT 4,D
			cpu.testBit(cpu.DE.Hi(), 4)
		},
		0x63: func() { // BIT 4,E
			cpu.testBit(cpu.DE.Lo(), 4)
		},
		0x64: func() { // BIT 4,H
			cpu.testBit(cpu.HL.Hi(), 4)
		},
		0x65: func() { // BIT 4,L
			cpu.testBit(cpu.HL.Lo(), 4)
		},
		0x66: func() { // BIT 4,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 4)
		},

		0x6F: func() { // BIT 5,A
			cpu.testBit(cpu.AF.Hi(), 5)
		},
		0x68: func() { // BIT 5,B
			cpu.testBit(cpu.BC.Hi(), 5)
		},
		0x69: func() { // BIT 5,C
			cpu.testBit(cpu.BC.Lo(), 5)
		},
		0x6A: func() { // BIT 5,D
			cpu.testBit(cpu.DE.Hi(), 5)
		},
		0x6B: func() { // BIT 5,E
			cpu.testBit(cpu.DE.Lo(), 5)
		},
		0x6C: func() { // BIT 5,H
			cpu.testBit(cpu.HL.Hi(), 5)
		},
		0x6D: func() { // BIT 5,L
			cpu.testBit(cpu.HL.Lo(), 5)
		},
		0x6E: func() { // BIT 5,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 5)
		},

		0x77: func() { // BIT 6,A
			cpu.testBit(cpu.AF.Hi(), 6)
		},
		0x70: func() { // BIT 6,B
			cpu.testBit(cpu.BC.Hi(), 6)
		},
		0x71: func() { // BIT 6,C
			cpu.testBit(cpu.BC.Lo(), 6)
		},
		0x72: func() { // BIT 6,D
			cpu.testBit(cpu.DE.Hi(), 6)
		},
		0x73: func() { // BIT 6,E
			cpu.testBit(cpu.DE.Lo(), 6)
		},
		0x74: func() { // BIT 6,H
			cpu.testBit(cpu.HL.Hi(), 6)
		},
		0x75: func() { // BIT 6,L
			cpu.testBit(cpu.HL.Lo(), 6)
		},
		0x76: func() { // BIT 6,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 6)
		},

		0x7F: func() { // BIT 7,A
			cpu.testBit(cpu.AF.Hi(), 7)
		},
		0x78: func() { // BIT 7,B
			cpu.testBit(cpu.BC.Hi(), 7)
		},
		0x79: func() { // BIT 7,C
			cpu.testBit(cpu.BC.Lo(), 7)
		},
		0x7A: func() { // BIT 7,D
			cpu.testBit(cpu.DE.Hi(), 7)
		},
		0x7B: func() { // BIT 7,E
			cpu.testBit(cpu.DE.Lo(), 7)
		},
		0x7C: func() { // BIT 7,H
			cpu.testBit(cpu.HL.Hi(), 7)
		},
		0x7D: func() { // BIT 7,L
			cpu.testBit(cpu.HL.Lo(), 7)
		},
		0x7E: func() { // BIT 7,(HL)
			cpu.testBit(cpu.mmu.Read(cpu.HL.HiLo()), 7)
		},

		0xC7: func() { // SET 0,A
			cpu.setBit(cpu.AF.Hi(), 0, cpu.AF.SetHi)
		},
		0xC0: func() { // SET 0,B
			cpu.setBit(cpu.BC.Hi(), 0, cpu.BC.SetHi)
		},
		0xC1: func() { // SET 0,C
			cpu.setBit(cpu.BC.Lo(), 0, cpu.BC.SetLo)
		},
		0xC2: func() { // SET 0,D
			cpu.setBit(cpu.DE.Hi(), 0, cpu.DE.SetHi)
		},
		0xC3: func() { // SET 0,E
			cpu.setBit(cpu.DE.Lo(), 0, cpu.DE.SetLo)
		},
		0xC4: func() { // SET 0,H
			cpu.setBit(cpu.HL.Hi(), 0, cpu.HL.SetHi)
		},
		0xC5: func() { // SET 0,L
			cpu.setBit(cpu.HL.Lo(), 0, cpu.HL.SetLo)
		},
		0xC6: func() { // SET 0,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 0, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xCF: func() { // SET 1,A
			cpu.setBit(cpu.AF.Hi(), 1, cpu.AF.SetHi)
		},
		0xC8: func() { // SET 1,B
			cpu.setBit(cpu.BC.Hi(), 1, cpu.BC.SetHi)
		},
		0xC9: func() { // SET 1,C
			cpu.setBit(cpu.BC.Lo(), 1, cpu.BC.SetLo)
		},
		0xCA: func() { // SET 1,D
			cpu.setBit(cpu.DE.Hi(), 1, cpu.DE.SetHi)
		},
		0xCB: func() { // SET 1,E
			cpu.setBit(cpu.DE.Lo(), 1, cpu.DE.SetLo)
		},
		0xCC: func() { // SET 1,H
			cpu.setBit(cpu.HL.Hi(), 1, cpu.HL.SetHi)
		},
		0xCD: func() { // SET 1,L
			cpu.setBit(cpu.HL.Lo(), 1, cpu.HL.SetLo)
		},
		0xCE: func() { // SET 1,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 1, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xD7: func() { // SET 2,A
			cpu.setBit(cpu.AF.Hi(), 2, cpu.AF.SetHi)
		},
		0xD0: func() { // SET 2,B
			cpu.setBit(cpu.BC.Hi(), 2, cpu.BC.SetHi)
		},
		0xD1: func() { // SET 2,C
			cpu.setBit(cpu.BC.Lo(), 2, cpu.BC.SetLo)
		},
		0xD2: func() { // SET 2,D
			cpu.setBit(cpu.DE.Hi(), 2, cpu.DE.SetHi)
		},
		0xD3: func() { // SET 2,E
			cpu.setBit(cpu.DE.Lo(), 2, cpu.DE.SetLo)
		},
		0xD4: func() { // SET 2,H
			cpu.setBit(cpu.HL.Hi(), 2, cpu.HL.SetHi)
		},
		0xD5: func() { // SET 2,L
			cpu.setBit(cpu.HL.Lo(), 2, cpu.HL.SetLo)
		},
		0xD6: func() { // SET 2,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 2, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xDF: func() { // SET 3,A
			cpu.setBit(cpu.AF.Hi(), 3, cpu.AF.SetHi)
		},
		0xD8: func() { // SET 3,B
			cpu.setBit(cpu.BC.Hi(), 3, cpu.BC.SetHi)
		},
		0xD9: func() { // SET 3,C
			cpu.setBit(cpu.BC.Lo(), 3, cpu.BC.SetLo)
		},
		0xDA: func() { // SET 3,D
			cpu.setBit(cpu.DE.Hi(), 3, cpu.DE.SetHi)
		},
		0xDB: func() { // SET 3,E
			cpu.setBit(cpu.DE.Lo(), 3, cpu.DE.SetLo)
		},
		0xDC: func() { // SET 3,H
			cpu.setBit(cpu.HL.Hi(), 3, cpu.HL.SetHi)
		},
		0xDD: func() { // SET 3,L
			cpu.setBit(cpu.HL.Lo(), 3, cpu.HL.SetLo)
		},
		0xDE: func() { // SET 3,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 3, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xE7: func() { // SET 4,A
			cpu.setBit(cpu.AF.Hi(), 4, cpu.AF.SetHi)
		},
		0xE0: func() { // SET 4,B
			cpu.setBit(cpu.BC.Hi(), 4, cpu.BC.SetHi)
		},
		0xE1: func() { // SET 4,C
			cpu.setBit(cpu.BC.Lo(), 4, cpu.BC.SetLo)
		},
		0xE2: func() { // SET 4,D
			cpu.setBit(cpu.DE.Hi(), 4, cpu.DE.SetHi)
		},
		0xE3: func() { // SET 4,E
			cpu.setBit(cpu.DE.Lo(), 4, cpu.DE.SetLo)
		},
		0xE4: func() { // SET 4,H
			cpu.setBit(cpu.HL.Hi(), 4, cpu.HL.SetHi)
		},
		0xE5: func() { // SET 4,L
			cpu.setBit(cpu.HL.Lo(), 4, cpu.HL.SetLo)
		},
		0xE6: func() { // SET 4,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 4, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xEF: func() { // SET 5,A
			cpu.setBit(cpu.AF.Hi(), 5, cpu.AF.SetHi)
		},
		0xE8: func() { // SET 5,B
			cpu.setBit(cpu.BC.Hi(), 5, cpu.BC.SetHi)
		},
		0xE9: func() { // SET 5,C
			cpu.setBit(cpu.BC.Lo(), 5, cpu.BC.SetLo)
		},
		0xEA: func() { // SET 5,D
			cpu.setBit(cpu.DE.Hi(), 5, cpu.DE.SetHi)
		},
		0xEB: func() { // SET 5,E
			cpu.setBit(cpu.DE.Lo(), 5, cpu.DE.SetLo)
		},
		0xEC: func() { // SET 5,H
			cpu.setBit(cpu.HL.Hi(), 5, cpu.HL.SetHi)
		},
		0xED: func() { // SET 5,L
			cpu.setBit(cpu.HL.Lo(), 5, cpu.HL.SetLo)
		},
		0xEE: func() { // SET 5,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 5, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xF7: func() { // SET 6,A
			cpu.setBit(cpu.AF.Hi(), 6, cpu.AF.SetHi)
		},
		0xF0: func() { // SET 6,B
			cpu.setBit(cpu.BC.Hi(), 6, cpu.BC.SetHi)
		},
		0xF1: func() { // SET 6,C
			cpu.setBit(cpu.BC.Lo(), 6, cpu.BC.SetLo)
		},
		0xF2: func() { // SET 6,D
			cpu.setBit(cpu.DE.Hi(), 6, cpu.DE.SetHi)
		},
		0xF3: func() { // SET 6,E
			cpu.setBit(cpu.DE.Lo(), 6, cpu.DE.SetLo)
		},
		0xF4: func() { // SET 6,H
			cpu.setBit(cpu.HL.Hi(), 6, cpu.HL.SetHi)
		},
		0xF5: func() { // SET 6,L
			cpu.setBit(cpu.HL.Lo(), 6, cpu.HL.SetLo)
		},
		0xF6: func() { // SET 6,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 6, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xFF: func() { // SET 7,A
			cpu.setBit(cpu.AF.Hi(), 7, cpu.AF.SetHi)
		},
		0xF8: func() { // SET 7,B
			cpu.setBit(cpu.BC.Hi(), 7, cpu.BC.SetHi)
		},
		0xF9: func() { // SET 7,C
			cpu.setBit(cpu.BC.Lo(), 7, cpu.BC.SetLo)
		},
		0xFA: func() { // SET 7,D
			cpu.setBit(cpu.DE.Hi(), 7, cpu.DE.SetHi)
		},
		0xFB: func() { // SET 7,E
			cpu.setBit(cpu.DE.Lo(), 7, cpu.DE.SetLo)
		},
		0xFC: func() { // SET 7,H
			cpu.setBit(cpu.HL.Hi(), 7, cpu.HL.SetHi)
		},
		0xFD: func() { // SET 7,L
			cpu.setBit(cpu.HL.Lo(), 7, cpu.HL.SetLo)
		},
		0xFE: func() { // SET 7,(HL)
			addr := cpu.HL.HiLo()
			cpu.setBit(cpu.mmu.Read(addr), 7, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x87: func() { // RES 0,A
			cpu.resetBit(cpu.AF.Hi(), 0, cpu.AF.SetHi)
		},
		0x80: func() { // RES 0,B
			cpu.resetBit(cpu.BC.Hi(), 0, cpu.BC.SetHi)
		},
		0x81: func() { // RES 0,C
			cpu.resetBit(cpu.BC.Lo(), 0, cpu.BC.SetLo)
		},
		0x82: func() { // RES 0,D
			cpu.resetBit(cpu.DE.Hi(), 0, cpu.DE.SetHi)
		},
		0x83: func() { // RES 0,E
			cpu.resetBit(cpu.DE.Lo(), 0, cpu.DE.SetLo)
		},
		0x84: func() { // RES 0,H
			cpu.resetBit(cpu.HL.Hi(), 0, cpu.HL.SetHi)
		},
		0x85: func() { // RES 0,L
			cpu.resetBit(cpu.HL.Lo(), 0, cpu.HL.SetLo)
		},
		0x86: func() { // RES 0,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 0, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x8F: func() { // RES 1,A
			cpu.resetBit(cpu.AF.Hi(), 1, cpu.AF.SetHi)
		},
		0x88: func() { // RES 1,B
			cpu.resetBit(cpu.BC.Hi(), 1, cpu.BC.SetHi)
		},
		0x89: func() { // RES 1,C
			cpu.resetBit(cpu.BC.Lo(), 1, cpu.BC.SetLo)
		},
		0x8A: func() { // RES 1,D
			cpu.resetBit(cpu.DE.Hi(), 1, cpu.DE.SetHi)
		},
		0x8B: func() { // RES 1,E
			cpu.resetBit(cpu.DE.Lo(), 1, cpu.DE.SetLo)
		},
		0x8C: func() { // RES 1,H
			cpu.resetBit(cpu.HL.Hi(), 1, cpu.HL.SetHi)
		},
		0x8D: func() { // RES 1,L
			cpu.resetBit(cpu.HL.Lo(), 1, cpu.HL.SetLo)
		},
		0x8E: func() { // RES 1,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 1, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x97: func() { // RES 2,A
			cpu.resetBit(cpu.AF.Hi(), 2, cpu.AF.SetHi)
		},
		0x90: func() { // RES 2,B
			cpu.resetBit(cpu.BC.Hi(), 2, cpu.BC.SetHi)
		},
		0x91: func() { // RES 2,C
			cpu.resetBit(cpu.BC.Lo(), 2, cpu.BC.SetLo)
		},
		0x92: func() { // RES 2,D
			cpu.resetBit(cpu.DE.Hi(), 2, cpu.DE.SetHi)
		},
		0x93: func() { // RES 2,E
			cpu.resetBit(cpu.DE.Lo(), 2, cpu.DE.SetLo)
		},
		0x94: func() { // RES 2,H
			cpu.resetBit(cpu.HL.Hi(), 2, cpu.HL.SetHi)
		},
		0x95: func() { // RES 2,L
			cpu.resetBit(cpu.HL.Lo(), 2, cpu.HL.SetLo)
		},
		0x96: func() { // RES 2,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 2, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x9F: func() { // RES 3,A
			cpu.resetBit(cpu.AF.Hi(), 3, cpu.AF.SetHi)
		},
		0x98: func() { // RES 3,B
			cpu.resetBit(cpu.BC.Hi(), 3, cpu.BC.SetHi)
		},
		0x99: func() { // RES 3,C
			cpu.resetBit(cpu.BC.Lo(), 3, cpu.BC.SetLo)
		},
		0x9A: func() { // RES 3,D
			cpu.resetBit(cpu.DE.Hi(), 3, cpu.DE.SetHi)
		},
		0x9B: func() { // RES 3,E
			cpu.resetBit(cpu.DE.Lo(), 3, cpu.DE.SetLo)
		},
		0x9C: func() { // RES 3,H
			cpu.resetBit(cpu.HL.Hi(), 3, cpu.HL.SetHi)
		},
		0x9D: func() { // RES 3,L
			cpu.resetBit(cpu.HL.Lo(), 3, cpu.HL.SetLo)
		},
		0x9E: func() { // RES 3,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 3, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xA7: func() { // RES 4,A
			cpu.resetBit(cpu.AF.Hi(), 4, cpu.AF.SetHi)
		},
		0xA0: func() { // RES 4,B
			cpu.resetBit(cpu.BC.Hi(), 4, cpu.BC.SetHi)
		},
		0xA1: func() { // RES 4,C
			cpu.resetBit(cpu.BC.Lo(), 4, cpu.BC.SetLo)
		},
		0xA2: func() { // RES 4,D
			cpu.resetBit(cpu.DE.Hi(), 4, cpu.DE.SetHi)
		},
		0xA3: func() { // RES 4,E
			cpu.resetBit(cpu.DE.Lo(), 4, cpu.DE.SetLo)
		},
		0xA4: func() { // RES 4,H
			cpu.resetBit(cpu.HL.Hi(), 4, cpu.HL.SetHi)
		},
		0xA5: func() { // RES 4,L
			cpu.resetBit(cpu.HL.Lo(), 4, cpu.HL.SetLo)
		},
		0xA6: func() { // RES 4,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 4, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xAF: func() { // RES 5,A
			cpu.resetBit(cpu.AF.Hi(), 5, cpu.AF.SetHi)
		},
		0xA8: func() { // RES 5,B
			cpu.resetBit(cpu.BC.Hi(), 5, cpu.BC.SetHi)
		},
		0xA9: func() { // RES 5,C
			cpu.resetBit(cpu.BC.Lo(), 5, cpu.BC.SetLo)
		},
		0xAA: func() { // RES 5,D
			cpu.resetBit(cpu.DE.Hi(), 5, cpu.DE.SetHi)
		},
		0xAB: func() { // RES 5,E
			cpu.resetBit(cpu.DE.Lo(), 5, cpu.DE.SetLo)
		},
		0xAC: func() { // RES 5,H
			cpu.resetBit(cpu.HL.Hi(), 5, cpu.HL.SetHi)
		},
		0xAD: func() { // RES 5,L
			cpu.resetBit(cpu.HL.Lo(), 5, cpu.HL.SetLo)
		},
		0xAE: func() { // RES 5,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 5, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xB7: func() { // RES 6,A
			cpu.resetBit(cpu.AF.Hi(), 6, cpu.AF.SetHi)
		},
		0xB0: func() { // RES 6,B
			cpu.resetBit(cpu.BC.Hi(), 6, cpu.BC.SetHi)
		},
		0xB1: func() { // RES 6,C
			cpu.resetBit(cpu.BC.Lo(), 6, cpu.BC.SetLo)
		},
		0xB2: func() { // RES 6,D
			cpu.resetBit(cpu.DE.Hi(), 6, cpu.DE.SetHi)
		},
		0xB3: func() { // RES 6,E
			cpu.resetBit(cpu.DE.Lo(), 6, cpu.DE.SetLo)
		},
		0xB4: func() { // RES 6,H
			cpu.resetBit(cpu.HL.Hi(), 6, cpu.HL.SetHi)
		},
		0xB5: func() { // RES 6,L
			cpu.resetBit(cpu.HL.Lo(), 6, cpu.HL.SetLo)
		},
		0xB6: func() { // RES 6,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 6, func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0xBF: func() { // RES 7,A
			cpu.resetBit(cpu.AF.Hi(), 7, cpu.AF.SetHi)
		},
		0xB8: func() { // RES 7,B
			cpu.resetBit(cpu.BC.Hi(), 7, cpu.BC.SetHi)
		},
		0xB9: func() { // RES 7,C
			cpu.resetBit(cpu.BC.Lo(), 7, cpu.BC.SetLo)
		},
		0xBA: func() { // RES 7,D
			cpu.resetBit(cpu.DE.Hi(), 7, cpu.DE.SetHi)
		},
		0xBB: func() { // RES 7,E
			cpu.resetBit(cpu.DE.Lo(), 7, cpu.DE.SetLo)
		},
		0xBC: func() { // RES 7,H
			cpu.resetBit(cpu.HL.Hi(), 7, cpu.HL.SetHi)
		},
		0xBD: func() { // RES 7,L
			cpu.resetBit(cpu.HL.Lo(), 7, cpu.HL.SetLo)
		},
		0xBE: func() { // RES 7,(HL)
			addr := cpu.HL.HiLo()
			cpu.resetBit(cpu.mmu.Read(addr), 7, func(val byte) { cpu.mmu.Write(addr, val) })
		},
	}

	for k, v := range cpu.instructions {
		if v == nil {
			opcode := k
			cpu.instructions[k] = func() {
				log.Printf("Encountered unknown CB instruction: %#2x", opcode)
				cpu.stop = true
			}
		}
	}

	cpu.cycles = [0x100]int{
		1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
		0, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1,
		2, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
		2, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		2, 2, 2, 2, 2, 2, 0, 2, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
		2, 3, 3, 4, 3, 4, 2, 4, 2, 4, 3, 0, 3, 6, 2, 4,
		2, 3, 3, 0, 3, 4, 2, 4, 2, 4, 3, 0, 3, 0, 2, 4,
		3, 3, 2, 0, 0, 4, 2, 4, 4, 1, 4, 0, 0, 0, 2, 4,
		3, 3, 2, 1, 0, 4, 2, 4, 3, 2, 4, 1, 0, 0, 2, 4,
	}

	cpu.cyclesCB = [0x100]int{
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
		2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	}
}
