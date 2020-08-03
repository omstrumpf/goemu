package audio

// ChanneledSample is a stereo audio sample composed of multiple audio channels
type ChanneledSample struct {
	Channels []Sample
}

// Combine combines the channeled data into a single stereo sample
func (s ChanneledSample) Combine() Sample {
	return Sample{s.L(), s.R()}
}

// L returns the sample's left component
func (s ChanneledSample) L() float64 {
	l := float64(0)

	for _, c := range s.Channels {
		l += c.L()
	}

	return l
}

// R returns the sample's right component
func (s ChanneledSample) R() float64 {
	r := float64(0)

	for _, c := range s.Channels {
		r += c.R()
	}

	return r
}

// M returns the average of the sample's left and right components (mono)
func (s ChanneledSample) M() float64 {
	return (s.L() + s.R()) / 2
}
