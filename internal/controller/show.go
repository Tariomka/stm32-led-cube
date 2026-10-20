package controller

import (
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/common"
	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/blob"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/render"
)

const saveDelay = 2 * time.Second // settings are written once the user stops changing them

// ShowController plays the light show of the current mode on the display: it runs the player,
// feeds Audio cartoons the microphone level, places this cube in multi-cube shows, runs key
// functions and keeps the player settings across power cycles.
type ShowController struct {
	display    *DisplayController
	settings   *SettingsController
	mode       *ModeController
	cascade    *CascadeController
	microphone *component.Microphone // PA0

	player *lightshow.Player
	canvas render.Canvas
	saveAt time.Time // when to persist the settings, zero if nothing changed
	// pauseAfterFrame restores a pause once the first frame of a show is on the cube, so it isn't
	// blank or still showing the previous show
	pauseAfterFrame bool
	changed         bool // the show was replaced, the current frame wait should end
}

// NewShowController starts the show of the current mode with the saved settings.
func NewShowController(
	display *DisplayController,
	settings *SettingsController,
	mode *ModeController,
	cascade *CascadeController) *ShowController {
	this := &ShowController{
		display:    display,
		settings:   settings,
		mode:       mode,
		cascade:    cascade,
		microphone: component.NewMicrophone(),
		canvas:     render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}},
	}
	this.player = lightshow.New(mode.Show(), &this.canvas)
	this.player.Log = func(msg string) { println(msg) }
	this.player.AudioFrame = this.microphone.Level
	this.player.ADC = func(channel int, raw bool) int {
		// Unverified: no original show calls GetAdc/GetRealAdc
		if raw {
			return this.microphone.Raw()
		}

		return this.microphone.Level(0)
	}

	if saved, ok := settings.Load(); ok {
		this.player.SetSpeed(int(saved.Speed))
		this.player.PlayAll = saved.PlayAll
		this.player.DisplayOff = saved.DisplayOff
		this.pauseAfterFrame = saved.Paused
		if saved.Mode == mode.Mode {
			this.player.Select(int(saved.Cartoon))
		}
	}
	return this
}

// Step plays one step of the show at show time now (ms) and returns how long to wait before the
// next step.
func (this *ShowController) Step(now uint32) int {
	rendered, waitMs := this.player.Step(now)
	if rendered {
		this.display.Show(this.canvas.Displays[0].Voxels)
		if this.pauseAfterFrame {
			this.pauseAfterFrame = false
			this.player.Paused = true
		}
	}
	this.display.SetBlank(this.player.DisplayOff)
	return waitMs
}

// Poll follows DIP switch changes and saves the settings when due; call it regularly.
func (this *ShowController) Poll() {
	if this.mode.Update() {
		println("Mode=" + this.mode.Mode.String())
		this.setShow(this.mode.Show())
		this.saveAt = time.Now().Add(saveDelay)
	}
	if !this.saveAt.IsZero() && time.Now().After(this.saveAt) {
		this.save()
	}
}

// Changed reports (once) that the show was replaced, so the current frame wait should end.
func (this *ShowController) Changed() bool {
	changed := this.changed
	this.changed = false
	return changed
}

// Place positions this cube in the show canvas by its cascade number, like the factory firmware.
// The built-in shows are for one cube at the origin; uploaded show data has its own table.
func (this *ShowController) Place() {
	count, number := 1, this.cascade.Number()
	var display blob.Display
	if uploaded := this.mode.Uploaded(); uploaded != nil {
		count = uploaded.DisplayCount()
		if number < count {
			display = uploaded.Display(number)
		}
	}
	if number >= count {
		println("Warn: Pos > Cubes, " + strconv.Itoa(number) + "," + strconv.Itoa(count))
	}
	if display.Dir > 1 {
		println("Warn: display dir " + strconv.Itoa(display.Dir) + " is not supported")
	}
	this.display.SetUpsideDown(display.Dir == 1 && !this.settings.IgnoresDisplayDir())
	println("Total=" + strconv.Itoa(count))
	println("Run x=" + strconv.Itoa(display.X) + ",y=" + strconv.Itoa(display.Y) + ",z=" + strconv.Itoa(display.Z) +
		",dir=" + strconv.Itoa(display.Dir) + ",type=" + strconv.Itoa(this.settings.CubeTypeIndex()))
	this.canvas.Displays[0] = render.NewDisplay(display.X, display.Y, display.Z, 8, 8, 8)
}

// setShow switches the player to another show; a paused player still draws its first frame.
func (this *ShowController) setShow(next lightshow.Show) {
	if this.player.Paused {
		this.player.Paused, this.pauseAfterFrame = false, true
	}
	this.player.SetShow(next)
	this.Place()
	this.changed = true
}

// UploadStarted stops playing the show slot before an upload rewrites it.
func (this *ShowController) UploadStarted() {
	this.mode.SetUploading(true)
	if this.mode.Mode == Text {
		this.player.SetShow(this.mode.Show())
		this.changed = true
	}
}

// UploadFinished plays the new show data in Text mode.
func (this *ShowController) UploadFinished() {
	this.mode.SetUploading(false)
	if this.mode.Mode == Text {
		this.setShow(this.mode.Show())
	}
}

// RunKeys runs the key functions of the keys set in mask (bit n = KEY(n+1)).
func (this *ShowController) RunKeys(mask uint8) {
	for key := 0; key < KeyCount; key++ {
		if mask&(1<<key) != 0 {
			this.changeSettings(func() {
				println("DoKey" + strconv.Itoa(key+1))
				this.player.Key(key)
			})
		}
	}
}

// Action runs a key function code (IR remote, console Action command).
func (this *ShowController) Action(action int) {
	this.changeSettings(func() { this.player.Action(action) })
}

// Select restarts the show on a cartoon (console Cartoon command); false if out of range.
func (this *ShowController) Select(cartoon int) bool {
	if cartoon >= this.player.Count() {
		return false
	}

	this.changeSettings(func() {
		println("Change Cartoon " + strconv.Itoa(cartoon))
		this.player.Select(cartoon)
	})
	return true
}

// Status describes the player state for the console Status command.
func (this *ShowController) Status() string {
	return "Mode=" + this.mode.Mode.String() +
		",Cartoon=" + strconv.Itoa(this.player.Current()) +
		",Count=" + strconv.Itoa(this.player.Count()) +
		",Speed=" + this.player.SpeedText() +
		",PlayModel=" + common.Flag(this.player.PlayAll) +
		",Pause=" + common.Flag(this.player.Paused || this.pauseAfterFrame) +
		",DisplayOff=" + common.Flag(this.player.DisplayOff) +
		",DeviceNum=" + strconv.Itoa(this.cascade.Number())
}

// changeSettings runs a key function and schedules saving the settings.
func (this *ShowController) changeSettings(callback func()) {
	this.pauseAfterFrame = false
	callback()
	this.saveAt = time.Now().Add(saveDelay)
}

func (this *ShowController) save() {
	this.saveAt = time.Time{}
	err := this.settings.Save(Settings{
		Speed:      uint8(this.player.Speed()),
		PlayAll:    this.player.PlayAll,
		Paused:     this.player.Paused || this.pauseAfterFrame,
		DisplayOff: this.player.DisplayOff,
		Mode:       this.mode.Mode,
		Cartoon:    uint8(this.player.Current()),
	})
	if err != nil {
		println("Save settings error!")
	}
}
