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
type LedLayout [8][6]uint32

func (ll *LedLayout) LedColor(x, y, z uint8, c Color) error {
	if c == NoColor {
		return ll.LedOff(x, y, z)
	}
	if err := common.ErrIfOutOfBounds(x, "X"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] |= 1 << layoutOffsetShift(y) << x
	return nil
}

func (ll *LedLayout) LedRow(y, z uint8, c Color, rowValues byte) error {
	if c == NoColor {
		return ll.LedRowOff(y, z)
	}
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] &= ^shiftedValue(y, 255)
	ll[z][layoutOffsetIndex(y, c)] |= shiftedValue(y, rowValues)
	return nil
}

func (ll *LedLayout) LedLayer(z uint8, c Color) error {
	if c == NoColor {
		return ll.LedLayerOff(z)
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(0, c)] |= ^uint32(0)
	ll[z][layoutOffsetIndex(4, c)] |= ^uint32(0)
	return nil
}

func (ll *LedLayout) LedBlock(c Color) {
	if c == NoColor {
		ll.LedBlockOff()
		return
	}

	for _, layer := range ll {
		layer[layoutOffsetIndex(0, c)] |= ^uint32(0)
		layer[layoutOffsetIndex(4, c)] |= ^uint32(0)
	}
}

func (ll *LedLayout) LedOff(x, y, z uint8) error {
	if err := common.ErrIfOutOfBounds(x, "X"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, Green)] &= ^(1 << layoutOffsetShift(y) << x)
	ll[z][layoutOffsetIndex(y, Blue)] &= ^(1 << layoutOffsetShift(y) << x)
	ll[z][layoutOffsetIndex(y, Red)] &= ^(1 << layoutOffsetShift(y) << x)
	return nil
}

func (ll *LedLayout) LedRowOff(y, z uint8) error {
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, Green)] &= ^shiftedValue(y, 255)
	ll[z][layoutOffsetIndex(y, Blue)] &= ^shiftedValue(y, 255)
	ll[z][layoutOffsetIndex(y, Red)] &= ^shiftedValue(y, 255)
	return nil
}

func (ll *LedLayout) LedLayerOff(z uint8) error {
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	for _, states := range ll[z] {
		states &= 0
	}
	return nil
}

func (ll *LedLayout) LedBlockOff() {
	for _, layer := range ll {
		for _, states := range layer {
			states &= 0
		}
	}
}

func layoutOffsetIndex(index uint8, c Color) uint8 {
	return (uint8(c)-1)*2 + (index >> 2)
}

func layoutOffsetShift(index uint8) uint8 {
	return (index & ^(uint8(1) << 2)) * 8
}

func shiftedValue(index, value uint8) uint32 {
	return uint32(value) << layoutOffsetShift(index)
}
