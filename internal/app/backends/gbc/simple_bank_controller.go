package gbc

// SimpleBankController is a basic cartridge memory controller with a fixed 32K of ROM, and no RAM
type SimpleBankController struct {
	buf [0x8000]byte
}

// NewSimpleBankController constructs a valid SimpleBankController struct
func NewSimpleBankController(data []byte) *SimpleBankController {
	sbc := new(SimpleBankController)

	if len(data) > 0x8000 {
		logger.Warningf("SMC loading oversized ROM. Data will be truncated.")
	}

	copy(sbc.buf[:], data)

	return sbc
}

func (sbc *SimpleBankController) Read(addr uint16) byte {
	if addr > 0x8000 {
		logger.Errorf("SMC encountered read out of range: %#4x", addr)
		return 0
	}

	return sbc.buf[addr]
}

func (sbc *SimpleBankController) Write(addr uint16, val byte) {
	logger.Warningf("SMC ROM write encountered: %#4x", addr)
}
