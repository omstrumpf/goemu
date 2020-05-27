package audio

type timer struct {
	countdown int
	period    int

	callback func()
}

func newTimer(period int, callback func()) *timer {
	return &timer{countdown: period, period: period, callback: callback}
}

func (t *timer) runForClocks(clocks int) {
	for c := 0; c < clocks; c++ {
		t.countdown--

		if t.countdown == 0 {
			t.countdown = t.period
			t.callback()
		}
	}
}
