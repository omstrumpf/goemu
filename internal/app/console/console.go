package console

import (
	"image/color"
	"time"
)

// Console is the interface for a console backend
type Console interface {
	Tick()

	PressButton(Button)
	ReleaseButton(Button)

	IsStopped() bool

	GetFrameBuffer() []color.RGBA
	GetFrameTime() time.Duration
	GetAudioChannel() *chan AudioSample
	GetAudioBitrate() int
	GetScreenWidth() int
	GetScreenHeight() int
	GetConsoleName() string
	GetGameName() string
}

// Button represents a button on the console
type Button byte

const (
	// ButtonDown is the down button on the d pad
	ButtonDown Button = 0
	// ButtonUp is the up button on the d pad
	ButtonUp Button = 1
	// ButtonLeft is the left button on the d pad
	ButtonLeft Button = 2
	// ButtonRight is the right button on the d pad
	ButtonRight Button = 3
	// ButtonStart is the start button
	ButtonStart Button = 6
	// ButtonSelect is the select button
	ButtonSelect Button = 7
	// ButtonB is the b button
	ButtonB Button = 5
	// ButtonA is the a button
	ButtonA Button = 4
)
