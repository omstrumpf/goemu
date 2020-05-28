package audio

type envelope struct {
	timer *timer

	initVolume byte
	volume     byte

	sweepPeriod  byte
	sweepCounter byte

	mode bool
}

func newEnvelope() *envelope {
	envelope := envelope{}

	envelope.timer = newTimerByHz(64, envelope.tick)

	return &envelope
}

func (envelope *envelope) runForClocks(clocks int) {
	envelope.timer.runForClocks(clocks)
}

func (envelope *envelope) tick() {
	if envelope.sweepPeriod != 0 {
		envelope.sweepCounter--
		if envelope.sweepCounter == 0 {
			envelope.sweepCounter = envelope.sweepPeriod

			if envelope.volume > 0 && envelope.volume < 15 {
				if envelope.mode {
					envelope.volume++
				} else {
					envelope.volume--
				}
			}
		}
	}
}

func (envelope *envelope) sample() byte {
	return byte(envelope.volume)
}

func (envelope *envelope) trigger() {
	envelope.sweepCounter = envelope.sweepPeriod
	envelope.volume = envelope.initVolume
}
