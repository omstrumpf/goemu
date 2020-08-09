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

func (lc *lengthCounter) enable(trigger bool) {
	// Extra clock when enabled in the first half of timer period
	// may disable channel
	bonusTick := false
	if lc.timer.countdown > (lc.timer.period / 2) {
		bonusTick = true

		if lc.counter != 0 && !lc.enabled {
			lc.enabled = true
			lc.tick()
		}
	}

	lc.enabled = true

	// If LC is now 0, and channel was triggered, should clock again
	if trigger && lc.counter == 0 && bonusTick {
		lc.counter = lc.initCounter - 1
	}

}

func (lc *lengthCounter) disable() {
	lc.enabled = false
}
