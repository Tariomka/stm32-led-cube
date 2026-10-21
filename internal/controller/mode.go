package controller

import (
	"machine"
	"strconv"

	"github.com/Tariomka/led-common-lib/pkg/lightshow"
	"github.com/Tariomka/led-common-lib/pkg/lightshow/blob"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/builtin"
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

// ModeController follows the DIP switches and provides the show of the selected mode: the
// built-in shows, or in Text mode the show data stored by an upload.
type ModeController struct {
	// DIP switches, ON = low. SW3/SW4 aren't reachable on the base and only reported.
	switch1 component.InputPin // PA1
	switch2 component.InputPin // PC3
	switch3 component.InputPin // PA15
	switch4 component.InputPin // PA12

	Mode      Mode
	uploaded  *blob.Blob // show data in the upload slot, nil if none or while an upload runs
	uploading bool

	next  Mode
	reads int
}

// NewModeController reads the switches; it needs the debug port off (NewSystemController) for SW3.
func NewModeController() *ModeController {
	mode := &ModeController{
		switch1: component.NewInputPin(machine.PA1),
		switch2: component.NewInputPin(machine.PC3),
		switch3: component.NewInputPin(machine.PA15),
		switch4: component.NewInputPin(machine.PA12),
	}
	mode.Mode = mode.read()
	mode.next = mode.Mode
	return mode
}

func (this *ModeController) read() Mode {
	switch1, switch2 := this.switch1.Get(), this.switch2.Get()
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
func (this *ModeController) Update() bool {
	mode := this.read()
	if mode == this.Mode {
		this.reads = 0
		return false
	}

	if mode != this.next {
		this.next = mode
		this.reads = 0
	}
	if this.reads++; this.reads < switchReads {
		return false
	}

	this.Mode = mode
	this.reads = 0
	return true
}

// Config is the factory boot line: the mode number and SW3/SW4 (bit set = switch ON).
func (this *ModeController) Config() string {
	config2 := 0
	if !this.switch3.Get() {
		config2 |= 1
	}
	if !this.switch4.Get() {
		config2 |= 2
	}
	return "Config1=" + strconv.Itoa(int(this.Mode)) + ",Config2=" + strconv.Itoa(config2)
}

// Show returns the show of the current mode: built in for Dynamic, Audio and Test, the uploaded
// show data (all its groups) for Text.
func (this *ModeController) Show() lightshow.Show {
	this.uploaded = nil
	switch this.Mode {
	case Dynamic:
		return builtin.New(builtin.Dynamic)
	case Audio:
		return builtin.New(builtin.Audio)
	case Test:
		return builtin.New(builtin.Test)
	}

	if this.uploading {
		return builtin.New(nil)
	}
	if b, err := blob.Parse(component.ShowSlotData()); err == nil {
		this.uploaded = b
		show := lightshow.NewBlobShow(b, 0, 1, 2, 3)
		println("Uploaded show: " + strconv.Itoa(show.Count()) + " cartoons")
		return show
	}
	println("No uploaded show")
	return builtin.New(nil)
}

// Uploaded is the show data played in Text mode, nil if none.
func (this *ModeController) Uploaded() *blob.Blob { return this.uploaded }

// SetUploading marks the upload slot as being rewritten; Show then doesn't read it.
func (this *ModeController) SetUploading(uploading bool) {
	this.uploading = uploading
	if uploading {
		this.uploaded = nil
	}
}
