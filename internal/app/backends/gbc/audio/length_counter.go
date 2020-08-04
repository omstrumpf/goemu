package audio

type lengthCounter struct {
	channel *channel
	timer   *timer

	initCounter int
	counter     int

	enabled bool
}

func newLengthCounter(channel *channel, maxLength int) *lengthCounter {
	lc := lengthCounter{
		channel:     channel,
		initCounter: maxLength,
	}

	lc.timer = newTimerByHz(256, lc.tick)

	return &lc
}

func (lc *lengthCounter) runForClocks(clocks int) {
	lc.timer.runForClocks(clocks)
}

func (lc *lengthCounter) tick() {
	if lc.enabled && lc.counter > 0 {
		lc.counter--

		if lc.counter == 0 {
			lc.channel.enabled = false
		}
	}
}

func (lc *lengthCounter) trigger() {
	if lc.counter == 0 {
		lc.counter = lc.initCounter
	}
}

func (lc *lengthCounter) updateCounter(val byte) {
	lc.counter = lc.initCounter - int(val)
}
