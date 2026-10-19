package runner

import (
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/cascade"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/remote"
)

const syncPeriod = 8 * time.Millisecond // master Sync rate, like the original tick task

// input handles the DIP switches, keys, the IR remote, the console and the cube cascade.
//
// Like the original, key functions don't run straight from the key interrupts: the master
// turns key presses into the key mask of its Sync frames and every cube (master included)
// runs a key's function on the rising edge of its bit, so cascaded cubes stay in step. The
// master also provides the show clock. Slaves ignore their own keys and IR remote.
type input struct {
	link      *cascade.Link
	keyMask   cascade.KeyMask
	slaveKeys uint8 // last key mask received by a slave
	pressed   uint8 // key presses not sent yet
	showStart time.Time
	lastSync  time.Time
}

// clock is the show time in ms: the master's own, or the last one a slave received.
func (this *CubeRunner) clock() uint64 {
	if this.link.Master {
		return uint64(time.Since(this.showStart).Milliseconds())
	}
	return this.link.Clock
}

func (this *CubeRunner) pollCascade() (synced bool) {
	this.Board.CascadePoll(func(payload []byte) {
		if this.link.Handle(payload) {
			synced = true
		}
	})
	return synced
}

// discoverCascade finds this cube's place in the chain at power-up (about 1.1 s).
func (this *CubeRunner) discoverCascade() {
	this.link = cascade.NewLink(this.Board.CascadeSend)
	this.link.Discover(func(ms int) {
		until := time.Now().Add(time.Duration(ms) * time.Millisecond)
		for time.Now().Before(until) {
			this.pollCascade()
			this.Board.PollSerial()
			time.Sleep(time.Millisecond)
		}
	})
	deviceType := 0
	if this.link.Master {
		deviceType = 1
	}
	println("DeviceType=" + strconv.Itoa(deviceType) + ",DeviceNum=" + strconv.Itoa(this.link.Number))
	this.showStart = time.Now()
}

func (this *CubeRunner) heldKeys() uint8 {
	var held uint8
	for key := 0; key < controller.KeyCount; key++ {
		if this.Board.KeyHeld(key) {
			held |= 1 << key
		}
	}
	return held
}

func (this *CubeRunner) runKeys(edges uint8) {
	for key := 0; key < controller.KeyCount; key++ {
		if edges&(1<<key) != 0 {
			this.action(func() {
				println("DoKey" + strconv.Itoa(key+1))
				this.Player.Key(key)
			})
		}
	}
}

// handleInput processes the console, DIP switches, cascade, keys and IR; it returns true when
// the show was replaced.
func (this *CubeRunner) handleInput() bool {
	this.Board.PollSerial()
	if !this.saveAt.IsZero() && time.Now().After(this.saveAt) {
		this.saveSettings()
	}
	if this.Tracker.Update() {
		this.changeMode()
	}
	synced := this.pollCascade()
	this.pressed |= this.Board.TakeKeys()
	code := this.Board.TakeRemote()

	if !this.link.Master {
		if synced {
			this.runKeys(cascade.Edges(&this.slaveKeys, this.link.Keys))
		}
		this.pressed = 0
	} else {
		if time.Since(this.lastSync) >= syncPeriod {
			this.lastSync = time.Now()
			mask, edges := this.keyMask.Next(this.heldKeys(), this.pressed)
			this.pressed = 0
			this.link.SendSync(this.clock(), mask)
			this.runKeys(edges)
		}
		if code != 0 {
			println("IR:" + strconv.FormatInt(int64(code), 16))
			this.action(func() { this.Player.Action(remote.Action(code)) })
		}
	}

	changed := this.showChanged
	this.showChanged = false
	return changed
}
