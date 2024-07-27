package controller

import (
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/src/component"
	"tinygo.org/x/drivers/shiftregister"
)

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Demultiplexer component.Demultiplexer // (74HC154) Provides power to each layer of led cube
	LedDriver     *shiftregister.Device   // (MBI5024GP/GF) Shift register led driver

	UartOnboard *machine.UART // Onboard UART serial connection, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection, RX - PA3, TX - PA2

	// I2C *machine.I2C // (24C02) , SDA - PB7, SCL - PB6

	LedGreen component.OutputPin // Onboard Green Led - D1, pin PB9. Cathode control
	LedRed   component.OutputPin // Onboard Red Led - D2, pin PB8. Cathode control

	ledDriverOutputEnable component.OutputPin

	ButtonPrevious  component.InputPin // Main board KEY1, PC0
	ButtonNext      component.InputPin // Main board KEY2, PC1
	ButtonSpeedMore component.InputPin // Main board KEY3, PC2
	ButtonSpeedLess component.InputPin // Main board KEY4, PB3
	ButtonRunPause  component.InputPin // Main board KEY5, PA14
	ButtonCycle     component.InputPin // Main board KEY6, PA13
	ButtonOnOff     component.InputPin // Main board KEY7, PA11

	// InfraRed // ?, PC6
}

func NewYellowBoard() *YellowBoard {
	board := YellowBoard{
		Demultiplexer: component.NewDemultiplexer(
			machine.PB0,
			machine.PB1,
			machine.PB10,
			machine.PB11,
			machine.PA8,
			machine.PC7,
		),
		LedDriver: shiftregister.New(
			shiftregister.SIXTEEN_BITS,
			machine.PC4,          // (Latch) Data strobe input terminal
			machine.SPI0_SCK_PIN, // (SPI_CLOCK) Clock input terminal for data shift on rising edge
			machine.SPI0_SDO_PIN, // (SPI_MOSI)Serial-data input to the shift register, PA7.
		),
		ledDriverOutputEnable: component.NewOutputPin(machine.PC5),

		// I2C: NewOnBoardI2C(),

		UartOnboard: machine.UART1,

		LedGreen: component.NewOutputPin(machine.PB9),
		LedRed:   component.NewOutputPin(machine.PB8),

		ButtonPrevious:  component.NewInputPin(machine.PC0),
		ButtonNext:      component.NewInputPin(machine.PC1),
		ButtonSpeedMore: component.NewInputPin(machine.PC2),
		ButtonSpeedLess: component.NewInputPin(machine.PB3),
		ButtonRunPause:  component.NewInputPin(machine.PA14),
		ButtonCycle:     component.NewInputPin(machine.PA13),
		ButtonOnOff:     component.NewInputPin(machine.PA11),
	}

	board.LedDriver.Configure()
	board.ledDriverOutputEnable.Pin.High()
	board.UartOnboard.Configure(machine.UARTConfig{BaudRate: 38400})

	return &board
}

func (yb *YellowBoard) LightLeds(ll LedLayout) {
	yb.ledDriverOutputEnable.Pin.Low()
	for z, layer := range ll {
		yb.Demultiplexer.EnableLayer(uint8(z))

		for _, packet := range layer {
			yb.LedDriver.WriteMask(packet)
		}
	}
	yb.ledDriverOutputEnable.Pin.High()
}

func (yb *YellowBoard) blinkError() {
	for i := 0; i < 3; i++ {
		yb.LedRed.Pin.Low()
		time.Sleep(100 * time.Millisecond)

		yb.LedRed.Pin.High()
		time.Sleep(100 * time.Millisecond)
	}
}
