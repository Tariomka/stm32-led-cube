package runner

import "github.com/Tariomka/stm32-led-cube/internal/controller"

type Runner interface {
	Start()
}

type RunnerLoop interface {
}

type CubeRunner struct {
	Board        controller.Board
	LayoutWorker controller.LayoutWorker
	Tracker      *controller.StateTracker
}

func NewRunner(config RunnerConfic) Runner {
	// TODO: implement other cube sizes/types
	if config.BaseSize != Size8 || config.Height != Size8 || config.LedType != RGB {
		return nil
	}

	tracker := controller.NewStateTracker(config.LightShows)

	return &CubeRunner{
		LayoutWorker: &controller.LedLayout{},
		Board:        controller.NewYellowBoard(tracker),
		Tracker:      tracker,
	}
}

func (cr *CubeRunner) Start() {
	cr.Board.BlinkStartup()

	for {
		switch cr.Tracker.CurrentMode {
		case controller.OnboardMode:
			cr.runOnboardLoop()
		case controller.DebugMode:
			cr.runDebugLoop()
		case controller.SerialMode:
			cr.runSerialLoop()
		}
	}
}

func (cr *CubeRunner) runOnboardLoop() {
	for _, frameCallback := range cr.Tracker.CurrentLightShow() {
		cr.LayoutWorker.ResetBlock()
		frameCallback(cr.LayoutWorker)
		cr.Tracker.ExecuteFrame(func() {
			cr.Board.LightLeds(cr.LayoutWorker)
		})
	}
}

func (cr *CubeRunner) runDebugLoop() {
	// Placeholder
	cr.Board.BlinkDebug()
	cr.Board.EnableLeds()
}

func (cr *CubeRunner) runSerialLoop() {
	// Placeholder
	cr.Board.DisableLeds()
}
