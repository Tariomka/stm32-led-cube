package runner

import (
	"time"

	"github.com/Tariomka/led-common-lib/pkg/led"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
)

type Runner interface {
	Start()
}

type RunnerLoop interface {
}

type CubeRunner struct {
	Board        controller.Board
	LayoutWorker led.LayoutWorker
	Tracker      *controller.StateTracker
}

func NewRunner(config RunnerConfig) Runner {
	// TODO: implement other cube sizes/types
	if config.BaseSize != Size8 || config.Height != Size8 || config.LedType != RGB {
		return nil
	}

	tracker := controller.NewStateTracker(config.LightShows)
	return &CubeRunner{
		LayoutWorker: &led.LedLayout{},
		Board:        controller.NewYellowBoard(tracker),
		Tracker:      tracker,
	}
}

func (this *CubeRunner) Start() {
	this.Board.GetIndicator().BlinkStartup()

	for {
		switch this.Tracker.CurrentMode {
		case controller.OnboardMode:
			this.runOnboardLoop()
		case controller.DebugMode:
			this.runDebugLoop()
		case controller.SerialMode:
			this.runSerialLoop()
		}
	}
}

func (this *CubeRunner) runOnboardLoop() {
	for _, frameCallback := range this.Tracker.CurrentLightShow() {
		this.LayoutWorker.ResetBlock()
		frameCallback(this.LayoutWorker)
		this.Tracker.ExecuteFrame(func() {
			this.Board.LightLeds(this.LayoutWorker)
		})
	}
}

func (this *CubeRunner) runDebugLoop() {
	// Placeholder
	this.Board.GetIndicator().Blink(component.BlinkGreen, 100*time.Millisecond, false)
	this.Board.EnableLeds()
}

func (this *CubeRunner) runSerialLoop() {
	// Placeholder
	this.Board.DisableLeds()
	data := this.Board.Receive()
	if len(data) > 0 {
		this.Board.Send("Received:" + string(data))
	}
	this.Board.GetIndicator().Blink(component.BlinkGreen, 100*time.Millisecond, false)
	time.Sleep(2 * time.Second)
}
