// Package player runs the original iCubeSmart show engine: it walks the cartoons of a show
// group, renders their frames and executes key functions. It is a non-blocking port of the
// vendor player loop (Simulator.dll FUN_10004ca0); see .agent/re/spec-engine.md.
//
// Player is not safe for concurrent use: call Key/Action from the same loop that calls Step
// (queue key presses from interrupts and apply them before Step).
package player

import (
	"strconv"

	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

// Key function codes stored in the show config (KEY1..KEY7).
const (
	ActionNone        = 0
	ActionReplay      = 1
	ActionPrevious    = 2
	ActionNext        = 3
	ActionSpeedUp     = 4
	ActionSpeedReset  = 5
	ActionSpeedDown   = 6
	ActionPlayModel   = 7
	ActionPause       = 8
	ActionDisplayOff  = 9
	ActionJumpCartoon = 10 // 10..19 jump to cartoon 0..9 of the group
)

// Speed factors in tenths: 0.1, 0.2, 0.5, 0.7, 1.0, 1.5, 2.0, 3.0
var speedTenths = [8]int{1, 2, 5, 7, 10, 15, 20, 30}

const defaultSpeed = 4

// pollMs is how long to wait before calling Step again when nothing was drawn.
const pollMs = 100

type Player struct {
	blob     *blob.Blob
	canvas   *render.Canvas
	assets   *assets
	cartoons []blob.Cartoon

	current int
	running bool // current cartoon has been started
	change  bool // a key selected another (or the same) cartoon
	loop    int
	frame   int
	drawn   bool
	until   uint32 // audio cartoon end time, 0 = forever

	vm *vm

	PlayAll    bool // play model: true = advance through the group, false = repeat one cartoon
	Paused     bool
	DisplayOff bool
	speed      int

	keyCapture  bool
	capturedKey int
	colorStep   int

	// AudioFrame selects the frame of Audio cartoons (0..63) for the cartoon's level mode
	// (header Loop field). The vendor PC simulator always uses 0; on the cube it follows the
	// microphone level (see package audio).
	AudioFrame func(mode int) int
	// ADC feeds the script functions GetAdc/GetRealAdc.
	ADC func(channel int, raw bool) int
	// Log receives the engine's console messages.
	Log func(msg string)
}

// New creates a player for show group 0..3 drawing into canvas.
func New(blob *blob.Blob, canvas *render.Canvas, group int) *Player {
	p := &Player{blob: blob, canvas: canvas, assets: newAssets(blob)}
	p.assets.Missing = func(kind, name string) { p.log("Not found " + kind + ":'" + name + "'") }
	p.SetGroup(group)
	return p
}

// SetGroup restarts the player on another show group with default settings.
func (this *Player) SetGroup(group int) {
	this.cartoons = this.blob.Group(group)
	this.current, this.running, this.change = 0, false, false
	this.PlayAll, this.Paused, this.DisplayOff = true, false, false
	this.speed = defaultSpeed
	this.keyCapture, this.capturedKey = false, 0
	this.vm = nil
}

func (this *Player) Current() int { return this.current }

func (this *Player) Speed() int { return this.speed }

// SetSpeed sets the speed index 0..7 (4 = normal).
func (this *Player) SetSpeed(speed int) {
	if speed >= 0 && speed < len(speedTenths) {
		this.speed = speed
	}
}

// Select restarts the show on cartoon n of the group.
func (this *Player) Select(n int) {
	if n >= 0 && n < len(this.cartoons) {
		this.current, this.change = n, true
	}
}

func (this *Player) log(msg string) {
	if this.Log != nil {
		this.Log(msg)
	}
}

// Step advances the show. When rendered is true a new frame is in the canvas and should be
// shown for waitMs; otherwise call Step again after waitMs. now is a millisecond clock.
func (this *Player) Step(now uint32) (rendered bool, waitMs int) {
	if this.Paused {
		return false, pollMs
	}

	if len(this.cartoons) == 0 {
		return false, pollMs
	}

	for tries := 0; tries <= 2*len(this.cartoons)+1; tries++ {
		if this.change {
			this.change, this.running, this.vm = false, false, nil
		}
		c := this.cartoons[this.current]
		if !this.running {
			this.start(c, now)
		}
		switch c.Type {
		case blob.Normal:
			if rendered, waitMs = this.stepNormal(c); rendered {
				return
			}
		case blob.Audio:
			if rendered, waitMs = this.stepAudio(c, now); rendered {
				return
			}
		case blob.Script:
			var done bool
			if rendered, waitMs, done = this.stepScript(c, now); !done {
				return
			}
		}
		this.finish()
	}

	return false, pollMs
}

func (this *Player) start(cartoon blob.Cartoon, now uint32) {
	this.running, this.loop, this.frame, this.drawn, this.until = true, 0, 0, false, 0
	if cartoon.Type == blob.Audio {
		if secs := cartoon.Header().FrameCount; secs != 0 {
			this.until = now + uint32(secs)*1000
		}
	}
	if cartoon.Type == blob.Script {
		this.vm = newVM(this, cartoon.Body)
	}
}

// finish ends the current cartoon; the next one follows the play model.
func (this *Player) finish() {
	this.running, this.vm = false, nil
	if this.PlayAll {
		this.current = (this.current + 1) % len(this.cartoons)
	}
}

func (this *Player) frameWait(delay int) int {
	return delay * 10 / speedTenths[this.speed]
}

func (this *Player) stepNormal(cartoon blob.Cartoon) (bool, int) {
	header := cartoon.Header()
	loops := int(uint16(header.Loop))
	if header.FrameCount <= 0 || this.loop >= loops {
		return false, 0
	}

	this.canvas.Clear(0)
	renderFrame(cartoon, this.frame, this.canvas, this.assets)
	this.frame++
	if this.frame >= header.FrameCount {
		this.frame = 0
		this.loop++
	}
	return true, this.frameWait(header.Delay)
}

func (this *Player) stepAudio(cartoon blob.Cartoon, now uint32) (bool, int) {
	if this.drawn && this.until != 0 && int32(now-this.until) >= 0 {
		return false, 0
	}

	frame := 0
	if this.AudioFrame != nil {
		frame = this.AudioFrame(int(cartoon.Header().Loop))
	}
	this.canvas.Clear(0)
	renderFrame(cartoon, frame, this.canvas, this.assets)
	this.drawn = true
	return true, this.frameWait(cartoon.Header().Delay)
}

// Key runs the configured function of KEY1..KEY7 (key 0..6).
func (this *Player) Key(key int) {
	this.Action(this.blob.KeyFunction(key))
}

func (this *Player) Action(action int) {
	if action < 1 {
		return
	}

	if this.keyCapture {
		this.capturedKey = action
		return
	}

	n := len(this.cartoons)
	if n == 0 {
		return
	}

	if action >= ActionJumpCartoon && action < ActionJumpCartoon+10 {
		target := action - ActionJumpCartoon
		if target >= n {
			return
		}

		this.log("Change Cartoon " + strconv.Itoa(target))
		this.current, this.change = target, true
		return
	}

	// Messages match the original firmware's console output
	switch action {
	case ActionReplay:
		this.log("Replay Cartoon")
		this.change = true
	case ActionPrevious:
		this.log("Previous Cartoon")
		this.current = (this.current + n - 1) % n
		this.change = true
	case ActionNext:
		this.log("Next Cartoon")
		this.current = (this.current + 1) % n
		this.change = true
	case ActionSpeedUp:
		if this.speed < len(speedTenths)-1 {
			this.speed++
			this.log("Add speed to " + this.speedText())
		}
	case ActionSpeedReset:
		if this.speed != defaultSpeed {
			this.speed = defaultSpeed
			this.log("Speed reset to " + this.speedText())
		}
	case ActionSpeedDown:
		if this.speed > 0 {
			this.speed--
			this.log("Sub speed to " + this.speedText())
		}
	case ActionPlayModel:
		this.PlayAll = !this.PlayAll
		this.log("Play Model=" + boolText(this.PlayAll))
	case ActionPause:
		this.Paused = !this.Paused
		this.log("DeviceSuppend:" + boolText(this.Paused))
	case ActionDisplayOff:
		this.DisplayOff = !this.DisplayOff
		this.log("Display off=" + boolText(this.DisplayOff))
	}
}

// speedText formats the speed factor like "%.2f".
func (this *Player) speedText() string {
	tenths := speedTenths[this.speed]
	return strconv.Itoa(tenths/10) + "." + strconv.Itoa(tenths%10) + "0"
}

func boolText(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

// stepScript runs the script until it waits, ends or uses up its instruction budget.
func (this *Player) stepScript(cartoon blob.Cartoon, now uint32) (rendered bool, waitMs int, done bool) {
	v := this.vm
	if v == nil {
		return false, 0, true
	}
	for i := 0; i < maxSteps; i++ {
		if v.phase == 2 || v.failed {
			return false, 0, true
		}
		if wait, yielded := v.exec(now); yielded {
			return true, this.frameWait(wait), false
		}
	}
	return false, 0, false
}
