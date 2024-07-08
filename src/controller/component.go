package controller

import (
	"errors"
	"machine"
)

// Driver for controlling led colors in a layer
type LedDriver struct {
	SpiClk  OutputPin // SPI Clock - Clock input terminal for data shift, PA5
	SpiMosi OutputPin // SPI MOSI - Serial-data input to the shift register, PA7
	LE      OutputPin // LE - Data strobe input terminal, PC4
	OE      OutputPin // OE - Output enable terminal, PC5
}

func NewLedDriver(clk, mosi, le, oe machine.Pin) LedDriver {
	return LedDriver{
		SpiClk:  NewOutputPin(clk),
		SpiMosi: NewOutputPin(mosi),
		LE:      NewOutputPin(le),
		OE:      NewOutputPin(oe),
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

	third := index >> 2
	second := (index - third*4) >> 1
	first := index % 2

	demux.MultiA0.Pin.Set(first == 1)
	demux.MultiA1.Pin.Set(second == 1)
	demux.MultiA2.Pin.Set(third == 1)

	return nil
}
