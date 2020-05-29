package audio

import "github.com/omstrumpf/goemu/internal/app/log"

type noiseWave struct {
	timer *timer

	divisorCode byte
	clockShift  byte
	widthMode   bool

	countdown uint16

	lsfr uint16
}

func newNoiseWave() *noiseWave {
	nw := noiseWave{}

	nw.timer = newTimerByHz(2097152, nw.tick)

	return &nw
}

func (nw *noiseWave) runForClocks(clocks int) {
	nw.timer.runForClocks(clocks)
}

func (nw *noiseWave) tick() {
	nw.countdown--

	if nw.countdown == 0 {
		nw.countdown = 1 << nw.clockShift

		nw.tickLSFR()
	}
}

func (nw *noiseWave) tickLSFR() {
	xorResult := ((nw.lsfr & 0b10) >> 1) ^ (nw.lsfr & 0b01)

	nw.lsfr >>= 1

	if nw.widthMode {
		nw.lsfr = (nw.lsfr & 0b0011_1111_1011_1111) | (xorResult << 14) | (xorResult << 5)
	} else {

		nw.lsfr = (nw.lsfr & 0b0011_1111_1111_1111) | (xorResult << 14)
	}
}

func (nw *noiseWave) sample() byte {
	return byte(nw.lsfr&1) ^ 1
}

func (nw *noiseWave) trigger() {
	nw.lsfr = 0xFFFF
}

func (nw *noiseWave) updateDivisor(code byte) {
	var duration int
	switch code {
	case 0:
		duration = 8
	case 1:
		duration = 16
	case 2:
		duration = 32
	case 3:
		duration = 48
	case 4:
		duration = 64
	case 5:
		duration = 80
	case 6:
		duration = 96
	case 7:
		duration = 112
	default:
		duration = 8
		log.Errorf("Noise wave divisor got invalid code: %#02x", code)
	}

	nw.timer.resetDuration(duration)
}
