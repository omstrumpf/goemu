package audio

import "github.com/omstrumpf/goemu/internal/app/log"

type sweep struct {
	timer *timer

	squareWave *squareWave
	channel    *channel

	period  byte
	counter byte

	shift   byte
	enabled bool
	negate  bool

	freqShadow uint16
}

func newSweep(channel *channel, squareWave *squareWave) *sweep {
	sweep := sweep{channel: channel, squareWave: squareWave}

	sweep.timer = newTimerByHz(128, sweep.tick)

	return &sweep
}

func (sweep *sweep) runForClocks(clocks int) {
	sweep.timer.runForClocks(clocks)
}

func (sweep *sweep) tick() {
	if sweep.period != 0 {
		sweep.counter--

		if sweep.counter == 0 {
			sweep.counter = sweep.period

			if sweep.enabled {
				newFreq := sweep.calculateNewFrequency()

				if !sweep.checkOverflow(newFreq) && sweep.shift != 0 {
					sweep.writeFreq(newFreq)

					sweep.checkOverflow(sweep.calculateNewFrequency())
				}
			}
		}
	}
}

func (sweep *sweep) trigger() {
	sweep.freqShadow = sweep.squareWave.frequency
	sweep.counter = sweep.period

	sweep.enabled = (sweep.period != 0 || sweep.shift != 0)

	if sweep.shift != 0 {
		sweep.checkOverflow(sweep.calculateNewFrequency())
	}
}

func (sweep *sweep) calculateNewFrequency() uint16 {
	newFreq := sweep.freqShadow >> sweep.shift

	if sweep.negate {
		newFreq = ^newFreq
	}

	newFreq += sweep.freqShadow

	return newFreq
}

func (sweep *sweep) checkOverflow(freq uint16) bool {
	if freq > 2047 {
		log.Debugf("Sweep unit disabling channel")
		sweep.channel.disable()
		return true
	}

	return false
}

func (sweep *sweep) writeFreq(freq uint16) {
	sweep.squareWave.frequency = freq
	sweep.freqShadow = freq
}
