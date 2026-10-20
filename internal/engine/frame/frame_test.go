package frame

import "testing"

func TestEncode(t *testing.T) {
	var v [512]byte
	v[3+2*8+5*64] = 0b11_01_10 // red 3, green 1, blue 2 at x=3, y=2, z=5
	var b Bits
	Encode(v[:], Mapping{}, &b)
	bit := byte(1 << 3)
	checks := []struct {
		plane, layer, index int
		want                byte
	}{
		{0, 2, redOffset + 5, bit}, {1, 2, redOffset + 5, bit}, {2, 2, redOffset + 5, bit},
		{0, 2, greenOffset + 5, bit}, {1, 2, greenOffset + 5, 0}, {2, 2, greenOffset + 5, 0},
		{0, 2, blueOffset + 5, bit}, {1, 2, blueOffset + 5, bit}, {2, 2, blueOffset + 5, 0},
	}
	for _, c := range checks {
		if got := b[c.plane][c.layer][c.index]; got != c.want {
			t.Errorf("plane %d layer %d byte %d = %08b, want %08b", c.plane, c.layer, c.index, got, c.want)
		}
	}

	Encode(v[:], Mapping{FlipX: true, FlipZ: true, FlipLayer: true, SwapXZ: true}, &b)
	if got := b[0][5][redOffset+4]; got != 1<<2 {
		t.Errorf("mapped voxel = %08b at layer 5 row 4, want bit 2", got)
	}
}
