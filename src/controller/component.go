package controller

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

func NewLedDriver(le, oe machine.Pin) LedDriver {
	return LedDriver{
		spi: NewSpiOutput(),

		latch: NewOutputPin(le),
		blank: NewOutputPin(oe),
	}
}

func (ld LedDriver) ClearLayer() {
	for i := 0; i < 24; i++ {
		ld.latch.Pin.Low()
		ld.spi.Tx([]byte{0}, nil)
		ld.latch.Pin.High()
	}
	ld.latch.Pin.Low()
}

func (ld LedDriver) LightLayer(data []byte) error {
	ld.latch.Pin.Low()
	err := ld.spi.Tx(data, nil)
	ld.latch.Pin.High()
	ld.latch.Pin.Low()
	return err
}

// (74HC138)
// Line decoder demultiplexer for providing power to each layer of led cube (Anode control).
type Demultiplexer struct {
	MultiA0     OutputPin // Demultiplexer address input 0, PB0
	MultiA1     OutputPin // Demultiplexer address input 1, PB1
	MultiA2     OutputPin // Demultiplexer address input 2, PB10
	MultiEnable OutputPin // Demultiplexer Enable pin, PB11
}

func NewDemultiplexer(a0, a1, a2, en machine.Pin) Demultiplexer {
	demux := Demultiplexer{
		MultiA0:     NewOutputPin(a0),
		MultiA1:     NewOutputPin(a1),
		MultiA2:     NewOutputPin(a2),
		MultiEnable: NewOutputPin(en),
	}
	demux.MultiEnable.Pin.High()
	return demux
}

func (demux Demultiplexer) EnableLayer(index uint8) error {
	err := ErrIfOutOfBounds(index, "Z")
	if err != nil {
		return err
	}

	demux.MultiA0.Pin.Set(index&1 == 1)
	demux.MultiA1.Pin.Set(index>>1&1 == 1)
	demux.MultiA2.Pin.Set(index>>2&1 == 1)

	return nil
}
