package interrupts

import "github.com/omstrumpf/goemu/internal/app/log"

// InterruptDevice manages enabled/flagged interrupts
type InterruptDevice struct {
	enable byte
	flag   byte
}

// Interrupt bit indices
const (
	VBlankBit = 0
	LCDBit    = 1
	TimerBit  = 2
	SerialBit = 3
	JoypadBit = 4
)

// NewInterruptDevice constructs a valid InterruptDevice struct
func NewInterruptDevice() *InterruptDevice {
	id := new(InterruptDevice)

	return id
}

// Request requests the given interrupt bit
func (id *InterruptDevice) Request(bit uint8) {
	id.flag = id.flag | (1 << bit)
}

// Reset resets the given interrupt bit
func (id *InterruptDevice) Reset(bit uint8) {
	id.flag = id.flag & ^(1 << bit)
}

func (id *InterruptDevice) Read(addr uint16) byte {
	switch addr {
	case 0xFF0F:
		return id.flag
	case 0xFFFF:
		return id.enable
	}

	log.Warningf("Encountered unexpected interrupt location: %#4x", addr)
	return 0
}

func (id *InterruptDevice) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF0F:
		id.flag = val
		return
	case 0xFFFF:
		id.enable = val
		return
	}

	log.Warningf("Encountered unexpected interrupt location: %#4x", addr)
}
