package audio

type channel struct {
	squareWave    *squareWave
	lengthCounter *lengthCounter
}

func newChannel() *channel {
	c := channel{}

	c.squareWave = newSquareWave()
	c.lengthCounter = newLengthCounter()

	return &c
}

func (c *channel) runForClocks(clocks int) {
	c.squareWave.runForClocks(clocks)
	c.lengthCounter.runForClocks(clocks)
}

func (c *channel) sample() float64 {
	if c.lengthCounter.channelEnabled() {
		return c.squareWave.sample()
	}
	return 0
}
