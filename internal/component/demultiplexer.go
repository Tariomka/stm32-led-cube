package component

import "machine"

// Line decoder demultiplexer for providing power to each layer of led cube (Anode control).
type Demultiplexer struct {
	MultiA0     OutputPin // Demultiplexer address input 0, PB0
	MultiA1     OutputPin // Demultiplexer address input 1, PB1
	MultiA2     OutputPin // Demultiplexer address input 2, PB10
	MultiEnable OutputPin // Demultiplexer Enable pin, PA8. High turns off, Low turns on
}

func NewDemultiplexer() Demultiplexer {
	demux := Demultiplexer{
		MultiA0:     NewOutputPin(machine.PB0),
		MultiA1:     NewOutputPin(machine.PB1),
		MultiA2:     NewOutputPin(machine.PB10),
		MultiEnable: NewOutputPin(machine.PA8),
	}

	NewOutputPin(machine.PB11).Low() // For 8x8x8 cube, the fourth address input needs to always be Low
	NewOutputPin(machine.PC7).High()
	demux.MultiEnable.High()
	return demux
}

// EnableLayer powers layer 0..7.
func (this Demultiplexer) EnableLayer(index uint8) {
	this.MultiEnable.Low()
	this.MultiA0.Set(index&1 == 1)
	this.MultiA1.Set(index>>1&1 == 1)
	this.MultiA2.Set(index>>2&1 == 1)
}

func (this Demultiplexer) Disable() { this.MultiEnable.High() }

func (this Demultiplexer) Enable() { this.MultiEnable.Low() }
