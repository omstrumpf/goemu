package gba

// CPU represents the central processing unit. Stores the state and instruction map.
type CPU struct {
	AF Register // Accumulator and flags
	BC Register // General use
	DE Register // General use
	HL Register // General use

	SP Register // Stack pointer
	PC Register // Program counter

	halt bool
	stop bool
	ime  bool // Interrupt disable

	clock int

	mmu *MMU // Memory Management Unit

	instructions [0x100]func() // Instruction map
	cycles       [0x100]int    // Cycle cost map
}

// NewCPU constructs a valid CPU struct
func NewCPU(mmu *MMU) *CPU {
	cpu := new(CPU)

	cpu.mmu = mmu

	cpu.AF.mask = 0xFFF0 // Last four flag bits not used (always 0)
	cpu.halt = false
	cpu.ime = false

	cpu.PopulateInstructions()

	return cpu
}

// ProcessNextInstruction fetches the next instruction, executes it, and increments the clock accordingly
func (cpu *CPU) ProcessNextInstruction() {
	// Fetch the next instruction and increment PC
	opcode := cpu.mmu.Read(cpu.PC.Inc())

	// Execute the instruction
	cpu.instructions[opcode]()

	// Increment the clock accordingly
	cpu.clock += cpu.cycles[opcode]
}

// IsHalted returns true if the CPU is halted
func (cpu *CPU) IsHalted() bool {
	return cpu.halt
}

// IsStopped returns true if the CPU is stopped
func (cpu *CPU) IsStopped() bool {
	return cpu.stop
}

// Flags
func (cpu *CPU) z() bool {
	return cpu.AF.HiLo()>>7&1 == 1
}

func (cpu *CPU) n() bool {
	return cpu.AF.HiLo()>>6&1 == 1
}

func (cpu *CPU) h() bool {
	return cpu.AF.HiLo()>>5&1 == 1
}

func (cpu *CPU) c() bool {
	return cpu.AF.HiLo()>>4&1 == 1
}

func (cpu *CPU) setZ(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<7)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & uint16(^byte(1<<7)))
	}
}

func (cpu *CPU) setN(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<6)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & uint16(^byte(1<<6)))
	}
}

func (cpu *CPU) setH(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<5)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & uint16(^byte(1<<5)))
	}
}

func (cpu *CPU) setC(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<4)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & uint16(^byte(1<<4)))
	}
}
