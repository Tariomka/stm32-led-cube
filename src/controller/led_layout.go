package controller

import "github.com/Tariomka/stm32-led-cube/src/common"

type Color uint8

const (
	NoColor Color = iota
	Green
	Blue
	Red
)

// State representation of all led colors of the cube.
//
// There are 8 layers, each having 64 leds, each having 3 contacts (GBR),
// making a single layer hold 192 bits of information.
//
// 12 shift registers control all 192 activations, each shift register having 16 activations
// and all shift registers connected in series.
// Shift registers receive data in little-endian (Least significant bit first)
//
// First 8 bytes [0-7] control the Green color, next 8 [8-15] - Blue, last 8 [16-23] - Red
// First byte in a color block controlls the back side of the layer, last (8'th) byte - the front.
// First bit (0b00000001) presumably controlls the left led, the last (0b10000000) - the right led.
type LedLayout [8][24]byte

func (ll *LedLayout) LedColor(x, y, z uint8, c Color) error {
	if c == NoColor {
		return ll.LedOff(x, y, z)
	}
	err := common.ErrIfOutOfBounds(x, "X")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] |= 1 << x
	return nil
}

func (ll *LedLayout) LedRow(y, z uint8, c Color, rowValues byte) error {
	if c == NoColor {
		return ll.LedRowOff(y, z)
	}
	err := common.ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] = rowValues
	return nil
}

func (ll *LedLayout) LedOff(x, y, z uint8) error {
	err := common.ErrIfOutOfBounds(x, "X")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, Green)] &= ^(1 << x)
	ll[z][layoutOffsetIndex(y, Blue)] &= ^(1 << x)
	ll[z][layoutOffsetIndex(y, Red)] &= ^(1 << x)
	return nil
}

func (ll *LedLayout) LedRowOff(y, z uint8) error {
	err := common.ErrIfOutOfBounds(y, "Y")
	if err != nil {
		return err
	}
	err = common.ErrIfOutOfBounds(z, "Z")
	if err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, Green)] = 0
	ll[z][layoutOffsetIndex(y, Blue)] = 0
	ll[z][layoutOffsetIndex(y, Red)] = 0
	return nil
}

func layoutOffsetIndex(index uint8, c Color) uint8 {
	return (uint8(c)-1)*8 + index
}
