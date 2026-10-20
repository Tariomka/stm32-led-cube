package controller

import (
	"strconv"
	"strings"

	"github.com/Tariomka/stm32-led-cube/internal/common"
)

// Console commands of the factory firmware (with its replies) plus the show control commands Key,
// Action, Cartoon and Status for the Pico and the PC tools.
type command struct {
	name, help string
	run        func(this *TerminalController, args []string)
}

var commands []command

func init() {
	commands = []command{
		{"Help", "Commands explain", (*TerminalController).help},
		{"ReBoot", ";System reset", (*TerminalController).rebootCommand},
		{"ReadMem", "ReadMem [address] [offset] [count] ;Read memory", (*TerminalController).readMem},
		{"SetReg", "SetReg [key] [value] ;Write registry", (*TerminalController).setReg},
		{"GetReg", "GetReg [key] ;Read registry", (*TerminalController).getReg},
		{"RegList", ";List all registry entries", (*TerminalController).regList},
		{"DelReg", "DelReg [key] ;Delete one entry", (*TerminalController).delReg},
		{"ClrReg", ";Delete all registry entries", (*TerminalController).clrReg},
		{"SetType", "SetType [type] ;Set the cube type", (*TerminalController).setType},
		{"SetSerial", "SetSerial [serial] ;Set the cube serial number", (*TerminalController).setSerial},
		{"Test", "Test", (*TerminalController).test},
		{"Key", "Key [1-7] ;Press a key", (*TerminalController).key},
		{"Action", "Action [code] ;Run a key function", (*TerminalController).action},
		{"Cartoon", "Cartoon [index] ;Play a cartoon", (*TerminalController).cartoon},
		{"Status", ";Show the player state", (*TerminalController).status},
	}
}

// run executes one command line ("name arg1 arg2 ...", name is case insensitive).
func (this *TerminalController) run(line string) {
	var args []string
	for _, arg := range strings.Split(line, " ") {
		if arg != "" {
			args = append(args, arg)
		}
	}
	if len(args) == 0 {
		return
	}

	for _, cmd := range commands {
		if common.EqualFold(cmd.name, args[0]) {
			cmd.run(this, args[1:])
			return
		}
	}
	this.println(helpLine(commands[0]))
}

func (this *TerminalController) section(name string) { this.println(name + ":-----------------") }

func helpLine(cmd command) string {
	if len(cmd.name) < 8 {
		return cmd.name + "\t\t:" + cmd.help
	}
	return cmd.name + "\t:" + cmd.help
}

func (this *TerminalController) help(args []string) {
	this.println("Commands explain:-----------------")
	for _, cmd := range commands {
		this.println(helpLine(cmd))
	}
}

func (this *TerminalController) rebootCommand(args []string) { this.reboot() }

// parseNumber reads "0x..." as hex and anything else as decimal, like the factory firmware.
func parseNumber(s string) (uint32, bool) {
	base := 10
	if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s, base = s[2:], 16
	}
	v, err := strconv.ParseUint(s, base, 32)
	return uint32(v), err == nil
}

func hex(v uint32, digits int) string {
	const hexDigits = "0123456789ABCDEF"
	buf := make([]byte, digits)
	for i := digits - 1; i >= 0; i-- {
		buf[i] = hexDigits[v&0xF]
		v >>= 4
	}
	return string(buf)
}

func hexDigit(b byte) byte {
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10
	}
	return 0
}

func (this *TerminalController) readMem(args []string) {
	this.section("DoReadMemery")
	if len(args) != 3 {
		this.println("Param Number Error!")
		return
	}
	address, ok1 := parseNumber(args[0])
	count, ok2 := parseNumber(args[2])
	offset, _ := parseNumber(args[1])
	if !ok1 || !ok2 {
		this.println("Param Number Error!")
		return
	}

	this.println("Address=0x" + hex(address, 8) + ",offset=" + strconv.Itoa(int(int32(offset))) + ",count=" + strconv.Itoa(int(count)))
	var row strings.Builder
	for i := uint32(0); i < count; i++ {
		if i%16 == 0 && i != 0 {
			this.println(row.String())
			row.Reset()
		}
		b, ok := this.system.ReadMemory(address + offset + i)
		if !ok {
			this.println(row.String())
			this.println("Address Error!")
			return
		}
		row.WriteString(hex(uint32(b), 2) + " ")
	}
	this.println(row.String())
	this.println("ReadMemery End!")
}

