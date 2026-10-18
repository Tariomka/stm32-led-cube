package component

import (
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/console"
)

// SerialPort shares UART1 between the vendor-style console and the Pico link: console
// commands are executed as they arrive, everything else is buffered for the link reader.
type SerialPort struct {
	uart    *machine.UART
	demux   console.Demux
	start   time.Time
	pending [1024]byte
	head    int
	size    int
}

func NewSerialPort(uart *machine.UART, command func(line string)) *SerialPort {
	uart.Configure(machine.UARTConfig{})
	return &SerialPort{uart: uart, demux: console.Demux{Command: command}, start: time.Now()}
}

// Poll moves received bytes through the demultiplexer; call it regularly.
func (this *SerialPort) Poll() {
	now := uint32(time.Since(this.start).Milliseconds())
	for this.uart.Buffered() > 0 {
		b, err := this.uart.ReadByte()
		if err != nil {
			return
		}
		if this.demux.Feed(b, now) && this.size < len(this.pending) {
			this.pending[(this.head+this.size)%len(this.pending)] = b
			this.size++
		}
	}
}

// Read returns buffered Pico link bytes (never blocks).
func (this *SerialPort) Read(p []byte) (int, error) {
	this.Poll()
	n := 0
	for n < len(p) && this.size > 0 {
		p[n] = this.pending[this.head]
		this.head = (this.head + 1) % len(this.pending)
		this.size--
		n++
	}
	return n, nil
}

func (this *SerialPort) Write(p []byte) (int, error) {
	return this.uart.Write(p)
}

func (this *SerialPort) WriteString(s string) {
	this.uart.Write([]byte(s))
}
