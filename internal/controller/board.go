package controller

import (
	"machine"
	"time"

	"github.com/Tariomka/led-common-lib/pkg/led"
	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/component/registers"
)

type Board interface {
	LightLeds(s led.Slicer)            // Lights up a single frame
	Send(message string)               // Sends message over UART
	Receive() []byte                   // Receives data from UART
	GetIndicator() component.Indicator // Gets onboard led indicator

	DisableLeds()
	EnableLeds()
}

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Demultiplexer component.Demultiplexer // (74HC154) Provides power to each layer of led cube
	LedDriver     *component.LedDriver    // (MBI5024GP/GF) Shift register led driver

	UartOnboard *component.UART // Onboard UART serial connection, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection (3.5 jacks), RX - PA3, TX - PA2
	// I2C *machine.I2C // (24C02) , SDA - PB7, SCL - PB6

	// Onboard leds:
	//
	// Green led -> D1, pin PB9
	//
	// Red led -> D2, pin PB8
	LedsOnboard component.OnboardLeds

	ButtonPrevious  component.InputPin // Main board KEY1, PC0
	ButtonNext      component.InputPin // Main board KEY2, PC1
	ButtonSpeedMore component.InputPin // Main board KEY3, PC2
	ButtonSpeedLess component.InputPin // Main board KEY4, PB3
	ButtonRunPause  component.InputPin // Main board KEY5, PA14
	ButtonCycle     component.InputPin // Main board KEY6, PA13
	ButtonOnOff     component.InputPin // Main board KEY7, PA11

	Switch1 component.InputPin // Main board switch 1, PA1
	Switch2 component.InputPin // Main board switch 2, PC3

	// InfraRed // ?, PC6
}

func NewYellowBoard(tracker *StateTracker) Board {
	registers.PrintAndResetCrashLog()
	registers.UpdateRegisters()

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
		UartOnboard: component.NewUart(machine.UART1),
		// I2C: NewOnBoardI2C(),
		LedsOnboard:     component.NewOnboardLeds(machine.PB9, machine.PB8),
		ButtonPrevious:  component.NewInputPin(machine.PC0),
		ButtonNext:      component.NewInputPin(machine.PC1),
		ButtonSpeedMore: component.NewInputPin(machine.PC2),
		ButtonSpeedLess: component.NewInputPin(machine.PB3),
		ButtonRunPause:  component.NewInputPin(machine.PA14),
		ButtonCycle:     component.NewInputPin(machine.PA13),
		ButtonOnOff:     component.NewInputPin(machine.PA11),
		Switch1:         component.NewInputPin(machine.PA1),
		Switch2:         component.NewInputPin(machine.PC3),
	}
	board.setInterrupts(tracker)

	global.HandleFatal = func() {
		board.LedsOnboard.LedGreen.High()
		board.LedsOnboard.LedRed.Low()
		board.Demultiplexer.Disable()
		for {
			time.Sleep(1 * time.Second)
		}
	}
	global.HandleError = board.LedsOnboard.BlinkError
	global.HandleStartup = func() {
		println("Led Cube is starting up")
		board.LedsOnboard.BlinkStartup()
	}

	// registers.PrintButtonStates()
	// registers.PrintRegisterValues()
	// registers.PrintInputOutputRegisterValues()
	return &board
}

func (this *YellowBoard) LightLeds(s led.Slicer) {
	for index, slice := range s.IterateSlices() {
		if err := this.LedDriver.LightLayer(slice); err != nil {
			this.LedsOnboard.BlinkError()
		}

		if err := this.Demultiplexer.EnableLayer(index); err != nil {
			this.LedsOnboard.BlinkError()
		}
	}
}

func (this *YellowBoard) Send(message string) {
	this.UartOnboard.Send(message)
}

func (this *YellowBoard) Receive() []byte {
	return this.UartOnboard.Receive()
}

func (this *YellowBoard) GetIndicator() component.Indicator {
	return this.LedsOnboard
}

// ===============================================
// To Be Deleted ?
func (this *YellowBoard) DisableLeds() {
	this.Demultiplexer.Disable()
}

func (this *YellowBoard) EnableLeds() {
	this.Demultiplexer.Enable()
}

// ===============================================

func (this *YellowBoard) setInterrupts(tracker *StateTracker) {
	this.ButtonPrevious.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 1 pressed") // TODO: Remove after debugging
		registers.PrintButtonStates()
		tracker.PrevLightShow()
	})
	this.ButtonNext.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 2 pressed") // TODO: Remove after debugging
		tracker.NextLightShow()
	})
	this.ButtonSpeedMore.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 3 pressed") // TODO: Remove after debugging
		tracker.IncreaseSpeed()
	})
	this.ButtonSpeedLess.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 4 pressed") // TODO: Remove after debugging
		tracker.DecreadeSpeed()
	})
	this.ButtonRunPause.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 5 pressed") // TODO: Remove after debugging
		tracker.SwitchRunPause()
	})
	this.ButtonCycle.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 6 pressed") // TODO: Remove after debugging
		tracker.CycleMode()
	})
	this.ButtonOnOff.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		println("Key 7 pressed") // TODO: Remove after debugging
		// TODO: add sleep mode logic
		tracker.SwitchRunPause() // TODO: remove after sleep logic is implemented
	})
}
