package audio

import "github.com/omstrumpf/goemu/internal/app/log"

type volumeShifter struct {
	volumeCode byte
}

func newVolumeShifter() *volumeShifter {
	return &volumeShifter{}
}

func (vol *volumeShifter) runForClocks(clocks int) {}

func (vol *volumeShifter) sample() float64 {
	switch vol.volumeCode {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 0.5
	case 3:
		return 0.25
	default:
		log.Errorf("APU volume shifter has invalid volume code")
		return 0
	}
}

func (vol *volumeShifter) trigger() {}
