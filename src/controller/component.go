package controller

import (
	"errors"
	"machine"
)

// Shift register for controlling led colors in a layer
type LedDriver struct {
	SPI machine.SPI // SPI interface for communicating with shift register

	// clock OutputPin // (SPI_CLOCK)Clock input terminal for data shift on rising edge, PA5
	// data  OutputPin // (SPI_MOSI)Serial-data input to the shift register, PA7

	latch OutputPin // (Latch)Data strobe input terminal, when LOW - data is latched, when HIGH - data is transfered, PC4
	blank OutputPin // (Blank)Output enable terminal, when LOW - output drivers are enabled, when HIGH - output is blanked, PC5
}

func NewLedDriver(le, oe machine.Pin) LedDriver {
	return LedDriver{
		SPI: NewSpiOutput(),

		// clock: NewOutputPin(clk),
		// data:  NewOutputPin(mosi),
		latch: NewOutputPin(le),
		blank: NewOutputPin(oe),
	}
}

// TODO: Make led rendering function

// Line decoder demultiplexer for providing power to each layer of led cube
type Demultiplexer struct {
	MultiA0     OutputPin // Demultiplexer address input 0, PB0
	MultiA1     OutputPin // Demultiplexer address input 1, PB1
	MultiA2     OutputPin // Demultiplexer address input 2, PB10
	MultiEnable OutputPin // Demultiplexer Enable pin, PB11
}

func NewDemultiplexer(a0, a1, a2, en machine.Pin) Demultiplexer {
	return Demultiplexer{
		MultiA0:     NewOutputPin(a0),
		MultiA1:     NewOutputPin(a1),
		MultiA2:     NewOutputPin(a2),
		MultiEnable: NewOutputPin(en),
	}
}

func (demux Demultiplexer) LightLayer(index uint8) error {
	if index > 7 {
		println("Index is out of range[0-7]. Received layer index:", index)
		return errors.New("index out of bounds")
	}

	demux.MultiA0.Pin.Set(index&1 == 1)
	demux.MultiA1.Pin.Set(index>>1&1 == 1)
	demux.MultiA2.Pin.Set(index>>2&1 == 1)

	return nil
}
