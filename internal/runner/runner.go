// Package runner starts the cube: it creates the controllers and runs the main loop.
package runner

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/controller"
)

const tick = 10 * time.Millisecond // input polling period while a frame is shown

type Runner struct {
	system   *controller.SystemController
	settings *controller.SettingsController
	keys     *controller.KeyController
	remote   *controller.RemoteController
	mode     *controller.ModeController
	cascade  *controller.CascadeController
	show     *controller.ShowController
	terminal *controller.TerminalController
}

func New() *Runner {
	system := controller.NewSystemController() // first: frees the debug port pins, starts the watchdog
	settings := controller.NewSettingsController()
	display := controller.NewDisplayController()
	keys := controller.NewKeyController()
	mode := controller.NewModeController()
	cascade := controller.NewCascadeController()
	show := controller.NewShowController(display, settings, mode, cascade)
	return &Runner{
		system:   system,
		settings: settings,
		keys:     keys,
		remote:   controller.NewRemoteController(),
		mode:     mode,
		cascade:  cascade,
		show:     show,
		terminal: controller.NewTerminalController(system, settings, show, keys),
	}
}

func (this *Runner) Start() {
	this.system.Startup(this.settings.CubeType())
	println(this.mode.Config())
	this.cascade.Discover(this.idle)
	this.show.Place()

	for {
		this.step()
	}
}

// idle keeps the watchdog fed and the terminal served.
func (this *Runner) idle() {
	this.system.Feed()
	this.terminal.Poll()
}

// step plays one step of the show and handles input until its frame time is over on the show
// clock (shared by all cascaded cubes) or the show is replaced.
func (this *Runner) step() {
	start := this.cascade.Clock()
	deadline := start + uint64(this.show.Step(uint32(start)))
	for {
		this.idle()
		this.show.Poll()
		this.show.RunKeys(this.cascade.Keys(this.keys.Take(), this.keys.Held()))
		if this.cascade.Master() {
			if action := this.remote.Take(); action != 0 {
				this.show.Action(action)
			}
		}
		if this.show.Changed() {
			return
		}

		now := this.cascade.Clock()
		if now >= deadline {
			return
		}
		wait := tick
		if this.cascade.Master() {
			wait = min(time.Duration(deadline-now)*time.Millisecond, tick)
		}
		time.Sleep(wait)
	}
}
