package audio

type channel struct {
	source        signalSource
	lengthCounter *lengthCounter
	volume        volumeUnit
	dac           *dac

	enabled bool
}

func newChannel(maxLength int, source signalSource, volume volumeUnit) *channel {
	c := &channel{
		source: source,
		volume: volume,
		dac:    newDAC(),
	}

	c.lengthCounter = newLengthCounter(c, maxLength)

	return c
}

func (c *channel) runForClocks(clocks int) {
	c.source.runForClocks(clocks)
	c.lengthCounter.runForClocks(clocks)
	c.volume.runForClocks(clocks)
}

func (c *channel) sample() float64 {
	if !c.enabled { // Channel is disabled, output 0
		return 0
	}

	// digital output from source is multiplied by volume unit, and converted to analog voltage.

	src := c.source.sample()

	volume := c.volume.sample()

	return c.dac.convert(byte(float64(src) * volume))
}

func (c *channel) trigger() {
	if c.dac.enabled {
		c.enable()
	}

	c.lengthCounter.trigger()
	c.volume.trigger()
	c.source.trigger()
}

func (c *channel) disableDAC() {
	c.enabled = false
	c.dac.enabled = false
}

func (c *channel) enableDAC() {
	c.dac.enabled = true
}

func (c *channel) disable() {
	c.enabled = false
}

func (c *channel) enable() {
	c.enabled = true
}

func (c *channel) active() bool {
	return c.enabled && c.dac.enabled
}
