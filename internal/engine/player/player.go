// Package player runs the original iCubeSmart show engine: it walks the cartoons of a show
// group, renders their frames and executes key functions. It is a non-blocking port of the
// vendor player loop (Simulator.dll FUN_10004ca0); see .agent/re/spec-engine.md.
//
// Player is not safe for concurrent use: call Key/Action from the same loop that calls Step
// (queue key presses from interrupts and apply them before Step).
package player

import (
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

	// AudioFrame selects the frame of Audio cartoons (0..63). The vendor PC simulator always
	// uses 0; on the cube it follows the microphone level.
	AudioFrame func() int
	// ADC feeds the script functions GetAdc/GetRealAdc.
	ADC func(channel int, raw bool) int
	// Log receives the engine's console messages.
	Log func(msg string)
}

// New creates a player for show group 0..3 drawing into canvas.
func New(b *blob.Blob, canvas *render.Canvas, group int) *Player {
	p := &Player{blob: b, canvas: canvas, assets: newAssets(b)}
	p.assets.Missing = func(kind, name string) { p.log("Not found " + kind + ":'" + name + "'") }
	p.SetGroup(group)
	return p
}

// SetGroup restarts the player on another show group with default settings.
func (p *Player) SetGroup(group int) {
	p.cartoons = p.blob.Group(group)
	p.current, p.running, p.change = 0, false, false
	p.PlayAll, p.Paused, p.DisplayOff = true, false, false
	p.speed = defaultSpeed
	p.keyCapture, p.capturedKey = false, 0
	p.vm = nil
}

func (p *Player) Current() int { return p.current }

func (p *Player) Speed() int { return p.speed }

func (p *Player) log(msg string) {
	if p.Log != nil {
		p.Log(msg)
	}
}

// Step advances the show. When rendered is true a new frame is in the canvas and should be
// shown for waitMs; otherwise call Step again after waitMs. now is a millisecond clock.
func (p *Player) Step(now uint32) (rendered bool, waitMs int) {
	if p.Paused {
		return false, pollMs
	}
	if len(p.cartoons) == 0 {
		return false, pollMs
	}
	for tries := 0; tries <= 2*len(p.cartoons)+1; tries++ {
		if p.change {
			p.change, p.running, p.vm = false, false, nil
		}
		c := p.cartoons[p.current]
		if !p.running {
			p.start(c, now)
		}
		switch c.Type {
		case blob.Normal:
			if rendered, waitMs = p.stepNormal(c); rendered {
				return
			}
		case blob.Audio:
			if rendered, waitMs = p.stepAudio(c, now); rendered {
				return
			}
		case blob.Script:
			var done bool
			if rendered, waitMs, done = p.stepScript(c, now); !done {
				return
			}
		}
		p.finish()
	}
	return false, pollMs
}

func (p *Player) start(c blob.Cartoon, now uint32) {
	p.running, p.loop, p.frame, p.drawn, p.until = true, 0, 0, false, 0
	if c.Type == blob.Audio {
		if secs := c.Header().FrameCount; secs != 0 {
			p.until = now + uint32(secs)*1000
		}
	}
	if c.Type == blob.Script {
		p.vm = newVM(p, c.Body)
	}
}

// finish ends the current cartoon; the next one follows the play model.
func (p *Player) finish() {
	p.running, p.vm = false, nil
	if p.PlayAll {
		p.current = (p.current + 1) % len(p.cartoons)
	}
}

func (p *Player) frameWait(delay int) int {
	return delay * 10 / speedTenths[p.speed]
}

func (p *Player) stepNormal(c blob.Cartoon) (bool, int) {
	h := c.Header()
	loops := int(uint16(h.Loop))
	if h.FrameCount <= 0 || p.loop >= loops {
		return false, 0
	}
	p.canvas.Clear(0)
	renderFrame(c, p.frame, p.canvas, p.assets)
	p.frame++
	if p.frame >= h.FrameCount {
		p.frame = 0
		p.loop++
	}
	return true, p.frameWait(h.Delay)
}

func (p *Player) stepAudio(c blob.Cartoon, now uint32) (bool, int) {
	if p.drawn && p.until != 0 && int32(now-p.until) >= 0 {
		return false, 0
	}
	frame := 0
	if p.AudioFrame != nil {
		frame = p.AudioFrame()
	}
	p.canvas.Clear(0)
	renderFrame(c, frame, p.canvas, p.assets)
	p.drawn = true
	return true, p.frameWait(c.Header().Delay)
}

// Key runs the configured function of KEY1..KEY7 (key 0..6).
func (p *Player) Key(key int) {
	p.Action(p.blob.KeyFunction(key))
}

func (p *Player) Action(action int) {
	if action < 1 {
		return
	}
	if p.keyCapture {
		p.capturedKey = action
		return
	}
	n := len(p.cartoons)
	if n == 0 {
		return
	}
	if action >= ActionJumpCartoon && action < ActionJumpCartoon+10 {
		target := action - ActionJumpCartoon
		if target >= n {
			return
		}
		p.log("Change Cartoon")
		p.current, p.change = target, true
		return
	}
	switch action {
	case ActionReplay:
		p.log("Replay Cartoon")
		p.change = true
	case ActionPrevious:
		p.log("Previous Cartoon")
		p.current = (p.current + n - 1) % n
		p.change = true
	case ActionNext:
		p.log("Next Cartoon")
		p.current = (p.current + 1) % n
		p.change = true
	case ActionSpeedUp:
		if p.speed < len(speedTenths)-1 {
			p.speed++
			p.log("Add speed")
		}
	case ActionSpeedReset:
		if p.speed != defaultSpeed {
			p.speed = defaultSpeed
			p.log("Speed reset")
		}
	case ActionSpeedDown:
		if p.speed > 0 {
			p.speed--
			p.log("Sub speed")
		}
	case ActionPlayModel:
		p.PlayAll = !p.PlayAll
		p.log("Play Model")
	case ActionPause:
		p.Paused = !p.Paused
		p.log("Pause")
	case ActionDisplayOff:
		p.DisplayOff = !p.DisplayOff
		p.log("Display off")
	}
}
