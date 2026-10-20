package controller

import (
	"device/arm"
	"unsafe"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// SystemController owns the parts that concern the whole microcontroller: the debug port pins,
// the watchdog, the status LEDs, resets and memory reads.
type SystemController struct {
	leds          component.OnboardLeds // green D1 - PB9, red D2 - PB8
	watchdogReset bool                  // the watchdog caused the last reset
}

// NewSystemController frees the debug port pins (KEY4-6 and SW3 share them) and starts the
// watchdog; create it before the other controllers.
func NewSystemController() *SystemController {
	component.DisableDebugPort()
	return &SystemController{
		leds:          component.NewOnboardLeds(),
		watchdogReset: component.StartWatchdog(),
	}
}

// Startup prints the boot banner and blinks the green LED.
func (this *SystemController) Startup(cubeType string) {
	println("Led Cube is starting up")
	if this.watchdogReset {
		println("Watchdog reset!")
	}
	println("T='" + cubeType + "'")
	this.leds.BlinkStartup()
}

// Feed keeps the watchdog from resetting the microcontroller; call it at least every second.
func (this *SystemController) Feed() { component.FeedWatchdog() }

func (this *SystemController) Reset() { arm.SystemReset() }

// ReadMemory reads flash, SRAM and system memory; peripherals are refused because reads can have
// side effects.
func (this *SystemController) ReadMemory(address uint32) (byte, bool) {
	isAddressValid := (address >= 0x08000000 && address < 0x08080000) ||
		(address >= 0x20000000 && address < 0x20010000) ||
		(address >= 0x1FFFF000 && address < 0x1FFFF810)
	if !isAddressValid {
		return 0, false
	}

	return *(*byte)(unsafe.Pointer(uintptr(address))), true
}
