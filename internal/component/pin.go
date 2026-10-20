package component

import (
	"device/arm"
	"device/stm32"
	"machine"
)

// Pin that's setup for Input mode
type InputPin struct{ machine.Pin }

// NewInputPin creates a new InputPin with input pull-up configuration
func NewInputPin(pin machine.Pin) InputPin {
	input := InputPin{Pin: pin}
	// This is a workaround for setting up PinInputPullup because the TinyGo compiler has a single bit leakage bug
	input.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	input.High() // ODR=1 selects pull-up
	return input
}

// DisableDebugPort frees PA13, PA14, PA15, PB3 and PB4 by enabling AFIO Clock and disabling the SWJ debug port
func DisableDebugPort() {
	stm32.RCC.SetAPB2ENR_AFIOEN(stm32.RCC_APB2ENR_AFIOEN_Enabled)
	stm32.AFIO.SetMAPR_SWJ_CFG(0b100) // AFIO_REMAP_SWJ_DISABLE
}

// Pin that's setup for Output mode
type OutputPin struct{ machine.Pin }

func NewOutputPin(pin machine.Pin) OutputPin {
	output := OutputPin{Pin: pin}
	output.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return output
}

func NewUart(uart *machine.UART, tx machine.Pin, rx machine.Pin) *machine.UART {
	uart.Configure(machine.UARTConfig{
		BaudRate: 115200,
		TX:       tx,
		RX:       rx,
	})
	return uart
}

func NewSpiOutput() *machine.SPI {
	spi := machine.SPI0
	spi.Configure(machine.SPIConfig{
		Frequency: 25 * machine.MHz, // Max should be around 32mhz
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
		LSBFirst:  true,
		Mode:      1,
	})

	return spi
}

func NewOnBoardI2C() *machine.I2C {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 100 * machine.KHz,
		SCL:       machine.I2C0_SCL_PIN,
		SDA:       machine.I2C0_SDA_PIN,
	})
	return i2c
}

// lowerPinInterruptPriority lets the runtime tick timer (priority 0xC1) preempt the pin's EXTI handler.
// EXTI defaults to the highest priority, and time.Now() inside it spins forever when the tick timer overflows meanwhile.
func lowerPinInterruptPriority(pin machine.Pin) {
	irq := uint32(stm32.IRQ_EXTI15_10)
	switch line := uint32(pin % 16); {
	case line <= 4:
		irq = stm32.IRQ_EXTI0 + line
	case line <= 9:
		irq = stm32.IRQ_EXTI9_5
	}
	arm.SetPriority(irq, 0xE0)
}
