package component

import (
	"machine"

	"github.com/Tariomka/stm32-led-cube/internal/cascade"
)

// CascadePort is the 3.5 mm jack link to the neighboring cubes: USART2, TX - PA2, RX - PA3.
type CascadePort struct {
	uart   *machine.UART
	parser cascade.Parser
}

func NewCascadePort() *CascadePort {
	port := &CascadePort{
		uart: NewUart(machine.UART2, machine.PA2, machine.PA3),
	}
	NewInputPin(machine.PA3) // pull-up: an unplugged jack must not produce noise
	return port
}

func (this *CascadePort) Send(payload []byte) { this.uart.Write(cascade.Encode(payload)) }

// Poll passes every complete received payload to handle.
func (this *CascadePort) Poll(handle func(payload []byte)) {
	for this.uart.Buffered() > 0 {
		b, err := this.uart.ReadByte()
		if err != nil {
			return
		}

		if payload, ok := this.parser.Feed(b); ok {
			handle(payload)
		}
	}
}
