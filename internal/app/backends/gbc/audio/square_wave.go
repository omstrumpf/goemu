package audio

type squareWave struct {
	timer *timer

	duty        byte
	dutyCounter int

	frequency uint32

	value float64
}

func newSquareWave() *squareWave {
	sw := squareWave{}

	sw.timer = newTimerByClocks(1, sw.tick)

	sw.duty = 2

	return &sw
}

func (sw *squareWave) runForClocks(clocks int) {
	sw.timer.runForClocks(clocks)
}

func (sw *squareWave) tick() {
	sw.value = float64(dutyMap[sw.duty][sw.dutyCounter])

	sw.dutyCounter = (sw.dutyCounter + 1) % 8
}

func (sw *squareWave) sample() float64 {
	return sw.value
}

var dutyMap [4][8]float64 = [4][8]float64{
	{-1, -1, -1, -1, -1, -1, -1, 1},
	{1, -1, -1, -1, -1, -1, -1, 1},
	{1, -1, -1, -1, -1, 1, 1, 1},
	{-1, 1, 1, 1, 1, 1, 1, -1},
}

func (sw *squareWave) updateFrequency(freq uint32) {
	sw.frequency = freq
	sw.timer.resetDuration(2048 - int(freq))
}
