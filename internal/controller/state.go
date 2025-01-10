package controller

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
	lightShows           []LightShow
	pause                bool
	previousMode         Mode
}

func NewStateTracker(ls []LightShow) *StateTracker {
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

func (st *StateTracker) CurrentLightShow() LightShow {
	if len(st.lightShows) < 1 {
		return []Frame{}
	}

	return st.lightShows[st.lightShowIndex]
}

func (st *StateTracker) ExecuteFrame(frameCallback func()) {
	// TODO: implement colorDepth as well.
	// Most likely this function will need to be refactored to func(Board, LayoutWorker) signature
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
	st.frameRepetitionCount -= 2
}

func (st *StateTracker) DecreadeSpeed() {
	st.frameRepetitionCount += 2
}

func (st *StateTracker) SwitchRunPause() {
	st.pause = !st.pause
}
