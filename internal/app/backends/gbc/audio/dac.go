package audio

import "github.com/omstrumpf/goemu/internal/app/log"

type dac struct {
	enabled bool
}

func newDAC() *dac {
	return &dac{}
}

// convert converts a 4-bit digital signal from the volume unit to an analog voltage
func (dac *dac) convert(val byte) float64 {
	if !dac.enabled {
		return 0
	}

	if val >= 0 && val <= 15 {
		return (float64(val) * 2 / 15) - 1
	}

	log.Errorf("DAC got input value out of range: %d", val)
	return 0
}
