package audio

type dataWave struct {
	timer *timer

	data            [32]byte
	positionCounter byte

	frequency uint32

	value byte // Current value of the wave
}

func newDataWave() *dataWave {
	dw := dataWave{}

	dw.timer = newTimerByClocks(1, dw.tick)

	return &dw
}

func (dw *dataWave) runForClocks(clocks int) {
	dw.timer.runForClocks(clocks)
}

func (dw *dataWave) tick() {
	dw.positionCounter = (dw.positionCounter + 1) % 32

	dw.value = dw.data[dw.positionCounter]
}

func (dw *dataWave) sample() byte {
	return dw.value
}

func (dw *dataWave) trigger() {
	dw.positionCounter = 0
	dw.updateFrequency(dw.frequency)
}

func (dw *dataWave) updateFrequency(freq uint32) {
	dw.frequency = freq
	dw.timer.resetDuration((2048 - int(freq)) / 2)
}
