package controller

import (
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

type Board interface {
	LightLeds(s Slicer)                          // Lights up a single frame
	Run(lw LayoutWorker, lightShows []LightShow) // Main loop
	Blinker
}

type Blinker interface {
	BlinkStartup()
	BlinkError()
}

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Demultiplexer component.Demultiplexer // (74HC154) Provides power to each layer of led cube
	LedDriver     *component.LedDriver    // (MBI5024GP/GF) Shift register led driver

	UartOnboard *machine.UART // Onboard UART serial connection, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection, RX - PA3, TX - PA2
	// I2C *machine.I2C // (24C02) , SDA - PB7, SCL - PB6

	LedGreen component.OutputPin // Onboard Green Led - D1, pin PB9. Cathode control
	LedRed   component.OutputPin // Onboard Red Led - D2, pin PB8. Cathode control

	ButtonPrevious  component.InputPin // Main board KEY1, PC0
	ButtonNext      component.InputPin // Main board KEY2, PC1
	ButtonSpeedMore component.InputPin // Main board KEY3, PC2
	ButtonSpeedLess component.InputPin // Main board KEY4, PB3
	ButtonRunPause  component.InputPin // Main board KEY5, PA14
	ButtonCycle     component.InputPin // Main board KEY6, PA13
	ButtonOnOff     component.InputPin // Main board KEY7, PA11

	// InfraRed // ?, PC6

	speed        uint32
	programIndex uint32
}

// func NewYellowBoard() *YellowBoard {
func NewYellowBoard() Board {
	board := YellowBoard{
		Demultiplexer: component.NewDemultiplexer(
			machine.PB0,
			machine.PB1,
			machine.PB10,
			machine.PB11,
			machine.PA8,
			machine.PC7,
		),
		LedDriver: component.NewLedDriver(
			machine.SPI0,
			machine.SPI0_SCK_PIN,
			machine.SPI0_SDO_PIN,
			machine.PC4,
			machine.PC5,
		),

		UartOnboard: machine.UART1,
		// I2C: NewOnBoardI2C(),

		LedGreen: component.NewOutputPin(machine.PB9),
		LedRed:   component.NewOutputPin(machine.PB8),

		ButtonPrevious:  component.NewInputPin(machine.PC0),
		ButtonNext:      component.NewInputPin(machine.PC1),
		ButtonSpeedMore: component.NewInputPin(machine.PC2),
		ButtonSpeedLess: component.NewInputPin(machine.PB3),
		ButtonRunPause:  component.NewInputPin(machine.PA14),
		ButtonCycle:     component.NewInputPin(machine.PA13),
		ButtonOnOff:     component.NewInputPin(machine.PA11),

		speed:        10,
		programIndex: 0,
	}

	board.UartOnboard.Configure(machine.UARTConfig{BaudRate: 38400})

	// board.ButtonPrevious.Pin.SetInterrupt(machine.PinFalling)
	// board.ButtonNext
	// board.ButtonSpeedMore
	// board.ButtonSpeedLess
	// board.ButtonRunPause
	// board.ButtonCycle
	// board.ButtonOnOff

	return &board
}

func (yb *YellowBoard) Run(lw LayoutWorker, lightShows []LightShow) {
	if len(lightShows) < 1 {
		yb.BlinkError()
		return
	}

	for {
		for _, lightShow := range lightShows[yb.programIndex] {
			lw.ResetBlock()
			lightShow(lw)
			for i := uint32(0); i < yb.speed; i++ {
				yb.LightLeds(lw)
			}
		}
	}
}

func (yb *YellowBoard) LightLeds(s Slicer) {
	for index, slice := range s.IterateSlices() {
		if err := yb.LedDriver.LightLayer(slice); err != nil {
			yb.BlinkError()
		}

		if err := yb.Demultiplexer.EnableLayer(index); err != nil {
			yb.BlinkError()
		}
	}
}

func (yb *YellowBoard) BlinkStartup() {
	for i := 0; i < 3; i++ {
		yb.LedRed.Pin.Low()
		yb.LedGreen.Pin.Low()
		time.Sleep(100 * time.Millisecond)

		yb.LedRed.Pin.High()
		yb.LedGreen.Pin.High()
		time.Sleep(500 * time.Millisecond)
	}
}

func (yb *YellowBoard) BlinkError() {
	for i := 0; i < 5; i++ {
		yb.LedRed.Pin.Low()
		time.Sleep(100 * time.Millisecond)

		yb.LedRed.Pin.High()
		time.Sleep(100 * time.Millisecond)
	}
}
