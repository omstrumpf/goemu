package gbc

import (
	"testing"

	"github.com/omstrumpf/goemu/internal/app/console"
)

func TestInputs(t *testing.T) {
	imd := newInputMemoryDevice()

	if imd.Read(0) != 0xEF {
		t.Errorf("Expected inputs to initialize to 0xEF, got %#04x", imd.Read(0))
	}

	imd.PressButton(console.ButtonDown)
	imd.PressButton(console.ButtonA)

	if imd.Read(0) != 0xE7 {
		t.Errorf("Expected inputs to read 0xE7, got %#04x", imd.Read(0))
	}

	imd.Write(0, 1<<4)

	if imd.Read(0) != 0xDE {
		t.Errorf("Expected inputs to read 0xDE, got %#04x", imd.Read(0))
	}

	imd.ReleaseButton(console.ButtonA)
	imd.PressButton(console.ButtonB)

	if imd.Read(0) != 0xDD {
		t.Errorf("Expected inputs to read 0xDD, got %#04x", imd.Read(0))
	}
}
