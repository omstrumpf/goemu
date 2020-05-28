package audio

import (
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Bitrate is the number of samples output per second
const Bitrate int = 44100

const bufferLength int = Bitrate // 1 second worth of buffer

// APU is the gameboy's Audio Processing Unit
type APU struct {
	outchan chan console.AudioSample

	sampleTimer *timer

	channel1 *channel
	channel2 *channel

	dropping bool
}

// NewAPU constructs a valid APU struct
func NewAPU() *APU {
	apu := new(APU)

	apu.outchan = make(chan console.AudioSample, bufferLength)

	apu.sampleTimer = newTimerByHz(Bitrate, apu.takeSample)

	apu.channel1 = newChannel()
	apu.channel2 = newChannel()

	apu.initDefaults()

	return apu
}

// RunForClocks runs the APU for the given number of clock cycles
func (apu *APU) RunForClocks(clocks int) {
	apu.channel1.runForClocks(clocks)
	apu.channel2.runForClocks(clocks)
	apu.sampleTimer.runForClocks(clocks)

	// TODO combine the lengthCounter, envelope, etc. timers into one frame sequencer
}

// GetOutputChannel returns the channel that the APU writes to.
func (apu *APU) GetOutputChannel() *chan console.AudioSample {
	return &apu.outchan
}

func (apu *APU) takeSample() {
	s := float64(0)

	s += apu.channel1.sample()
	s += apu.channel2.sample()

	apu.enqueueSample(s, s)
}

func (apu *APU) enqueueSample(l float64, r float64) {
	sample := console.StereoSample(l, r)

	select {
	case apu.outchan <- sample:
		log.Tracef("APU produced audio sample: %v", sample)
		apu.dropping = false
	default:
		if !apu.dropping {
			log.Warningf("APU output buffer full. Dropping samples.")
		} else {
			apu.dropping = true
		}
	}
}

func (apu *APU) Read(addr uint16) byte {

	// log.Warningf("Encountered read with unknown APU control address: %#04x", addr)
	return 0xFF
}

func (apu *APU) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF10:
	case 0xFF11:
		apu.channel1.squareWave.duty = (val & 0b1100_0000) >> 6
		apu.channel1.lengthCounter.counter = (val & 0b0011_1111)
	case 0xFF12:
		apu.channel1.volumeEnvelope.initVolume = (val & 0b1111_0000) >> 4
		apu.channel1.volumeEnvelope.mode = (val & 0b0000_1000) != 0
		apu.channel1.volumeEnvelope.sweepPeriod = (val & 0b0000_0111)
	case 0xFF13:
		apu.channel1.squareWave.updateFrequency((apu.channel1.squareWave.frequency & 0x700) | uint32(val))
	case 0xFF14:
		if val&(1<<7) != 0 {
			log.Tracef("APU CH1 Trigger. Frequency: %d", apu.channel1.squareWave.frequency)
			apu.channel1.trigger()
		}
		apu.channel1.squareWave.updateFrequency((uint32(val&0x7) << 8) | (apu.channel1.squareWave.frequency & 0xFF))
	case 0xFF15:
	case 0xFF16:
		apu.channel2.squareWave.duty = (val & 0b1100_0000) >> 6
		apu.channel2.lengthCounter.counter = (val & 0b0011_1111)
	case 0xFF17:
		apu.channel2.volumeEnvelope.initVolume = (val & 0b1111_0000) >> 4
		apu.channel2.volumeEnvelope.mode = (val & 0b0000_1000) != 0
		apu.channel2.volumeEnvelope.sweepPeriod = (val & 0b0000_0111)
	case 0xFF18:
		apu.channel2.squareWave.updateFrequency((apu.channel2.squareWave.frequency & 0x700) | uint32(val))
	case 0xFF19:
		if val&(1<<7) != 0 {
			log.Tracef("APU CH2 Trigger. Frequency: %d", apu.channel2.squareWave.frequency)
			apu.channel2.trigger()
		}
		apu.channel2.squareWave.updateFrequency((uint32(val&0x7) << 8) | (apu.channel2.squareWave.frequency & 0xFF))
	}

	// log.Warningf("Encountered write with unknown APU control address: %#4x", addr)
}

func (apu *APU) resetDefaults() {
	apu.Write(0xFF10, 0x80)
	apu.Write(0xFF11, 0xBF)
	apu.Write(0xFF12, 0xF3)
	apu.Write(0xFF14, 0xBF)
	apu.Write(0xFF16, 0x3F)
	apu.Write(0xFF17, 0x00)
	apu.Write(0xFF19, 0xBF)
	apu.Write(0xFF1A, 0x7F)
	apu.Write(0xFF1B, 0xFF)
	apu.Write(0xFF1C, 0x9F)
	apu.Write(0xFF1E, 0xBF)
	apu.Write(0xFF20, 0xFF)
	apu.Write(0xFF21, 0x00)
	apu.Write(0xFF22, 0x00)
	apu.Write(0xFF23, 0xBF)
	apu.Write(0xFF24, 0x77)
	apu.Write(0xFF25, 0xF3)
	apu.Write(0xFF26, 0xF1)
}

func (apu *APU) initDefaults() {
	apu.resetDefaults()
	apu.Write(0xFF13, 0xFF)
	apu.Write(0xFF15, 0xFF)
	apu.Write(0xFF18, 0xFF)
	apu.Write(0xFF1D, 0xFF)
	apu.Write(0xFF1F, 0xFF)
}
