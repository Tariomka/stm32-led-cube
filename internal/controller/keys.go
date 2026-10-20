package controller

import (
	"machine"
	"runtime/interrupt"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

const KeyCount = 7

// KeyController reads KEY1..KEY7 on the base board: KEY1 PC0 Previous, KEY2 PC1 Next, KEY3 PC2
// Speed+, KEY4 PB3 Speed-, KEY5 PA14 Run/Pause, KEY6 PA13 Play model, KEY7 PA11 On/Off.
// Presses are caught by interrupts and collected until Take.
type KeyController struct {
	buttons [KeyCount]*component.Button
	pending uint8 // bit n = KEY(n+1) pressed since the last Take
}

// NewKeyController needs the debug port off (NewSystemController) for KEY4-KEY6.
func NewKeyController() *KeyController {
	keys := &KeyController{}
	pins := [KeyCount]machine.Pin{machine.PC0, machine.PC1, machine.PC2, machine.PB3, machine.PA14, machine.PA13, machine.PA11}
	for index, pin := range pins {
		keys.buttons[index] = component.NewButton(pin)
		bit := uint8(1) << index
		if err := keys.buttons[index].OnPress(func() { keys.pending |= bit }); err != nil {
			println("KEY", index+1, "OnPress error:", err)
		}
	}
	return keys
}

// Take returns the keys pressed since the last call (bit n = KEY(n+1)).
func (this *KeyController) Take() uint8 {
	mask := interrupt.Disable()
	keys := this.pending
	this.pending = 0
	interrupt.Restore(mask)
	return keys
}

// Press queues a press of KEY(key+1) as if the key was pressed (console Key command).
func (this *KeyController) Press(key int) {
	mask := interrupt.Disable()
	this.pending |= 1 << key
	interrupt.Restore(mask)
}

// Held returns the keys held down right now.
func (this *KeyController) Held() uint8 {
	var held uint8
	for index, button := range this.buttons {
		if button.Pressed() {
			held |= 1 << index
		}
	}
	return held
}
