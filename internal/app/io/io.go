package io

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/omstrumpf/goemu/internal/app/backends"
)

// IO manages the graphical output of the emulator
type IO struct {
	console backends.Console

	win *pixelgl.Window
	pic *pixel.PictureData
}

// NewIO constructs a valid IO struct
func NewIO(console backends.Console) *IO {
	io := new(IO)

	io.console = console

	io.setupWindow()
	io.setupPicture()

	return io
}

// ProcessInput reads input and writes it to the gameboy
func (io *IO) ProcessInput() {
	// TODO
}

// Render renders the gba's frame buffer to the display
func (io *IO) Render() {
	io.win.Clear(color.RGBA{R: 255, G: 0, B: 0, A: 0xFF})

	picture := pixel.Picture(io.pic)
	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(io.win, pixel.IM)

	shift := io.win.Bounds().Size().Scaled(0.5).Sub(pixel.ZV)
	mat := pixel.IM.Scaled(pixel.ZV, 1).Moved(shift)
	io.win.SetMatrix(mat)

	io.win.Update()
}

// ShouldExit returns true if the emulator should exit
func (io *IO) ShouldExit() bool {
	return io.win.Closed()
}

func (io *IO) setupWindow() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "GoEmu Emulator (" + io.console.GetConsoleName() + ")",
		Bounds: pixel.R(0, 0, float64(io.console.GetScreenWidth()), float64(io.console.GetScreenHeight())), // TODO scaling
	})
	if err != nil {
		panic(err)
	}

	io.win = win
}

func (io *IO) setupPicture() {
	io.pic = &pixel.PictureData{
		Pix:    io.console.GetFrameBuffer(),
		Stride: io.console.GetScreenWidth(),
		Rect:   pixel.R(0, 0, float64(io.console.GetScreenWidth()), float64(io.console.GetScreenHeight())),
	}
}
