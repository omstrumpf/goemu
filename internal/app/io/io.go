package io

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/omstrumpf/goemu/internal/app/console"
)

// IO manages the graphical output of the emulator
type IO struct {
	console console.Console

	win *pixelgl.Window
	pic *pixel.PictureData

	paused bool
}

// NewIO constructs a valid IO struct
func NewIO(console console.Console) *IO {
	io := new(IO)

	io.console = console

	io.setupWindow()
	io.setupPicture()

	return io
}

// ProcessInput reads input and writes it to the console
func (io *IO) ProcessInput() {
	for key, val := range functionKeys {
		if io.win.JustPressed(key) {
			val(io)
		}
	}

	for key, val := range buttonKeys {
		if io.win.JustPressed(key) {
			io.console.PressButton(val)
		}
		if io.win.JustReleased(key) {
			io.console.ReleaseButton(val)
		}
	}
}

// Render renders the console's frame buffer to the display
func (io *IO) Render() {
	io.win.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 0xFF})

	picture := pixel.Picture(io.pic)
	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(io.win, pixel.IM)

	shift := io.win.Bounds().Size().Scaled(0.5).Sub(pixel.ZV)
	mat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(io.getScaleFactor(), io.getScaleFactor()*-1)).Moved(shift)
	io.win.SetMatrix(mat)

	io.win.Update()
}

// ShouldEmulate returns true if the emulator should emulate (not paused)
func (io *IO) ShouldEmulate() bool {
	return !io.paused
}

// ShouldExit returns true if the emulator should exit
func (io *IO) ShouldExit() bool {
	return io.win.Closed()
}

func (io *IO) setupWindow() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "GoEmu Emulator (" + io.console.GetConsoleName() + ")",
		Bounds:    pixel.R(0, 0, float64(io.console.GetScreenWidth()), float64(io.console.GetScreenHeight())),
		Resizable: true,
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

func (io *IO) getScaleFactor() float64 {
	scaleWidth := io.win.Bounds().W() / float64(io.console.GetScreenWidth())
	scaleHeight := io.win.Bounds().H() / float64(io.console.GetScreenHeight())

	if scaleWidth < scaleHeight {
		return scaleWidth
	}
	return scaleHeight
}
