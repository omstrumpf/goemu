package audio

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
	return (byte(nw.lsfr) ^ 0xFF) & 1
}

func (nw *noiseWave) trigger() {
	nw.lsfr = 0xFFFF
	nw.updateDivisor(nw.divisorCode)
	nw.countdown = 1 << nw.clockShift
}

func (nw *noiseWave) updateDivisor(code byte) {
	nw.timer.resetDuration(int(code) + 1)
}
