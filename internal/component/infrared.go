package component

import (
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/remote"
)

// InfraredReceiver decodes the IR remote receiver output (active low) on an EXTI pin.
type InfraredReceiver struct {
	InputPin
	decoder  *remote.Decoder
	lastEdge time.Time
	lastMark int // µs, 0 = unknown
	code     int // last decoded code, 0 = none
}

func NewInfraredReceiver() *InfraredReceiver {
	ir := &InfraredReceiver{
		InputPin: NewInputPin(machine.PC6),
		decoder:  remote.NewDecoder(),
	}
	err := ir.SetInterrupt(machine.PinToggle, ir.edge)
	if err != nil {
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
		if code, ok := this.decoder.Pulse(this.lastMark, interval); ok {
			this.code = code
		}
	}
	this.lastMark = interval
}

// Take returns the last decoded code and clears it, 0 when nothing was received.
func (this *InfraredReceiver) Take() int {
	mask := interrupt.Disable()
	code := this.code
	this.code = 0
	interrupt.Restore(mask)
	return code
}
