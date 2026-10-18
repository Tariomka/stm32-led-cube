package controller

type Mode uint8

const (
	OnboardMode Mode = iota // Original iCubeSmart shows
	SerialMode              // Frames from the Raspberry Pi Pico link
	DebugMode
	modeCount
)

type StateTracker struct {
	CurrentMode Mode
}

func NewStateTracker() *StateTracker {
	return &StateTracker{CurrentMode: OnboardMode}
}

func (this *StateTracker) CycleMode() {
	this.CurrentMode = (this.CurrentMode + 1) % modeCount
}
