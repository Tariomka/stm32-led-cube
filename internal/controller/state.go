package controller

import (
	"machine"

	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/player"
	"github.com/Tariomka/stm32-led-cube/internal/show"
)

// Mode is the show mode selected with DIP switches SW1/SW2 (ON = 0).
type Mode uint8

const (
	Dynamic Mode = iota // SW1 = 1, SW2 = 1: the original animations
	Audio               // SW1 = 0, SW2 = 1: cartoons follow the microphone
	Test                // SW1 = 1, SW2 = 0: LED test patterns
	Text                // SW1 = 0, SW2 = 0: the show uploaded from the PC or the Pico
)

func (this Mode) String() string {
	switch this {
	case Dynamic:
		return "Dynamic"
	case Audio:
		return "Audio"
	case Test:
		return "Test"
	case Text:
		return "Text"
	default:
		return "Unknown"
	}
}

// switchReads is how many equal readings in a row change the mode (debounce).
const switchReads = 3

// StateTracker follows the DIP switches and provides the show of the selected mode.
type StateTracker struct {
	// DIP switches, ON = low. SW3/SW4 aren't reachable on the base and only reported.
	Switch1 component.InputPin // PA1
	Switch2 component.InputPin // PC3
	Switch3 component.InputPin // PA15
	Switch4 component.InputPin // PA12

	Mode Mode
	// Uploaded is the show data in the upload slot, nil if none or while an upload runs.
	Uploaded  *blob.Blob
	uploading bool

	next  Mode
	reads int
}

// NewStateTracker reads the switches; create it after the board (PA15 needs the debug port off).
func NewStateTracker() *StateTracker {
	tracker := &StateTracker{
		Switch1: component.NewInputPin(machine.PA1),
		Switch2: component.NewInputPin(machine.PC3),
		Switch3: component.NewInputPin(machine.PA15),
		Switch4: component.NewInputPin(machine.PA12),
	}
	tracker.Mode = tracker.read()
	tracker.next = tracker.Mode
	return tracker
}

func (this *StateTracker) read() Mode {
	switch1, switch2 := this.Switch1.Get(), this.Switch2.Get()
	switch {
	case switch1 && switch2:
		return Dynamic
	case !switch1 && switch2:
		return Audio
	case switch1 && !switch2:
		return Test
	default:
		return Text
	}
}

// Update reads the switches and reports whether the mode changed.
func (this *StateTracker) Update() bool {
	mode := this.read()
	if mode == this.Mode {
		this.reads = 0
		return false
	}

	if mode != this.next {
		this.next = mode
		this.reads = 0
	}
	if this.reads++; this.reads < 3 {
		return false
	}

	this.Mode = mode
	this.reads = 0
	return true
}

// Config2 reads SW3/SW4 like the original (bit set = switch ON); it is only reported.
func (this *StateTracker) Config2() int {
	value := 0
	if !this.Switch3.Get() {
		value |= 1
	}
	if !this.Switch4.Get() {
		value |= 2
	}
	return value
}

// Show returns the show of the current mode: built in for Dynamic, Audio and Test, the uploaded
// show data (all its groups) for Text.
func (this *StateTracker) Show() player.Show {
	this.Uploaded = nil
	switch this.Mode {
	case Dynamic:
		return show.New(show.Dynamic)
	case Audio:
		return show.New(show.Audio)
	case Test:
		return show.New(show.Test)
	}
	if !this.uploading {
		if b, err := blob.Parse(component.ShowSlotData()); err == nil {
			this.Uploaded = b
			return player.NewBlobShow(b, 0, 1, 2, 3)
		}
	}
	return show.New(nil)
}

// SetUploading marks the upload slot as being rewritten; Show then doesn't read it.
func (this *StateTracker) SetUploading(uploading bool) {
	this.uploading = uploading
	if uploading {
		this.Uploaded = nil
	}
}
