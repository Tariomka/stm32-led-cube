package controller

// Show uploads use the protocol of the vendor "Download Software" (the iCubeSmart PC app's update
// tool), answered like the factory bootloader. Only the show data of the uploaded firmware image
// (Update.hex) is kept, in the show slot. Payloads of A5 05 len payload check 55 frames:
//
//	PC -> cube  [17, md5[16], size u32]                       connect
//	cube -> PC  [16, 0, 0, 0, 0, 0]                            ready (non-zero byte 1 = nothing to do)
//	PC -> cube  [32, n, line u16 BE, address u32 BE, data[n]]  one hex record
//	cube -> PC  [33, n, line u16 BE]                           record stored
//	PC -> cube  [37, 0]                                        end (sent twice)

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

const (
	// Show data region of the vendor firmware image (its show compiler writes it at 0x0800B000).
	showStart = 0x0800B000
	showEnd   = 0x08032100

	msgConnect = 17
	msgReady   = 16
	msgRecord  = 32
	msgStored  = 33
	msgEnd     = 37
)

type upload struct {
	slot   component.ShowSlotStorage
	active bool
	failed bool
}

// uploadFrame handles the payload of one received frame.
func (this *TerminalController) uploadFrame(payload []byte) {
	if len(payload) == 0 {
		return
	}

	switch payload[0] {
	case msgConnect:
		this.rebootAt = time.Time{}
		println("Show upload started")
		this.show.UploadStarted()
		this.upload.active = true
		this.upload.failed = this.upload.slot.Begin() != nil
		this.reply([]byte{msgReady, 0, 0, 0, 0, 0})
	case msgRecord:
		if !this.upload.active || len(payload) < 8 || len(payload) != 8+int(payload[1]) {
			return
		}
		address := int(payload[4])<<24 | int(payload[5])<<16 | int(payload[6])<<8 | int(payload[7])
		if !this.storeRecord(address, payload[8:]) {
			this.upload.failed = true
			return
		}
		this.reply([]byte{msgStored, payload[1], payload[2], payload[3]})
	case msgEnd:
		if !this.upload.active {
			return
		}
		this.upload.active = false
		if this.upload.failed {
			println("Show upload failed")
		} else {
			println("Show upload OK")
		}
		this.show.UploadFinished()
	}
}

// storeRecord writes the part of a record inside the show region; other records are acknowledged
// without storing them.
func (this *TerminalController) storeRecord(address int, data []byte) bool {
	start, end := max(address, showStart), min(address+len(data), showEnd)
	if start >= end {
		return true
	}
	return this.upload.slot.Write(start-showStart, data[start-address:end-address]) == nil
}

// reply sends a payload to the PC: A5 05 len payload check 55 (check = XOR of all other bytes).
func (this *TerminalController) reply(payload []byte) {
	frame := make([]byte, 0, len(payload)+5)
	frame = append(frame, frameStart, 0x05, byte(len(payload)))
	frame = append(frame, payload...)
	check := byte(frameStart ^ 0x05 ^ 0x55 ^ byte(len(payload)))
	for _, b := range payload {
		check ^= b
	}
	this.port.Write(append(frame, check, 0x55))
}
