package gba

import (
	"io/ioutil"
	"testing"
)

func TestMMUInit(t *testing.T) {
	mmu := NewMMU()

	if mmu.biosEnable != true {
		t.Error("MMU biosEnable flag should init to true")
	}
}

func TestMMUReadWrite(t *testing.T) {
	mmu := NewMMU()

	got := mmu.Read(0)
	if got != BIOS[0] {
		t.Errorf("Expected initial value to be %#2x, got %#2x", BIOS[0], got)
	}

	mmu.Write(0, 0x12)
	mmu.Write(1, 0x34)
	mmu.Write16(2, 0x7856)

	got = mmu.Read(0)
	if got != 0x12 {
		t.Errorf("Expected to read written value of 0x12, got %#2x", got)
	}
	got16 := mmu.Read16(1)
	if got16 != 0x5634 {
		t.Errorf("Expected to read written value of 0x5634, got %#4x", got16)
	}
	got = mmu.Read(3)
	if got != 0x78 {
		t.Errorf("Expected to read writte nvalue of 0x78, got %#2x", got)
	}
}

func TestMMUBios(t *testing.T) {
	mmu := NewMMU()

	mmu.Write(0x099, 0x12)
	mmu.Write(0x100, 0x34)

	got := mmu.Read(0x099)
	if got != 0x12 {
		t.Errorf("Expected to read written value of 0x12, got %#2x", got)
	}
	got = mmu.Read(0x100)
	if got != 0x34 {
		t.Errorf("Expected to read written value of 0x34, got %#2x", got)
	}

	mmu.DisableBios()

	got = mmu.Read(0x099)
	if got != 0 {
		t.Errorf("Expected to read 0 after disabling bios, got %#2x", got)
	}
	got = mmu.Read(0x100)
	if got != 0x34 {
		t.Errorf("Expected value 0x34 to be preserved after disabling bios, got %#2x", got)
	}
}

func TestMMUMappings(t *testing.T) {
	mmu := NewMMU()

	mmu.Write(0x0000, 0x1) // BIOS
	mmu.Write(0x0099, 0x2) // BIOS
	mmu.Write(0x0100, 0x3) // ROM
	mmu.Write(0x7FFF, 0x4) // ROM
	mmu.Write(0x8000, 0x5) // VRAM
	mmu.Write(0x9FFF, 0x6) // VRAM
	mmu.Write(0xA000, 0x7) // ERAM
	mmu.Write(0xBFFF, 0x8) // ERAM
	mmu.Write(0xC000, 0x9) // WRAM
	mmu.Write(0xDFFF, 0xA) // WRAM
	mmu.Write(0xFE00, 0xB) // GOAM
	mmu.Write(0xFF80, 0xC) // ZRAM
	mmu.Write(0xFFFF, 0xD) // ZRAM

	if mmu.bios[0x0000] != 0x1 {
		t.Errorf("Expected to read 0x1 from bios memory, got %#2x", mmu.bios[0x0000])
	}
	if mmu.bios[0x0099] != 0x2 {
		t.Errorf("Expected to read 0x2 from bios memory, got %#2x", mmu.bios[0x0099])
	}
	if mmu.rom[0x0100] != 0x3 {
		t.Errorf("Expected to read 0x3 from rom, got %#2x", mmu.rom[0x0100])
	}
	if mmu.rom[0x7FFF] != 0x4 {
		t.Errorf("Expected to read 0x4 from rom, got %#2x", mmu.rom[0x7FFF])
	}
	if mmu.vram[0x0000] != 0x5 {
		t.Errorf("Expected to read 0x5 from vram, got %#2x", mmu.vram[0x0000])
	}
	if mmu.vram[0x1FFF] != 0x6 {
		t.Errorf("Expected to read 0x6 from vram, got %#2x", mmu.vram[0x1FFF])
	}
	if mmu.eram[0x0000] != 0x7 {
		t.Errorf("Expected to read 0x7 from eram, got %#2x", mmu.eram[0x0000])
	}
	if mmu.eram[0x1FFF] != 0x8 {
		t.Errorf("Expected to read 0x8 from eram, got %#2x", mmu.eram[0x1FFF])
	}
	if mmu.wram[0x0000] != 0x9 {
		t.Errorf("Expected to read 0x9 from wram, got %#2x", mmu.wram[0x0000])
	}
	if mmu.wram[0x1FFF] != 0xA {
		t.Errorf("Expected to read 0xA from wram, got %#2x", mmu.wram[0x1FFF])
	}
	if mmu.goam[0x0000] != 0xB {
		t.Errorf("Expected to read 0xB from goam, got %#2x", mmu.goam[0x0000])
	}
	if mmu.zram[0x0000] != 0xC {
		t.Errorf("Expected to read 0xC from zram, got %#2x", mmu.zram[0x0000])
	}
	if mmu.zram[0x007F] != 0xD {
		t.Errorf("Expected to read 0xD from zram, got %#2x", mmu.zram[0x007F])
	}

	if mmu.Read(0xE000) != 0x9 {
		t.Errorf("Expected to read 0x9 from wram shadow, got %#2x", mmu.Read(0xE000))
	}
	if mmu.Read(0xF000) != 0x0 {
		t.Errorf("Expected to read 0x0 from wram shadow, got %#2x", mmu.Read(0xF000))
	}
}

func TestMMUZeros(t *testing.T) {
	mmu := NewMMU()

	mmu.Write(0xFEFF, 0xFF) // Should be forced to 0

	if mmu.Read(0xFEFF) != 0 {
		t.Errorf("Expected to read 0, got %#2x", mmu.Read(0xFEFF))
	}
}

func TestMMULoadRom(t *testing.T) {
	mmu := NewMMU()

	romfile := "../../../roms/cpu_instrs.gb"

	mmu.LoadROM(romfile)

	buf, err := ioutil.ReadFile(romfile)
	if err != nil {
		panic(err)
	}

	mmu.DisableBios()

	if buf[0] != mmu.Read(0) {
		t.Errorf("Expected to read %#2x, got %#2x", buf[0], mmu.Read(0))
	}

}
