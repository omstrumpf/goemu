package audio

import (
	"github.com/omstrumpf/goemu/internal/app/console"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// Bitrate is the number of samples output per second
// const Bitrate int = 44100
const Bitrate int = 43690

const bufferLength int = Bitrate // 1 second worth of buffer

// APU is the gameboy's Audio Processing Unit
type APU struct {
	outchan chan console.AudioSample

	sampleTimer *timer

	squareWave1    *squareWave
	squareWave2    *squareWave
	dataWave       *dataWave
	noiseWave      *noiseWave
	lengthCounter1 *lengthCounter
	lengthCounter2 *lengthCounter
	lengthCounter3 *lengthCounter
	lengthCounter4 *lengthCounter
	envelope1      *envelope
	envelope2      *envelope
	envelope4      *envelope
	volumeShifter  *volumeShifter
	sweep          *sweep

	channel1 *channel
	channel2 *channel
	channel3 *channel
	channel4 *channel
}

// NewAPU constructs a valid APU struct
func NewAPU() *APU {
	apu := new(APU)

	apu.outchan = make(chan console.AudioSample, bufferLength)

	apu.sampleTimer = newTimerByHz(Bitrate, apu.takeSample)

	apu.squareWave1 = newSquareWave()
	apu.sweep = newSweep(apu.squareWave1)
	apu.squareWave2 = newSquareWave()
	apu.dataWave = newDataWave()
	apu.noiseWave = newNoiseWave()
	apu.lengthCounter1 = newLengthCounter(0x3F)
	apu.lengthCounter2 = newLengthCounter(0x3F)
	apu.lengthCounter3 = newLengthCounter(0xFF)
	apu.lengthCounter4 = newLengthCounter(0x3F)
	apu.envelope1 = newEnvelope()
	apu.envelope2 = newEnvelope()
	apu.envelope4 = newEnvelope()
	apu.volumeShifter = newVolumeShifter()

	apu.channel1 = newChannel(apu.squareWave1, apu.lengthCounter1, apu.envelope1)
	apu.channel2 = newChannel(apu.squareWave2, apu.lengthCounter2, apu.envelope2)
	apu.channel3 = newChannel(apu.dataWave, apu.lengthCounter3, apu.volumeShifter)
	apu.channel4 = newChannel(apu.noiseWave, apu.lengthCounter4, apu.envelope4)

	apu.initDefaults()

	return apu
}

// RunForClocks runs the APU for the given number of clock cycles
func (apu *APU) RunForClocks(clocks int) {
	apu.channel1.runForClocks(clocks)
	apu.channel2.runForClocks(clocks)
	apu.channel3.runForClocks(clocks)
	apu.channel4.runForClocks(clocks)
	apu.sweep.runForClocks(clocks)
	apu.sampleTimer.runForClocks(clocks)
}

// GetOutputChannel returns the channel that the APU writes to.
func (apu *APU) GetOutputChannel() *chan console.AudioSample {
	return &apu.outchan
}

func (apu *APU) takeSample() {
	s := float64(0)

	s += apu.channel1.sample()
	s += apu.channel2.sample()
	s += apu.channel3.sample()
	s += apu.channel4.sample()

	apu.enqueueSample(s, s)
}

func (apu *APU) enqueueSample(l float64, r float64) {
	sample := console.StereoSample(l, r)

	select {
	case apu.outchan <- sample:
		log.Tracef("APU produced audio sample: %v", sample)
	default:
		log.Warningf("APU output buffer full. Dropping samples.")
	}
}

func (apu *APU) Read(addr uint16) byte {

	log.Warningf("Encountered read with unknown APU control address: %#04x", addr)
	return 0xFF
}

func (apu *APU) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF10: // CH1 Sweep
		apu.sweep.period = (val & 0b0111_0000) >> 4
		apu.sweep.negate = (val & 0b0000_1000) != 0
		apu.sweep.shift = (val & 0b0000_0111)
	case 0xFF11: // CH1 Duty and Length Counter
		apu.squareWave1.duty = (val & 0b1100_0000) >> 6
		apu.lengthCounter1.counter = (val & 0b0011_1111)
	case 0xFF12: // CH1 Envelope
		apu.envelope1.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope1.mode = (val & 0b0000_1000) != 0
		apu.envelope1.sweepPeriod = (val & 0b0000_0111)
		apu.channel1.dac.enabled = (val&0b1111_1000 != 0)
	case 0xFF13: // CH1 Frequency LSB
		apu.squareWave1.frequency = (apu.squareWave1.frequency & 0x700) | uint32(val)
	case 0xFF14: // CH1 Trigger, Length Enable, Frequency MSB
		apu.lengthCounter1.enabled = (val & (1 << 6)) != 0
		apu.squareWave1.frequency = (uint32(val&0x7) << 8) | (apu.squareWave1.frequency & 0xFF)
		if val&(1<<7) != 0 {
			log.Tracef("APU CH1 Trigger. Frequency: %d", apu.squareWave1.frequency)
			apu.sweep.trigger()
			apu.channel1.trigger()
		}
	case 0xFF15: // Unused
	case 0xFF16: // CH2 Duty and Length Counter
		apu.squareWave2.duty = (val & 0b1100_0000) >> 6
		apu.lengthCounter2.counter = (val & 0b0011_1111)
	case 0xFF17: // CH2 Envelope
		apu.envelope2.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope2.mode = (val & 0b0000_1000) != 0
		apu.envelope2.sweepPeriod = (val & 0b0000_0111)
		apu.channel2.dac.enabled = (val&0b1111_1000 != 0)
	case 0xFF18: // CH2 Frequency LSB
		apu.squareWave2.frequency = (apu.squareWave2.frequency & 0x700) | uint32(val)
	case 0xFF19: // CH2 Trigger, Length Enable, Frequency MSB
		apu.lengthCounter2.enabled = (val & (1 << 6)) != 0
		apu.squareWave2.frequency = (uint32(val&0x7) << 8) | (apu.squareWave2.frequency & 0xFF)
		if val&(1<<7) != 0 {
			log.Tracef("APU CH2 Trigger. Frequency: %d", apu.squareWave2.frequency)
			apu.channel2.trigger()
		}
	case 0xFF1A: // CH3 DAC power
		apu.channel3.dac.enabled = (val&0b1000_0000 != 0)
	case 0xFF1B: // CH3 length counter
		apu.lengthCounter3.counter = val
	case 0xFF1C: // CH3 volume code
		apu.volumeShifter.volumeCode = (val & 0b0110_0000) >> 5
	case 0xFF1D: // CH3 frequency LSB
		apu.dataWave.frequency = (apu.dataWave.frequency & 0x700) | uint32(val)
	case 0xFF1E: // CH3 trigger, length enable, frequency MSB
		apu.lengthCounter3.enabled = (val & (1 << 6)) != 0
		apu.dataWave.frequency = (uint32(val&0x7) << 8) | (apu.dataWave.frequency & 0xFF)
		if val&(1<<7) != 0 {
			log.Tracef("APU CH3 Trigger. Frequency: %d", apu.dataWave.frequency)
			apu.channel3.trigger()
		}
	case 0xFF1F: // Unused
	case 0xFF20: // CH4 length counter
		apu.lengthCounter4.counter = (val & 0b0011_1111)
	case 0xFF21: // CH4 envelope
		apu.envelope4.initVolume = (val & 0b1111_0000) >> 4
		apu.envelope4.mode = (val & 0b0000_1000) != 0
		apu.envelope4.sweepPeriod = (val & 0b0000_0111)
		apu.channel4.dac.enabled = (val&0b1111_1000 != 0)
	case 0xFF22: // CH4 clock shift, width, divisor code
		apu.noiseWave.clockShift = (val & 0b1111_0000) >> 4
		apu.noiseWave.widthMode = (val & 0b0000_1000) != 0
		apu.noiseWave.divisorCode = (val & 0b0000_0111)
		apu.noiseWave.updateDivisor(apu.noiseWave.divisorCode)
	case 0xFF23: // CH4 trigger, length enable
		apu.lengthCounter4.enabled = (val & (1 << 6)) != 0
		if (val & (1 << 7)) != 0 {
			log.Tracef("APU CH4 Trigger.")
			apu.channel4.trigger()
		}
	case 0xFF24: // Controls
	case 0xFF25: // Controls
	case 0xFF26: // Controls
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37,
		0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		// Wave data
		tableOffset := addr - 0xFF30
		apu.dataWave.data[tableOffset] = (val & 0xF0) >> 4
		apu.dataWave.data[tableOffset+1] = val & 0x0F
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
