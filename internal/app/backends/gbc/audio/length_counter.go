package audio

import "github.com/omstrumpf/goemu/internal/app/log"

type lengthCounter struct {
	timer *timer

	volumeCounter byte

	enabled bool
}

func newLengthCounter() *lengthCounter {
	lc := lengthCounter{}

	lc.timer = newTimer(4096, lc.tick)

	return &lc
}

func (lc *lengthCounter) runForClocks(clocks int) {
	lc.timer.runForClocks(clocks)
}

func (lc *lengthCounter) tick() {
	if lc.enabled && lc.volumeCounter > 0 {
		lc.volumeCounter--
		if lc.volumeCounter == 0 {
			log.Debugf("LC1 hit zero")
		}
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
	log.Debugf("Setting LC1 enabled: %t", val)
	lc.enabled = val
}
