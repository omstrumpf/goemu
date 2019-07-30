package backends

import (
	"image/color"
	"time"
)

// Console is the interface for a console backend
type Console interface {
	Tick()
	GetFrameBuffer() []color.RGBA
	GetFrameTime() time.Duration
	GetScreenWidth() int
	GetScreenHeight() int
	GetConsoleName() string
}
