package controller

import (
	"encoding/binary"
	"strconv"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// Cascade messages
const (
	msgHello  = 0x01 // [1]
	msgNumber = 0x02 // [2, number of the receiving cube]
	msgSync   = 0x20 // [0x20, clock ms (u64 LE), key mask]

	syncLength        = 10
	cascadeMaxPayload = 15
	syncPeriod        = 8 * time.Millisecond // master Sync rate, like the factory tick task
	keyPulseFrames    = 3
)

const (
	discovering = iota // listening for an upstream Hello
	numbering          // slave waiting for its number
	running
)

// CascadeController chains cubes over the 3.5 mm jacks (UART2, TX of one cube to the RX of the
// next) with the factory protocol: frames A5 55 len payload xor. At power-up every cube sends
// Hello; a cube that hears one is a slave, the master numbers the chain. From then on the master
// broadcasts Sync with its show clock and key mask, so all cubes play in step and run the same
// key functions. Slaves ignore their own keys.
type CascadeController struct {
	port *component.SerialPort

	// Frame parser
	parseState int
	parseSize  int
	parseCheck byte
	parseBuf   [cascadeMaxPayload + 1]byte

	master bool
	number int // position in the chain, 0 = master
	state  int

	showStart time.Time // master clock origin
	clock     uint64    // slave: last master clock (ms)
	lastSync  time.Time // master: when the last Sync was sent

	pressed   uint8 // master: presses not sent yet
	latch     [8]uint8
	held      [8]uint8
	lastMask  uint8 // key mask of the previous Sync
	slaveKeys uint8 // slave: key mask of the last Sync received
}

func NewCascadeController() *CascadeController {
	return &CascadeController{port: component.NewCascadePort(), master: true}
}

func (this *CascadeController) Master() bool { return this.master }

func (this *CascadeController) Number() int { return this.number }

// Clock is the show time in ms shared by all cubes: the master's own, or the last one a slave
// received.
func (this *CascadeController) Clock() uint64 {
	if this.master {
		return uint64(time.Since(this.showStart).Milliseconds())
	}
	return this.clock
}

// Discover finds this cube's place in the chain at power-up (about 1.1 s); idle is called every
// millisecond meanwhile.
func (this *CascadeController) Discover(idle func()) {
	wait := func(ms int) {
		until := time.Now().Add(time.Duration(ms) * time.Millisecond)
		for time.Now().Before(until) {
			idle()
			this.poll()
			time.Sleep(time.Millisecond)
		}
	}

	this.state = discovering
	wait(100)
	for range 5 {
		this.send([]byte{msgHello})
		wait(100)
	}
	if !this.master {
		this.state = numbering
		wait(500)
	} else {
		for range 5 {
			this.send([]byte{msgNumber, 1})
			wait(100)
		}
	}
	this.state = running

	deviceType := 0
	if this.master {
		deviceType = 1
	}
	println("DeviceType=" + strconv.Itoa(deviceType) + ",DeviceNum=" + strconv.Itoa(this.number))
	this.showStart = time.Now()
}

// Keys exchanges key states with the chain and returns the keys whose function runs now (bit n =
// KEY(n+1)). The master sends its presses and held keys with every Sync (8 ms); a slave returns
// the keys of the Syncs it receives.
func (this *CascadeController) Keys(pressed, held uint8) uint8 {
	synced := this.poll()
	if !this.master {
		if !synced {
			return 0
		}

		edges := this.slaveKeys &^ this.lastMask
		this.lastMask = this.slaveKeys
		return edges
	}

	this.pressed |= pressed
	if time.Since(this.lastSync) < syncPeriod {
		return 0
	}

	this.lastSync = time.Now()
	mask := this.keyMask(held, this.pressed)
	this.pressed = 0
	payload := make([]byte, syncLength)
	payload[0] = msgSync
	binary.LittleEndian.PutUint64(payload[1:9], this.Clock())
	payload[9] = mask
	this.send(payload)

	edges := mask &^ this.lastMask
	this.lastMask = mask
	return edges
}

// keyMask builds the key level mask of a Sync. The factory firmware sets a bit after 15 stable
// key polls; here a press holds its bit for a few Syncs so short taps reach the slaves, and a key
// held for 3 Syncs keeps its bit.
func (this *CascadeController) keyMask(held, pressed uint8) uint8 {
	var mask uint8
	for index := range 8 {
		bit := uint8(1) << index
		if pressed&bit != 0 {
			this.latch[index] = keyPulseFrames
		}
		if held&bit != 0 {
			this.held[index] = min(this.held[index]+1, keyPulseFrames)
		} else {
			this.held[index] = 0
		}
		if this.latch[index] > 0 {
			this.latch[index]--
			mask |= bit
		} else if this.held[index] >= keyPulseFrames {
			mask |= bit
		}
	}
	return mask
}

// poll handles the received payloads and reports whether a Sync arrived.
func (this *CascadeController) poll() (synced bool) {
	for {
		b, ok := this.port.Read()
		if !ok {
			return synced
		}
		if payload, ok := this.parse(b); ok && this.handle(payload) {
			synced = true
		}
	}
}

func (this *CascadeController) handle(payload []byte) (synced bool) {
	switch {
	case len(payload) == 1 && payload[0] == msgHello:
		if this.state == discovering {
			this.master = false
		}
	case len(payload) == 2 && payload[0] == msgNumber:
		if this.state == numbering {
			this.number = int(payload[1])
			this.send([]byte{msgNumber, payload[1] + 1})
		}
	case len(payload) == syncLength && payload[0] == msgSync:
		if this.master {
			return false
		}

		this.send(payload)
		this.clock = binary.LittleEndian.Uint64(payload[1:9])
		this.slaveKeys = payload[9]
		return true
	}
	return false
}

// send transmits a payload downstream: A5 55 len payload xor (xor = len ^ payload bytes).
func (this *CascadeController) send(payload []byte) {
	frame := make([]byte, 0, len(payload)+4)
	frame = append(frame, 0xA5, 0x55, byte(len(payload)))
	check := byte(len(payload))
	for _, b := range payload {
		check ^= b
	}
	frame = append(frame, payload...)
	this.port.Write(append(frame, check))
}

// parse returns a payload (valid until the next call) when b completes a frame.
func (this *CascadeController) parse(b byte) ([]byte, bool) {
	switch {
	case this.parseState == 0:
		if b == 0xA5 {
			this.parseState = 1
		}
	case this.parseState == 1:
		this.parseState = 0
		if b == 0x55 {
			this.parseState = 2
		}
	case this.parseState == 2:
		this.parseState = 0
		if b >= 1 && b <= cascadeMaxPayload+1 {
			this.parseSize, this.parseCheck, this.parseState = int(b), b, 3
		}
	case this.parseState < 3+this.parseSize:
		this.parseBuf[this.parseState-3] = b
		this.parseCheck ^= b
		this.parseState++
	default:
		this.parseState = 0
		if b == this.parseCheck {
			return this.parseBuf[:this.parseSize], true
		}
	}
	return nil, false
}
