package runner

import (
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/common/global"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/player"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
	"github.com/Tariomka/stm32-led-cube/internal/settings"
	"github.com/Tariomka/stm32-led-cube/internal/show"
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

	canvas render.Canvas
	saveAt time.Time // when to persist the settings, zero if nothing changed
	// pauseAfterFrame restores a pause once the first frame of a show is on the cube, so it isn't
	// blank or still showing the previous show
	pauseAfterFrame bool
	showChanged     bool // the show was replaced, stop waiting for the current frame

	input
}

func NewRunner(config RunnerConfig) Runner {
	// TODO: implement other cube sizes/types
	if config.BaseSize != Size8 || config.Height != Size8 || config.LedType != RGB {
		return nil
	}

	board := controller.NewYellowBoard()
	runner := &CubeRunner{
		Board:   board,
		Tracker: controller.NewStateTracker(),
		canvas:  render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}},
	}
	show.SetMissing(func(kind, name string) { println("Not found " + kind + ":'" + name + "'") })
	runner.Player = player.New(runner.modeShow(), &runner.canvas)
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
		runner.Player.DisplayOff = saved.DisplayOff
		runner.pauseAfterFrame = saved.Paused
		if controller.Mode(saved.Mode) == runner.Tracker.Mode {
			runner.Player.Select(int(saved.Cartoon))
		}
	}
	board.OnUpload(runner.uploadStarted, runner.uploadFinished)
	board.SetRemote(runner)
	return runner
}

// modeShow is the show of the DIP mode: built in, or the uploaded show data in Text mode.
func (this *CubeRunner) modeShow() player.Show {
	modeShow := this.Tracker.Show()
	if this.Tracker.Mode == controller.Text {
		if this.Tracker.Uploaded == nil {
			println("No uploaded show")
		} else {
			println("Uploaded show: " + strconv.Itoa(modeShow.Count()) + " cartoons")
		}
	}
	return modeShow
}

// setShow switches the player to another show; a paused player still draws its first frame.
func (this *CubeRunner) setShow(next player.Show) {
	if this.Player.Paused {
		this.Player.Paused, this.pauseAfterFrame = false, true
	}
	this.Player.SetShow(next)
	this.placeDisplay()
	this.showChanged = true
}

// changeMode follows a DIP switch change without a restart.
func (this *CubeRunner) changeMode() {
	println("Mode=" + this.Tracker.Mode.String())
	this.setShow(this.modeShow())
	this.saveAt = time.Now().Add(saveDelay)
}

// uploadStarted stops playing the show slot before the upload rewrites it.
func (this *CubeRunner) uploadStarted() {
	this.Tracker.SetUploading(true)
	if this.Tracker.Mode == controller.Text {
		this.Player.SetShow(this.Tracker.Show())
		this.showChanged = true
	}
}

func (this *CubeRunner) uploadFinished(ok bool) {
	this.Tracker.SetUploading(false)
	if ok {
		println("Show upload OK")
	} else {
		println("Show upload failed")
	}
	if this.Tracker.Mode == controller.Text {
		this.setShow(this.modeShow())
	}
}

func (this *CubeRunner) saveSettings() {
	this.saveAt = time.Time{}
	err := this.Board.SaveSettings(settings.Settings{
		Speed:      uint8(this.Player.Speed()),
		PlayAll:    this.Player.PlayAll,
		Paused:     this.Player.Paused || this.pauseAfterFrame,
		DisplayOff: this.Player.DisplayOff,
		Mode:       uint8(this.Tracker.Mode),
		Cartoon:    uint8(this.Player.Current()),
	})
	if err != nil {
		println("Save settings error!")
	}
}

// action runs a key function of the shows and schedules saving the settings.
func (this *CubeRunner) action(run func()) {
	this.pauseAfterFrame = false
	run()
	this.saveAt = time.Now().Add(saveDelay)
}

