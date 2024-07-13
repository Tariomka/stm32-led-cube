package controller

import "machine"

// Pin that's setup for Input mode
type InputPin struct{ Pin machine.Pin }

func NewInputPin(pin machine.Pin) InputPin {
	input := InputPin{Pin: pin}
	input.Pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	return input
}

// Pin that's setup for Output mode
type OutputPin struct{ Pin machine.Pin }

func NewOutputPin(pin machine.Pin) OutputPin {
	output := OutputPin{Pin: pin}
	output.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return output
}

func NewSpiOutput() machine.SPI {
	spi := machine.SPI0
	spi.Configure(machine.SPIConfig{
		Frequency: 100_000,
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
		LSBFirst:  true,
		Mode:      1,
		// DataBits:  16,
	})

	return spi
}
