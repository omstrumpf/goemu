package io

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/io/audio"
	audio_inspector "github.com/omstrumpf/goemu/internal/app/io/audio/inspector"
)

// IO manages the graphical and audio output of the emulator
type IO struct {
	console console.Console

	win *pixelgl.Window
	pic *pixel.PictureData

	audioInspector *audio_inspector.AudioInspector
	audioPlayer    *audio.Player

	paused bool
	muted  bool
}

// NewIO constructs a valid IO struct
func NewIO(console console.Console) *IO {
	io := new(IO)

	io.console = console

	io.audioInspector = audio_inspector.NewAudioInspector() // TODO only open this at user request
	io.audioPlayer = audio.NewPlayer(io.console.GetAudioBitrate())

	io.setupWindow()

	go io.distributeAudio()

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
	picture := pixel.Picture(io.pic)
	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(io.win, pixel.IM)

	shift := io.win.Bounds().Size().Scaled(0.5).Sub(pixel.ZV)
	mat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(io.getScaleFactor(), io.getScaleFactor()*-1)).Moved(shift)
	io.win.SetMatrix(mat)

	io.win.Update()

	io.audioInspector.Render()
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
		Title:     "GoEmu Emulator (" + io.console.GetConsoleName() + " - " + io.console.GetGameName() + ")",
		Bounds:    pixel.R(0, 0, float64(io.console.GetScreenWidth()), float64(io.console.GetScreenHeight())),
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}

	io.win = win

	io.pic = &pixel.PictureData{
		Pix:    io.console.GetFrameBuffer(),
		Stride: io.console.GetScreenWidth(),
		Rect:   pixel.R(0, 0, float64(io.console.GetScreenWidth()), float64(io.console.GetScreenHeight())),
	}
}

func (io *IO) distributeAudio() {
	channel := io.console.GetAudioChannel()

	for {
		if io == nil {
			return
		}

		select {
		case sample := <-*channel:
			io.audioPlayer.InputChannel <- sample.Combine()
			// TODO this assumes 4 channels
			io.audioInspector.InputChannel <- [4]float64{sample.Channels[0].M(), sample.Channels[1].M(), sample.Channels[2].M(), sample.Channels[3].M()}
		default:
			time.Sleep(100 * time.Millisecond)
		}
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

func (io *IO) mute() {
	fmt.Println("Muting audio.")

	io.audioPlayer.SetMute(true)
	io.muted = true
}

func (io *IO) unmute() {
	fmt.Println("Unmuting audio.")

	io.audioPlayer.SetMute(false)
	io.muted = false
}
