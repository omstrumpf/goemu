package gbc

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

	instructions   [0x100]func() // Instruction map
	instructionsCB [0x100]func() // CB instruction map
	cycles         [0x100]int    // Cycle cost map
	cyclesCB       [0x100]int    // CB cycle cost map
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
	if cpu.stop {
		return
	}

	if cpu.halt {
		cpu.clock++
	} else {
		// Fetch the next instruction and increment PC
		opcode := cpu.mmu.Read(cpu.PC.Inc())

		// Execute the instruction
		cpu.instructions[opcode]()

		// Increment the clock accordingly
		cpu.clock += cpu.cycles[opcode]
	}

	// Check for interrupts
	cpu.handleInterrupts()
}

// IsHalted returns true if the CPU is halted
func (cpu *CPU) IsHalted() bool {
	return cpu.halt
}

// IsStopped returns true if the CPU is stopped
func (cpu *CPU) IsStopped() bool {
	return cpu.stop
}

// Interrupts
func (cpu *CPU) handleInterrupts() {
	if !cpu.ime {
		return
	}

	interruptByte := cpu.mmu.Read(0xFFFF) & cpu.mmu.Read(0xFF0F)

	if interruptByte&0x1F != 0 {
		cpu.ime = false
		cpu.halt = false
	}

	if interruptByte&1 != 0 { // V-Blank
		cpu.mmu.interrupts.ResetInterrupt(interruptVBlankBit)
		cpu.call(0x40)
	} else if interruptByte&2 != 0 { // LCD STAT
		cpu.mmu.interrupts.ResetInterrupt(interruptLCDBit)
		cpu.call(0x48)
	} else if interruptByte&4 != 0 { // Timer
		cpu.mmu.interrupts.ResetInterrupt(interruptTimerBit)
		cpu.call(0x50)
	} else if interruptByte&8 != 0 { // Serial
		cpu.mmu.interrupts.ResetInterrupt(interruptSerialBit)
		cpu.call(0x58)
	} else if interruptByte&16 != 0 { // Joypad
		cpu.mmu.interrupts.ResetInterrupt(interruptJoypadBit)
		cpu.call(0x60)
	}
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
		cpu.AF.Set(cpu.AF.HiLo() & ^uint16(1<<7))
	}
}

func (cpu *CPU) setN(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<6)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & ^uint16(1<<6))
	}
}

func (cpu *CPU) setH(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<5)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & ^uint16(1<<5))
	}
}

func (cpu *CPU) setC(set bool) {
	if set {
		cpu.AF.Set(cpu.AF.HiLo() | 1<<4)
	} else {
		cpu.AF.Set(cpu.AF.HiLo() & ^uint16(1<<4))
	}
}
