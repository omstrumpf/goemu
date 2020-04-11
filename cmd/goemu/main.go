package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/faiface/pixel/pixelgl" // I/O
	"github.com/omstrumpf/goemu/internal/app/backends/gbc"
	"github.com/omstrumpf/goemu/internal/app/io"
)

// Runs a temporary version of the GBC emulator. Will have a global entrypoint later that allows selecting another backend.
func main() {
	pixelgl.Run(_main)
}

func _main() {
	fmt.Println("Welcome to goemu!")

	gameboy := gbc.NewGBC()

	buf, err := ioutil.ReadFile("roms/tetris.gb")
	if err != nil {
		panic(err)
	}
	gameboy.LoadROM(buf)

	io := io.NewIO(gameboy)

	ticker := time.NewTicker(gameboy.GetFrameTime())

	// Game loop
	for range ticker.C {
		if io.ShouldExit() || gameboy.IsStopped() {
			return
		}

		io.ProcessInput()

		if io.ShouldEmulate() {
			gameboy.Tick()
		}

		io.Render()
	}
}
