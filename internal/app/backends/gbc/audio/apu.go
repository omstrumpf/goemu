package audio

import (
	"github.com/omstrumpf/goemu/internal/app/console/audio"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Bitrate is the number of samples output per second
// const Bitrate int = 44100
const Bitrate int = 43690

const bufferLength int = Bitrate // 1 second worth of buffer

// APU is the gameboy's Audio Processing Unit
type APU struct {
	outchan chan audio.ChanneledSample

	sampleTimer *timer

	squareWave1   *squareWave
	squareWave2   *squareWave
	dataWave      *dataWave
	noiseWave     *noiseWave
	envelope1     *envelope
	envelope2     *envelope
	envelope4     *envelope
	volumeShifter *volumeShifter
	sweep         *sweep

	volumeLeft   byte
	volumeRight  byte
	outputSelect byte
	enabled      bool

	channel1 *channel
	channel2 *channel
	channel3 *channel
	channel4 *channel

	lastWrites map[uint16]byte
}

// NewAPU constructs a valid APU struct
func NewAPU(speedFactor float64) *APU {
	apu := new(APU)

	apu.outchan = make(chan audio.ChanneledSample, bufferLength)

	apu.sampleTimer = newTimerByHz(int(float64(Bitrate)/speedFactor), apu.takeSample)

	apu.squareWave1 = newSquareWave()
	apu.squareWave2 = newSquareWave()
	apu.dataWave = newDataWave()
	apu.noiseWave = newNoiseWave()
	apu.envelope1 = newEnvelope()
	apu.envelope2 = newEnvelope()
	apu.envelope4 = newEnvelope()
	apu.volumeShifter = newVolumeShifter()

	apu.channel1 = newChannel(64, apu.squareWave1, apu.envelope1)
	apu.channel2 = newChannel(64, apu.squareWave2, apu.envelope2)
	apu.channel3 = newChannel(256, apu.dataWave, apu.volumeShifter)
	apu.channel4 = newChannel(64, apu.noiseWave, apu.envelope4)

	apu.sweep = newSweep(apu.channel1, apu.squareWave1)

	apu.lastWrites = make(map[uint16]byte)

	apu.initDefaults()

	return apu
}

// RunForClocks runs the APU for the given number of clock cycles
func (apu *APU) RunForClocks(clocks int) {
	if apu.enabled {
		apu.channel1.runForClocks(clocks)
		apu.channel2.runForClocks(clocks)
		apu.channel3.runForClocks(clocks)
		apu.channel4.runForClocks(clocks)
		apu.sweep.runForClocks(clocks)
	}

	apu.sampleTimer.runForClocks(clocks)
}

// GetOutputChannel returns the channel that the APU writes to.
func (apu *APU) GetOutputChannel() *chan audio.ChanneledSample {
	return &apu.outchan
}

func (apu *APU) takeSample() {
	sample := audio.ChanneledSample{Channels: []audio.Sample{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}}

	if !apu.enabled {
		apu.enqueueSample(sample)
		return
	}

	volumeLeft := float64(apu.volumeLeft+1) / 64
	volumeRight := float64(apu.volumeRight+1) / 64

	if apu.outputSelect&0b1000_0000 > 0 {
		sample.Channels[3][0] += apu.channel4.sample() * volumeLeft
	}
	if apu.outputSelect&0b0100_0000 > 0 {
		sample.Channels[2][0] += apu.channel3.sample() * volumeLeft
	}
	if apu.outputSelect&0b0010_0000 > 0 {
		sample.Channels[1][0] += apu.channel2.sample() * volumeLeft
	}
	if apu.outputSelect&0b0001_0000 > 0 {
		sample.Channels[0][0] += apu.channel1.sample() * volumeLeft
	}
	if apu.outputSelect&0b0000_1000 > 0 {
		sample.Channels[3][1] += apu.channel4.sample() * volumeRight
	}
	if apu.outputSelect&0b0000_0100 > 0 {
		sample.Channels[2][1] += apu.channel3.sample() * volumeRight
	}
	if apu.outputSelect&0b0000_0010 > 0 {
		sample.Channels[1][1] += apu.channel2.sample() * volumeRight
	}
	if apu.outputSelect&0b0000_0001 > 0 {
		sample.Channels[0][1] += apu.channel1.sample() * volumeRight
	}

	apu.enqueueSample(sample)
}

func (apu *APU) enqueueSample(sample audio.ChanneledSample) {
	select {
	case apu.outchan <- sample:
		log.Tracef("APU produced audio sample: %v", sample)
	default:
		log.Warningf("APU output buffer full. Dropping samples.")
	}
}

func (apu *APU) enable() {
	apu.enabled = true
}

func (apu *APU) disable() {
	// Write 0 to all APU registers
	for addr := uint16(0xFF10); addr < 0xFF26; addr++ {
		apu.Write(addr, 0)
	}

	apu.channel1.disableDAC()
	apu.channel2.disableDAC()
	apu.channel3.disableDAC()
	apu.channel4.disableDAC()

	apu.enabled = false
}

func (apu *APU) Read(addr uint16) byte {
	switch addr {
	case 0xFF10:
		return apu.lastWrites[addr] | 0x80
	case 0xFF11:
		return apu.lastWrites[addr] | 0x3F
	case 0xFF12:
		return apu.lastWrites[addr] | 0x00
	case 0xFF13:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF14:
		return apu.lastWrites[addr] | 0xBF
	case 0xFF15:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF16:
		return apu.lastWrites[addr] | 0x3F
	case 0xFF17:
		return apu.lastWrites[addr] | 0x00
	case 0xFF18:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF19:
		return apu.lastWrites[addr] | 0xBF
	case 0xFF1A:
		return apu.lastWrites[addr] | 0x7F
	case 0xFF1B:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF1C:
		return apu.lastWrites[addr] | 0x9F
	case 0xFF1D:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF1E:
		return apu.lastWrites[addr] | 0xBF
	case 0xFF1F:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF20:
		return apu.lastWrites[addr] | 0xFF
	case 0xFF21:
		return apu.lastWrites[addr] | 0x00
	case 0xFF22:
		return apu.lastWrites[addr] | 0x00
	case 0xFF23:
		return apu.lastWrites[addr] | 0xBF
	case 0xFF24:
		return apu.lastWrites[addr] | 0x00
	case 0xFF25:
		return apu.lastWrites[addr] | 0x00
	case 0xFF26:
		if !apu.enabled {
			return 0b0111_0000
		}

		ret := byte(0b1111_0000)

		if apu.channel4.active() {
			ret |= 0b0000_1000
		}
		if apu.channel3.active() {
			ret |= 0b0000_0100
		}
		if apu.channel2.active() {
			ret |= 0b0000_0010
		}
		if apu.channel1.active() {
			ret |= 0b0000_0001
		}

		return ret
	case 0xFF27, 0xFF28, 0xFF29, 0xFF2A, 0xFF2B, 0xFF2C, 0xFF2D, 0xFF2E, 0xFF2F:
		return 0xFF
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37,
		0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		// Wave data
		tableOffset := (addr - 0xFF30) * 2
		return (apu.dataWave.data[tableOffset] << 4) | (apu.dataWave.data[tableOffset+1])
	}

	log.Warningf("Encountered read with unknown APU control address: %#04x", addr)
	return 0xFF
}

func (apu *APU) Write(addr uint16, val byte) {
	if !apu.enabled && addr < 0xFF26 {
		return
	}

	apu.lastWrites[addr] = val

	switch addr {
	case 0xFF10: // CH1 Sweep
		apu.sweep.period = (val & 0b0111_0000) >> 4
		apu.sweep.negate = (val & 0b0000_1000) != 0
		apu.sweep.shift = (val & 0b0000_0111)
	case 0xFF11: // CH1 Duty and Length Counter
		apu.squareWave1.duty = (val & 0b1100_0000) >> 6
		apu.channel1.lengthCounter.updateCounter(val & 0b0011_1111)
	case 0xFF12: // CH1 Envelope
		apu.envelope1.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope1.mode = (val & 0b0000_1000) != 0
		apu.envelope1.sweepPeriod = (val & 0b0000_0111)
		if val&0b1111_1000 != 0 {
			apu.channel1.enableDAC()
		} else {
			apu.channel1.disableDAC()
		}
	case 0xFF13: // CH1 Frequency LSB
		apu.squareWave1.frequency = (apu.squareWave1.frequency & 0x700) | uint16(val)
		apu.squareWave1.updateFrequency()
	case 0xFF14: // CH1 Trigger, Length Enable, Frequency MSB
		trigger := val&0b1000_0000 != 0
		if (val & 0b0100_0000) != 0 {
			apu.channel1.lengthCounter.enable(trigger)
		} else {
			apu.channel1.lengthCounter.disable()
		}
		apu.squareWave1.frequency = (uint16(val&0x7) << 8) | (apu.squareWave1.frequency & 0xFF)
		apu.squareWave1.updateFrequency()
		if trigger {
			log.Tracef("APU CH1 Trigger. Frequency: %d", apu.squareWave1.frequency)
			apu.channel1.trigger()
			apu.sweep.trigger()
		}
	case 0xFF15: // Unused
	case 0xFF16: // CH2 Duty and Length Counter
		apu.squareWave2.duty = (val & 0b1100_0000) >> 6
		apu.channel2.lengthCounter.updateCounter(val & 0b0011_1111)
	case 0xFF17: // CH2 Envelope
		apu.envelope2.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope2.mode = (val & 0b0000_1000) != 0
		apu.envelope2.sweepPeriod = (val & 0b0000_0111)
		if val&0b1111_1000 != 0 {
			apu.channel2.enableDAC()
		} else {
			apu.channel2.disableDAC()
		}
	case 0xFF18: // CH2 Frequency LSB
		apu.squareWave2.frequency = (apu.squareWave2.frequency & 0x700) | uint16(val)
		apu.squareWave2.updateFrequency()
	case 0xFF19: // CH2 Trigger, Length Enable, Frequency MSB
		trigger := val&0b1000_0000 != 0
		if (val & 0b0100_0000) != 0 {
			apu.channel2.lengthCounter.enable(trigger)
		} else {
			apu.channel2.lengthCounter.disable()
		}
		apu.squareWave2.frequency = (uint16(val&0x7) << 8) | (apu.squareWave2.frequency & 0xFF)
		apu.squareWave2.updateFrequency()
		if trigger {
			log.Tracef("APU CH2 Trigger. Frequency: %d", apu.squareWave2.frequency)
			apu.channel2.trigger()
		}
	case 0xFF1A: // CH3 DAC power
		if val&0b1000_0000 != 0 {
			apu.channel3.enableDAC()
		} else {
			apu.channel3.disableDAC()
		}
	case 0xFF1B: // CH3 length counter
		apu.channel3.lengthCounter.updateCounter(val)
	case 0xFF1C: // CH3 volume code
		apu.volumeShifter.volumeCode = (val & 0b0110_0000) >> 5
	case 0xFF1D: // CH3 frequency LSB
		apu.dataWave.frequency = (apu.dataWave.frequency & 0x700) | uint16(val)
		apu.dataWave.updateFrequency()
	case 0xFF1E: // CH3 trigger, length enable, frequency MSB
		trigger := val&0b1000_0000 != 0
		if (val & 0b0100_0000) != 0 {
			apu.channel3.lengthCounter.enable(trigger)
		} else {
			apu.channel3.lengthCounter.disable()
		}
		apu.dataWave.frequency = (uint16(val&0x7) << 8) | (apu.dataWave.frequency & 0xFF)
		apu.dataWave.updateFrequency()
		if trigger {
			log.Tracef("APU CH3 Trigger. Frequency: %d", apu.dataWave.frequency)
			apu.channel3.trigger()
		}
	case 0xFF1F: // Unused
	case 0xFF20: // CH4 length counter
		apu.channel4.lengthCounter.updateCounter(val & 0b0011_1111)
	case 0xFF21: // CH4 envelope
		apu.envelope4.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope4.mode = (val & 0b0000_1000) != 0
		apu.envelope4.sweepPeriod = (val & 0b0000_0111)
		if val&0b1111_1000 != 0 {
			apu.channel4.enableDAC()
		} else {
			apu.channel4.disableDAC()
		}
	case 0xFF22: // CH4 clock shift, width, divisor code
		apu.noiseWave.clockShift = (val & 0b1111_0000) >> 4
		apu.noiseWave.widthMode = (val & 0b0000_1000) != 0
		apu.noiseWave.divisorCode = (val & 0b0000_0111)
		apu.noiseWave.updateDivisor(apu.noiseWave.divisorCode)
	case 0xFF23: // CH4 trigger, length enable
		trigger := val&0b1000_0000 != 0
		if (val & 0b0100_0000) != 0 {
			apu.channel4.lengthCounter.enable(trigger)
		} else {
			apu.channel4.lengthCounter.disable()
		}
		if trigger {
			log.Tracef("APU CH4 Trigger.")
			apu.channel4.trigger()
		}
	case 0xFF24: // Controls
		apu.volumeLeft = (val & 0b0111_0000) >> 4
		apu.volumeRight = (val & 0b0000_0111)
	case 0xFF25: // Controls
		apu.outputSelect = val
	case 0xFF26: // Controls
		if (val & 0b1000_0000) != 0 {
			apu.enable()
		} else {
			apu.disable()
		}
	case 0xFF27, 0xFF28, 0xFF29, 0xFF2A, 0xFF2B, 0xFF2C, 0xFF2D, 0xFF2E, 0xFF2F:
		// Unused
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37,
		0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		// Wave data
		tableOffset := (addr - 0xFF30) * 2
		apu.dataWave.data[tableOffset] = (val & 0xF0) >> 4
		apu.dataWave.data[tableOffset+1] = val & 0x0F
	default:
		log.Warningf("Encountered write with unknown APU control address: %#4x", addr)
	}
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
