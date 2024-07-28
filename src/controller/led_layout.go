package controller

import "github.com/Tariomka/stm32-led-cube/src/common"

type Color uint8

const (
	NoColor Color = iota
	Green
	Blue
	Red
)

const (
	none byte = 0b00000000
	all  byte = 0b11111111
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
// Bits from right to left control each led from right to left led,
// i.e. 0b00000001 turns on the right most led, 0b10000000 - left most led.
type LedLayout [8][24]byte

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

	ll[z][layoutOffsetIndex(y, c)] |= 1 << x
	return nil
}

func (ll *LedLayout) LedRowIndividual(y, z uint8, c Color, rowValues byte) error {
	if c == NoColor {
		return ll.LedRowOff(y, z)
	}
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] |= rowValues
	return nil
}

func (ll *LedLayout) LedRow(y, z uint8, c Color) error {
	if c == NoColor {
		return ll.LedRowOff(y, z)
	}
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	ll[z][layoutOffsetIndex(y, c)] |= all
	return nil
}

func (ll *LedLayout) LedLayer(z uint8, c Color) error {
	if c == NoColor {
		return ll.LedLayerOff(z)
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	offset := layoutOffsetIndex(0, c)
	for i := offset; i < offset+8; i++ {
		ll[z][i] |= all
	}
	return nil
}

func (ll *LedLayout) LedBlock(c Color) {
	if c == NoColor {
		ll.LedBlockOff()
		return
	}

	for z := range ll {
		offset := layoutOffsetIndex(0, c)
		for i := offset; i < offset+8; i++ {
			ll[z][i] |= all
		}
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

	offset := layoutOffsetIndex(y, Green)
	for i := offset; i < 24; i += 8 {
		ll[z][i] &= ^(1 << x)
	}
	return nil
}

func (ll *LedLayout) LedRowIndividualOff(y, z uint8, rowValues byte) error {
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	offset := layoutOffsetIndex(y, Green)
	for i := offset; i < 24; i += 8 {
		ll[z][i] &= ^rowValues
	}
	return nil
}

func (ll *LedLayout) LedRowOff(y, z uint8) error {
	if err := common.ErrIfOutOfBounds(y, "Y"); err != nil {
		return err
	}
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	offset := layoutOffsetIndex(y, Green)
	for i := offset; i < 24; i += 8 {
		ll[z][i] &= none
	}
	return nil
}

func (ll *LedLayout) LedLayerOff(z uint8) error {
	if err := common.ErrIfOutOfBounds(z, "Z"); err != nil {
		return err
	}

	for i := range ll[z] {
		ll[z][i] &= none
	}
	return nil
}

func (ll *LedLayout) LedBlockOff() {
	for z := range ll {
		for i := range ll[z] {
			ll[z][i] &= none
		}
	}
}

func layoutOffsetIndex(index uint8, c Color) uint8 {
	return (uint8(c)-1)*8 + index
}
