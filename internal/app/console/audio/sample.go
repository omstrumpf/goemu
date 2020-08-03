package audio

// Sample is a stereo audio sample
type Sample [2]float64

// L returns the sample's left component
func (s Sample) L() float64 {
	return s[0]
}

// R returns the sample's right component
func (s Sample) R() float64 {
	return s[1]
}

// M returns the average of the sample's left and right components (mono)
func (s Sample) M() float64 {
	return (s.L() + s.R()) / 2
}
