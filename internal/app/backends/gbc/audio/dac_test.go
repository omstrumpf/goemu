package audio

import (
	"testing"
)

const EPSILON = 0.001

func floatEq(a float64, b float64) bool {
	return ((a - b) < EPSILON) && ((b - a) < EPSILON)
}

func TestDAC(t *testing.T) {
	dac := newDAC()

	if !floatEq(dac.convert(0), 0) {
		t.Errorf("Expected dac(0) to return 0 when disabled, got %f", dac.convert(0))
	}

	dac.enabled = true

	if !floatEq(dac.convert(0), -1) {
		t.Errorf("Expected dac(0) to return -1, got %f", dac.convert(0))
	}
	if !floatEq(dac.convert(1), -0.866) {
		t.Errorf("Expected dac(1) to return -0.866, got %f", dac.convert(1))
	}
	if !floatEq(dac.convert(2), -0.733) {
		t.Errorf("Expected dac(2) to return -0.733, got %f", dac.convert(1))
	}
	if !floatEq(dac.convert(14), 0.866) {
		t.Errorf("Expected dac(14) to return 0.866, got %f", dac.convert(14))
	}
	if !floatEq(dac.convert(15), 1) {
		t.Errorf("Expected dac(15) to return 1, got %f", dac.convert(15))
	}
}
