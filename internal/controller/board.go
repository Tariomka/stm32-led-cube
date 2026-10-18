package controller

import (
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/Tariomka/led-common-lib/pkg/led"
	"github.com/Tariomka/led-common-lib/pkg/network"
	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/engine/frame"
)

const KeyCount = 7

type Board interface {
	Show(voxels []byte)                // Presents an 8x8x8 engine frame (index x + y*8 + z*64, 0bRRGGBB)
	ShowLayout(s led.Slicer)           // Presents a led-common-lib layout at full brightness
	SetBlank(blank bool)               // Turns all LEDs off / back on
	TakeKeys() uint8                   // Keys pressed since the last call, bit n = KEY(n+1)
	KeyHeld(key int) bool              // Whether KEY(key+1) is currently held down
	ShowGroup() int                    // Show group 0..3 selected with DIP SW1/SW2
	Send(message string)               // Sends message over UART
	Receive() []byte                   // Receives data from UART
	GetIndicator() component.Indicator // Gets onboard led indicator
}

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Scanner *component.Scanner // MBI5024 chain (SPI1 PA5/PA7, LE PC4, OE PC5) + 74HC154 layer select
	Mapping frame.Mapping      // Orientation of the show coordinates on the physical cube

	UartProcessor *network.UartProcessor // Onboard UART serial connection processor, RX - PA10, TX - PA9
	// UartMainBoard *machine.UART // Main board UART serial connection (3.5 jacks), RX - PA3, TX - PA2
	// I2C *machine.I2C // (24C02) , SDA - PB7, SCL - PB6

	// Onboard leds:
	//
	// Green led -> D1, pin PB9
	//
	// Red led -> D2, pin PB8
	LedsOnboard component.OnboardLeds

	// KEY1 PC0 Previous, KEY2 PC1 Next, KEY3 PC2 Speed+, KEY4 PB3 Speed-,
	// KEY5 PA14 Run/Pause, KEY6 PA13 Cycle, KEY7 PA11 On/Off
	Keys [KeyCount]*component.Button

	// DIP switches, ON = Low. SW1+SW2 select the show group (applied at reset)
	Switch1 component.InputPin // PA1
	Switch2 component.InputPin // PC3
	Switch3 component.InputPin // PA15
	Switch4 component.InputPin // PA12

	// InfraRed // ?, PC6
}

// Keys pressed in interrupt context, consumed by TakeKeys
var pendingKeys uint8

func NewYellowBoard() Board {
	// KEY4-6 and SW3 share pins with the JTAG/SWD port
	component.DisableDebugPort()

	board := YellowBoard{
		Scanner: component.NewScanner(
			machine.SPI0,
			machine.SPI0_SCK_PIN,
			machine.SPI0_SDO_PIN,
			machine.PC4,
			machine.PC5,
			component.NewDemultiplexer(
				machine.PB0,
				machine.PB1,
				machine.PB10,
				machine.PB11,
				machine.PA8,
				machine.PC7,
			),
		),
		UartProcessor: network.NewUartProcessor(component.NewConfiguredUart(machine.UART1)),
		// I2C: NewOnBoardI2C(),
		LedsOnboard: component.NewOnboardLeds(machine.PB9, machine.PB8),
		Switch1:     component.NewInputPin(machine.PA1),
		Switch2:     component.NewInputPin(machine.PC3),
		Switch3:     component.NewInputPin(machine.PA15),
		Switch4:     component.NewInputPin(machine.PA12),
	}
	for i, pin := range [KeyCount]machine.Pin{machine.PC0, machine.PC1, machine.PC2, machine.PB3, machine.PA14, machine.PA13, machine.PA11} {
		board.Keys[i] = component.NewButton(pin)
		bit := uint8(1) << i
		board.Keys[i].OnPress(func() { pendingKeys |= bit })
	}
	board.Scanner.Start()

	global.HandleFatal = func() {
		board.LedsOnboard.LedGreen.High()
		board.LedsOnboard.LedRed.Low()
		board.Scanner.SetBlank(true)
		for {
			time.Sleep(1 * time.Second)
		}
	}
	global.HandleError = board.LedsOnboard.BlinkError
	global.HandleStartup = func() {
		println("Led Cube is starting up")
		board.LedsOnboard.BlinkStartup()
	}
	global.HandlePrint = board.Send

	return &board
}

func (this *YellowBoard) Show(voxels []byte) {
	frame.Encode(voxels, this.Mapping, this.Scanner.Back())
	this.Scanner.Present()
}

func (this *YellowBoard) ShowLayout(s led.Slicer) {
	back := this.Scanner.Back()
	for index, slice := range s.IterateSlices() {
		if int(index) < frame.Layers && len(slice) >= frame.BytesPerLayer {
			frame.FromLayer(back, int(index), slice)
		}
	}
	this.Scanner.Present()
}

func (this *YellowBoard) SetBlank(blank bool) {
	this.Scanner.SetBlank(blank)
}

func (this *YellowBoard) TakeKeys() uint8 {
	mask := interrupt.Disable()
	keys := pendingKeys
	pendingKeys = 0
	interrupt.Restore(mask)
	return keys
}

func (this *YellowBoard) KeyHeld(key int) bool {
	return key >= 0 && key < KeyCount && this.Keys[key].Pressed()
}

// ShowGroup follows the base board label: SW1/SW2 = 1/1 dynamic, 0/1 audio, 1/0 test, 0/0 text,
// where 1 means the pin reads high (switch off).
func (this *YellowBoard) ShowGroup() int {
	sw1, sw2 := this.Switch1.Get(), this.Switch2.Get()
	switch {
	case sw1 && sw2:
		return 0
	case !sw1 && sw2:
		return 1
	case sw1 && !sw2:
		return 2
	}
	return 3
}

func (this *YellowBoard) Send(message string) {
	this.UartProcessor.WriteMessage(message)
}

func (this *YellowBoard) Receive() []byte {
	dtype, content, err := this.UartProcessor.Read()
	if err != nil {
		global.HandleError()
		return nil
	}

	switch dtype {
	case network.UartPing:
		this.UartProcessor.SendPong()
	}

	return content
}

func (this *YellowBoard) GetIndicator() component.Indicator {
	return this.LedsOnboard
}
