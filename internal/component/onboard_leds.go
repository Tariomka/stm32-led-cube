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
		this.LedRed.Low()
	case BlinkGreen:
		this.LedGreen.Low()
	case BlinkBoth:
		this.LedGreen.Low()
		this.LedRed.Low()
	case BlinkStaggerRed:
		this.LedRed.Low()
		time.Sleep(duration)
		this.LedRed.High()
		this.LedGreen.Low()
	case BlinkStaggerGreen:
		this.LedGreen.Low()
		time.Sleep(duration)
		this.LedGreen.High()
		this.LedRed.Low()
	}
	time.Sleep(duration)

	this.LedGreen.High()
	this.LedRed.High()
	if delay {
		time.Sleep(duration)
	}
}
