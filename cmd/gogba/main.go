package main

import (
	"fmt"

	"github.com/omstrumpf/gogba/internal/app/gba"
)

func main() {
	fmt.Println("Welcome to GoGBA!")

	gameboy := gba.NewGBA()

	gameboy.Tick()
}
