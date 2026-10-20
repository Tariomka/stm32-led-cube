package component

import "machine"

// SerialPort is a UART at 115200 8N1.
type SerialPort struct {
	uart *machine.UART
}

// NewConsolePort is UART1 on the core board header (TX - PA9, RX - PA10): console, uploads and log
// output (println writes here too).
func NewConsolePort() *SerialPort {
	return &SerialPort{uart: NewUart(machine.UART1, machine.UART_TX_PIN, machine.UART_RX_PIN)}
}

// NewCascadePort is UART2 on the 3.5 mm jacks (TX - PA2, RX - PA3) to the neighboring cubes.
func NewCascadePort() *SerialPort {
	port := &SerialPort{uart: NewUart(machine.UART2, machine.PA2, machine.PA3)}
	NewInputPin(machine.PA3) // pull-up: an unplugged jack must not produce noise
	return port
}

// Read returns the next received byte, ok = false when none is waiting.
func (this *SerialPort) Read() (b byte, ok bool) {
	if this.uart.Buffered() == 0 {
		return 0, false
	}

	b, err := this.uart.ReadByte()
	return b, err == nil
}

func (this *SerialPort) Write(data []byte) { this.uart.Write(data) }

func (this *SerialPort) WriteString(s string) { this.uart.Write([]byte(s)) }
