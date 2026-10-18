// Package cascade is the original firmware's protocol for chaining cubes over the 3.5 mm jacks
// (USART2, 115200 baud, TX of one cube to RX of the next). See .agent/re/firmware-map.md.
//
// Frames are A5 55 len payload xor, with xor = len ^ payload bytes and len 1..15. At power-up
// every cube sends Hello 5 times; a cube that hears a Hello is a slave. The master then numbers
// the chain (Number), and from then on broadcasts Sync with its show clock and key state.
package cascade

import "encoding/binary"

const (
	MsgHello  = 0x01 // [1]
	MsgNumber = 0x02 // [2, number of the receiving cube]
	MsgSync   = 0x20 // [0x20, clock ms (u64 LE), key mask]

	SyncLength = 10
	maxPayload = 15
)

// Encode wraps a payload in a frame.
func Encode(payload []byte) []byte {
	frame := make([]byte, 0, len(payload)+4)
	frame = append(frame, 0xA5, 0x55, byte(len(payload)))
	check := byte(len(payload))
	for _, b := range payload {
		check ^= b
	}
	frame = append(frame, payload...)
	return append(frame, check)
}

// Parser extracts payloads from the received byte stream.
type Parser struct {
	state int
	size  int
	check byte
	buf   [maxPayload + 1]byte
}

// Feed returns a payload (valid until the next call) when b completes a frame.
func (this *Parser) Feed(b byte) ([]byte, bool) {
	switch {
	case this.state == 0:
		if b == 0xA5 {
			this.state = 1
		}
	case this.state == 1:
		this.state = 0
		if b == 0x55 {
			this.state = 2
		}
	case this.state == 2:
		this.state = 0
		if b >= 1 && b <= 16 {
			this.size, this.check, this.state = int(b), b, 3
		}
	case this.state < 3+this.size:
		this.buf[this.state-3] = b
		this.check ^= b
		this.state++
	default:
		this.state = 0
		if b == this.check {
			return this.buf[:this.size], true
		}
	}
	return nil, false
}

const (
	stateDiscover = iota // listening for an upstream Hello
	stateNumber          // slave waiting for its number
	stateRunning
)

// Link is one cube's view of the chain.
type Link struct {
	// Send transmits a payload downstream.
	Send func(payload []byte)

	Master bool   // no cube upstream
	Number int    // position in the chain, 0 = master
	Clock  uint64 // slave: last master clock (ms)
	Keys   uint8  // slave: last master key mask
	state  int
}

func NewLink(send func([]byte)) *Link {
	return &Link{Send: send, Master: true}
}

// Handle processes a received payload and reports whether it was a Sync.
func (this *Link) Handle(payload []byte) bool {
	switch {
	case len(payload) == 1 && payload[0] == MsgHello:
		if this.state == stateDiscover {
			this.Master = false
		}
	case len(payload) == 2 && payload[0] == MsgNumber:
		if this.state == stateNumber {
			this.Number = int(payload[1])
			this.Send([]byte{MsgNumber, payload[1] + 1})
		}
	case len(payload) == SyncLength && payload[0] == MsgSync:
		if this.Master {
			return false
		}
		this.Send(payload)
		this.Clock = binary.LittleEndian.Uint64(payload[1:9])
		this.Keys = payload[9]
		return true
	}
	return false
}

// Discover runs the power-up sequence. wait(ms) must keep feeding received payloads to Handle.
func (this *Link) Discover(wait func(ms int)) {
	this.state = stateDiscover
	wait(100)
	for range 5 {
		this.Send([]byte{MsgHello})
		wait(100)
	}
	if !this.Master {
		this.state = stateNumber
		wait(500)
	} else {
		for range 5 {
			this.Send([]byte{MsgNumber, 1})
			wait(100)
		}
	}
	this.state = stateRunning
}

// SendSync broadcasts the master's clock and key mask.
func (this *Link) SendSync(clock uint64, keys uint8) {
	payload := make([]byte, SyncLength)
	payload[0] = MsgSync
	binary.LittleEndian.PutUint64(payload[1:9], clock)
	payload[9] = keys
	this.Send(payload)
}

// KeyMask builds the key level mask of Sync frames. The original polls keys and sets a bit
// after 15 stable polls; here a press edge holds its bit for a few frames so short taps reach
// the slaves, and a key held for 3 frames keeps its bit. Every cube runs a key's function on
// the rising edge of its bit, so the master and the slaves act on the same frame.
type KeyMask struct {
	latch [8]uint8
	held  [8]uint8
	last  uint8
}

const pulseFrames = 3

// Next returns the mask for this frame and its rising edges.
func (this *KeyMask) Next(held, pressed uint8) (mask, edges uint8) {
	for index := range 8 {
		bit := uint8(1) << index
		if pressed&bit != 0 {
			this.latch[index] = pulseFrames
		}
		if held&bit != 0 {
			this.held[index] = min(this.held[index]+1, pulseFrames)
		} else {
			this.held[index] = 0
		}
		if this.latch[index] > 0 {
			this.latch[index]--
			mask |= bit
		} else if this.held[index] >= pulseFrames {
			mask |= bit
		}
	}
	return mask, Edges(&this.last, mask)
}

// Edges returns the bits that rose since the previous mask and stores mask as previous.
func Edges(previous *uint8, mask uint8) uint8 {
	edges := mask &^ *previous
	*previous = mask
	return edges
}
