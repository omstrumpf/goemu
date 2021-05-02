package banking

// Controller is a memory bank controller
type Controller interface {
	Read(uint16) byte
	Write(uint16, byte)

	RunForClocks(int)

	// TODO add save interactions
}
