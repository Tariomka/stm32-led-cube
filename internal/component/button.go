package component

import (
	"device/arm"
	"device/stm32"
	"machine"
	"time"
)

const debounceTime = 30 * time.Millisecond

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
		if now.Sub(this.lastPress) < debounceTime {
			return
		}
		this.lastPress = now
		callback()
	})
	lowerPinInterruptPriority(this.Pin)
	return err
}

// lowerPinInterruptPriority lets the runtime tick timer (priority 0xC1) preempt the pin's EXTI
// handler. EXTI defaults to the highest priority, and time.Now() inside it spins forever when
// the tick timer overflows meanwhile.
func lowerPinInterruptPriority(pin machine.Pin) {
	irq := uint32(stm32.IRQ_EXTI15_10)
	switch line := uint32(pin % 16); {
	case line <= 4:
		irq = stm32.IRQ_EXTI0 + line
	case line <= 9:
		irq = stm32.IRQ_EXTI9_5
	}
	arm.SetPriority(irq, 0xE0)
}
