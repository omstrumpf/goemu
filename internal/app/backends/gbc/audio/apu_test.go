package audio

import (
	"testing"
)

func TestAPUPower(t *testing.T) {
	apu := NewAPU(1)

	if !apu.enabled {
		t.Errorf("Expected APU to start enabled")
	}

	apu.Write(0xFF26, 0b0000_0000)

	if apu.enabled {
		t.Errorf("Expected reset of FF26 to disable APU")
	}

	apu.Write(0xFF26, 0b1000_0000)

	if !apu.enabled {
		t.Errorf("Expected write to FF26 to enable APU")
	}
}

func TestAPUWaveRAM(t *testing.T) {
	apu := NewAPU(1)

	apu.Write(0xFF26, 0b1000_0000)

	apu.Write(0xFF30, 0xA7)
	apu.Write(0xFF31, 0xB7)
	apu.Write(0xFF3E, 0xC7)
	apu.Write(0xFF3F, 0xD7)

	if apu.Read(0xFF30) != 0xA7 {
		t.Errorf("Expected WAVE ram to persist write of 0xA7, got %#02X", apu.Read(0xFF30))
	}

	if apu.Read(0xFF31) != 0xB7 {
		t.Errorf("Expected WAVE ram to persist write of 0xB7, got %#02X", apu.Read(0xFF31))
	}

	if apu.Read(0xFF3E) != 0xC7 {
		t.Errorf("Expected WAVE ram to persist write of 0xC7, got %#02X", apu.Read(0xFF3E))
	}

	if apu.Read(0xFF3F) != 0xD7 {
		t.Errorf("Expected WAVE ram to persist write of 0xD7, got %#02X", apu.Read(0xFF3F))
	}

	apu.Write(0xFF26, 0b0000_0000)

	if apu.Read(0xFF30) != 0xA7 {
		t.Errorf("Expected WAVE ram to persist write of 0xA7 after power off, got %#02X", apu.Read(0xFF30))
	}

	if apu.Read(0xFF31) != 0xB7 {
		t.Errorf("Expected WAVE ram to persist write of 0xB7 after power off, got %#02X", apu.Read(0xFF31))
	}

	apu.Write(0xFF26, 0b1000_0000)

	if apu.Read(0xFF3E) != 0xC7 {
		t.Errorf("Expected WAVE ram to persist write of 0xC7 after power off, got %#02X", apu.Read(0xFF3E))
	}

	if apu.Read(0xFF3F) != 0xD7 {
		t.Errorf("Expected WAVE ram to persist write of 0xD7 after power off, got %#02X", apu.Read(0xFF3F))
	}

}
