package audio

type squareWave struct {
	timer *timer

	duty        byte
	dutyCounter int

	frequency uint16

	value byte // Current binary value of the wave
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
	sw.value = dutyMap[sw.duty][sw.dutyCounter]

	sw.dutyCounter = (sw.dutyCounter + 1) % 8
}

func (sw *squareWave) sample() byte {
	return sw.value
}

func (sw *squareWave) trigger() {
	sw.updateFrequency()

	// TODO:
	// Upon the channel INIT trigger bit being set for either channel 1
	// or 2, the wave position's incrementing will be delayed by 1/12 of a full cycle.
	// IT WILL *NOT* BE RESET TO 0 BY A CHANNEL INIT. 'Gauntlet II' does a very slick
	// job of timing? itself around this fact.
}

func (sw *squareWave) updateFrequency() {
	sw.timer.period = 2048 - int(sw.frequency)
}

var dutyMap [4][8]byte = [4][8]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 1, 1},
	{0, 1, 1, 1, 1, 1, 1, 0},
}
