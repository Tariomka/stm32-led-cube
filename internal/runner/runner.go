package runner

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/player"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

type Runner interface {
	Start()
}

const (
	cycleKey      = 5 // KEY6: short press = show key function, long press = cycle runner mode
	longPressTime = time.Second
	tick          = 10 * time.Millisecond
)

type CubeRunner struct {
	Board   controller.Board
	Player  *player.Player
	Tracker *controller.StateTracker

	canvas     render.Canvas
	start      time.Time
	cycleSince time.Time // when the pending KEY6 press started, zero if none
}

func NewRunner(config RunnerConfig) Runner {
	// TODO: implement other cube sizes/types
	if config.BaseSize != Size8 || config.Height != Size8 || config.LedType != RGB {
		return nil
	}

	shows, err := blob.Parse(config.Shows)
	if err != nil {
		println("show data:", err.Error())
		global.HandleFatal()
	}

	board := controller.NewYellowBoard()
	runner := &CubeRunner{
		Board:   board,
		Tracker: controller.NewStateTracker(),
		canvas:  render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}},
		start:   time.Now(),
	}
	runner.Player = player.New(shows, &runner.canvas, board.ShowGroup())
	runner.Player.Log = func(msg string) { println(msg) }
	return runner
}

func (this *CubeRunner) millis() uint32 {
	return uint32(time.Since(this.start).Milliseconds())
}

func (this *CubeRunner) Start() {
	global.HandleStartup()

	for {
		switch this.Tracker.CurrentMode {
		case controller.OnboardMode:
			this.runOnboard()
		case controller.DebugMode:
			this.runDebugLoop()
		case controller.SerialMode:
			this.runSerialLoop()
		}
	}
}

// runOnboard plays one step of the original shows and waits for its frame time,
// handling keys every tick.
func (this *CubeRunner) runOnboard() {
	rendered, waitMs := this.Player.Step(this.millis())
	if rendered {
		this.Board.Show(this.canvas.Displays[0].Voxels)
	}
	this.Board.SetBlank(this.Player.DisplayOff)

	deadline := time.Now().Add(time.Duration(waitMs) * time.Millisecond)
	for {
		if this.handleKeys() {
			return
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return
		}
		time.Sleep(min(remaining, tick))
	}
}

// handleKeys applies pressed keys; it returns true when the runner mode changed.
func (this *CubeRunner) handleKeys() bool {
	keys := this.Board.TakeKeys()
	for key := 0; key < controller.KeyCount; key++ {
		if keys&(1<<key) == 0 {
			continue
		}
		if key == cycleKey {
			this.cycleSince = time.Now()
			continue
		}
		if this.Tracker.CurrentMode == controller.OnboardMode {
			this.Player.Key(key)
		}
	}

	if this.cycleSince.IsZero() {
		return false
	}
	if this.Board.KeyHeld(cycleKey) {
		if time.Since(this.cycleSince) < longPressTime {
			return false
		}
		this.cycleSince = time.Time{}
		this.Tracker.CycleMode()
		this.Board.SetBlank(false)
		return true
	}
	this.cycleSince = time.Time{}
	if this.Tracker.CurrentMode == controller.OnboardMode {
		this.Player.Key(cycleKey)
	}
	return false
}

func (this *CubeRunner) runDebugLoop() {
	// Placeholder
	this.Board.GetIndicator().Blink(component.BlinkGreen, 100*time.Millisecond, false)
	this.handleKeys()
}

func (this *CubeRunner) runSerialLoop() {
	// Placeholder
	this.Board.SetBlank(true)
	data := this.Board.Receive()
	if len(data) > 0 {
		this.Board.Send("Received:" + string(data))
	}
	this.Board.GetIndicator().Blink(component.BlinkGreen, 100*time.Millisecond, false)
	for range 200 {
		if this.handleKeys() {
			return
		}
		time.Sleep(tick)
	}
}
