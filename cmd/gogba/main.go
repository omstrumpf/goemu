package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl" // I/O
	"github.com/omstrumpf/gogba/internal/app/gba"
	"github.com/omstrumpf/gogba/internal/app/io"
)

func main() {
	pixelgl.Run(_main)
}

func _main() {
	fmt.Println("Welcome to GoGBA!")

	gameboy := gba.NewGBA()
	io := io.NewIO(gameboy)

	ticker := time.NewTicker(gba.SecondsPerFrame)

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
