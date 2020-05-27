package io

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/omstrumpf/goemu/internal/app/console"
)

var buttonKeys = map[pixelgl.Button]console.Button{
	pixelgl.KeyDown:      console.ButtonDown,
	pixelgl.KeyUp:        console.ButtonUp,
	pixelgl.KeyLeft:      console.ButtonLeft,
	pixelgl.KeyRight:     console.ButtonRight,
	pixelgl.KeyEnter:     console.ButtonStart,
	pixelgl.KeyBackspace: console.ButtonSelect,
	pixelgl.KeyZ:         console.ButtonA,
	pixelgl.KeyX:         console.ButtonB,
}

var functionKeys = map[pixelgl.Button]func(*IO){
	pixelgl.KeyEscape: func(io *IO) {
		io.paused = !io.paused
	},
	pixelgl.KeyP: func(io *IO) {
		io.paused = !io.paused
	},
	pixelgl.KeyM: func(io *IO) {
		if io.muted {
			io.unmute()
		} else {
			io.mute()
		}
	},
}
