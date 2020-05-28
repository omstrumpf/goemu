package audio

import "github.com/omstrumpf/goemu/internal/app/log"

// Converts a 4-bit digital signal from the volume unit to an analog voltage
func dac(val byte) float64 {
	if val >= 0 && val <= 15 {
		return (float64(val) * 2 / 15) - 1
	}

	log.Errorf("DAC got input value out of range: %d", val)
	return 0
}
