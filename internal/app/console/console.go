package console

import "time"

// Console is the interface for a console backend
type Console interface {
	GetFrameTime() time.Duration
	Tick()
}
