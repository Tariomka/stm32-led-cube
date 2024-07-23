package component

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

// (MBI5024GP/GF)
func NewSpiOutput() machine.SPI {
	spi := machine.SPI1
	spi.Configure(machine.SPIConfig{
		Frequency: 100_000,
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
		LSBFirst:  true,
		Mode:      1,
	})

	return spi
}

// (24C02)
func NewOnBoardI2C() *machine.I2C {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		Frequency: 100_000,
		SCL:       machine.I2C0_SCL_PIN,
		SDA:       machine.I2C0_SDA_PIN,
	})
	return i2c
}
