package gbc

type memoryDevice interface {
	Read(uint16) byte
	Write(uint16, byte)
}

type standardMemoryDevice struct {
	buf []byte
}

func newStandardMemoryDevice(size int) *standardMemoryDevice {
	smd := new(standardMemoryDevice)
	smd.buf = make([]byte, size)

	return smd
}

func (smd *standardMemoryDevice) Read(addr uint16) byte {
	return smd.buf[addr]
}

func (smd *standardMemoryDevice) Write(addr uint16, val byte) {
	smd.buf[addr] = val
}

type zeroMemoryDevice struct {
}

func newZeroMemoryDevice() *zeroMemoryDevice {
	return new(zeroMemoryDevice)
}

func (zmd *zeroMemoryDevice) Read(addr uint16) byte {
	return 0
}

func (zmd *zeroMemoryDevice) Write(addr uint16, val byte) {}
