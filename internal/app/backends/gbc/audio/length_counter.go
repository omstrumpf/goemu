package audio

type lengthCounter struct {
	timer *timer

	volumeCounter byte

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
	if lc.enabled && lc.volumeCounter > 0 {
		lc.volumeCounter--
	}
}

func (lc *lengthCounter) channelEnabled() bool {
	return lc.volumeCounter > 0
}

func (lc *lengthCounter) getCounter() byte {
	return lc.volumeCounter
}

func (lc *lengthCounter) setCounter(val byte) {
	lc.volumeCounter = val
}

func (lc *lengthCounter) getEnabled() bool {
	return lc.enabled
}

func (lc *lengthCounter) setEnabled(val bool) {
	lc.enabled = val
}
