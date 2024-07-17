package controller

import (
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/src/component"
)

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Demultiplexer component.Demultiplexer // Provides power to each layer of led cube
	LedDriver     *component.LedDriver    // Colors leds in a layer

	// I2C *machine.I2C // (24C02) , SDA - PB7, SCL - PB6

	// UartOnboard *machine.UART // Onboard UART serial connection, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection, RX - PA3, TX - PA2

	LedGreen component.OutputPin // Onboard Green Led - D1, pin PB9. Cathode control
	LedRed   component.OutputPin // Onboard Red Led - D2, pin PB8. Cathode control

	ButtonPrevious  component.InputPin // Main board KEY1, PC0
	ButtonNext      component.InputPin // Main board KEY2, PC1
	ButtonSpeedMore component.InputPin // Main board KEY3, PC2
	ButtonSpeedLess component.InputPin // Main board KEY4, PB3
	ButtonRunPause  component.InputPin // Main board KEY5, PA14
	ButtonCycle     component.InputPin // Main board KEY6, PA13
	ButtonOnOff     component.InputPin // Main board KEY7, PA11
}

func NewYellowBoard() *YellowBoard {
	board := YellowBoard{
		Demultiplexer: component.NewDemultiplexer(
			machine.PB0,
			machine.PB1,
			machine.PB10,
			machine.PB11,
			machine.PA8,
		),
		LedDriver: component.NewLedDriver(machine.PC4, machine.PC5),

		// I2C: NewOnBoardI2C(),

		// UartOnboard: machine.UART1,
		// UartMainBoard: machine.UART2,

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

	// board.UartOnboard.Configure(machine.UARTConfig{BaudRate: 38400})
	// board.UartMainBoard.Configure(machine.UARTConfig{
	// 	BaudRate: 38400,
	// 	TX:       machine.PA3,
	// 	RX:       machine.PA2,
	// })

	return &board
}

func (yb *YellowBoard) LightLeds(ll LedLayout) {
	for z, layer := range ll {
		err := yb.Demultiplexer.EnableLayer(uint8(z))
		if err != nil {
			yb.blinkError()
		}

		yb.LedDriver.LightLayer(layer[:])
	}
}

func (yb *YellowBoard) blinkError() {
	for i := 0; i < 3; i++ {
		yb.LedRed.Pin.Low()
		time.Sleep(100 * time.Millisecond)

		yb.LedRed.Pin.High()
		time.Sleep(100 * time.Millisecond)
	}
}
