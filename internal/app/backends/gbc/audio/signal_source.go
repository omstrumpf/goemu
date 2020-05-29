package audio

type signalSource interface {
	runForClocks(int)
	sample() byte
	trigger()
}
