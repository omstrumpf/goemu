package audio

import (
	"math"

	"github.com/omstrumpf/goemu/internal/app/backends/gbc/constants"
)

type timer struct {
	countdown int
	period    int

	callback func()
}

func newTimerByClocks(clocks int, callback func()) *timer {
	return &timer{countdown: clocks, period: clocks, callback: callback}
}

func newTimerByHz(period int, callback func()) *timer {
	clocks := int(math.Round(float64(constants.ClockSpeed / float64(period))))
	return newTimerByClocks(clocks, callback)
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

func (t *timer) resetDuration(clocks int) {
	t.countdown = clocks
	t.period = clocks
}
