package component

import "machine"

type UART struct {
	*machine.UART
	buffer      []byte
	handleError func()
}

func NewUart(uart *machine.UART, handleErrorCallback func()) *UART {
	device := UART{
		UART:        uart,
		buffer:      make([]byte, 4*1024),
		handleError: handleErrorCallback,
	}
	device.Configure(machine.UARTConfig{})

	return &device
}

func (this *UART) Send(message string) {
	if _, err := this.Write([]byte(message)); err != nil {
		this.handleError()
	}
}

func (this *UART) Receive() []byte {
	n, err := this.Read(this.buffer)
	if err != nil {
		this.handleError()
		return nil
	}

	return this.buffer[:n]
}
