package gba

// Register represents a 16 bit GBA register
type Register struct {
	val uint16 // The value in the register

	mask uint16 // A mask over the register, for bits that cannot be set
}

// Hi returns the high byte
func (reg *Register) Hi() byte {
	return byte(reg.val >> 8)
}

// Lo returns the low byte
func (reg *Register) Lo() byte {
	return byte(reg.val & 0xFF)
}

// HiLo returns the entire 16-bit value
func (reg *Register) HiLo() uint16 {
	return reg.val
}

// SetHi sets the high byte
func (reg *Register) SetHi(val byte) {
	reg.val = uint16(val)<<8 | (uint16(reg.val) & 0xFF)
	reg.updateMask()
}

// SetLo sets the low byte
func (reg *Register) SetLo(val byte) {
	reg.val = uint16(val) | (uint16(reg.val) & 0xFF00)
	reg.updateMask()
}

// Set sets the entire 16-bit value
func (reg *Register) Set(val uint16) {
	reg.val = val
	reg.updateMask()
}

// Inc increments and returns original value
func (reg *Register) Inc() uint16 {
	ret := reg.val
	reg.val++
	return ret
}

// Inc2 increments by 2 and returns original value
func (reg *Register) Inc2() uint16 {
	ret := reg.val
	reg.val += 2
	return ret
}

// Dec decrements and returns original value
func (reg *Register) Dec() uint16 {
	ret := reg.val
	reg.val--
	return ret
}

// Dec2 decrements by 2 and returns original value
func (reg *Register) Dec2() uint16 {
	ret := reg.val
	reg.val -= 2
	return ret
}

// updateMask applies the mask if present
func (reg *Register) updateMask() {
	if reg.mask != 0 {
		reg.val &= reg.mask
	}
}
