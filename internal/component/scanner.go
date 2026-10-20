package component

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"
)

// Plane times in microseconds: 8 layers x (120+360+720) us = 9.6 ms per frame (~104 Hz)
var planeMicros = [framePlanes]uint32{120, 360, 720}

// Scanner refreshes the cube from a timer interrupt: one layer at a time, each layer shown
// once per plane, the planes weighted 1:3:6 (color levels 1, 2, 3 = 10%, 40%, 100%).
type Scanner struct {
	spi   *machine.SPI
	latch OutputPin // LE, PC4 (MBI5024 latches on the high pulse)
	blank OutputPin // OE, PC5 (High = outputs off)
	demux Demultiplexer
	timer *machine.TIM

	buffers [2]frameBits
	front   int
	pending bool

	layer, plane int
	blanked      bool
}

func NewScanner() *Scanner {
	scanner := &Scanner{
		spi:   NewSpiOutput(),
		latch: NewOutputPin(machine.PC4),
		blank: NewOutputPin(machine.PC5),
		demux: NewDemultiplexer(),
		timer: &machine.TIM2,
		// The first interrupt advances to plane 0 of layer 0
		layer: frameLayers - 1,
		plane: framePlanes - 1,
	}
	scanner.blank.High()
	scanner.latch.Low()

	// A 65 ms period makes TinyGo choose prescaler 72 -> 1 MHz timer clock
	// Configure the timer for 1 us ticks and set up the wraparound interrupt for refreshing the scanner
	scanner.timer.Configure(machine.PWMConfig{Period: 65_000_000})
	scanner.timer.Device.ARR.Set(planeMicros[0] - 1)
	scanner.timer.Device.EGR.SetBits(stm32.TIM_EGR_UG)
	scanner.timer.SetWraparoundInterrupt(scanner.refresh)

	return scanner
}

// Show encodes an 8x8x8 voxel buffer (index x + y*8 + z*64, 0bRRGGBB) into the back buffer and
// waits until the refresh switches to it, so a frame is never shown half drawn.
func (this *Scanner) Show(voxels []byte, mapping Mapping) {
	encodeFrame(voxels, mapping, &this.buffers[1-this.front])
	mask := interrupt.Disable()
	this.pending = true
	interrupt.Restore(mask)
	for this.isPending() {
	}
}

func (this *Scanner) isPending() bool {
	mask := interrupt.Disable()
	isPending := this.pending
	interrupt.Restore(mask)
	return isPending
}

// SetBlank turns all LEDs off (true) or resumes the refresh (false).
func (this *Scanner) SetBlank(blank bool) {
	mask := interrupt.Disable()
	this.blanked = blank
	interrupt.Restore(mask)
}

func (this *Scanner) refresh() {
	this.blank.High()

	this.plane++
	if this.plane == framePlanes {
		this.plane = 0
		this.layer++
		if this.layer == frameLayers {
			this.layer = 0
			if this.pending {
				this.front = 1 - this.front
				this.pending = false
			}
		}
	}
	next := this.plane + 1
	if next == framePlanes {
		next = 0
	}
	// ARR is preloaded: this value applies to the period after the one that starts now
	this.timer.Device.ARR.Set(planeMicros[next] - 1)

	if this.blanked {
		return
	}

	this.spi.Tx(this.buffers[this.front][this.plane][this.layer][:], nil)
	this.latch.High()
	this.latch.Low()
	this.demux.EnableLayer(uint8(this.layer))
	this.blank.Low()
}
