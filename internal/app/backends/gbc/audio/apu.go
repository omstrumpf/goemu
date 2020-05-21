package audio

import (
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Bitrate is the number of samples output per second
const Bitrate int = 44100

// BufferLength is the number of samples to buffer in the APU's output channel
const BufferLength int = Bitrate // 1 second worth of buffer

// APU is the gameboy's Audio Processing Unit
type APU struct {
	outchan chan console.AudioSample

	dropping bool
}

// NewAPU constructs a valid APU struct
func NewAPU() *APU {
	apu := new(APU)

	apu.outchan = make(chan console.AudioSample, BufferLength)

	return apu
}

// RunForClocks runs the APU for the given number of clock cycles
func (apu *APU) RunForClocks(clocks int) {
	// TODO
}

// GetOutputChannel returns the channel that the APU writes to.
func (apu *APU) GetOutputChannel() chan console.AudioSample {
	return apu.outchan
}

func (apu *APU) enqueueSample(val float64) {
	sample := console.MonoSample(val)

	select {
	case apu.outchan <- sample:
		log.Tracef("APU produced audio sample: %f", val)
		apu.dropping = false
	default:
		if !apu.dropping {
			log.Warningf("APU output buffer full. Dropping samples.")
		} else {
			apu.dropping = true
		}
	}
}
