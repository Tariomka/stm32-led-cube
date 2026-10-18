package component

import (
	"machine"

	"github.com/Tariomka/stm32-led-cube/internal/common/utils"
)

// Line decoder demultiplexer for providing power to each layer of led cube (Anode control).
type Demultiplexer struct {
	MultiA0     OutputPin // Demultiplexer address input 0, PB0
	MultiA1     OutputPin // Demultiplexer address input 1, PB1
	MultiA2     OutputPin // Demultiplexer address input 2, PB10
	MultiEnable OutputPin // Demultiplexer Enable pin, PA8. High turns off, Low turns on
}

func NewDemultiplexer(a0, a1, a2, a3, en1, en2 machine.Pin) Demultiplexer {
	demux := Demultiplexer{
		MultiA0:     NewOutputPin(a0),
		MultiA1:     NewOutputPin(a1),
		MultiA2:     NewOutputPin(a2),
		MultiEnable: NewOutputPin(en1),
	}

	NewOutputPin(a3).Low() // For 8x8x8 cube, the fourth address input needs to always be Low
	NewOutputPin(en2).High()

	demux.MultiEnable.High()
	return demux
}

func (this Demultiplexer) EnableLayer(index uint8) error {
	if err := utils.ErrIfOutOfBounds(index); err != nil {
		return err
	}

	this.MultiEnable.Low()
	this.MultiA0.Set(index&1 == 1)
	this.MultiA1.Set(index>>1&1 == 1)
	this.MultiA2.Set(index>>2&1 == 1)

	return nil
}

func (this Demultiplexer) Disable() {
	this.MultiEnable.High()
}

func (this Demultiplexer) Enable() {
	this.MultiEnable.Low()
}
