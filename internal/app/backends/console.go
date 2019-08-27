package backends

import (
	"image/color"
	"time"

	"github.com/omstrumpf/goemu/internal/app/buttons"
)

// Console is the interface for a console backend
type Console interface {
	Tick()

	LoadROM([]byte)

	PressButton(buttons.Button)
	ReleaseButton(buttons.Button)

	IsStopped() bool

	GetFrameBuffer() []color.RGBA
	GetFrameTime() time.Duration
	GetScreenWidth() int
	GetScreenHeight() int
	GetConsoleName() string
}
