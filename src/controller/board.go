package controller

import "machine"

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Demultiplexer Demultiplexer // Provides power to each layer of led cube
	LedDriver     LedDriver     // Colors leds in a layer

	UartOnboard *machine.UART // Onboard UART serial connection, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection, RX - PA3, TX - PA2

	LedGreen OutputPin // Onboard Green Led - D1, pin PB9
	LedRed   OutputPin // Onboard Red Led - D2, pin PB8

	ButtonPrevious  InputPin // Main board KEY1, PC0
	ButtonNext      InputPin // Main board KEY2, PC1
	ButtonSpeedMore InputPin // Main board KEY3, PC2
	ButtonSpeedLess InputPin // Main board KEY4, PB3
	ButtonRunPause  InputPin // Main board KEY5, PA14
	ButtonCycle     InputPin // Main board KEY6, PA13
	ButtonOnOff     InputPin // Main board KEY7, PA11
}

func NewYellowBoard() YellowBoard {
	board := YellowBoard{
		LedDriver:     NewLedDriver(machine.PC4, machine.PC5),
		Demultiplexer: NewDemultiplexer(machine.PB0, machine.PB1, machine.PB10, machine.PB11),

		UartOnboard: machine.UART1,
		// UartMainBoard: machine.UART2,

		LedGreen: NewOutputPin(machine.PB9),
		LedRed:   NewOutputPin(machine.PB8),

		ButtonPrevious:  NewInputPin(machine.PC0),
		ButtonNext:      NewInputPin(machine.PC1),
		ButtonSpeedMore: NewInputPin(machine.PC2),
		ButtonSpeedLess: NewInputPin(machine.PB3),
		ButtonRunPause:  NewInputPin(machine.PA14),
		ButtonCycle:     NewInputPin(machine.PA13),
		ButtonOnOff:     NewInputPin(machine.PA11),
	}

	board.UartOnboard.Configure(machine.UARTConfig{})
	// board.UartMainBoard.Configure(machine.UARTConfig{
	// 	BaudRate: 38400,
	// 	TX:       machine.PA3,
	// 	RX:       machine.PA2,
	// })

	return board
}
