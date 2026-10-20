package component

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"

	"github.com/Tariomka/stm32-led-cube/internal/engine/frame"
)

const basePlaneMicros = 400 // Shortest plane time in microseconds: 8 layers x (1+2) x 400 us = 9.6 ms per frame (~104 Hz)

// Scanner refreshes the cube from a timer interrupt: one layer at a time, each layer shown
// once per bit plane with binary weighted durations (2-bit color depth per channel).
type Scanner struct {
	spi   *machine.SPI
	latch OutputPin // LE, PC4 (MBI5024 latches on the high pulse)
	blank OutputPin // OE, PC5 (High = outputs off)
	demux Demultiplexer
	timer *machine.TIM

	buffers [2]frame.Bits
	front   int
	pending bool

	layer, plane int
	blanked      bool
	planeTicks   [frame.Planes]uint32
}

func NewScanner() *Scanner {
	scanner := &Scanner{
		spi:   NewSpiOutput(),
		latch: NewOutputPin(machine.PC4),
		blank: NewOutputPin(machine.PC5),
		demux: NewDemultiplexer(),
		timer: &machine.TIM2,
		// The first interrupt advances to plane 0 of layer 0
		layer: frame.Layers - 1,
		plane: frame.Planes - 1,
	}
	scanner.blank.High()
	scanner.latch.Low()
	for p := range scanner.planeTicks {
		scanner.planeTicks[p] = basePlaneMicros << p
	}

	// A 65 ms period makes TinyGo choose prescaler 72 -> 1 MHz timer clock
	// Configure the timer for 1 us ticks and set up the wraparound interrupt for refreshing the scanner
	scanner.timer.Configure(machine.PWMConfig{Period: 65_000_000})
	scanner.timer.Device.ARR.Set(scanner.planeTicks[0] - 1)
	scanner.timer.Device.EGR.SetBits(stm32.TIM_EGR_UG)
	scanner.timer.SetWraparoundInterrupt(scanner.refresh)

	return scanner
}

// Back returns the buffer to draw the next frame into. It must not be used again until Present.
func (this *Scanner) Back() *frame.Bits { return &this.buffers[1-this.front] }

// Present shows the back buffer from the next full refresh cycle on.
func (this *Scanner) Present() {
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
	if this.plane == frame.Planes {
		this.plane = 0
		this.layer++
		if this.layer == frame.Layers {
			this.layer = 0
			if this.pending {
				this.front = 1 - this.front
				this.pending = false
			}
		}
	}
	next := this.plane + 1
	if next == frame.Planes {
		next = 0
	}
	// ARR is preloaded: this value applies to the period after the one that starts now
	this.timer.Device.ARR.Set(this.planeTicks[next] - 1)

	if this.blanked {
		return
	}

	this.spi.Tx(this.buffers[this.front][this.plane][this.layer][:], nil)
	this.latch.High()
	this.latch.Low()
	this.demux.EnableLayer(uint8(this.layer))
	this.blank.Low()
}
