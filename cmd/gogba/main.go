package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl" // I/O
	"github.com/omstrumpf/gogba/internal/app/gba"
	"github.com/omstrumpf/gogba/internal/app/io"
)

// Runs a temporary version of the GBA emulator. Will have a global entrypoint later that allows selecting another backend.
func main() {
	pixelgl.Run(_main)
}

func _main() {
	fmt.Println("Welcome to GoGBA!")

	gameboy := gba.NewGBA()
	io := io.NewIO(gameboy)

	ticker := time.NewTicker(gameboy.GetFrameTime())

	// Game loop
	for range ticker.C {
		if io.ShouldExit() {
			return
		}

		io.ProcessInput()

		gameboy.Tick()

		io.Render()
	}
}
