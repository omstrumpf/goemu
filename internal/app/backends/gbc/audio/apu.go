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

	apu.sampleTimer = newTimer(24, apu.takeSample) // TODO all these magic numbers should be calculated based on sample rate / GBC clock speed

	apu.channel1 = newChannel()
	apu.channel2 = newChannel()

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
	switch addr {
	case 0xFF10:
	case 0xFF11:
	case 0xFF12:
	case 0xFF13:
		return byte(apu.channel1.squareWave.frequency)
	case 0xFF14:
		triggerBit := byte(1 << 7)
		lengthEnableBit := byte(0)
		if apu.channel1.lengthCounter.getEnabled() {
			lengthEnableBit = 1 << 6
		}
		frequencyBits := byte((apu.channel1.squareWave.frequency >> 8) & 0b0000_0111)
		emptyBits := byte(0b0011_1000)
		return triggerBit | lengthEnableBit | emptyBits | frequencyBits
	}

	// log.Warningf("Encountered read with unknown APU control address: %#04x", addr)
	return 0xFF
}

func (apu *APU) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF10:
	case 0xFF11:
		apu.channel1.squareWave.duty = (val & 0b1100_0000) >> 6
		apu.channel1.lengthCounter.setCounter(val & 0b0011_0000)
	case 0xFF12:
	case 0xFF13:
		apu.channel1.squareWave.updateFrequency((apu.channel1.squareWave.frequency & 0x700) | uint32(val))
	case 0xFF14:
		if val&(1<<7) != 0 {
			if apu.channel1.lengthCounter.volumeCounter == 0 {
				apu.channel1.lengthCounter.setCounter(64)
			}
			apu.channel1.lengthCounter.setEnabled(true)
			apu.channel1.squareWave.updateFrequency(apu.channel1.squareWave.frequency)

			// TODO trigger
		}
		apu.channel1.squareWave.updateFrequency((uint32(val&0x7) << 8) | (apu.channel1.squareWave.frequency & 0xFF))
	case 0xFF15:
	case 0xFF16:
		apu.channel2.squareWave.duty = (val & 0b1100_0000) >> 6
		apu.channel2.lengthCounter.setCounter(val & 0b0011_0000)
	case 0xFF17:
	case 0xFF18:
		apu.channel2.squareWave.updateFrequency((apu.channel2.squareWave.frequency & 0x700) | uint32(val))
	case 0xFF19:
		if val&(1<<7) != 0 {
			if apu.channel2.lengthCounter.volumeCounter == 0 {
				apu.channel2.lengthCounter.setCounter(64)
			}
			apu.channel2.lengthCounter.setEnabled(true)
			apu.channel2.squareWave.updateFrequency(apu.channel2.squareWave.frequency)

			// TODO trigger
		}
		apu.channel2.squareWave.updateFrequency((uint32(val&0x7) << 8) | (apu.channel2.squareWave.frequency & 0xFF))
	}

	// log.Warningf("Encountered write with unknown APU control address: %#4x", addr)
}
