package console

// Demux splits the UART1 byte stream between the original firmware's text console
// ("#?command args?!", like the vendor PC tools send), the vendor tools' binary frames
// (A5 05 len data check 55) and the Raspberry Pi Pico link (led-common-lib frames:
// 0xAA type len_hi len_lo payload crc32 0x55, plus the ASCII handshake).
type Demux struct {
	state     int
	header    [4]byte
	headerLen int
	frameLeft int
	line      []byte
	last      uint32
	started   bool
	check     byte

	// Command receives each complete console line.
	Command func(line string)
	// Frame receives the payload of each valid vendor binary frame.
	Frame func(payload []byte)
}

const (
	idle      = iota
	hash      // '#' received
	capture   // "#?" received, collecting the command line
	question  // closing '?' received, waiting for '!'
	frameHead // inside a Pico frame header
	frameBody // inside a Pico frame payload/trailer
	vendor05  // 0xA5 received, waiting for 0x05
	vendorLen
	vendorData
	vendorCheck
	vendorEnd

	maxLine     = 180
	lineTimeout = 1000 // ms without bytes resets a partial command, like the original

	frameStart   = 0xAA
	frameTrailer = 4 + 1 // crc32 + end marker
	vendorStart  = 0xA5
)

func (this *Demux) partial() bool {
	return (this.state >= hash && this.state <= question) || this.state >= vendor05
}

// Feed processes one received byte and reports whether it belongs to the Pico link.
func (this *Demux) Feed(b byte, nowMs uint32) bool {
	if this.started && nowMs-this.last > lineTimeout && this.partial() {
		this.state = idle
	}
	this.started, this.last = true, nowMs

	switch this.state {
	case hash:
		if b == '?' {
			this.state, this.line = capture, this.line[:0]
			return false
		}
		this.state = idle
	case capture:
		if b == '?' {
			this.state = question
		} else if len(this.line) < maxLine {
			this.line = append(this.line, b)
		} else {
			this.state = idle
		}
		return false
	case question:
		this.state = idle
		if b == '!' && len(this.line) > 0 && this.Command != nil {
			this.Command(string(this.line))
		}
		return false
	case frameHead:
		this.header[this.headerLen] = b
		this.headerLen++
		if this.headerLen == len(this.header) {
			this.state = frameBody
			this.frameLeft = (int(this.header[2])<<8 | int(this.header[3])) + frameTrailer
		}
		return true
	case frameBody:
		this.frameLeft--
		if this.frameLeft == 0 {
			this.state = idle
		}
		return true
	case vendor05:
		if b != 0x05 {
			this.state = idle
			return this.Feed(b, nowMs)
		}
		this.state, this.check = vendorLen, vendorStart^b
		return false
	case vendorLen:
		if b == 0 || b >= maxLine {
			this.state = idle
			return false
		}
		this.state, this.frameLeft, this.check, this.line = vendorData, int(b), this.check^b, this.line[:0]
		return false
	case vendorData:
		this.line = append(this.line, b)
		this.check ^= b
		if this.frameLeft--; this.frameLeft == 0 {
			this.state = vendorCheck
		}
		return false
	case vendorCheck:
		this.state, this.header[0] = vendorEnd, b
		return false
	case vendorEnd:
		this.state = idle
		if b == 0x55 && this.header[0] == this.check^0x55 && this.Frame != nil {
			this.Frame(this.line)
		}
		return false
	}

	switch b {
	case '#':
		this.state = hash
		return false
	case vendorStart:
		this.state = vendor05
		return false
	case frameStart:
		this.state, this.header[0], this.headerLen = frameHead, b, 1
	}
	return true
}
