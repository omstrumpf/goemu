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
		0x87: func() { // ADD A,A
			cpu.addInstruction(cpu.AF.Hi(), false)
		},
		0x80: func() { // ADD A,B
			cpu.addInstruction(cpu.BC.Hi(), false)
		},
		0x81: func() { // ADD A,C
			cpu.addInstruction(cpu.BC.Lo(), false)
		},
		0x82: func() { // ADD A,D
			cpu.addInstruction(cpu.DE.Hi(), false)
		},
		0x83: func() { // ADD A,E
			cpu.addInstruction(cpu.DE.Lo(), false)
		},
		0x84: func() { // ADD A,H
			cpu.addInstruction(cpu.HL.Hi(), false)
		},
		0x85: func() { // ADD A,L
			cpu.addInstruction(cpu.HL.Lo(), false)
		},
		0xC6: func() { // ADD A,n
			cpu.addInstruction(cpu.mmu.Read(cpu.PC.Inc()), false)
		},
		0x86: func() { // ADD A,(HL)
			cpu.addInstruction(cpu.mmu.Read(cpu.HL.HiLo()), false)
		},

		0x8F: func() { // ADC A,A
			cpu.addInstruction(cpu.AF.Hi(), true)
		},
		0x88: func() { // ADC A,B
			cpu.addInstruction(cpu.BC.Hi(), true)
		},
		0x89: func() { // ADC A,C
			cpu.addInstruction(cpu.BC.Lo(), true)
		},
		0x8A: func() { // ADC A,D
			cpu.addInstruction(cpu.DE.Hi(), true)
		},
		0x8B: func() { // ADC A,E
			cpu.addInstruction(cpu.DE.Lo(), true)
		},
		0x8C: func() { // ADC A,H
			cpu.addInstruction(cpu.HL.Hi(), true)
		},
		0x8D: func() { // ADC A,L
			cpu.addInstruction(cpu.HL.Lo(), true)
		},
		0xCE: func() { // ADC A,n
			cpu.addInstruction(cpu.mmu.Read(cpu.PC.Inc()), true)
		},
		0x8E: func() { // ADC A,(HL)
			cpu.addInstruction(cpu.mmu.Read(cpu.HL.HiLo()), true)
		},

		0x97: func() { // SUB A,A
			cpu.subInstruction(cpu.AF.Hi(), false)
		},
		0x90: func() { // SUB A,B
			cpu.subInstruction(cpu.BC.Hi(), false)
		},
		0x91: func() { // SUB A,C
			cpu.subInstruction(cpu.BC.Lo(), false)
		},
		0x92: func() { // SUB A,D
			cpu.subInstruction(cpu.DE.Hi(), false)
		},
		0x93: func() { // SUB A,E
			cpu.subInstruction(cpu.DE.Lo(), false)
		},
		0x94: func() { // SUB A,H
			cpu.subInstruction(cpu.HL.Hi(), false)
		},
		0x95: func() { // SUB A,L
			cpu.subInstruction(cpu.HL.Lo(), false)
		},
		0xD6: func() { // SUB A,n
			cpu.subInstruction(cpu.mmu.Read(cpu.PC.Inc()), false)
		},
		0x96: func() { // SUB A,(HL)
			cpu.subInstruction(cpu.mmu.Read(cpu.HL.HiLo()), false)
		},

		0x9F: func() { // SBC A,A
			cpu.subInstruction(cpu.AF.Hi(), true)
		},
		0x98: func() { // SBC A,B
			cpu.subInstruction(cpu.BC.Hi(), true)
		},
		0x99: func() { // SBC A,C
			cpu.subInstruction(cpu.BC.Lo(), true)
		},
		0x9A: func() { // SBC A,D
			cpu.subInstruction(cpu.DE.Hi(), true)
		},
		0x9B: func() { // SBC A,E
			cpu.subInstruction(cpu.DE.Lo(), true)
		},
		0x9C: func() { // SBC A,H
			cpu.subInstruction(cpu.HL.Hi(), true)
		},
		0x9D: func() { // SBC A,L
			cpu.subInstruction(cpu.HL.Lo(), true)
		},
		0xDE: func() { // SBC A,n
			cpu.subInstruction(cpu.mmu.Read(cpu.PC.Inc()), true)
		},
		0x9E: func() { // SBC A,(HL)
			cpu.subInstruction(cpu.mmu.Read(cpu.HL.HiLo()), true)
		},

		0xA7: func() { // AND A
			cpu.andInstruction(cpu.AF.Hi())
		},
		0xA0: func() { // AND B
			cpu.andInstruction(cpu.BC.Hi())
		},
		0xA1: func() { // AND C
			cpu.andInstruction(cpu.BC.Lo())
		},
		0xA2: func() { // AND D
			cpu.andInstruction(cpu.DE.Hi())
		},
		0xA3: func() { // AND E
			cpu.andInstruction(cpu.DE.Lo())
		},
		0xA4: func() { // AND H
			cpu.andInstruction(cpu.HL.Hi())
		},
		0xA5: func() { // AND L
			cpu.andInstruction(cpu.HL.Lo())
		},
		0xE6: func() { // AND n
			cpu.andInstruction(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xA6: func() { // AND (HL)
			cpu.andInstruction(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xAF: func() { // XOR A
			cpu.xorInstruction(cpu.AF.Hi())
		},
		0xA8: func() { // XOR B
			cpu.xorInstruction(cpu.BC.Hi())
		},
		0xA9: func() { // XOR C
			cpu.xorInstruction(cpu.BC.Lo())
		},
		0xAA: func() { // XOR D
			cpu.xorInstruction(cpu.DE.Hi())
		},
		0xAB: func() { // XOR E
			cpu.xorInstruction(cpu.DE.Lo())
		},
		0xAC: func() { // XOR H
			cpu.xorInstruction(cpu.HL.Hi())
		},
		0xAD: func() { // XOR L
			cpu.xorInstruction(cpu.HL.Lo())
		},
		0xEE: func() { // XOR n
			cpu.xorInstruction(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xAE: func() { // XOR (HL)
			cpu.xorInstruction(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xB7: func() { // OR A
			cpu.orInstruction(cpu.AF.Hi())
		},
		0xB0: func() { // OR B
			cpu.orInstruction(cpu.BC.Hi())
		},
		0xB1: func() { // OR C
			cpu.orInstruction(cpu.BC.Lo())
		},
		0xB2: func() { // OR D
			cpu.orInstruction(cpu.DE.Hi())
		},
		0xB3: func() { // OR E
			cpu.orInstruction(cpu.DE.Lo())
		},
		0xB4: func() { // OR H
			cpu.orInstruction(cpu.HL.Hi())
		},
		0xB5: func() { // OR L
			cpu.orInstruction(cpu.HL.Lo())
		},
		0xF6: func() { // OR n
			cpu.orInstruction(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xB6: func() { // OR (HL)
			cpu.orInstruction(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0xBF: func() { // CP A
			cpu.cpInstruction(cpu.AF.Hi())
		},
		0xB8: func() { // CP B
			cpu.cpInstruction(cpu.BC.Hi())
		},
		0xB9: func() { // CP C
			cpu.cpInstruction(cpu.BC.Lo())
		},
		0xBA: func() { // CP D
			cpu.cpInstruction(cpu.DE.Hi())
		},
		0xBB: func() { // CP E
			cpu.cpInstruction(cpu.DE.Lo())
		},
		0xBC: func() { // CP H
			cpu.cpInstruction(cpu.HL.Hi())
		},
		0xBD: func() { // CP L
			cpu.cpInstruction(cpu.HL.Lo())
		},
		0xFE: func() { // CP n
			cpu.cpInstruction(cpu.mmu.Read(cpu.PC.Inc()))
		},
		0xBE: func() { // CP (HL)
			cpu.cpInstruction(cpu.mmu.Read(cpu.HL.HiLo()))
		},

		0x3C: func() { // INC A
			cpu.incInstruction(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x04: func() { // INC B
			cpu.incInstruction(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x0C: func() { // INC C
			cpu.incInstruction(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x14: func() { // INC D
			cpu.incInstruction(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x1C: func() { // INC E
			cpu.incInstruction(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x24: func() { // INC H
			cpu.incInstruction(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x2C: func() { // INC L
			cpu.incInstruction(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x34: func() { // INC (HL)
			addr := cpu.HL.HiLo()
			cpu.incInstruction(cpu.mmu.Read(cpu.HL.HiLo()), func(val byte) { cpu.mmu.Write(addr, val) })
		},

		0x3D: func() { // DEC A
			cpu.decInstruction(cpu.AF.Hi(), cpu.AF.SetHi)
		},
		0x05: func() { // DEC B
			cpu.decInstruction(cpu.BC.Hi(), cpu.BC.SetHi)
		},
		0x0D: func() { // DEC C
			cpu.decInstruction(cpu.BC.Lo(), cpu.BC.SetLo)
		},
		0x15: func() { // DEC D
			cpu.decInstruction(cpu.DE.Hi(), cpu.DE.SetHi)
		},
		0x1D: func() { // DEC E
			cpu.decInstruction(cpu.DE.Lo(), cpu.DE.SetLo)
		},
		0x25: func() { // DEC H
			cpu.decInstruction(cpu.HL.Hi(), cpu.HL.SetHi)
		},
		0x2D: func() { // DEC L
			cpu.decInstruction(cpu.HL.Lo(), cpu.HL.SetLo)
		},
		0x35: func() { // DEC (HL)
			addr := cpu.HL.HiLo()
			cpu.decInstruction(cpu.mmu.Read(cpu.HL.HiLo()), func(val byte) { cpu.mmu.Write(addr, val) })
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
			cpu.add16Instruction(cpu.BC.HiLo())
		},
		0x19: func() { // ADD HL,DE
			cpu.add16Instruction(cpu.DE.HiLo())
		},
		0x29: func() { // ADD HL,HL
			cpu.add16Instruction(cpu.HL.HiLo())
		},
		0x39: func() { // ADD HL,SP
			cpu.add16Instruction(cpu.SP.HiLo())
		},

		0x03: func() { // INC BC
			cpu.BC.Set(cpu.BC.HiLo() + 1)
		},
		0x13: func() { // INC DE
			cpu.BC.Set(cpu.DE.HiLo() + 1)
		},
		0x23: func() { // INC HL
			cpu.BC.Set(cpu.HL.HiLo() + 1)
		},
		0x33: func() { // INC SP
			cpu.BC.Set(cpu.SP.HiLo() + 1)
		},

		0x0B: func() { // DEC BC
			cpu.BC.Set(cpu.BC.HiLo() - 1)
		},
		0x1B: func() { // DEC DE
			cpu.BC.Set(cpu.DE.HiLo() - 1)
		},
		0x2B: func() { // DEC HL
			cpu.BC.Set(cpu.HL.HiLo() - 1)
		},
		0x3B: func() { // DEC SP
			cpu.BC.Set(cpu.SP.HiLo() - 1)
		},

		0xE8: func() { // ADD SP,d
			cpu.add16SignedInstruction(cpu.SP.HiLo(), int8(cpu.PC.Inc()), cpu.SP.Set)
		},
		0xF8: func() { // LD HL,SP,d
			cpu.add16SignedInstruction(cpu.SP.HiLo(), int8(cpu.PC.Inc()), cpu.HL.Set)
		},

		//// Rotate / Shift ////
		0x07: func() { // RLCA
			original := cpu.AF.Hi()
			result := byte(original<<1) | (original >> 7)

			cpu.AF.SetHi(result)

			cpu.setZ(false)
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(original > 0x7F)
		},
		0x17: func() { // RLA
			original := cpu.AF.Hi()
			carry := 0
			if cpu.c() {
				carry = 1
			}
			result := byte(original<<1) + byte(carry)

			cpu.AF.SetHi(result)

			cpu.setZ(false)
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(original > 0x7F)
		},
		0x0F: func() { // RRCA
			original := cpu.AF.Hi()
			result := byte(original>>1) | (original << 7)

			cpu.AF.SetHi(result)

			cpu.setZ(false)
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(original > 0x7F)
		},
		0x1F: func() { // RRA
			original := cpu.AF.Hi()
			carry := 0
			if cpu.c() {
				carry = 1 << 7
			}
			result := byte(original>>1) | byte(carry)

			cpu.AF.SetHi(result)

			cpu.setZ(false)
			cpu.setN(false)
			cpu.setH(false)
			cpu.setC(1&original == 1)
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
				// TODO pause execution (notify user)
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
}

func (cpu *CPU) addInstruction(operand byte, useCarry bool) {
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

func (cpu *CPU) subInstruction(operand byte, useCarry bool) {
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

func (cpu *CPU) andInstruction(operand byte) {
	result := cpu.AF.Hi() & operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(true)
	cpu.setC(false)
}

func (cpu *CPU) xorInstruction(operand byte) {
	result := cpu.AF.Hi() ^ operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

func (cpu *CPU) orInstruction(operand byte) {
	result := cpu.AF.Hi() | operand

	cpu.AF.SetHi(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

func (cpu *CPU) cpInstruction(operand byte) {
	original := cpu.AF.Hi()
	result := original - operand

	cpu.setZ(result == 0)
	cpu.setN(true)
	cpu.setH((original & 0x0F) > (operand & 0x0F))
	cpu.setC(original > operand)
}

func (cpu *CPU) incInstruction(val byte, setter func(byte)) {
	result := val + 1

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(false)
	cpu.setH(((val & 0xF) + 1) > 0xF)
}

func (cpu *CPU) decInstruction(val byte, setter func(byte)) {
	result := val - 1

	setter(result)

	cpu.setZ(result == 0)
	cpu.setN(true)
	cpu.setH(val&0xF == 0)
}

func (cpu *CPU) add16Instruction(val uint16) {
	result := uint32(cpu.HL.HiLo()) + uint32(val)

	cpu.HL.Set(uint16(result))
	cpu.setN(false)
	cpu.setH(uint32(val&0xFFF) > (result & 0xFFF))
	cpu.setC(result > 0xFFFF)
}

func (cpu *CPU) add16SignedInstruction(original uint16, operand int8, setter func(uint16)) {
	result := int32(original) + int32(operand)

	setter(uint16(result))

	// xor operands and result to determine carries
	xor := original ^ uint16(operand) ^ uint16(result)

	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH((xor & 0x10) == 0x10)
	cpu.setC((xor & 0x100) == 0x100)
}
