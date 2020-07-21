package audio

type lengthCounter struct {
	timer *timer

	initCounter byte
	counter     byte

	enabled bool
}

func newLengthCounter(maxLength byte) *lengthCounter {
	lc := lengthCounter{initCounter: maxLength}

	lc.timer = newTimerByHz(256, lc.tick)

	return &lc
}

func (lc *lengthCounter) runForClocks(clocks int) {
	lc.timer.runForClocks(clocks)
}

func (lc *lengthCounter) tick() {
	if lc.enabled && lc.counter > 0 {
		lc.counter--
	}
}

func (lc *lengthCounter) channelEnabled() bool {
	return lc.counter > 0
}

func (lc *lengthCounter) trigger() {
	if lc.counter == 0 {
		lc.counter = lc.initCounter
	}
	lc.enabled = true
}

func (lc *lengthCounter) updateCounter(val byte) {
	lc.counter = lc.initCounter - val + 1
}
