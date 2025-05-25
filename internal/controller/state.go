package controller

import "github.com/Tariomka/led-common-lib/pkg/led"

type Mode uint8

const (
	OnboardMode Mode = iota + 1
	SerialMode
	DebugMode
	StandbyMode
)

type StateTracker struct {
	CurrentMode          Mode
	colorDepth           uint8
	frameRepetitionCount uint32
	lightShowIndex       uint32
	lightShows           []led.LightShow
	pause                bool
	previousMode         Mode
}

func NewStateTracker(ls []led.LightShow) *StateTracker {
	return &StateTracker{
		CurrentMode:          OnboardMode,
		colorDepth:           1,
		frameRepetitionCount: 1,
		lightShowIndex:       0,
		lightShows:           ls,
		pause:                false,
		previousMode:         StandbyMode,
	}
}

func (st *StateTracker) CurrentLightShow() led.LightShow {
	if len(st.lightShows) < 1 || st.lightShowIndex >= uint32(len(st.lightShows)) {
		return nil
	}

	return st.lightShows[st.lightShowIndex]
}

func (st *StateTracker) ExecuteFrame(frameCallback func()) {
	// TODO: implement colorDepth as well.
	// Most likely this function will need to be refactored to func(Board, LayoutWorker) signature
	// For the time being colorDepth is implicit - frames switch so fast, it works as colorDepth

	// Put this to frame loop to preserve frame loop possition and not need to wait for frame to finish?
	for st.pause {
		frameCallback()
	}

	for i := uint32(0); i < st.frameRepetitionCount; i++ {
		frameCallback()
	}
}

func (st *StateTracker) CycleMode() {
	switch st.CurrentMode {
	case OnboardMode:
		st.previousMode = st.CurrentMode
		st.CurrentMode = SerialMode
	case SerialMode:
		st.previousMode = st.CurrentMode
		st.CurrentMode = DebugMode
	case DebugMode:
		st.previousMode = st.CurrentMode
		st.CurrentMode = OnboardMode
	}
}

func (st *StateTracker) NextLightShow() {
	st.lightShowIndex++
	if st.lightShowIndex >= uint32(len(st.lightShows)) || st.lightShowIndex == ^uint32(0) {
		st.lightShowIndex = 0
	}
}

func (st *StateTracker) PrevLightShow() {
	if st.lightShowIndex == 0 {
		st.lightShowIndex = uint32(len(st.lightShows))
	}
	st.lightShowIndex--
}

func (st *StateTracker) IncreaseSpeed() {
	if st.frameRepetitionCount > 2 {
		st.frameRepetitionCount -= 2
	} else {
		st.frameRepetitionCount = 1
	}
}

func (st *StateTracker) DecreadeSpeed() {
	if st.frameRepetitionCount < ^uint32(2) {
		st.frameRepetitionCount += 2
	} else {
		st.frameRepetitionCount = ^uint32(0)
	}
}

func (st *StateTracker) SwitchRunPause() {
	st.pause = !st.pause
}
