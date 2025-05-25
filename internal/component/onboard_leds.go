package component

import (
	"machine"
	"time"
)

type BlinkType byte

const (
	BlinkRed BlinkType = iota
	BlinkGreen
	BlinkBoth
	BlinkStaggerRed
	BlinkStaggerGreen
)

type Indicator interface {
	BlinkStartup() // Startup indicator
	BlinkError()   // Error indicator
	Blink(bType BlinkType, duration time.Duration, delay bool)
}

type OnboardLeds struct {
	LedGreen OutputPin // Onboard Green Led - D1, pin PB9. Cathode control
	LedRed   OutputPin // Onboard Red Led - D2, pin PB8. Cathode control
}

func NewOnboardLeds(greenPin machine.Pin, redPin machine.Pin) OnboardLeds {
	return OnboardLeds{
		LedGreen: NewOutputPin(greenPin),
		LedRed:   NewOutputPin(redPin),
	}
}

func (this OnboardLeds) BlinkStartup() {
	for range 3 {
		this.Blink(BlinkStaggerGreen, 100*time.Millisecond, true)
	}
}

func (this OnboardLeds) BlinkError() {
	for range 5 {
		this.Blink(BlinkRed, 100*time.Millisecond, true)
	}
}

func (this OnboardLeds) Blink(bType BlinkType, duration time.Duration, delay bool) {
	switch bType {
	case BlinkRed:
		this.LedRed.Pin.Low()
	case BlinkGreen:
		this.LedGreen.Pin.Low()
	case BlinkBoth:
		this.LedGreen.Pin.Low()
		this.LedRed.Pin.Low()
	case BlinkStaggerRed:
		this.LedRed.Pin.Low()
		time.Sleep(duration)
		this.LedRed.Pin.High()
		this.LedGreen.Pin.Low()
	case BlinkStaggerGreen:
		this.LedGreen.Pin.Low()
		time.Sleep(duration)
		this.LedGreen.Pin.High()
		this.LedRed.Pin.Low()
	}
	time.Sleep(duration)

	this.LedGreen.Pin.High()
	this.LedRed.Pin.High()
	if delay {
		time.Sleep(duration)
	}
}
