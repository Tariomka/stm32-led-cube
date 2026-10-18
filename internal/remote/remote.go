// Package remote decodes the NEC infrared remote of the original kit and maps its buttons to
// show key actions, like the original firmware (FUN_08033eb4 decoder, FUN_08034ddc table).
package remote

// Decoder turns mark/space durations (µs) into 16-bit codes: address << 8 | command.
// Bits are stored MSB first, so the codes match the ones printed for common 21-key remotes.
type Decoder struct {
	bit  int // -1 = waiting for a leader
	data [4]byte
}

func NewDecoder() *Decoder { return &Decoder{bit: -1} }

func within(value, from, width int) bool { return value >= from && value < from+width }

// Pulse feeds one mark (IR burst) and the space after it. ok is true when a full frame with
// valid complements has been received. Repeat frames are ignored.
func (this *Decoder) Pulse(mark, space int) (code int, ok bool) {
	switch {
	case within(mark, 7001, 3999) && within(space, 3501, 1999): // leader
		this.bit, this.data = 0, [4]byte{}
		return 0, false
	case within(mark, 411, 299) && within(space, 411, 299): // 0
	case within(mark, 411, 299) && within(space, 1281, 799): // 1
		if this.bit >= 0 {
			this.data[this.bit/8] |= 0x80 >> (this.bit % 8)
		}
	default: // repeat leader or noise
		this.bit = -1
		return 0, false
	}

	if this.bit < 0 {
		return 0, false
	}

	this.bit++
	if this.bit < 32 {
		return 0, false
	}

	this.bit = -1
	if this.data[0] != ^this.data[1] || this.data[2] != ^this.data[3] {
		return 0, false
	}

	return int(this.data[0])<<8 | int(this.data[2]), true
}

// Key actions of the remote buttons (address 0x00). Labels are for the common 21-key remote.
var actions = map[int]int{
	0xA2: 9,  // CH-: display off
	0xE2: 7,  // CH+: play model
	0x02: 4,  // >>|: speed up
	0xC2: 5,  // >||: speed reset
	0xE0: 2,  // VOL-: previous
	0xA8: 8,  // VOL+: pause
	0x90: 3,  // EQ: next
	0x98: 6,  // 100+: speed down
	0xB0: 1,  // 200+: replay
	0x68: 10, // 0..9: jump to cartoon 0..9
	0x30: 11,
	0x18: 12,
	0x7A: 13,
	0x10: 14,
	0x38: 15,
	0x5A: 16,
	0x42: 17,
	0x4A: 18,
	0x52: 19,
}

// Action returns the key action for a decoded code, 0 when the button has none.
func Action(code int) int { return actions[code] }
