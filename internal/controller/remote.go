package controller

import (
	"runtime/interrupt"
	"strconv"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// Key functions of the kit's 21-key remote (address 0x00), like the factory firmware.
var remoteActions = map[int]int{
	0xA2: 9,  // CH-: display off
	0xE2: 7,  // CH+: play model
	0x02: 4,  // >>|: speed up
	0xC2: 5,  // >||: speed reset
	0xE0: 2,  // VOL-: previous
	0xA8: 8,  // VOL+: pause
	0x90: 3,  // EQ: next
	0x98: 6,  // 100+: speed down
	0xB0: 1,  // 200+: replay
	0x68: 10, // 0..9: play cartoon 0..9
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

// RemoteController decodes the kit's NEC infrared remote (receiver on PC6) into key function codes.
type RemoteController struct {
	receiver *component.InfraredReceiver
	bit      int // next bit of the frame, -1 = waiting for a leader
	data     [4]byte
	code     int // last decoded code (address << 8 | command), 0 = none
}

func NewRemoteController() *RemoteController {
	remote := &RemoteController{bit: -1}
	remote.receiver = component.NewInfraredReceiver(remote.pulse)
	return remote
}

// Take returns the key function of the last button received and clears it, 0 when none. It logs
// the code like the factory firmware.
func (this *RemoteController) Take() int {
	mask := interrupt.Disable()
	code := this.code
	this.code = 0
	interrupt.Restore(mask)
	if code == 0 {
		return 0
	}

	println("IR:" + strconv.FormatInt(int64(code), 16))
	return remoteActions[code]
}

func within(value, from, width int) bool { return value >= from && value < from+width }

// pulse decodes one mark and the space after it (µs); runs in interrupt context. Bits are stored
// MSB first so the codes match the ones printed for common 21-key remotes; repeats are ignored.
func (this *RemoteController) pulse(mark, space int) {
	switch {
	case within(mark, 7001, 3999) && within(space, 3501, 1999): // leader
		this.bit, this.data = 0, [4]byte{}
		return
	case within(mark, 411, 299) && within(space, 411, 299): // 0
	case within(mark, 411, 299) && within(space, 1281, 799): // 1
		if this.bit >= 0 {
			this.data[this.bit/8] |= 0x80 >> (this.bit % 8)
		}
	default: // repeat leader or noise
		this.bit = -1
		return
	}

	if this.bit < 0 {
		return
	}
	if this.bit++; this.bit < 32 {
		return
	}

	this.bit = -1
	if this.data[0] == ^this.data[1] && this.data[2] == ^this.data[3] {
		this.code = int(this.data[0])<<8 | int(this.data[2])
	}
}
