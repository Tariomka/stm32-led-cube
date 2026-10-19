package console

// Demux splits the UART1 byte stream between the text console ("#?command args?!") and the
// binary frames of the upload protocol (A5 05 len data check 55). Other bytes are ignored.
type Demux struct {
	state     int
	frameLeft int
	line      []byte
	last      uint32
	started   bool
	check     byte
	received  byte

	// Command receives each complete console line.
	Command func(line string)
	// Frame receives the payload of each valid binary frame.
	Frame func(payload []byte)
}

const (
	idle     = iota
	hash     // '#' received
	capture  // "#?" received, collecting the command line
	question // closing '?' received, waiting for '!'
	vendor05 // 0xA5 received, waiting for 0x05
	vendorLen
	vendorData
	vendorCheck
	vendorEnd

	maxLine     = 180
	lineTimeout = 1000 // ms without bytes resets a partial command or frame

	vendorStart = 0xA5
)

// Feed processes one received byte.
func (this *Demux) Feed(b byte, nowMs uint32) {
	if this.started && nowMs-this.last > lineTimeout {
		this.state = idle
	}
	this.started, this.last = true, nowMs

	switch this.state {
	case hash:
		if b == '?' {
			this.state, this.line = capture, this.line[:0]
			return
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
		return
	case question:
		this.state = idle
		if b == '!' && len(this.line) > 0 && this.Command != nil {
			this.Command(string(this.line))
		}
		return
	case vendor05:
		this.state = idle
		if b == 0x05 {
			this.state, this.check = vendorLen, vendorStart^b
			return
		}
	case vendorLen:
		if b == 0 || b >= maxLine {
			this.state = idle
			return
		}
		this.state, this.frameLeft, this.check, this.line = vendorData, int(b), this.check^b, this.line[:0]
		return
	case vendorData:
		this.line = append(this.line, b)
		this.check ^= b
		if this.frameLeft--; this.frameLeft == 0 {
			this.state = vendorCheck
		}
		return
	case vendorCheck:
		this.state, this.received = vendorEnd, b
		return
	case vendorEnd:
		this.state = idle
		if b == 0x55 && this.received == this.check^0x55 && this.Frame != nil {
			this.Frame(this.line)
		}
		return
	}

	switch b {
	case '#':
		this.state = hash
	case vendorStart:
		this.state = vendor05
	}
}
