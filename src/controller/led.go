package controller

type Color uint8

const (
	NoColor Color = iota
	Green
	Blue
	Red
)

// State representation of all led colors of the cube.
//
// There are 8 layers, each having 64 leds, each having 3 contacts (RGB),
// making a single layer hold 192 bits of information.
//
// 12 shift registers control all 192 activations, each shift register having 16 activations
// and all shift registers connected in series.
// Shift registers receive data in little-endian (Least significant bit first)
//
// If a byte is set to 0, the entire row is off.
// If a byte is 1, the leftmost LED is on.
// If a byte is 255, the entire row is on.
type LedLayout [8][24]byte

func (ll *LedLayout) LedColor(x, y, z uint8, c Color) error {
	if c == NoColor {
		return ll.LedOff(x, y, z)
	}
	err := ErrIfOutOfBounds(x, "X")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, uint8(c))] |= 1 << x
	return nil
}

func (ll *LedLayout) LedRow(y, z uint8, c Color, rowValues byte) error {
	if c == NoColor {
		return ll.LedRowOff(y, z)
	}
	err := ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, uint8(c))] = rowValues
	return nil
}

func (ll *LedLayout) LedOff(x, y, z uint8) error {
	err := ErrIfOutOfBounds(x, "X")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, uint8(Green))] &= ^(1 << x)
	ll[z][layoutOffsetIndex(y, uint8(Blue))] &= ^(1 << x)
	ll[z][layoutOffsetIndex(y, uint8(Red))] &= ^(1 << x)
	return nil
}

func (ll *LedLayout) LedRowOff(y, z uint8) error {
	err := ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, uint8(Green))] = 0
	ll[z][layoutOffsetIndex(y, uint8(Blue))] = 0
	ll[z][layoutOffsetIndex(y, uint8(Red))] = 0
	return nil
}

func layoutOffsetIndex(index, colorValue uint8) uint8 {
	// Entire sequence is reversed.
	//
	// Each pair of bytes corresponds to a single shift register.
	// Each shift register in the chain alternates R->B->G.
	//
	// [I / 2 * 6] is required for cutting off the floating point,
	// making the index a multiple of 2 times 3.
	//
	// Finally, checking if even to distinguish the 2 parts of the register.
	return 23 - ((colorValue-1)*2 + index/2*6 + index%2)
}
