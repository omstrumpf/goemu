package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/faiface/pixel/pixelgl" // I/O
	"github.com/juju/loggo"
	"github.com/omstrumpf/goemu/internal/app/backends/gbc"
	"github.com/omstrumpf/goemu/internal/app/io"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Runs a temporary version of the GBC emulator. Will have a global entrypoint later that allows selecting another backend.
func main() {
	pixelgl.Run(_main)
}

func _main() {
	fmt.Println("Welcome to goemu!")

	// testmode := flag.Bool("testmode", false, "Run the emulator in testmode, and produce a hash of the steady-state screen output")
	loglevel := flag.String("loglevel", "ERROR", "Log level. ERROR, WARNING, DEBUG, TRACE.")
	skiplogo := flag.Bool("skiplogo", false, "Skip the logo scroll sequence")
	fastmode := flag.Bool("fastmode", false, "Don't limit emulation speed")
	frames := flag.Uint64("frames", 0, "Number of frames to emulate. 0 is infinite.")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Please specify a romfile.")
		return
	}
	romfile := flag.Arg(0)

	switch *loglevel {
	default:
		fmt.Printf("Unsupported log level: %s.\n", *loglevel)
		fallthrough
	case "ERROR":
		loggo.ConfigureLoggers(`<root>=ERROR`)
	case "WARNING":
		loggo.ConfigureLoggers(`<root>=WARNING`)
	case "DEBUG":
		loggo.ConfigureLoggers(`<root>=DEBUG`)
	case "TRACE":
		loggo.ConfigureLoggers(`<root>=TRACE`)
	}

	log.Tracef("Loading romfile")

	buf, err := ioutil.ReadFile(romfile)
	if err != nil {
		panic(err)
	}

	log.Tracef("Initializing gameboy")

	gameboy := gbc.NewGBC(*skiplogo, buf)

	io := io.NewIO(gameboy)

	var ticker *time.Ticker
	if *fastmode {
		ticker = time.NewTicker(time.Nanosecond)
	} else {
		ticker = time.NewTicker(gameboy.GetFrameTime())
	}

	// Game loop
	frame := uint64(0)
	maxFrame := *frames - 1
	for range ticker.C {
		if maxFrame != 0 && frame > *frames-1 {
			break
		}
		log.Tracef("Emulating frame %d", frame)
		frame++

		if io.ShouldExit() {
			return
		}

		io.ProcessInput()

		if io.ShouldEmulate() {
			gameboy.Tick()
		}

		io.Render()
	}
}
