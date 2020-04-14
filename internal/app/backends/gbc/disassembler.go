package gbc

import "fmt"

// Disassemble returns a disassembled string and the next unparsed PC
func (cpu *CPU) Disassemble(pc uint16) (uint16, string) {
	opcode := cpu.mmu.Read(pc)

	switch opcode {
	case 0x7F:
		return pc + 1, "LD A,A"
	case 0x78:
		return pc + 1, "LD A,B"
	case 0x79:
		return pc + 1, "LD A,C"
	case 0x7A:
		return pc + 1, "LD A,D"
	case 0x7B:
		return pc + 1, "LD A,E"
	case 0x7C:
		return pc + 1, "LD A,H"
	case 0x7D:
		return pc + 1, "LD A,L"
	case 0x47:
		return pc + 1, "LD B,A"
	case 0x40:
		return pc + 1, "LD B,B"
	case 0x41:
		return pc + 1, "LD B,C"
	case 0x42:
		return pc + 1, "LD B,D"
	case 0x43:
		return pc + 1, "LD B,E"
	case 0x44:
		return pc + 1, "LD B,H"
	case 0x45:
		return pc + 1, "LD B,L"
	case 0x4F:
		return pc + 1, "LD C,A"
	case 0x48:
		return pc + 1, "LD C,B"
	case 0x49:
		return pc + 1, "LD C,C"
	case 0x4A:
		return pc + 1, "LD C,D"
	case 0x4B:
		return pc + 1, "LD C,E"
	case 0x4C:
		return pc + 1, "LD C,H"
	case 0x4D:
		return pc + 1, "LD C,L"
	case 0x57:
		return pc + 1, "LD D,A"
	case 0x50:
		return pc + 1, "LD D,B"
	case 0x51:
		return pc + 1, "LD D,C"
	case 0x52:
		return pc + 1, "LD D,D"
	case 0x53:
		return pc + 1, "LD D,E"
	case 0x54:
		return pc + 1, "LD D,H"
	case 0x55:
		return pc + 1, "LD D,L"
	case 0x5F:
		return pc + 1, "LD E,A"
	case 0x58:
		return pc + 1, "LD E,B"
	case 0x59:
		return pc + 1, "LD E,C"
	case 0x5A:
		return pc + 1, "LD E,D"
	case 0x5B:
		return pc + 1, "LD E,E"
	case 0x5C:
		return pc + 1, "LD E,H"
	case 0x5D:
		return pc + 1, "LD E,L"
	case 0x67:
		return pc + 1, "LD H,A"
	case 0x60:
		return pc + 1, "LD H,B"
	case 0x61:
		return pc + 1, "LD H,C"
	case 0x62:
		return pc + 1, "LD H,D"
	case 0x63:
		return pc + 1, "LD H,E"
	case 0x64:
		return pc + 1, "LD H,H"
	case 0x65:
		return pc + 1, "LD H,L"
	case 0x6F:
		return pc + 1, "LD L,A"
	case 0x68:
		return pc + 1, "LD L,B"
	case 0x69:
		return pc + 1, "LD L,C"
	case 0x6A:
		return pc + 1, "LD L,D"
	case 0x6B:
		return pc + 1, "LD L,E"
	case 0x6C:
		return pc + 1, "LD L,H"
	case 0x6D:
		return pc + 1, "LD L,L"
	case 0x3E:
		return pc + 2, fmt.Sprintf("LD A,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x06:
		return pc + 2, fmt.Sprintf("LD B,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x0E:
		return pc + 2, fmt.Sprintf("LD C,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x16:
		return pc + 2, fmt.Sprintf("LD D,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x1E:
		return pc + 2, fmt.Sprintf("LD E,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x26:
		return pc + 2, fmt.Sprintf("LD H,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x2E:
		return pc + 2, fmt.Sprintf("LD L,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x7E:
		return pc + 1, "LD A,(HL)"
	case 0x46:
		return pc + 1, "LD B,(HL)"
	case 0x4E:
		return pc + 1, "LD C,(HL)"
	case 0x56:
		return pc + 1, "LD D,(HL)"
	case 0x5E:
		return pc + 1, "LD E,(HL)"
	case 0x66:
		return pc + 1, "LD H,(HL)"
	case 0x6E:
		return pc + 1, "LD L,(HL)"
	case 0x77:
		return pc + 1, "LD (HL),A"
	case 0x70:
		return pc + 1, "LD (HL),B"
	case 0x71:
		return pc + 1, "LD (HL),C"
	case 0x72:
		return pc + 1, "LD (HL),D"
	case 0x73:
		return pc + 1, "LD (HL),E"
	case 0x74:
		return pc + 1, "LD (HL),H"
	case 0x75:
		return pc + 1, "LD (HL),L"
	case 0x36:
		return pc + 2, fmt.Sprintf("LD (HL),[%#02x]", cpu.mmu.Read(pc+1))
	case 0x0A:
		return pc + 1, "LD A,(BC)"
	case 0x1A:
		return pc + 1, "LD A,(DE)"
	case 0xFA:
		return pc + 3, fmt.Sprintf("LD A,(%#04x)", cpu.mmu.Read16(pc+1))
	case 0x02:
		return pc + 1, "LD (BC),A"
	case 0x12:
		return pc + 1, "LD (DE),A"
	case 0xEA:
		return pc + 3, fmt.Sprintf("LD (%#04x),A", cpu.mmu.Read16(pc+1))
	case 0x08:
		return pc + 3, fmt.Sprintf("LD (%#04x),SP", cpu.mmu.Read16(pc+1))
	case 0xF2:
		return pc + 1, "LD A,(FF00+C)"
	case 0xE2:
		return pc + 1, "LD (FF00+C),A"
	case 0xF0:
		return pc + 2, fmt.Sprintf("LD A,(FF00+[%#02x])", cpu.mmu.Read(pc+1))
	case 0xE0:
		return pc + 2, fmt.Sprintf("LD (FF00+[%#02x]),A", cpu.mmu.Read(pc+1))
	case 0x22:
		return pc + 1, "LDI (HL),A"
	case 0x2A:
		return pc + 1, "LDI A,(HL)"
	case 0x32:
		return pc + 1, "LDD (HL),A"
	case 0x3A:
		return pc + 1, "LDD A,(HL)"
	case 0x01:
		return pc + 3, fmt.Sprintf("LD BC,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0x11:
		return pc + 3, fmt.Sprintf("LD DE,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0x21:
		return pc + 3, fmt.Sprintf("LD HL,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0x31:
		return pc + 3, fmt.Sprintf("LD SP,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xF9:
		return pc + 1, "LD SP,HL"
	case 0xC5:
		return pc + 1, "PUSH BC"
	case 0xD5:
		return pc + 1, "PUSH DE"
	case 0xE5:
		return pc + 1, "PUSH HL"
	case 0xF5:
		return pc + 1, "PUSH AF"
	case 0xC1:
		return pc + 1, "POP BC"
	case 0xD1:
		return pc + 1, "POP DE"
	case 0xE1:
		return pc + 1, "POP HL"
	case 0xF1:
		return pc + 1, "POP AF"
	case 0x87:
		return pc + 1, "ADD A,A"
	case 0x80:
		return pc + 1, "ADD A,B"
	case 0x81:
		return pc + 1, "ADD A,C"
	case 0x82:
		return pc + 1, "ADD A,D"
	case 0x83:
		return pc + 1, "ADD A,E"
	case 0x84:
		return pc + 1, "ADD A,H"
	case 0x85:
		return pc + 1, "ADD A,L"
	case 0xC6:
		return pc + 2, fmt.Sprintf("ADD A,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x86:
		return pc + 1, "ADD A,(HL)"
	case 0x8F:
		return pc + 1, "ADC A,A"
	case 0x88:
		return pc + 1, "ADC A,B"
	case 0x89:
		return pc + 1, "ADC A,C"
	case 0x8A:
		return pc + 1, "ADC A,D"
	case 0x8B:
		return pc + 1, "ADC A,E"
	case 0x8C:
		return pc + 1, "ADC A,H"
	case 0x8D:
		return pc + 1, "ADC A,L"
	case 0xCE:
		return pc + 2, fmt.Sprintf("ADC A,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x8E:
		return pc + 1, "ADC A,(HL)"
	case 0x97:
		return pc + 1, "SUB A,A"
	case 0x90:
		return pc + 1, "SUB A,B"
	case 0x91:
		return pc + 1, "SUB A,C"
	case 0x92:
		return pc + 1, "SUB A,D"
	case 0x93:
		return pc + 1, "SUB A,E"
	case 0x94:
		return pc + 1, "SUB A,H"
	case 0x95:
		return pc + 1, "SUB A,L"
	case 0xD6:
		return pc + 2, fmt.Sprintf("SUB A,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x96:
		return pc + 1, "SUB A,(HL)"
	case 0x9F:
		return pc + 1, "SBC A,A"
	case 0x98:
		return pc + 1, "SBC A,B"
	case 0x99:
		return pc + 1, "SBC A,C"
	case 0x9A:
		return pc + 1, "SBC A,D"
	case 0x9B:
		return pc + 1, "SBC A,E"
	case 0x9C:
		return pc + 1, "SBC A,H"
	case 0x9D:
		return pc + 1, "SBC A,L"
	case 0xDE:
		return pc + 2, fmt.Sprintf("SBC A,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x9E:
		return pc + 1, "SBC A,(HL)"
	case 0xA7:
		return pc + 1, "AND A"
	case 0xA0:
		return pc + 1, "AND B"
	case 0xA1:
		return pc + 1, "AND C"
	case 0xA2:
		return pc + 1, "AND D"
	case 0xA3:
		return pc + 1, "AND E"
	case 0xA4:
		return pc + 1, "AND H"
	case 0xA5:
		return pc + 1, "AND L"
	case 0xE6:
		return pc + 2, fmt.Sprintf("AND [%#02x]", cpu.mmu.Read(pc+1))
	case 0xA6:
		return pc + 1, "AND (HL)"
	case 0xAF:
		return pc + 1, "XOR A"
	case 0xA8:
		return pc + 1, "XOR B"
	case 0xA9:
		return pc + 1, "XOR C"
	case 0xAA:
		return pc + 1, "XOR D"
	case 0xAB:
		return pc + 1, "XOR E"
	case 0xAC:
		return pc + 1, "XOR H"
	case 0xAD:
		return pc + 1, "XOR L"
	case 0xEE:
		return pc + 2, fmt.Sprintf("XOR [%#02x]", cpu.mmu.Read(pc+1))
	case 0xAE:
		return pc + 1, "XOR (HL)"
	case 0xB7:
		return pc + 1, "OR A"
	case 0xB0:
		return pc + 1, "OR B"
	case 0xB1:
		return pc + 1, "OR C"
	case 0xB2:
		return pc + 1, "OR D"
	case 0xB3:
		return pc + 1, "OR E"
	case 0xB4:
		return pc + 1, "OR H"
	case 0xB5:
		return pc + 1, "OR L"
	case 0xF6:
		return pc + 2, fmt.Sprintf("OR [%#02x]", cpu.mmu.Read(pc+1))
	case 0xB6:
		return pc + 1, "OR (HL)"
	case 0xBF:
		return pc + 1, "CP A"
	case 0xB8:
		return pc + 1, "CP B"
	case 0xB9:
		return pc + 1, "CP C"
	case 0xBA:
		return pc + 1, "CP D"
	case 0xBB:
		return pc + 1, "CP E"
	case 0xBC:
		return pc + 1, "CP H"
	case 0xBD:
		return pc + 1, "CP L"
	case 0xFE:
		return pc + 2, fmt.Sprintf("CP [%#02x]", cpu.mmu.Read(pc+1))
	case 0xBE:
		return pc + 1, "CP (HL)"
	case 0x3C:
		return pc + 1, "INC A"
	case 0x04:
		return pc + 1, "INC B"
	case 0x0C:
		return pc + 1, "INC C"
	case 0x14:
		return pc + 1, "INC D"
	case 0x1C:
		return pc + 1, "INC E"
	case 0x24:
		return pc + 1, "INC H"
	case 0x2C:
		return pc + 1, "INC L"
	case 0x34:
		return pc + 1, "INC (HL)"
	case 0x3D:
		return pc + 1, "DEC A"
	case 0x05:
		return pc + 1, "DEC B"
	case 0x0D:
		return pc + 1, "DEC C"
	case 0x15:
		return pc + 1, "DEC D"
	case 0x1D:
		return pc + 1, "DEC E"
	case 0x25:
		return pc + 1, "DEC H"
	case 0x2D:
		return pc + 1, "DEC L"
	case 0x35:
		return pc + 1, "DEC (HL)"
	case 0x27:
		return pc + 1, "DAA"
	case 0x2F:
		return pc + 1, "CPL"
	case 0x09:
		return pc + 1, "ADD HL,BC"
	case 0x19:
		return pc + 1, "ADD HL,DE"
	case 0x29:
		return pc + 1, "ADD HL,HL"
	case 0x39:
		return pc + 1, "ADD HL,SP"
	case 0x03:
		return pc + 1, "INC BC"
	case 0x13:
		return pc + 1, "INC DE"
	case 0x23:
		return pc + 1, "INC HL"
	case 0x33:
		return pc + 1, "INC SP"
	case 0x0B:
		return pc + 1, "DEC BC"
	case 0x1B:
		return pc + 1, "DEC DE"
	case 0x2B:
		return pc + 1, "DEC HL"
	case 0x3B:
		return pc + 1, "DEC SP"
	case 0xE8:
		return pc + 2, fmt.Sprintf("ADD SP,(%d)", int8(cpu.mmu.Read(pc+1)))
	case 0xF8:
		return pc + 2, fmt.Sprintf("LD HL,SP,(%d)", int8(cpu.mmu.Read(pc+1)))
	case 0x07:
		return pc + 1, "RLCA"
	case 0x17:
		return pc + 1, "RLA"
	case 0x0F:
		return pc + 1, "RRCA"
	case 0x1F:
		return pc + 1, "RRA"
	case 0x3F:
		return pc + 1, "CCF"
	case 0x37:
		return pc + 1, "SCF"
	case 0x00:
		return pc + 1, "NOP"
	case 0x76:
		return pc + 1, "HALT"
	case 0x10:
		return pc + 1, "STOP"
	case 0xF3:
		return pc + 1, "DI"
	case 0xFB:
		return pc + 1, "EI"
	case 0xC3:
		return pc + 3, fmt.Sprintf("JP [%#04x]", cpu.mmu.Read16(pc+1))
	case 0xE9:
		return pc + 1, "JP HL"
	case 0xC2:
		return pc + 3, fmt.Sprintf("JP NZ,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xCA:
		return pc + 3, fmt.Sprintf("JP Z,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xD2:
		return pc + 3, fmt.Sprintf("JP NC,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xDA:
		return pc + 3, fmt.Sprintf("JP C,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0x18:
		return pc + 2, fmt.Sprintf("JR [%#02x]", cpu.mmu.Read(pc+1))
	case 0x20:
		return pc + 2, fmt.Sprintf("JR NZ,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x28:
		return pc + 2, fmt.Sprintf("JR Z,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x30:
		return pc + 2, fmt.Sprintf("JR NC,[%#02x]", cpu.mmu.Read(pc+1))
	case 0x38:
		return pc + 2, fmt.Sprintf("JR C,[%#02x]", cpu.mmu.Read(pc+1))
	case 0xCD:
		return pc + 3, fmt.Sprintf("CALL [%#04x]", cpu.mmu.Read16(pc+1))
	case 0xC4:
		return pc + 3, fmt.Sprintf("CALL NZ,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xCC:
		return pc + 3, fmt.Sprintf("CALL Z,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xD4:
		return pc + 3, fmt.Sprintf("CALL NC,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xDC:
		return pc + 3, fmt.Sprintf("CALL C,[%#04x]", cpu.mmu.Read16(pc+1))
	case 0xC9:
		return pc + 1, "RET"
	case 0xC0:
		return pc + 1, "RET NZ"
	case 0xC8:
		return pc + 1, "RET Z"
	case 0xD0:
		return pc + 1, "RET NC"
	case 0xD8:
		return pc + 1, "RET C"
	case 0xD9:
		return pc + 1, "RETI"
	case 0xC7:
		return pc + 1, "RST 0x00"
	case 0xCF:
		return pc + 1, "RST 0x08"
	case 0xD7:
		return pc + 1, "RST 0x10"
	case 0xDF:
		return pc + 1, "RST 0x18"
	case 0xE7:
		return pc + 1, "RST 0x20"
	case 0xEF:
		return pc + 1, "RST 0x28"
	case 0xF7:
		return pc + 1, "RST 0x30"
	case 0xFF:
		return pc + 1, "RST 0x38"
	case 0xCB:
		return cpu.disassembleCB(pc + 1)
	default:
		return pc + 1, fmt.Sprintf("Unknown instruction: %#02x", opcode)
	}
}

func (cpu *CPU) disassembleCB(pc uint16) (uint16, string) {
	opcode := cpu.mmu.Read(pc)

	switch opcode {
	case 0x07:
		return pc + 1, "RLC A"
	case 0x00:
		return pc + 1, "RLC B"
	case 0x01:
		return pc + 1, "RLC C"
	case 0x02:
		return pc + 1, "RLC D"
	case 0x03:
		return pc + 1, "RLC E"
	case 0x04:
		return pc + 1, "RLC H"
	case 0x05:
		return pc + 1, "RLC L"
	case 0x06:
		return pc + 1, "RLC (HL)"
	case 0x17:
		return pc + 1, "RL A"
	case 0x10:
		return pc + 1, "RL B"
	case 0x11:
		return pc + 1, "RL C"
	case 0x12:
		return pc + 1, "RL D"
	case 0x13:
		return pc + 1, "RL E"
	case 0x14:
		return pc + 1, "RL H"
	case 0x15:
		return pc + 1, "RL L"
	case 0x16:
		return pc + 1, "RL (HL)"
	case 0x0F:
		return pc + 1, "RRC A"
	case 0x08:
		return pc + 1, "RRC B"
	case 0x09:
		return pc + 1, "RRC C"
	case 0x0A:
		return pc + 1, "RRC D"
	case 0x0B:
		return pc + 1, "RRC E"
	case 0x0C:
		return pc + 1, "RRC H"
	case 0x0D:
		return pc + 1, "RRC L"
	case 0x0E:
		return pc + 1, "RRC (HL)"
	case 0x1F:
		return pc + 1, "RR A"
	case 0x18:
		return pc + 1, "RR B"
	case 0x19:
		return pc + 1, "RR C"
	case 0x1A:
		return pc + 1, "RR D"
	case 0x1B:
		return pc + 1, "RR E"
	case 0x1C:
		return pc + 1, "RR H"
	case 0x1D:
		return pc + 1, "RR L"
	case 0x1E:
		return pc + 1, "RR (HL)"
	case 0x27:
		return pc + 1, "SLA A"
	case 0x20:
		return pc + 1, "SLA B"
	case 0x21:
		return pc + 1, "SLA C"
	case 0x22:
		return pc + 1, "SLA D"
	case 0x23:
		return pc + 1, "SLA E"
	case 0x24:
		return pc + 1, "SLA H"
	case 0x25:
		return pc + 1, "SLA L"
	case 0x26:
		return pc + 1, "SLA (HL)"
	case 0x37:
		return pc + 1, "SWAP A"
	case 0x30:
		return pc + 1, "SWAP B"
	case 0x31:
		return pc + 1, "SWAP C"
	case 0x32:
		return pc + 1, "SWAP D"
	case 0x33:
		return pc + 1, "SWAP E"
	case 0x34:
		return pc + 1, "SWAP H"
	case 0x35:
		return pc + 1, "SWAP L"
	case 0x36:
		return pc + 1, "SWAP (HL)"
	case 0x2F:
		return pc + 1, "SRA A"
	case 0x28:
		return pc + 1, "SRA B"
	case 0x29:
		return pc + 1, "SRA C"
	case 0x2A:
		return pc + 1, "SRA D"
	case 0x2B:
		return pc + 1, "SRA E"
	case 0x2C:
		return pc + 1, "SRA H"
	case 0x2D:
		return pc + 1, "SRA L"
	case 0x2E:
		return pc + 1, "SRA (HL)"
	case 0x3F:
		return pc + 1, "SRL A"
	case 0x38:
		return pc + 1, "SRL B"
	case 0x39:
		return pc + 1, "SRL C"
	case 0x3A:
		return pc + 1, "SRL D"
	case 0x3B:
		return pc + 1, "SRL E"
	case 0x3C:
		return pc + 1, "SRL H"
	case 0x3D:
		return pc + 1, "SRL L"
	case 0x3E:
		return pc + 1, "SRL (HL)"
	case 0x47:
		return pc + 1, "BIT 0,A"
	case 0x40:
		return pc + 1, "BIT 0,B"
	case 0x41:
		return pc + 1, "BIT 0,C"
	case 0x42:
		return pc + 1, "BIT 0,D"
	case 0x43:
		return pc + 1, "BIT 0,E"
	case 0x44:
		return pc + 1, "BIT 0,H"
	case 0x45:
		return pc + 1, "BIT 0,L"
	case 0x46:
		return pc + 1, "BIT 0,(HL)"
	case 0x4F:
		return pc + 1, "BIT 1,A"
	case 0x48:
		return pc + 1, "BIT 1,B"
	case 0x49:
		return pc + 1, "BIT 1,C"
	case 0x4A:
		return pc + 1, "BIT 1,D"
	case 0x4B:
		return pc + 1, "BIT 1,E"
	case 0x4C:
		return pc + 1, "BIT 1,H"
	case 0x4D:
		return pc + 1, "BIT 1,L"
	case 0x4E:
		return pc + 1, "BIT 1,(HL)"
	case 0x57:
		return pc + 1, "BIT 2,A"
	case 0x50:
		return pc + 1, "BIT 2,B"
	case 0x51:
		return pc + 1, "BIT 2,C"
	case 0x52:
		return pc + 1, "BIT 2,D"
	case 0x53:
		return pc + 1, "BIT 2,E"
	case 0x54:
		return pc + 1, "BIT 2,H"
	case 0x55:
		return pc + 1, "BIT 2,L"
	case 0x56:
		return pc + 1, "BIT 2,(HL)"
	case 0x5F:
		return pc + 1, "BIT 3,A"
	case 0x58:
		return pc + 1, "BIT 3,B"
	case 0x59:
		return pc + 1, "BIT 3,C"
	case 0x5A:
		return pc + 1, "BIT 3,D"
	case 0x5B:
		return pc + 1, "BIT 3,E"
	case 0x5C:
		return pc + 1, "BIT 3,H"
	case 0x5D:
		return pc + 1, "BIT 3,L"
	case 0x5E:
		return pc + 1, "BIT 3,(HL)"
	case 0x67:
		return pc + 1, "BIT 4,A"
	case 0x60:
		return pc + 1, "BIT 4,B"
	case 0x61:
		return pc + 1, "BIT 4,C"
	case 0x62:
		return pc + 1, "BIT 4,D"
	case 0x63:
		return pc + 1, "BIT 4,E"
	case 0x64:
		return pc + 1, "BIT 4,H"
	case 0x65:
		return pc + 1, "BIT 4,L"
	case 0x66:
		return pc + 1, "BIT 4,(HL)"
	case 0x6F:
		return pc + 1, "BIT 5,A"
	case 0x68:
		return pc + 1, "BIT 5,B"
	case 0x69:
		return pc + 1, "BIT 5,C"
	case 0x6A:
		return pc + 1, "BIT 5,D"
	case 0x6B:
		return pc + 1, "BIT 5,E"
	case 0x6C:
		return pc + 1, "BIT 5,H"
	case 0x6D:
		return pc + 1, "BIT 5,L"
	case 0x6E:
		return pc + 1, "BIT 5,(HL)"
	case 0x77:
		return pc + 1, "BIT 6,A"
	case 0x70:
		return pc + 1, "BIT 6,B"
	case 0x71:
		return pc + 1, "BIT 6,C"
	case 0x72:
		return pc + 1, "BIT 6,D"
	case 0x73:
		return pc + 1, "BIT 6,E"
	case 0x74:
		return pc + 1, "BIT 6,H"
	case 0x75:
		return pc + 1, "BIT 6,L"
	case 0x76:
		return pc + 1, "BIT 6,(HL)"
	case 0x7F:
		return pc + 1, "BIT 7,A"
	case 0x78:
		return pc + 1, "BIT 7,B"
	case 0x79:
		return pc + 1, "BIT 7,C"
	case 0x7A:
		return pc + 1, "BIT 7,D"
	case 0x7B:
		return pc + 1, "BIT 7,E"
	case 0x7C:
		return pc + 1, "BIT 7,H"
	case 0x7D:
		return pc + 1, "BIT 7,L"
	case 0x7E:
		return pc + 1, "BIT 7,(HL)"
	case 0xC7:
		return pc + 1, "SET 0,A"
	case 0xC0:
		return pc + 1, "SET 0,B"
	case 0xC1:
		return pc + 1, "SET 0,C"
	case 0xC2:
		return pc + 1, "SET 0,D"
	case 0xC3:
		return pc + 1, "SET 0,E"
	case 0xC4:
		return pc + 1, "SET 0,H"
	case 0xC5:
		return pc + 1, "SET 0,L"
	case 0xC6:
		return pc + 1, "SET 0,(HL)"
	case 0xCF:
		return pc + 1, "SET 1,A"
	case 0xC8:
		return pc + 1, "SET 1,B"
	case 0xC9:
		return pc + 1, "SET 1,C"
	case 0xCA:
		return pc + 1, "SET 1,D"
	case 0xCB:
		return pc + 1, "SET 1,E"
	case 0xCC:
		return pc + 1, "SET 1,H"
	case 0xCD:
		return pc + 1, "SET 1,L"
	case 0xCE:
		return pc + 1, "SET 1,(HL)"
	case 0xD7:
		return pc + 1, "SET 2,A"
	case 0xD0:
		return pc + 1, "SET 2,B"
	case 0xD1:
		return pc + 1, "SET 2,C"
	case 0xD2:
		return pc + 1, "SET 2,D"
	case 0xD3:
		return pc + 1, "SET 2,E"
	case 0xD4:
		return pc + 1, "SET 2,H"
	case 0xD5:
		return pc + 1, "SET 2,L"
	case 0xD6:
		return pc + 1, "SET 2,(HL)"
	case 0xDF:
		return pc + 1, "SET 3,A"
	case 0xD8:
		return pc + 1, "SET 3,B"
	case 0xD9:
		return pc + 1, "SET 3,C"
	case 0xDA:
		return pc + 1, "SET 3,D"
	case 0xDB:
		return pc + 1, "SET 3,E"
	case 0xDC:
		return pc + 1, "SET 3,H"
	case 0xDD:
		return pc + 1, "SET 3,L"
	case 0xDE:
		return pc + 1, "SET 3,(HL)"
	case 0xE7:
		return pc + 1, "SET 4,A"
	case 0xE0:
		return pc + 1, "SET 4,B"
	case 0xE1:
		return pc + 1, "SET 4,C"
	case 0xE2:
		return pc + 1, "SET 4,D"
	case 0xE3:
		return pc + 1, "SET 4,E"
	case 0xE4:
		return pc + 1, "SET 4,H"
	case 0xE5:
		return pc + 1, "SET 4,L"
	case 0xE6:
		return pc + 1, "SET 4,(HL)"
	case 0xEF:
		return pc + 1, "SET 5,A"
	case 0xE8:
		return pc + 1, "SET 5,B"
	case 0xE9:
		return pc + 1, "SET 5,C"
	case 0xEA:
		return pc + 1, "SET 5,D"
	case 0xEB:
		return pc + 1, "SET 5,E"
	case 0xEC:
		return pc + 1, "SET 5,H"
	case 0xED:
		return pc + 1, "SET 5,L"
	case 0xEE:
		return pc + 1, "SET 5,(HL)"
	case 0xF7:
		return pc + 1, "SET 6,A"
	case 0xF0:
		return pc + 1, "SET 6,B"
	case 0xF1:
		return pc + 1, "SET 6,C"
	case 0xF2:
		return pc + 1, "SET 6,D"
	case 0xF3:
		return pc + 1, "SET 6,E"
	case 0xF4:
		return pc + 1, "SET 6,H"
	case 0xF5:
		return pc + 1, "SET 6,L"
	case 0xF6:
		return pc + 1, "SET 6,(HL)"
	case 0xFF:
		return pc + 1, "SET 7,A"
	case 0xF8:
		return pc + 1, "SET 7,B"
	case 0xF9:
		return pc + 1, "SET 7,C"
	case 0xFA:
		return pc + 1, "SET 7,D"
	case 0xFB:
		return pc + 1, "SET 7,E"
	case 0xFC:
		return pc + 1, "SET 7,H"
	case 0xFD:
		return pc + 1, "SET 7,L"
	case 0xFE:
		return pc + 1, "SET 7,(HL)"
	case 0x87:
		return pc + 1, "RES 0,A"
	case 0x80:
		return pc + 1, "RES 0,B"
	case 0x81:
		return pc + 1, "RES 0,C"
	case 0x82:
		return pc + 1, "RES 0,D"
	case 0x83:
		return pc + 1, "RES 0,E"
	case 0x84:
		return pc + 1, "RES 0,H"
	case 0x85:
		return pc + 1, "RES 0,L"
	case 0x86:
		return pc + 1, "RES 0,(HL)"
	case 0x8F:
		return pc + 1, "RES 1,A"
	case 0x88:
		return pc + 1, "RES 1,B"
	case 0x89:
		return pc + 1, "RES 1,C"
	case 0x8A:
		return pc + 1, "RES 1,D"
	case 0x8B:
		return pc + 1, "RES 1,E"
	case 0x8C:
		return pc + 1, "RES 1,H"
	case 0x8D:
		return pc + 1, "RES 1,L"
	case 0x8E:
		return pc + 1, "RES 1,(HL)"
	case 0x97:
		return pc + 1, "RES 2,A"
	case 0x90:
		return pc + 1, "RES 2,B"
	case 0x91:
		return pc + 1, "RES 2,C"
	case 0x92:
		return pc + 1, "RES 2,D"
	case 0x93:
		return pc + 1, "RES 2,E"
	case 0x94:
		return pc + 1, "RES 2,H"
	case 0x95:
		return pc + 1, "RES 2,L"
	case 0x96:
		return pc + 1, "RES 2,(HL)"
	case 0x9F:
		return pc + 1, "RES 3,A"
	case 0x98:
		return pc + 1, "RES 3,B"
	case 0x99:
		return pc + 1, "RES 3,C"
	case 0x9A:
		return pc + 1, "RES 3,D"
	case 0x9B:
		return pc + 1, "RES 3,E"
	case 0x9C:
		return pc + 1, "RES 3,H"
	case 0x9D:
		return pc + 1, "RES 3,L"
	case 0x9E:
		return pc + 1, "RES 3,(HL)"
	case 0xA7:
		return pc + 1, "RES 4,A"
	case 0xA0:
		return pc + 1, "RES 4,B"
	case 0xA1:
		return pc + 1, "RES 4,C"
	case 0xA2:
		return pc + 1, "RES 4,D"
	case 0xA3:
		return pc + 1, "RES 4,E"
	case 0xA4:
		return pc + 1, "RES 4,H"
	case 0xA5:
		return pc + 1, "RES 4,L"
	case 0xA6:
		return pc + 1, "RES 4,(HL)"
	case 0xAF:
		return pc + 1, "RES 5,A"
	case 0xA8:
		return pc + 1, "RES 5,B"
	case 0xA9:
		return pc + 1, "RES 5,C"
	case 0xAA:
		return pc + 1, "RES 5,D"
	case 0xAB:
		return pc + 1, "RES 5,E"
	case 0xAC:
		return pc + 1, "RES 5,H"
	case 0xAD:
		return pc + 1, "RES 5,L"
	case 0xAE:
		return pc + 1, "RES 5,(HL)"
	case 0xB7:
		return pc + 1, "RES 6,A"
	case 0xB0:
		return pc + 1, "RES 6,B"
	case 0xB1:
		return pc + 1, "RES 6,C"
	case 0xB2:
		return pc + 1, "RES 6,D"
	case 0xB3:
		return pc + 1, "RES 6,E"
	case 0xB4:
		return pc + 1, "RES 6,H"
	case 0xB5:
		return pc + 1, "RES 6,L"
	case 0xB6:
		return pc + 1, "RES 6,(HL)"
	case 0xBF:
		return pc + 1, "RES 7,A"
	case 0xB8:
		return pc + 1, "RES 7,B"
	case 0xB9:
		return pc + 1, "RES 7,C"
	case 0xBA:
		return pc + 1, "RES 7,D"
	case 0xBB:
		return pc + 1, "RES 7,E"
	case 0xBC:
		return pc + 1, "RES 7,H"
	case 0xBD:
		return pc + 1, "RES 7,L"
	case 0xBE:
		return pc + 1, "RES 7,(HL)"
	default:
		return pc + 1, fmt.Sprintf("Unknown CB instruction: %#02x", opcode)
	}
}
