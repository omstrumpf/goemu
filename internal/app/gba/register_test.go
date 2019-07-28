package gba

import "testing"

func TestRegisterGetSet(t *testing.T) {
	AF := Register{0, 0}
	BC := Register{0x1234, 0}
	DE := Register{0xABCD, 0}

	if AF.HiLo() != 0 {
		t.Errorf("Expected register AF to contain 0, got %#4x", AF.HiLo())
	}
	if BC.HiLo() != 0x1234 {
		t.Errorf("Expected register BC to contain 0x1234, got %#4x", BC.HiLo())
	}
	if DE.HiLo() != 0xABCD {
		t.Errorf("Expected register DE to contain 0xABCD, got %#4x", DE.HiLo())
	}
	if BC.Hi() != 0x12 {
		t.Errorf("Expected hi byte of BC to contain 0x12, got %#2x", BC.Hi())
	}
	if BC.Lo() != 0x34 {
		t.Errorf("Expected low byte of BC to contain 0x34, got %#2x", BC.Lo())
	}
	if DE.Hi() != 0xAB {
		t.Errorf("Expected hi byte of DE to contain 0xAB, got %#2x", DE.Hi())
	}
	if DE.Lo() != 0xCD {
		t.Errorf("Expected low byte of DE to contain 0xCD, got %#2x", DE.Lo())
	}

	AF.Set(0x5678)
	BC.SetHi(0xAB)
	DE.SetLo(0x34)

	if AF.HiLo() != 0x5678 {
		t.Errorf("Expected register AF to contain 0x5678, got %#4x", AF.HiLo())
	}
	if BC.HiLo() != 0xAB34 {
		t.Errorf("Expected register BC to contain 0xAB34, got %#4x", BC.HiLo())
	}
	if DE.HiLo() != 0xAB34 {
		t.Errorf("Expected register DE to contain 0xAB34, got %#4x", DE.HiLo())
	}
}

func TestRegisterIncDec(t *testing.T) {
	AF := Register{0x0010, 0}

	before := AF.Inc()
	after := AF.HiLo()

	if before != 0x10 {
		t.Errorf("Expected Inc to return original value of 0x0010, got %#4x", before)
	}
	if after != 0x11 {
		t.Errorf("Expected Inc to increment the value by one, to 0x0011, got %#4x", after)
	}

	before = AF.Dec()
	after = AF.HiLo()

	if before != 0x11 {
		t.Errorf("Expected Dec to return original value of 0x0011, got %#4x", before)
	}
	if after != 0x10 {
		t.Errorf("Expected Dec to decrement the value by one, to 0x0010, got %#4x", after)
	}

	before = AF.Inc2()
	after = AF.HiLo()

	if before != 0x10 {
		t.Errorf("Expected Inc2 to return original value of 0x0010, got %#4x", before)
	}
	if after != 0x12 {
		t.Errorf("Expected Inc2 to increment the value by two, to 0x0012, got %#4x", after)
	}

	before = AF.Dec2()
	after = AF.HiLo()

	if before != 0x12 {
		t.Errorf("Expected Dec2 to return original value of 0x0012, got %#4x", before)
	}
	if after != 0x10 {
		t.Errorf("Expected Dec2 to decrement the value by two, to 0x0010, got %#4x", after)
	}

}

func TestRegisterMask(t *testing.T) {
	AF := Register{0x0000, 0xFFF0}
	BC := Register{0x0000, 0x0000}

	AF.Set(0xABCD)
	BC.Set(0xABCD)

	if AF.HiLo() != 0xABC0 {
		t.Errorf("Expected AF to be masked to 0xABC0, got %#4x", AF.HiLo())
	}
	if BC.HiLo() != 0xABCD {
		t.Errorf("Expected BC to contain 0xABCD, got %#4x", BC.HiLo())
	}
}
