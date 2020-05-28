package audio

type channel struct {
	squareWave     *squareWave
	lengthCounter  *lengthCounter
	volumeEnvelope *envelope
}

func newChannel() *channel {
	c := channel{}

	c.squareWave = newSquareWave()
	c.lengthCounter = newLengthCounter(64) // TODO this should be 256 for channel 3
	c.volumeEnvelope = newEnvelope()

	return &c
}

func (c *channel) runForClocks(clocks int) {
	c.squareWave.runForClocks(clocks)
	c.lengthCounter.runForClocks(clocks)
}

func (c *channel) sample() float64 {
	if !c.enabled() { // Channel is disabled, output 0
		return 0
	}

	// 1-bit output from wave source is multiplied by the 4-bit volume, and converted to analog voltage.

	bit := c.squareWave.sample()

	volume := c.volumeEnvelope.sample()

	return dac(bit * volume)
}

func (c *channel) enabled() bool {
	return c.lengthCounter.channelEnabled() && true // TODO master volume control switch
}

func (c *channel) trigger() {
	c.lengthCounter.trigger()
	c.volumeEnvelope.trigger()
	c.squareWave.trigger()
}
