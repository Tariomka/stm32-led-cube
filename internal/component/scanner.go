package component

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"

	"github.com/Tariomka/stm32-led-cube/internal/engine/frame"
)

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

// Shortest plane time in microseconds: 8 layers x (1+2) x 400 us = 9.6 ms per frame (~104 Hz)
const basePlaneMicros = 400

func NewScanner(spi *machine.SPI, sck, sdo, le, oe machine.Pin, demux Demultiplexer) *Scanner {
	s := &Scanner{
		spi:   NewSpiOutput(spi, sck, sdo),
		latch: NewOutputPin(le),
		blank: NewOutputPin(oe),
		demux: demux,
		timer: &machine.TIM2,
		// The first interrupt advances to plane 0 of layer 0
		layer: frame.Layers - 1,
		plane: frame.Planes - 1,
	}
	s.blank.High()
	s.latch.Low()
	for p := range s.planeTicks {
		s.planeTicks[p] = basePlaneMicros << p
	}
	return s
}

// Start configures TIM2 for 1 us ticks and begins refreshing.
func (s *Scanner) Start() {
	// A 65 ms period makes TinyGo choose prescaler 72 -> 1 MHz timer clock
	s.timer.Configure(machine.PWMConfig{Period: 65_000_000})
	s.timer.Device.ARR.Set(s.planeTicks[0] - 1)
	s.timer.Device.EGR.SetBits(stm32.TIM_EGR_UG)
	s.timer.SetWraparoundInterrupt(s.refresh)
}

// Back returns the buffer to draw the next frame into. It must not be used again until Present.
func (s *Scanner) Back() *frame.Bits {
	return &s.buffers[1-s.front]
}

// Present shows the back buffer from the next full refresh cycle on.
func (s *Scanner) Present() {
	mask := interrupt.Disable()
	s.pending = true
	interrupt.Restore(mask)
	for s.isPending() {
	}
}

func (s *Scanner) isPending() bool {
	mask := interrupt.Disable()
	p := s.pending
	interrupt.Restore(mask)
	return p
}

// SetBlank turns all LEDs off (true) or resumes the refresh (false).
func (s *Scanner) SetBlank(blank bool) {
	mask := interrupt.Disable()
	s.blanked = blank
	interrupt.Restore(mask)
}

func (s *Scanner) refresh() {
	s.blank.High()

	s.plane++
	if s.plane == frame.Planes {
		s.plane = 0
		s.layer++
		if s.layer == frame.Layers {
			s.layer = 0
			if s.pending {
				s.front = 1 - s.front
				s.pending = false
			}
		}
	}
	next := s.plane + 1
	if next == frame.Planes {
		next = 0
	}
	// ARR is preloaded: this value applies to the period after the one that starts now
	s.timer.Device.ARR.Set(s.planeTicks[next] - 1)

	if s.blanked {
		return
	}
	s.spi.Tx(s.buffers[s.front][s.plane][s.layer][:], nil)
	s.latch.High()
	s.latch.Low()
	s.demux.EnableLayer(uint8(s.layer))
	s.blank.Low()
}
