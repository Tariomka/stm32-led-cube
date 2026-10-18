package controller

import (
	"machine"
	"runtime/interrupt"
	"strconv"
	"time"

	"github.com/Tariomka/led-common-lib/pkg/led"
	"github.com/Tariomka/led-common-lib/pkg/network"
	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/console"
	"github.com/Tariomka/stm32-led-cube/internal/engine/frame"
	"github.com/Tariomka/stm32-led-cube/internal/settings"
)

const KeyCount = 7

type Board interface {
	Show(voxels []byte)              // Presents an 8x8x8 engine frame (index x + y*8 + z*64, 0bRRGGBB)
	ShowLayout(s led.Slicer)         // Presents a led-common-lib layout at full brightness
	SetBlank(blank bool)             // Turns all LEDs off / back on
	TakeKeys() uint8                 // Keys pressed since the last call, bit n = KEY(n+1)
	KeyHeld(key int) bool            // Whether KEY(key+1) is currently held down
	ShowGroup() int                  // Show group 0..3 selected with DIP SW1/SW2
	AudioLevel(mode int) int         // Microphone loudness 0..63 (audio cartoon level mode)
	AudioRaw() int                   // Latest 12-bit microphone sample
	TakeRemote() int                 // Last IR remote code since the previous call, 0 = none
	PollSerial()                     // Runs pending console commands, buffers Pico link data
	CascadeSend(payload []byte)      // Sends a payload to the next cube
	CascadePoll(handle func([]byte)) // Passes payloads received from the previous cube
	LoadSettings() (settings.Settings, bool)
	SaveSettings(s settings.Settings) error
	Send(message string)               // Sends message over UART
	Receive() []byte                   // Receives data from UART
	GetIndicator() component.Indicator // Gets onboard led indicator
}

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Scanner *component.Scanner // MBI5024 chain (SPI1 PA5/PA7, LE PC4, OE PC5) + 74HC154 layer select
	Mapping frame.Mapping      // Orientation of the show coordinates on the physical cube

	Serial        *component.SerialPort  // UART1, RX - PA10, TX - PA9: console + Pico link
	UartProcessor *network.UartProcessor // Pico link protocol on Serial
	Console       *console.Console       // Vendor-style console on Serial
	EEPROM        component.EEPROM       // 24C02, SCL - PB6, SDA - PB7
	Cascade       *component.CascadePort // 3.5 mm jacks, TX - PA2, RX - PA3

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

	Microphone *component.Microphone       // PA0
	Infrared   *component.InfraredReceiver // PC6
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
		EEPROM:      component.NewEEPROM(),
		Cascade:     component.NewCascadePort(),
		LedsOnboard: component.NewOnboardLeds(machine.PB9, machine.PB8),
		Switch1:     component.NewInputPin(machine.PA1),
		Switch2:     component.NewInputPin(machine.PC3),
		Switch3:     component.NewInputPin(machine.PA15),
		Switch4:     component.NewInputPin(machine.PA12),
		Microphone:  component.NewMicrophone(),
		Infrared:    component.NewInfraredReceiver(machine.PC6),
	}
	board.Console = &console.Console{Device: &board, Registry: loadRegistry()}
	board.Serial = component.NewSerialPort(machine.UART1, board.Console.Run)
	board.Console.Write = board.Serial.WriteString
	board.UartProcessor = network.NewUartProcessor(board.Serial)

	for i, pin := range [KeyCount]machine.Pin{machine.PC0, machine.PC1, machine.PC2, machine.PB3, machine.PA14, machine.PA13, machine.PA11} {
		board.Keys[i] = component.NewButton(pin)
		bit := uint8(1) << i
		board.Keys[i].OnPress(func() { pendingKeys |= bit })
	}
	board.Scanner.Start()
	board.Microphone.Start()
	board.Infrared.Start()

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
		println("T='" + board.CubeType() + "'")
		println("Config1=" + strconv.Itoa(board.ShowGroup()) + ",Config2=" + strconv.Itoa(board.config2()))
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

// config2 reads SW3/SW4 like the original (bit set = switch ON); it only prints it.
func (this *YellowBoard) config2() int {
	value := 0
	if !this.Switch3.Get() {
		value |= 1
	}
	if !this.Switch4.Get() {
		value |= 2
	}
	return value
}

func (this *YellowBoard) AudioLevel(mode int) int {
	return this.Microphone.Level(mode)
}

func (this *YellowBoard) AudioRaw() int {
	return this.Microphone.Raw()
}

func (this *YellowBoard) TakeRemote() int {
	return this.Infrared.Take()
}

func (this *YellowBoard) CascadeSend(payload []byte) {
	this.Cascade.Send(payload)
}

func (this *YellowBoard) CascadePoll(handle func([]byte)) {
	this.Cascade.Poll(handle)
}

func (this *YellowBoard) Send(message string) {
	this.UartProcessor.WriteMessage(message)
}

func (this *YellowBoard) Receive() []byte {
	dataType, content, err := this.UartProcessor.Read()
	if err != nil {
		global.HandleError()
		return nil
	}

	switch dataType {
	case network.UartPing:
		this.UartProcessor.SendPong()
	}

	return content
}

func (this *YellowBoard) GetIndicator() component.Indicator {
	return this.LedsOnboard
}
