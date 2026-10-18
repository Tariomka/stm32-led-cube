package runner

import (
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/player"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
	"github.com/Tariomka/stm32-led-cube/internal/settings"
)

type Runner interface {
	Start()
}

const (
	tick      = 10 * time.Millisecond
	saveDelay = 2 * time.Second // settings are written once the user stops pressing keys
)

type CubeRunner struct {
	Board   controller.Board
	Player  *player.Player
	Tracker *controller.StateTracker

	shows  *blob.Blob
	canvas render.Canvas
	group  int
	saveAt time.Time // when to persist the settings, zero if nothing changed

	input
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
		shows:   shows,
		canvas:  render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}},
		group:   board.ShowGroup(),
	}
	runner.Player = player.New(shows, &runner.canvas, runner.group)
	runner.Player.Log = func(msg string) { println(msg) }
	runner.Player.AudioFrame = board.AudioLevel
	runner.Player.ADC = func(channel int, raw bool) int {
		// Unverified: no original show calls GetAdc/GetRealAdc
		if raw {
			return board.AudioRaw()
		}
		return board.AudioLevel(0)
	}
	if saved, ok := board.LoadSettings(); ok {
		runner.Player.SetSpeed(int(saved.Speed))
		runner.Player.PlayAll = saved.PlayAll
		if int(saved.Group) == runner.group {
			runner.Player.Select(int(saved.Cartoon))
		}
	}
	return runner
}

func (this *CubeRunner) saveSettings() {
	this.saveAt = time.Time{}
	err := this.Board.SaveSettings(settings.Settings{
		Speed:   uint8(this.Player.Speed()),
		PlayAll: this.Player.PlayAll,
		Group:   uint8(this.group),
		Cartoon: uint8(this.Player.Current()),
	})
	if err != nil {
		println("Save settings error!")
	}
}

// action runs a key function of the original shows and schedules saving the settings.
func (this *CubeRunner) action(run func()) {
	if this.Tracker.CurrentMode != controller.OnboardMode {
		return
	}
	run()
	this.saveAt = time.Now().Add(saveDelay)
}

func (this *CubeRunner) Start() {
	global.HandleStartup()
	this.discoverCascade()
	this.placeDisplay()

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

// placeDisplay positions this cube in the show canvas by its cascade number, like the original.
func (this *CubeRunner) placeDisplay() {
	count, number := this.shows.DisplayCount(), this.link.Number
	var display blob.Display
	if number < count {
		display = this.shows.Display(number)
	} else {
		println("Warn: Pos > Cubes, " + strconv.Itoa(number) + "," + strconv.Itoa(count))
	}
	if display.Dir != 0 {
		println("Warn: display dir " + strconv.Itoa(display.Dir) + " is not supported")
	}
	println("Total=" + strconv.Itoa(count))
	println("Run x=" + strconv.Itoa(display.X) + ",y=" + strconv.Itoa(display.Y) + ",z=" + strconv.Itoa(display.Z) +
		",dir=" + strconv.Itoa(display.Dir))
	this.canvas.Displays[0] = render.NewDisplay(display.X, display.Y, display.Z, 8, 8, 8)
}

// runOnboard plays one step of the original shows and waits for its frame time on the show
// clock (shared by all cascaded cubes), handling input every tick.
func (this *CubeRunner) runOnboard() {
	start := this.clock()
	rendered, waitMs := this.Player.Step(uint32(start))
	if rendered {
		this.Board.Show(this.canvas.Displays[0].Voxels)
	}
	this.Board.SetBlank(this.Player.DisplayOff)

	deadline := start + uint64(waitMs)
	for {
		if this.handleInput() {
			return
		}
		now := this.clock()
		if now >= deadline {
			return
		}
		wait := tick
		if this.link.Master {
			wait = min(time.Duration(deadline-now)*time.Millisecond, tick)
		}
		time.Sleep(wait)
	}
}

func (this *CubeRunner) runDebugLoop() {
	// Placeholder
	this.Board.GetIndicator().Blink(component.BlinkGreen, 100*time.Millisecond, false)
	this.handleInput()
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
		if this.handleInput() {
			return
		}
		time.Sleep(tick)
	}
}
