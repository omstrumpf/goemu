package gbc

type interruptMemoryDevice struct {
	enable byte
	flag   byte
}

const (
	interruptVBlankBit = 0
	interruptLCDBit    = 1
	interruptTimerBit  = 2
	interruptSerialBit = 3
	interruptJoypadBit = 1
)

func newInterruptMemoryDevice() *interruptMemoryDevice {
	imd := new(interruptMemoryDevice)

	return imd
}

func (imd *interruptMemoryDevice) Read(addr uint16) byte {
	switch addr {
	case 0xFF0F:
		return imd.flag
	case 0xFFFF:
		return imd.enable
	}

	logger.Warningf("Encountered unexpected interrupt location: %#4x", addr)
	return 0
}

func (imd *interruptMemoryDevice) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF0F:
		imd.flag = val
		return
	case 0xFFFF:
		imd.enable = val
		return
	}

	logger.Warningf("Encountered unexpected interrupt location: %#4x", addr)
}

func (imd *interruptMemoryDevice) RequestInterrupt(bit uint8) {
	imd.flag = imd.flag | (1 << bit)
}

func (imd *interruptMemoryDevice) ResetInterrupt(bit uint8) {
	imd.flag = imd.flag & ^(1 << bit)
}
