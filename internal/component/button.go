package component

import (
	"machine"
	"time"
)

// Active-low push button with a debounced press interrupt
type Button struct {
	InputPin
	lastPress time.Time
}

func NewButton(pin machine.Pin) *Button {
	return &Button{InputPin: NewInputPin(pin)}
}

func (this *Button) Pressed() bool {
	return !this.Get()
}

// OnPress registers a callback that runs in interrupt context once per debounced press
func (this *Button) OnPress(callback func()) error {
	err := this.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		now := time.Now()
		if now.Sub(this.lastPress) < 30*time.Millisecond {
			return
		}
		this.lastPress = now
		callback()
	})
	lowerPinInterruptPriority(this.Pin)
	return err
}
