package component

import (
	"device/stm32"
	"machine"
	"time"
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

const debounceTime = 30 * time.Millisecond

// Active-low push button with a debounced press interrupt
type Button struct {
	InputPin
	lastPress time.Time
}

func NewButton(pin machine.Pin) *Button {
	return &Button{InputPin: NewInputPin(pin)}
}

func (this *Button) Pressed() bool {
	return !this.Get()
}

// OnPress registers a callback that runs in interrupt context once per debounced press
func (this *Button) OnPress(callback func()) error {
	return this.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		now := time.Now()
		if now.Sub(this.lastPress) < debounceTime {
			return
		}
		this.lastPress = now
		callback()
	})
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

func NewSpiOutput(spi *machine.SPI, sck, sdo machine.Pin) *machine.SPI {
	spi.Configure(machine.SPIConfig{
		Frequency: 25 * machine.MHz, // Max should be around 32mhz
		SCK:       sck,
		SDO:       sdo,
		LSBFirst:  true,
		Mode:      1,
	})

	return spi
}

type UART struct{ *machine.UART }

func NewConfiguredUart(uart *machine.UART) UART {
	device := UART{UART: uart}
	device.Configure(machine.UARTConfig{})

	return device
}

// func NewOnBoardI2C() *machine.I2C {
// 	i2c := machine.I2C0
// 	i2c.Configure(machine.I2CConfig{
// 		Frequency: 100_000,
// 		SCL:       machine.I2C0_SCL_PIN,
// 		SDA:       machine.I2C0_SDA_PIN,
// 	})
// 	return i2c
// }
