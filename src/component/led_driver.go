package component

import "machine"

// (MBI5024GP/GF)
// Shift register for controlling led colors in a layer(Cathode control)
type LedDriver struct {
	// SPI interface for communicating with shift register. SPI uses:
	// (SPI_CLOCK)Clock input terminal for data shift on rising edge, PA5;
	// (SPI_MOSI)Serial-data input to the shift register, PA7.
	spi machine.SPI

	latch OutputPin // (Latch)Data strobe input terminal, when LOW - data is latched, when HIGH - data is transfered, PC4
	blank OutputPin // (Blank)Output enable terminal, when LOW - output drivers are enabled, when HIGH - output is blanked, PC5
}

func NewLedDriver(le, oe machine.Pin) *LedDriver {
	ld := LedDriver{
		spi:   NewSpiOutput(),
		latch: NewOutputPin(le),
		blank: NewOutputPin(oe),
	}
	ld.latch.Pin.High()
	ld.blank.Pin.Low()
	return &ld
}

func (ld *LedDriver) ClearLayer() error {
	ld.latch.Pin.Low()
	err := ld.spi.Tx(make([]byte, 24), nil)
	ld.latch.Pin.High()
	return err
}

func (ld *LedDriver) LightLayer(data []byte) error {
	ld.latch.Pin.Low()
	err := ld.spi.Tx(data, nil)
	ld.latch.Pin.High()
	return err
}
