// Package upload receives light shows from the vendor "Download Software" (the iCubeSmart PC
// app's update tool). It answers the tool's update protocol like the original IAP bootloader,
// but keeps only the show data of the uploaded firmware (Update.hex) and stores it in a slot.
//
// Frames on UART1 are A5 05 len payload check 55 (check = XOR of all other bytes):
//
//	PC -> cube  [17, md5[16], size u32]                 connect
//	cube -> PC  [16, 0, 0, 0, 0, 0]                      ready (non-zero byte 1 = nothing to do)
//	PC -> cube  [32, n, line u16 BE, address u32 BE, data[n]]  one hex record
//	cube -> PC  [33, n, line u16 BE]                     record stored
//	PC -> cube  [37, 0]                                  end (sent twice)
package upload

const (
	// Show data region of the vendor firmware image (Compile.cs writes it at 0x0800B000).
	ShowStart = 0x0800B000
	ShowEnd   = 0x08032100
	ShowSize  = ShowEnd - ShowStart

	msgConnect = 17
	msgReady   = 16
	msgRecord  = 32
	msgStored  = 33
	msgEnd     = 37
)

// Storage keeps the received show data.
type Storage interface {
	Begin() error
	Write(offset int, data []byte) error
}

type Receiver struct {
	Storage Storage
	// Reply sends a payload to the PC (wrapped in a frame by the caller).
	Reply func(payload []byte)
	// Started is called on connect, before anything is written.
	Started func()
	// Finished is called at the end with whether every record was stored.
	Finished func(ok bool)

	active bool
	failed bool
}

// Frame handles one received payload.
func (this *Receiver) Frame(payload []byte) {
	if len(payload) == 0 {
		return
	}
	switch payload[0] {
	case msgConnect:
		if this.Started != nil {
			this.Started()
		}
		this.active = true
		this.failed = this.Storage.Begin() != nil
		this.Reply([]byte{msgReady, 0, 0, 0, 0, 0})
	case msgRecord:
		if !this.active || len(payload) < 8 || len(payload) != 8+int(payload[1]) {
			return
		}
		address := int(payload[4])<<24 | int(payload[5])<<16 | int(payload[6])<<8 | int(payload[7])
		if !this.store(address, payload[8:]) {
			this.failed = true
			return
		}
		this.Reply([]byte{msgStored, payload[1], payload[2], payload[3]})
	case msgEnd:
		if !this.active {
			return
		}
		this.active = false
		if this.Finished != nil {
			this.Finished(!this.failed)
		}
	}
}

// store writes the part of a record inside the show region; other records are acknowledged
// without storing them.
func (this *Receiver) store(address int, data []byte) bool {
	start, end := max(address, ShowStart), min(address+len(data), ShowEnd)
	if start >= end {
		return true
	}
	return this.Storage.Write(start-ShowStart, data[start-address:end-address]) == nil
}

// Frame wraps a payload for the PC: A5 05 len payload check 55.
func Frame(payload []byte) []byte {
	frame := make([]byte, 0, len(payload)+5)
	frame = append(frame, 0xA5, 0x05, byte(len(payload)))
	frame = append(frame, payload...)
	check := byte(0xA5 ^ 0x05 ^ 0x55 ^ byte(len(payload)))
	for _, b := range payload {
		check ^= b
	}
	return append(frame, check, 0x55)
}
