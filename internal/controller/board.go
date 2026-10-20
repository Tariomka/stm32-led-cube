package controller

import (
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/console"
	"github.com/Tariomka/stm32-led-cube/internal/engine/frame"
	"github.com/Tariomka/stm32-led-cube/internal/settings"
	"github.com/Tariomka/stm32-led-cube/internal/upload"
)

const KeyCount = 7

type Board interface {
	Show(voxels []byte)              // Presents an 8x8x8 engine frame (index x + y*8 + z*64, 0bRRGGBB)
	SetBlank(blank bool)             // Turns all LEDs off / back on
	TakeKeys() uint8                 // Keys pressed since the last call, bit n = KEY(n+1)
	KeyHeld(key int) bool            // Whether KEY(key+1) is currently held down
	AudioLevel(mode int) int         // Microphone loudness 0..63 (audio cartoon level mode)
	AudioRaw() int                   // Latest 12-bit microphone sample
	TakeRemote() int                 // Last IR remote code since the previous call, 0 = none
	PollSerial()                     // Runs console commands and show uploads received on UART1
	CascadeSend(payload []byte)      // Sends a payload to the next cube
	CascadePoll(handle func([]byte)) // Passes payloads received from the previous cube
	LoadSettings() (settings.Settings, bool)
	SaveSettings(s settings.Settings) error
	CubeType() string // Vendor type string from the EEPROM, e.g. "8x8x8A"
	// OnUpload sets the callbacks of a show upload from the PC tool or the Pico
	OnUpload(started func(), finished func(ok bool))
	// SetRemote lets console commands drive the shows (Key, Action, Status)
	SetRemote(remote console.Remote)
}

// ICubeSmart controller board, powered by GD32F103RET6
type YellowBoard struct {
	Scanner *component.Scanner // MBI5024 chain (SPI1 PA5/PA7, LE PC4, OE PC5) + 74HC154 layer select
	Mapping frame.Mapping      // Orientation of the show coordinates on the physical cube

	Serial  *component.SerialPort  // UART1, RX - PA10, TX - PA9: console, uploads and log output
	Console *console.Console       // Console commands on Serial
	Upload  *upload.Receiver       // Show uploads on Serial
	EEPROM  component.EEPROM       // 24C02, SCL - PB6, SDA - PB7
	Cascade *component.CascadePort // 3.5 mm jacks, TX - PA2, RX - PA3

	// Onboard leds:
	//
	// Green led -> D1, pin PB9
	//
	// Red led -> D2, pin PB8
	LedsOnboard component.OnboardLeds

	// KEY1 PC0 Previous, KEY2 PC1 Next, KEY3 PC2 Speed+, KEY4 PB3 Speed-,
	// KEY5 PA14 Run/Pause, KEY6 PA13 Play model, KEY7 PA11 On/Off
	Keys [KeyCount]*component.Button

	Microphone *component.Microphone       // PA0
	Infrared   *component.InfraredReceiver // PC6

	rebootAt time.Time // pending console ReBoot, zero if none
}

// Keys pressed in interrupt context, consumed by TakeKeys
var pendingKeys uint8

func NewYellowBoard() Board {
	// KEY4-6 and SW3 share pins with the JTAG/SWD port
	component.DisableDebugPort()

	board := YellowBoard{
		Scanner:     component.NewScanner(),
		EEPROM:      component.NewEEPROM(),
		Cascade:     component.NewCascadePort(),
		LedsOnboard: component.NewOnboardLeds(),
		Microphone:  component.NewMicrophone(),
		Infrared:    component.NewInfraredReceiver(),
		Upload:      &upload.Receiver{Storage: &component.ShowSlotStorage{}},
	}
	board.Console = &console.Console{Device: &board, Registry: loadRegistry()}
	board.Serial = component.NewSerialPort(board.Console.Run, board.Upload.Frame)
	board.Upload.Reply = func(payload []byte) { board.Serial.Write(upload.Frame(payload)) }
	board.Console.Write = board.Serial.WriteString

	for index, pin := range [KeyCount]machine.Pin{machine.PC0, machine.PC1, machine.PC2, machine.PB3, machine.PA14, machine.PA13, machine.PA11} {
		board.Keys[index] = component.NewButton(pin)
		bit := uint8(1) << index
		err := board.Keys[index].OnPress(func() { pendingKeys |= bit })
		if err != nil {
			println("Button '", index, "' OnPress error:", err)
		}
	}

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
		board.LedsOnboard.BlinkStartup()
	}

	return &board
}

func (this *YellowBoard) Show(voxels []byte) {
	frame.Encode(voxels, this.Mapping, this.Scanner.Back())
	this.Scanner.Present()
}

func (this *YellowBoard) SetBlank(blank bool) { this.Scanner.SetBlank(blank) }

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

func (this *YellowBoard) AudioLevel(mode int) int { return this.Microphone.Level(mode) }

func (this *YellowBoard) AudioRaw() int { return this.Microphone.Raw() }

func (this *YellowBoard) TakeRemote() int { return this.Infrared.Take() }

func (this *YellowBoard) CascadeSend(payload []byte) { this.Cascade.Send(payload) }

func (this *YellowBoard) CascadePoll(handle func([]byte)) { this.Cascade.Poll(handle) }

func (this *YellowBoard) SetRemote(remote console.Remote) { this.Console.Remote = remote }
