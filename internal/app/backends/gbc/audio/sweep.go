package audio

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

	negativeSweepSinceLastTrigger bool // For some very obscure behaviour
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
	if sweep.counter > 0 {
		sweep.counter--
	}

	if sweep.counter == 0 {
		sweep.reloadCounter()

		if sweep.enabled && sweep.period > 0 {
			newFreq := sweep.calculateNewFrequency()

			if !sweep.checkOverflow(newFreq) && sweep.shift > 0 {
				sweep.writeFreq(newFreq)

				sweep.checkOverflow(sweep.calculateNewFrequency())
			}
		}
	}
}

func (sweep *sweep) trigger() {
	sweep.negativeSweepSinceLastTrigger = false

	sweep.freqShadow = sweep.squareWave.frequency

	sweep.reloadCounter()

	sweep.enabled = (sweep.period > 0 || sweep.shift > 0)

	if sweep.shift > 0 {
		sweep.checkOverflow(sweep.calculateNewFrequency())
	}
}

func (sweep *sweep) setNegate(negate bool) {
	if sweep.negativeSweepSinceLastTrigger && sweep.negate && !negate {
		// Clearing the sweep negate mode bit after at least one sweep calculation
		// has been made using the negate mode since the last trigger causes the
		// channel to be immediately disabled.
		sweep.channel.disable()
	}

	sweep.negate = negate
}

func (sweep *sweep) reloadCounter() {
	if sweep.period > 0 {
		sweep.counter = sweep.period
	} else {
		sweep.counter = 8
	}
}

func (sweep *sweep) calculateNewFrequency() uint16 {
	newFreq := sweep.freqShadow >> sweep.shift

	if sweep.negate {
		newFreq = sweep.freqShadow - newFreq
		sweep.negativeSweepSinceLastTrigger = true
	} else {
		newFreq = sweep.freqShadow + newFreq
	}

	return newFreq
}

func (sweep *sweep) checkOverflow(freq uint16) bool {
	if freq > 2047 {
		sweep.channel.disable()
		return true
	}

	return false
}

func (sweep *sweep) writeFreq(freq uint16) {
	sweep.squareWave.frequency = freq
	sweep.freqShadow = freq
}
