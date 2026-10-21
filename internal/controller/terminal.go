package controller

import (
	"time"

	"github.com/Tariomka/led-common-lib/pkg/communication/uart"
	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// Byte stream states
const (
	streamIdle     = iota
	streamHash     // '#' received
	streamCommand  // "#?" received, collecting the command line
	streamQuestion // closing '?' received, waiting for '!'
	streamFrame05  // 0xA5 received, waiting for 0x05
	streamFrameLen
	streamFrameData
	streamFrameCheck
	streamFrameEnd

	frameKind = 0x05
)

// rebootDelay lets the vendor upload tool's connect frame, sent right after its "reboot" command,
// cancel the reboot: the factory firmware rebooted into its bootloader, we take the upload.
const rebootDelay = 300 * time.Millisecond

// TerminalController runs UART1 (core board header, TX - PA9, RX - PA10, 115200 8N1), the link to
// a PC terminal, the vendor PC tools and the Raspberry Pi Pico. It splits the received bytes into
// console commands ("#?command args?!", see terminalCommands.go) and binary frames
// (A5 05 len payload check 55) of show uploads (see terminalUpload.go). Replies and the log share the port.
type TerminalController struct {
	port     *component.SerialPort
	system   *SystemController
	settings *SettingsController
	show     *ShowController
	keys     *KeyController

	registry *registry
	upload   upload
	rebootAt time.Time // pending ReBoot, zero if none

	// Byte stream splitter
	state     int
	line      []byte
	frameLeft int
	check     byte
	received  byte
	start     time.Time
	last      uint32
	started   bool
}

func NewTerminalController(system *SystemController, settings *SettingsController, show *ShowController, keys *KeyController) *TerminalController {
	return &TerminalController{
		port:     component.NewConsolePort(),
		system:   system,
		settings: settings,
		show:     show,
		keys:     keys,
		registry: loadRegistry(),
		start:    time.Now(),
	}
}

// Poll handles the received bytes and a pending reboot; call it regularly.
func (this *TerminalController) Poll() {
	now := uint32(time.Since(this.start).Milliseconds())
	for {
		b, ok := this.port.Read()
		if !ok {
			break
		}

		this.feed(b, now)
	}

	if !this.rebootAt.IsZero() && time.Now().After(this.rebootAt) {
		this.system.Reset()
	}
}

// reboot resets the cube shortly after the reply is sent, unless an upload starts first.
func (this *TerminalController) reboot() {
	this.println("System Reboot!")
	this.rebootAt = time.Now().Add(rebootDelay)
}

func (this *TerminalController) println(s string) { this.port.WriteString(s + "\r\n") }

// feed processes one received byte.
func (this *TerminalController) feed(b byte, nowMs uint32) {
	if this.started && nowMs-this.last > uart.Timeout {
		this.state = streamIdle
	}
	this.started, this.last = true, nowMs

	switch this.state {
	case streamHash:
		if b == '?' {
			this.state, this.line = streamCommand, this.line[:0]
			return
		}
		this.state = streamIdle
	case streamCommand:
		if b == '?' {
			this.state = streamQuestion
		} else if len(this.line) < uart.MaxCommandLine {
			this.line = append(this.line, b)
		} else {
			this.state = streamIdle
		}
		return
	case streamQuestion:
		this.state = streamIdle
		if b == '!' && len(this.line) > 0 {
			this.run(string(this.line))
		}
		return
	case streamFrame05:
		this.state = streamIdle
		if b == frameKind {
			this.state, this.check = streamFrameLen, uart.FrameStart^b
			return
		}
	case streamFrameLen:
		if b == 0 || b > uart.MaxPayload {
			this.state = streamIdle
			return
		}
		this.state, this.frameLeft, this.check, this.line = streamFrameData, int(b), this.check^b, this.line[:0]
		return
	case streamFrameData:
		this.line = append(this.line, b)
		this.check ^= b
		if this.frameLeft--; this.frameLeft == 0 {
			this.state = streamFrameCheck
		}
		return
	case streamFrameCheck:
		this.state, this.received = streamFrameEnd, b
		return
	case streamFrameEnd:
		this.state = streamIdle
		if b == uart.FrameEnd && this.received == this.check^uart.FrameEnd {
			this.uploadFrame(this.line)
		}
		return
	}

	switch b {
	case '#':
		this.state = streamHash
	case uart.FrameStart:
		this.state = streamFrame05
	}
}
