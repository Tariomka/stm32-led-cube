package controller

import "github.com/Tariomka/led-common-lib/pkg/led"

type Mode uint8

// TODO: double check states if everything works correctly
const (
	StandbyMode Mode = iota
	OnboardMode
	SerialMode
	DebugMode
	PauseMode
)

type StateTracker struct {
	lightShows     []led.LightShow
	lightShowIndex uint32

	CurrentMode  Mode
	previousMode Mode

	frameRepetitionCount uint32
}

func NewStateTracker(ls []led.LightShow) *StateTracker {
	return &StateTracker{
		lightShows:           ls,
		CurrentMode:          OnboardMode,
		previousMode:         StandbyMode,
		frameRepetitionCount: 1,
	}
}

func (this *StateTracker) CurrentLightShow() led.LightShow {
	if len(this.lightShows) < 1 || this.lightShowIndex >= uint32(len(this.lightShows)) {
		return nil
	}

	return this.lightShows[this.lightShowIndex]
}

func (this *StateTracker) ExecuteFrame(frameCallback func()) {
	// TODO: implement colorDepth as well.
	// Most likely this function will need to be refactored to func(Board, LayoutWorker) signature
	// For the time being colorDepth is implicit - frames switch so fast, it works as colorDepth

	// Put this to frame loop to preserve frame loop possition and not need to wait for frame to finish?
	for this.CurrentMode == PauseMode {
		frameCallback()
	}

	for range this.frameRepetitionCount {
		frameCallback()
	}
}

func (this *StateTracker) CycleMode() {
	switch this.CurrentMode {
	case OnboardMode:
		this.previousMode = this.CurrentMode
		this.CurrentMode = SerialMode
	case SerialMode:
		this.previousMode = this.CurrentMode
		this.CurrentMode = DebugMode
	case DebugMode:
		this.previousMode = this.CurrentMode
		this.CurrentMode = OnboardMode
	case PauseMode:
		this.CurrentMode = this.previousMode
	}
}

func (this *StateTracker) NextLightShow() {
	this.lightShowIndex++
	if this.lightShowIndex >= uint32(len(this.lightShows)) || this.lightShowIndex == ^uint32(0) {
		this.lightShowIndex = 0
	}
}

func (this *StateTracker) PrevLightShow() {
	if this.lightShowIndex == 0 {
		this.lightShowIndex = uint32(len(this.lightShows))
	}
	this.lightShowIndex--
}

func (this *StateTracker) IncreaseSpeed() {
	if this.frameRepetitionCount > 2 {
		this.frameRepetitionCount -= 2
		return
	}

	this.frameRepetitionCount = 1
}

func (this *StateTracker) DecreadeSpeed() {
	if this.frameRepetitionCount < ^uint32(2) {
		this.frameRepetitionCount += 2
		return
	}

	this.frameRepetitionCount = ^uint32(0)
}

func (this *StateTracker) SwitchRunPause() {
	if this.CurrentMode == PauseMode {
		this.CurrentMode = this.previousMode
		return
	}

	this.previousMode = this.CurrentMode
	this.CurrentMode = PauseMode
}

// TODO: implement
func (this *StateTracker) TurnOnOff() {

}
