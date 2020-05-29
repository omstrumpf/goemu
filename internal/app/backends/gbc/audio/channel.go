package audio

type channel struct {
	source        signalSource
	lengthCounter *lengthCounter
	volume        volumeUnit
	dac           *dac
}

func newChannel(source signalSource, lengthCounter *lengthCounter, volume volumeUnit) *channel {
	c := channel{}

	c.source = source
	c.lengthCounter = lengthCounter
	c.volume = volume
	c.dac = newDAC()

	return &c
}

func (c *channel) runForClocks(clocks int) {
	c.source.runForClocks(clocks)
	c.lengthCounter.runForClocks(clocks)
}

func (c *channel) sample() float64 {
	if !c.enabled() { // Channel is disabled, output 0
		return 0
	}

	// digital output from source is multiplied by volume unit, and converted to analog voltage.

	src := c.source.sample()

	volume := c.volume.sample()

	return c.dac.convert(byte(float64(src) * volume))
}

func (c *channel) enabled() bool {
	return c.lengthCounter.channelEnabled() && true // TODO master volume control switch
}

func (c *channel) trigger() {
	c.lengthCounter.trigger()
	c.volume.trigger()
	c.source.trigger()
}
