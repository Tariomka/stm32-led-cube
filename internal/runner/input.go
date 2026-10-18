package runner

import (
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/cascade"
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/remote"
)

const (
	cycleKey      = 5 // KEY6: short press = show key function, long press = cycle runner mode
	cycleBit      = 1 << cycleKey
	longPressTime = time.Second
	syncPeriod    = 8 * time.Millisecond // master Sync rate, like the original tick task
)

// input handles keys, the IR remote, the console and the cube cascade.
//
// Like the original, key functions don't run straight from the key interrupts: the master
// turns key presses into the key mask of its Sync frames and every cube (master included)
// runs a key's function on the rising edge of its bit, so cascaded cubes stay in step. The
// master also provides the show clock. Slaves ignore their own keys and IR remote; KEY6 held
// for 1 s still cycles the runner mode locally.
type input struct {
	link       *cascade.Link
	keyMask    cascade.KeyMask
	slaveKeys  uint8 // last key mask received by a slave
	pressed    uint8 // key presses not sent yet
	showStart  time.Time
	lastSync   time.Time
	cycleSince time.Time // when the pending KEY6 press started, zero if none
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
		if key != cycleKey && this.Board.KeyHeld(key) {
			held |= 1 << key
		}
	}
	return held
}

func (this *CubeRunner) runKeys(edges uint8) {
	for key := 0; key < controller.KeyCount; key++ {
		if edges&(1<<key) != 0 {
			this.action(func() { this.Player.Key(key) })
		}
	}
}

// handleInput processes console, cascade, keys and IR; it returns true when the runner mode
// changed.
func (this *CubeRunner) handleInput() bool {
	this.Board.PollSerial()
	if !this.saveAt.IsZero() && time.Now().After(this.saveAt) {
		this.saveSettings()
	}
	synced := this.pollCascade()

	events := this.Board.TakeKeys()
	if events&cycleBit != 0 {
		this.cycleSince = time.Now()
	}
	this.pressed |= events &^ cycleBit
	modeChanged, shortPress := this.checkCycleKey()
	if shortPress {
		this.pressed |= cycleBit
	}

	code := this.Board.TakeRemote()
	if !this.link.Master {
		if synced {
			this.runKeys(cascade.Edges(&this.slaveKeys, this.link.Keys))
		}
		this.pressed = 0
		return modeChanged
	}

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
	return modeChanged
}

// checkCycleKey tells a KEY6 long press (cycle the runner mode) from a short one (key function).
func (this *CubeRunner) checkCycleKey() (modeChanged, shortPress bool) {
	if this.cycleSince.IsZero() {
		return false, false
	}
	if this.Board.KeyHeld(cycleKey) {
		if time.Since(this.cycleSince) < longPressTime {
			return false, false
		}
		this.cycleSince = time.Time{}
		this.Tracker.CycleMode()
		this.Board.SetBlank(false)
		return true, false
	}
	this.cycleSince = time.Time{}
	return false, true
}
