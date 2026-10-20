package component

import (
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/console"
)

// SerialPort runs UART1: console commands and upload frames are handled as they arrive.
type SerialPort struct {
	uart  *machine.UART
	demux console.Demux
	start time.Time
}

func NewSerialPort(command func(line string), frame func(payload []byte)) *SerialPort {
	return &SerialPort{
		uart:  NewUart(machine.UART1, machine.UART_TX_PIN, machine.UART_RX_PIN),
		demux: console.Demux{Command: command, Frame: frame},
		start: time.Now(),
	}
}

// Poll moves received bytes through the demultiplexer; call it regularly.
func (this *SerialPort) Poll() {
	now := uint32(time.Since(this.start).Milliseconds())
	for this.uart.Buffered() > 0 {
		b, err := this.uart.ReadByte()
		if err != nil {
			return
		}
		this.demux.Feed(b, now)
	}
}

func (this *SerialPort) Write(p []byte) (int, error) {
	return this.uart.Write(p)
}

func (this *SerialPort) WriteString(s string) {
	this.uart.Write([]byte(s))
}
