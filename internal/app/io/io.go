package io

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/omstrumpf/gogba/internal/app/console"
)

// IO manages the graphical output of the emulator
type IO struct {
	console console.Console

	win *pixelgl.Window
}

// NewIO constructs a valid IO struct
func NewIO(console console.Console) *IO {
	io := new(IO)

	io.console = console

	io.setupWindow()

	return io
}

// ProcessInput reads input and writes it to the gameboy
func (io *IO) ProcessInput() {
	// TODO
}

// Render renders the gba's frame buffer to the display
func (io *IO) Render() {
	io.win.Update()
}

// ShouldExit returns true if the emulator should exit
func (io *IO) ShouldExit() bool {
	return io.win.Closed()
}

func (io *IO) setupWindow() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "GoGBA Emulator", // TODO  get game title from ROM
		Bounds: pixel.R(0, 0, 160, 144),
	})
	if err != nil {
		panic(err)
	}

	io.win = win
}
