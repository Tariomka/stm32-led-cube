package component

import (
	"machine"

	"github.com/Tariomka/stm32-led-cube/src/common"
)

// (74HC154)
// Line decoder demultiplexer for providing power to each layer of led cube (Anode control).
type Demultiplexer struct {
	MultiA0     OutputPin // Demultiplexer address input 0, PB0
	MultiA1     OutputPin // Demultiplexer address input 1, PB1
	MultiA2     OutputPin // Demultiplexer address input 2, PB10
	MultiEnable OutputPin // Demultiplexer Enable pin, PA8
}

func NewDemultiplexer(a0, a1, a2, a3, en machine.Pin) Demultiplexer {
	demux := Demultiplexer{
		MultiA0:     NewOutputPin(a0),
		MultiA1:     NewOutputPin(a1),
		MultiA2:     NewOutputPin(a2),
		MultiEnable: NewOutputPin(en),
	}
	demux.MultiEnable.Pin.Low()

	// For 8x8x8 cube, the fourth address input needs to always be Low
	a3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	a3.Low()

	return demux
}

func (demux Demultiplexer) EnableLayer(index uint8) error {
	err := common.ErrIfOutOfBounds(index, "Z")
	if err != nil {
		return err
	}

	demux.MultiA0.Pin.Set(index&1 == 1)
	demux.MultiA1.Pin.Set(index>>1&1 == 1)
	demux.MultiA2.Pin.Set(index>>2&1 == 1)

	return nil
}
