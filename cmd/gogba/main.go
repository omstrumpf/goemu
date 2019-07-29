package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl" // I/O
	"github.com/omstrumpf/goemu/internal/app/gbc"
	"github.com/omstrumpf/goemu/internal/app/io"
)

// Runs a temporary version of the GBC emulator. Will have a global entrypoint later that allows selecting another backend.
func main() {
	pixelgl.Run(_main)
}

func _main() {
	fmt.Println("Welcome to gomu!")

	gameboy := gbc.NewGBC()
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