func (this *CubeRunner) Start() {
	global.HandleStartup()
	println("Config1=" + strconv.Itoa(int(this.Tracker.Mode)) + ",Config2=" + strconv.Itoa(this.Tracker.Config2()))
	this.discoverCascade()
	this.placeDisplay()

	for {
		this.runOnboard()
	}
}

// placeDisplay positions this cube in the show canvas by its cascade number, like the original.
// The built-in shows are for one cube at the origin; uploaded show data has its own table.
func (this *CubeRunner) placeDisplay() {
	if this.link == nil {
		return // Start places it after the cascade discovery
	}
	count, number := 1, this.link.Number
	var display blob.Display
	if uploaded := this.Tracker.Uploaded; uploaded != nil {
		count = uploaded.DisplayCount()
		if number < count {
			display = uploaded.Display(number)
		}
	}
	if number >= count {
		println("Warn: Pos > Cubes, " + strconv.Itoa(number) + "," + strconv.Itoa(count))
	}
	if display.Dir != 0 {
		println("Warn: display dir " + strconv.Itoa(display.Dir) + " is not supported")
	}
	println("Total=" + strconv.Itoa(count))
	println("Run x=" + strconv.Itoa(display.X) + ",y=" + strconv.Itoa(display.Y) + ",z=" + strconv.Itoa(display.Z) +
		",dir=" + strconv.Itoa(display.Dir) + ",type=" + strconv.Itoa(cubeTypeIndex(this.Board.CubeType())))
	this.canvas.Displays[0] = render.NewDisplay(display.X, display.Y, display.Z, 8, 8, 8)
}

// cubeTypeIndex is the original's scan variant for the EEPROM type (compared case-insensitively),
// only printed: our scanner always drives one 8x8x8 cube.
func cubeTypeIndex(cubeType string) int {
	types := [...]string{"8x8x8A", "8x8x8B", "8x16x8A", "161616A", "161616B", "161616C", "8x32x8A", "8x32x8B"}
	for index, name := range types {
		if len(name) != len(cubeType) {
			continue
		}
		same := true
		for i := range len(name) {
			same = same && name[i]|0x20 == cubeType[i]|0x20
		}
		if same {
			return index
		}
	}
	return -1
}

// runOnboard plays one step of the shows and waits for its frame time on the show clock
// (shared by all cascaded cubes), handling input every tick.
func (this *CubeRunner) runOnboard() {
	start := this.clock()
	rendered, waitMs := this.Player.Step(uint32(start))
	if rendered {
		this.Board.Show(this.canvas.Displays[0].Voxels)
		if this.pauseAfterFrame {
			this.pauseAfterFrame = false
			this.Player.Paused = true
		}
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

// Key presses KEY1..KEY7 (key 0..6) for a console command, like the physical key.
func (this *CubeRunner) Key(key int) {
	this.pressed |= 1 << key
}

// Action runs a key function code for a console command, like the IR remote.
func (this *CubeRunner) Action(action int) {
	this.action(func() { this.Player.Action(action) })
}

// Select restarts the show on a cartoon for a console command.
func (this *CubeRunner) Select(cartoon int) bool {
	if cartoon >= this.Player.Count() {
		return false
	}
	this.action(func() {
		println("Change Cartoon " + strconv.Itoa(cartoon))
		this.Player.Select(cartoon)
	})
	return true
}

// Status describes the player state for the console Status command.
func (this *CubeRunner) Status() string {
	number := 0
	if this.link != nil {
		number = this.link.Number
	}
	return "Mode=" + this.Tracker.Mode.String() +
		",Cartoon=" + strconv.Itoa(this.Player.Current()) +
		",Count=" + strconv.Itoa(this.Player.Count()) +
		",Speed=" + this.Player.SpeedText() +
		",PlayModel=" + flag(this.Player.PlayAll) +
		",Pause=" + flag(this.Player.Paused || this.pauseAfterFrame) +
		",DisplayOff=" + flag(this.Player.DisplayOff) +
		",DeviceNum=" + strconv.Itoa(number)
}

func flag(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
