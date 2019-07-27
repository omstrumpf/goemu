package gba

import "log"

// PopulateInstructions populates the opcode map for instructions and cycle costs
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
			cpu.HL.SetHi(cpu.mmu.Read(cpu.HL.HiLo()))
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
			cpu.mmu.Write16(cpu.SP.Dec2(), cpu.BC.HiLo())
		},
		0xD5: func() { // PUSH DE
			cpu.mmu.Write16(cpu.SP.Dec2(), cpu.DE.HiLo())
		},
		0xE5: func() { // PUSH HL
			cpu.mmu.Write16(cpu.SP.Dec2(), cpu.HL.HiLo())
		},
		0xF5: func() { // PUSH AF
			cpu.mmu.Write16(cpu.SP.Dec2(), cpu.AF.HiLo())
		},

		0xC1: func() { // POP BC
			cpu.BC.Set(cpu.mmu.Read16(cpu.SP.Inc2()))
		},
		0xD1: func() { // POP DE
			cpu.DE.Set(cpu.mmu.Read16(cpu.SP.Inc2()))
		},
		0xE1: func() { // POP HL
			cpu.HL.Set(cpu.mmu.Read16(cpu.SP.Inc2()))
		},
		0xF1: func() { // POP AF
			cpu.AF.Set(cpu.mmu.Read16(cpu.SP.Inc2()))
		},

		//// 8-bit ALU ////

		//// 16-bit ALU ////

		//// Rotate / Shift ////

		//// Singlebit ////

		//// Control ////
		0x3F: func() { // CCF
			f := cpu.AF.Lo()
			if cpu.c() {
				f &= 0x8F // 0b10001111, set n, h, and c to 0
			} else {
				f &= 0x9F // 0b1001111, set n, h to 0
				f |= 0x10 // 0b00010000, set c to 1
			}
			cpu.AF.SetLo(f)
		},
		0x37: func() { // SCF
			f := cpu.AF.Lo()
			f &= 0x9F                       // 0b10011111, set n, h to 0
			f |= 0x10                       // 0b00010000, set c to 1
			cpu.AF.SetLo(cpu.AF.Lo() & 0x9) // 0b1001
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
			cpu.ime = true
		},
		0xFB: func() { // EI
			cpu.ime = false
		},

		//// Jump /////
	}

	for k, v := range cpu.instructions {
		if v == nil {
			opcode := k
			cpu.instructions[k] = func() {
				log.Printf("Encountered unknown instruction: %#2x", opcode)
			}
		}
	}

	cpu.cycles = [0x100]int{
		//// 8-bit loads ////
		0x7F: 4,
		0x78: 4,
		0x79: 4,
		0x7A: 4,
		0x7B: 4,
		0x7C: 4,
		0x7D: 4,
		0x47: 4,
		0x40: 4,
		0x41: 4,
		0x42: 4,
		0x43: 4,
		0x44: 4,
		0x45: 4,
		0x4F: 4,
		0x48: 4,
		0x49: 4,
		0x4A: 4,
		0x4B: 4,
		0x4C: 4,
		0x4D: 4,
		0x57: 4,
		0x50: 4,
		0x51: 4,
		0x52: 4,
		0x53: 4,
		0x54: 4,
		0x55: 4,
		0x5F: 4,
		0x58: 4,
		0x59: 4,
		0x5A: 4,
		0x5B: 4,
		0x5C: 4,
		0x5D: 4,
		0x67: 4,
		0x60: 4,
		0x61: 4,
		0x62: 4,
		0x63: 4,
		0x64: 4,
		0x65: 4,
		0x6F: 4,
		0x68: 4,
		0x69: 4,
		0x6A: 4,
		0x6B: 4,
		0x6C: 4,
		0x6D: 4,

		0x3E: 8,
		0x06: 8,
		0x0E: 8,
		0x16: 8,
		0x1E: 8,
		0x26: 8,
		0x2E: 8,

		0x7E: 8,
		0x46: 8,
		0x4E: 8,
		0x56: 8,
		0x5E: 8,
		0x66: 8,
		0x6E: 8,

		0x77: 8,
		0x70: 8,
		0x71: 8,
		0x72: 8,
		0x73: 8,
		0x74: 8,
		0x75: 8,

		0x36: 12,

		0x0A: 8,
		0x1A: 8,
		0xFA: 16,

		0x02: 8,
		0x12: 8,
		0xEA: 16,

		0xF2: 8,
		0xE2: 8,
		0xF0: 12,
		0xE0: 12,

		0x22: 8,
		0x2A: 8,
		0x32: 8,
		0x3A: 8,

		//// 16-bit loads ////
		0x01: 12,
		0x11: 12,
		0x21: 12,
		0x31: 12,

		0xF9: 8,

		0xC5: 16,
		0xD5: 16,
		0xE5: 16,
		0xF5: 16,

		0xC1: 12,
		0xD1: 12,
		0xE1: 12,
		0xF1: 12,

		//// 8-bit ALU ////

		//// 16-bit ALU ////

		//// Rotate / Shift ////

		//// Singlebit ////

		//// Control ////
		0x3F: 4,
		0x37: 4,
		0x00: 4,
		0x76: 4, // N * 4
		0x10: 4, // ?
		0xF3: 4,
		0xFB: 4,

		//// Jump /////
	}
}
