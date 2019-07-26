package gba

// CPU Represents the central processing unit. Stores the state and instruction map.
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
func (cpu *CPU) ProcessNextInstruction() int {
	// Fetch the next instruction and increment PC
	opcode := cpu.mmu.Read(cpu.PC.Inc())

	// Execute the instruction
	cpu.instructions[opcode]()

	// Return the consumed cycles
	return cpu.cycles[opcode]
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
