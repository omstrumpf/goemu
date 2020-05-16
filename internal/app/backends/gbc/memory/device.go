package memory

// Device is a device that supports a Read/Write memory interface
type Device interface {
	Read(uint16) byte
	Write(uint16, byte)
}

// Simple is a simple Device backed by a buffer
type Simple struct {
	buf []byte
}

// NewSimple creates a Simple of the given size
func NewSimple(size int) *Simple {
	s := new(Simple)
	s.buf = make([]byte, size)

	return s
}

// NewSimpleWithData creats a Simple backed by the given data slice
func NewSimpleWithData(data []byte) *Simple {
	s := new(Simple)
	s.buf = data

	return s
}

func (s *Simple) Read(addr uint16) byte {
	return s.buf[addr]
}

func (s *Simple) Write(addr uint16, val byte) {
	s.buf[addr] = val
}

// Zero is a Device that always contains 0
type Zero struct{}

// NewZero creates a Zero
func NewZero() *Zero {
	return new(Zero)
}

func (z *Zero) Read(addr uint16) byte {
	return 0
}

func (z *Zero) Write(addr uint16, val byte) {}

// High is a Device that always contains 0xFF
type High struct{}

// NewHigh creates a High
func NewHigh() *High {
	return new(High)
}

func (h *High) Read(addr uint16) byte {
	return 0xFF
}

func (h *High) Write(addr uint16, val byte) {}
