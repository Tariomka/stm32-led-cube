package component

import (
	"machine"
	"time"
)

// InfraredReceiver times the IR receiver output (PC6, active low, both edges) and passes each
// burst (mark) and the gap after it (space), in microseconds, to the pulse callback. The callback
// runs in interrupt context.
type InfraredReceiver struct {
	InputPin
	pulse    func(mark, space int)
	lastEdge time.Time
	lastMark int // µs, 0 = unknown
}

func NewInfraredReceiver(pulse func(mark, space int)) *InfraredReceiver {
	ir := &InfraredReceiver{InputPin: NewInputPin(machine.PC6), pulse: pulse}
	if err := ir.SetInterrupt(machine.PinToggle, ir.edge); err != nil {
		println("InfraredReceiver: SetInterrupt error:", err)
	}
	lowerPinInterruptPriority(ir.Pin)
	return ir
}

func (this *InfraredReceiver) edge(machine.Pin) {
	now := time.Now()
	elapsed := now.Sub(this.lastEdge)
	this.lastEdge = now
	interval := 0
	if elapsed < 20*time.Millisecond {
		interval = int(elapsed.Microseconds())
	}

	// A falling edge ends a space; the interval before it was the mark
	if !this.Get() && interval != 0 && this.lastMark != 0 {
		this.pulse(this.lastMark, interval)
	}
	this.lastMark = interval
}
