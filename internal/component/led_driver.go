package component

import "machine"

const bytesInLayerCount = 24

// (MBI5024GP/GF)
// Shift register for controlling led colors in a layer(Cathode control)
type LedDriver struct {
	// SPI interface for communicating with shift register. SPI uses:
	//
	// - (SPI_CLOCK)Clock input terminal for data shift on rising edge, PA5;
	//
	// - (SPI_MOSI)Serial-data input to the shift register, PA7.
	spi *machine.SPI

	// Data strobe input terminal, PC4. Functionality:
	//
	// - When LOW - data is latched,
	//
	// - When HIGH - data is transfered
	latch OutputPin
	// Output enable terminal, PC5. Functionality:
	//
	// When LOW - output drivers are enabled,
	//
	// When HIGH - output is blanked
	blank OutputPin
}

func NewLedDriver(spi *machine.SPI, sck, sdo, le, oe machine.Pin) *LedDriver {
	ld := LedDriver{
		spi:   NewSpiOutput(spi, sck, sdo),
		latch: NewOutputPin(le),
		blank: NewOutputPin(oe),
	}
	ld.latch.High()
	ld.blank.Low()
	return &ld
}

func (this *LedDriver) ClearLayer() error {
	this.latch.Low()
	err := this.spi.Tx(make([]byte, bytesInLayerCount), nil)
	this.latch.High()
	return err
}

func (this *LedDriver) LightLayer(data []byte) error {
	this.latch.Low()
	err := this.spi.Tx(data, nil)
	this.latch.High()
	return err
}
