package audio

type lengthCounter struct {
	timer *timer

	counter byte

	enabled bool
}

func newLengthCounter() *lengthCounter {
	lc := lengthCounter{}

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

func (lc *lengthCounter) getCounter() byte {
	return lc.counter
}

func (lc *lengthCounter) setCounter(val byte) {
	lc.counter = val
}

func (lc *lengthCounter) getEnabled() bool {
	return lc.enabled
}

func (lc *lengthCounter) setEnabled(val bool) {
	lc.enabled = val
}