func (this *TerminalController) saveRegistry() {
	if err := this.registry.save(); err != nil {
		this.println("Save Regedit Error!")
	}
}

func (this *TerminalController) setReg(args []string) {
	this.section("DoWriteRegedit")
	if len(args) != 2 {
		this.println("WriteReg Param Number Error!")
		return
	}
	if this.registry.set(args[0], args[1]) != nil {
		this.println("WriteReg Error!")
		return
	}
	this.saveRegistry()
	this.println("WriteReg:" + args[0] + " OK!")
}

func (this *TerminalController) getReg(args []string) {
	this.section("DoReadRegedit")
	if len(args) != 1 {
		this.println("ReadReg Param Number Error!")
		return
	}
	if value, ok := this.registry.get(args[0]); ok {
		this.println("ReadReg:" + args[0] + "=" + value)
		return
	}
	this.println("Not Found " + args[0])
}

func (this *TerminalController) regList(args []string) {
	this.section("DoGetRegeditList")
	this.println("Regedit Count = " + strconv.Itoa(len(this.registry.entries)))
	for _, entry := range this.registry.entries {
		this.println(entry.key + "=" + entry.value)
	}
}

func (this *TerminalController) delReg(args []string) {
	this.section("DoDeleteRegedit")
	if len(args) != 1 {
		this.println("DelReg Param number is error!")
		return
	}
	if this.registry.delete(args[0]) {
		this.saveRegistry()
	}
	this.println("Delete Regedit:" + args[0] + " OK!")
}

func (this *TerminalController) clrReg(args []string) {
	this.section("DoCleanRegeditList")
	if len(this.registry.entries) != 0 {
		this.registry.entries = nil
		this.saveRegistry()
	}
	this.println("Clean Regedit OK!")
}

func (this *TerminalController) setType(args []string) {
	this.println("DoSetType")
	if len(args) != 1 {
		this.println("Argc number Error!")
		return
	}
	var cubeType [8]byte
	copy(cubeType[:], args[0])
	if this.settings.SetCubeType(cubeType) != nil {
		this.println("SetType Error!")
		return
	}
	this.reboot()
}

func (this *TerminalController) setSerial(args []string) {
	this.println("DoSetSerial")
	if len(args) != 1 || len(args[0]) != 16 {
		this.println("Argc number Error!")
		return
	}
	var serial [8]byte
	for i := range serial {
		serial[i] = hexDigit(args[0][2*i])<<4 | hexDigit(args[0][2*i+1])
	}
	if this.settings.SetSerial(serial) != nil {
		this.println("SetSerial Error!")
		return
	}
	this.reboot()
}

func (this *TerminalController) test(args []string) { this.println("DoTest") }

// controlArg parses the single numeric argument of a control command within [low, high].
func (this *TerminalController) controlArg(args []string, low, high int) (int, bool) {
	if len(args) == 1 {
		if v, ok := parseNumber(args[0]); ok && int(v) >= low && int(v) <= high {
			return int(v), true
		}
	}
	this.println("Param Error!")
	return 0, false
}

// key presses a key like the physical one, so cascaded cubes follow it.
func (this *TerminalController) key(args []string) {
	if key, ok := this.controlArg(args, 1, KeyCount); ok {
		this.keys.Press(key - 1)
	}
}

func (this *TerminalController) action(args []string) {
	if action, ok := this.controlArg(args, 1, 19); ok {
		this.show.Action(action)
	}
}

func (this *TerminalController) cartoon(args []string) {
	if index, ok := this.controlArg(args, 0, 255); ok && !this.show.Select(index) {
		this.println("Param Error!")
	}
}

func (this *TerminalController) status(args []string) { this.println("Status:" + this.show.Status()) }
