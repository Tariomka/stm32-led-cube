package runner

import "github.com/Tariomka/stm32-led-cube/internal/controller"

type Runner interface {
	Start(ls []controller.LightShow)
}

type CubeRunner struct {
	Board        controller.Board
	LayoutWorker controller.LayoutWorker
	Tracker      *controller.StateTracker
	// LightShows   []controller.LightShow
}

func NewRunner(config RunnerConfic) Runner {
	if config.BaseSize != Size8 || config.Height != Size8 || config.LedType != RGB {
		return nil
	}

	tracker := controller.NewStateTracker()

	return &CubeRunner{
		LayoutWorker: &controller.LedLayout{},
		Board:        controller.NewYellowBoard(tracker),
		Tracker:      tracker,
	}
}

func (cr *CubeRunner) Start(ls []controller.LightShow) {
	cr.Board.BlinkStartup()

	if len(ls) < 1 {
		cr.Tracker.EnableOnboardMode(false)
	}

	for {
		switch cr.Tracker.CurrentMode {
		case controller.OnboardMode:
			cr.runOnboardLoop(ls)
		case controller.DebugMode:
			cr.runDebugLoop()
		case controller.SerialMode:
			cr.runSerialLoop()
		}

	}
}

func (cr *CubeRunner) runOnboardLoop(ls []controller.LightShow) {
	for {
		if cr.Tracker.DisableOnboard {
			return
		}

		if cr.Tracker.Pause {
			// Continue might work as well, but there still needs to be a return for when Mode is changed
			// This needs to be in the for loop below, so pausing can occur mid way through a show
			// ?and somehow a state needs to be tracked to continue from the last paused frame?
			return
		}

		for _, lightShow := range ls[cr.Tracker.LightShowIndex] {
			cr.LayoutWorker.ResetBlock()
			lightShow(cr.LayoutWorker)
			for i := uint32(0); i < cr.Tracker.AnimationSpeed; i++ {
				cr.Board.LightLeds(cr.LayoutWorker)
			}
		}
	}
}

func (cr *CubeRunner) runDebugLoop() {
	return
}

func (cr *CubeRunner) runSerialLoop() {
	return
}
