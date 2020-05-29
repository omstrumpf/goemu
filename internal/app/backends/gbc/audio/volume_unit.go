package audio

type volumeUnit interface {
	runForClocks(int)
	sample() float64
	trigger()
}
