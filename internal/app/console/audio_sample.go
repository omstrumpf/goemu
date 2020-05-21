package console

// AudioSample is a stereo audio sample
type AudioSample [2]float64

// MonoSample creates an AudioSample from a single sample
func MonoSample(val float64) AudioSample {
	return AudioSample{val, val}
}

// StereoSample creates an AudioSample from a pair of samples
func StereoSample(l float64, r float64) AudioSample {
	return AudioSample{l, r}
}

// L returns the sample's left component
func (s AudioSample) L() float64 {
	return s[0]
}

// R returns the sample's right component
func (s AudioSample) R() float64 {
	return s[1]
}
