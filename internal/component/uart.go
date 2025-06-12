package component

import (
	"machine"

	"github.com/Tariomka/stm32-led-cube/internal/common/global"
)

type UART struct {
	*machine.UART
	buffer []byte
}

func NewUart(uart *machine.UART) *UART {
	device := UART{
		UART:   uart,
		buffer: make([]byte, 4*1024),
	}
	device.Configure(machine.UARTConfig{})

	return &device
}

// WIP
func (this *UART) Send(message string) {
	if _, err := this.Write([]byte(message)); err != nil {
		global.HandleError()
	}
}

// WIP
func (this *UART) Receive() []byte {
	n, err := this.Read(this.buffer)
	if err != nil {
		global.HandleError()
		return nil
	}

	return this.buffer[:n]
}
