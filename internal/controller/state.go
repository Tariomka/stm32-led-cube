package controller

type Mode uint8

const (
	OnboardMode Mode = iota + 1
	SerialMode
	DebugMode
	StandbyMode
)

type StateTracker struct {
	AnimationSpeed uint32 // Speed denoted how many times a frame should be shown
	LightShowIndex uint32
	Pause          bool
	colorDepth     uint8
	CurrentMode    Mode
	previousMode   Mode
	DisableOnboard bool
}

func NewStateTracker() *StateTracker {
	return &StateTracker{
		AnimationSpeed: 1,
		LightShowIndex: 0,
		Pause:          false,
		colorDepth:     1,
		CurrentMode:    OnboardMode,
		previousMode:   StandbyMode,
		DisableOnboard: false,
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
	st.LightShowIndex++
}

func (st *StateTracker) PrevLightShow() {
	st.LightShowIndex--
}

func (st *StateTracker) IncreaseSpeed() {
	st.AnimationSpeed -= 2
}

func (st *StateTracker) DecreadeSpeed() {
	st.AnimationSpeed += 2
}

func (st *StateTracker) SwitchRunPause() {
	st.Pause = !st.Pause
}

func (st *StateTracker) EnableOnboardMode(enable bool) {
	st.DisableOnboard = !enable
}
