package audio

type sweep struct {
	timer *timer

	squareWave *squareWave

	period  byte
	counter byte

	shift   byte
	enabled bool
	negate  bool

	freqShadow uint16
}

func newSweep(squareWave *squareWave) *sweep {
	sweep := sweep{squareWave: squareWave}

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

			if sweep.shift != 0 {
				sweep.calculateNewFrequency()
			}
		}
	}
}

func (sweep *sweep) trigger() {
	sweep.freqShadow = sweep.squareWave.frequency
	sweep.counter = sweep.period

	sweep.enabled = (sweep.period != 0 || sweep.shift != 0)

	if sweep.shift != 0 {
		sweep.calculateNewFrequency()
	}
}

func (sweep *sweep) calculateNewFrequency() {
	newFreq := sweep.freqShadow >> sweep.shift

	if sweep.negate {
		newFreq = ^newFreq
	}

	newFreq += sweep.freqShadow

	if newFreq > 2047 {
		sweep.squareWave.enabled = false
	} else {
		sweep.squareWave.frequency = newFreq
		sweep.freqShadow = newFreq
	}
}
