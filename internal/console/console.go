// Package console implements the original firmware's UART console commands (Help, ReBoot,
// ReadMem, SetReg, GetReg, RegList, DelReg, ClrReg, SetType, SetSerial, Test) with the same
// replies, see .agent/re/firmware-map.md, plus show control commands (Key, Action, Cartoon,
// Status) for the Pico and PC tools.
package console

import (
	"strconv"
	"strings"

	"github.com/Tariomka/stm32-led-cube/internal/registry"
)

// Device is the hardware behind the commands.
type Device interface {
	Reboot()
	ReadMemory(address uint32) (byte, bool)
	SetType(cubeType [8]byte) error
	SetSerial(serial [8]byte) error
	SaveRegistry(r *registry.Registry) error
}

// Remote is the show player behind the control commands.
type Remote interface {
	Key(key int)             // presses KEY1..KEY7 (key 0..6)
	Action(action int)       // runs a key function code (player.Action*)
	Select(cartoon int) bool // restarts the show on a cartoon, false if out of range
	Status() string          // "Mode=...,Cartoon=..." line
}

type command struct {
	name, help string
	run        func(c *Console, args []string)
}

var commands []command

func init() {
	commands = []command{
		{"Help", "Commands explain", (*Console).help},
		{"ReBoot", ";System reset", (*Console).reboot},
		{"ReadMem", "ReadMem [address] [offset] [count] ;Read memory", (*Console).readMem},
		{"SetReg", "SetReg [key] [value] ;Write registry", (*Console).setReg},
		{"GetReg", "GetReg [key] ;Read registry", (*Console).getReg},
		{"RegList", ";List all registry entries", (*Console).regList},
		{"DelReg", "DelReg [key] ;Delete one entry", (*Console).delReg},
		{"ClrReg", ";Delete all registry entries", (*Console).clrReg},
		{"SetType", "SetType [type] ;Set the cube type", (*Console).setType},
		{"SetSerial", "SetSerial [serial] ;Set the cube serial number", (*Console).setSerial},
		{"Test", "Test", (*Console).test},
		{"Key", "Key [1-7] ;Press a key", (*Console).key},
		{"Action", "Action [code] ;Run a key function", (*Console).action},
		{"Cartoon", "Cartoon [index] ;Play a cartoon", (*Console).cartoon},
		{"Status", ";Show the player state", (*Console).status},
	}
}

type Console struct {
	Device   Device
	Registry *registry.Registry
	// Remote runs the show control commands; nil until the shows are running.
	Remote Remote
	// Write sends raw text to the UART.
	Write func(s string)
}

func (this *Console) println(s string) { this.Write(s + "\r\n") }

func (this *Console) section(name string) { this.println(name + ":-----------------") }

// Run executes one command line ("name arg1 arg2 ...", name is case insensitive).
func (this *Console) Run(line string) {
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
		if equalFold(cmd.name, args[0]) {
			cmd.run(this, args[1:])
			return
		}
	}
	this.println(helpLine(commands[0]))
}

// equalFold compares ASCII case-insensitively (strings.EqualFold pulls in the unicode tables).
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if 'A' <= x && x <= 'Z' {
			x += 'a' - 'A'
		}
		if 'A' <= y && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}

func helpLine(cmd command) string {
	if len(cmd.name) < 8 {
		return cmd.name + "\t\t:" + cmd.help
	}
	return cmd.name + "\t:" + cmd.help
}

func (this *Console) help(args []string) {
	this.println("Commands explain:-----------------")
	for _, cmd := range commands {
		this.println(helpLine(cmd))
	}
}

func (this *Console) reboot(args []string) {
	this.println("System Reboot!")
	this.Device.Reboot()
}

// parseNumber reads "0x..." as hex and anything else as decimal, like the original.
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

func (this *Console) readMem(args []string) {
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
		b, ok := this.Device.ReadMemory(address + offset + i)
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

func (this *Console) save() {
	if err := this.Device.SaveRegistry(this.Registry); err != nil {
		this.println("Save Regedit Error!")
	}
}

func (this *Console) setReg(args []string) {
	this.section("DoWriteRegedit")
	if len(args) != 2 {
		this.println("WriteReg Param Number Error!")
		return
	}
	if this.Registry.Set(args[0], args[1]) != nil {
		this.println("WriteReg Error!")
		return
	}
	this.save()
	this.println("WriteReg:" + args[0] + " OK!")
}

func (this *Console) getReg(args []string) {
	this.section("DoReadRegedit")
	if len(args) != 1 {
		this.println("ReadReg Param Number Error!")
		return
	}
	if value, ok := this.Registry.Get(args[0]); ok {
		this.println("ReadReg:" + args[0] + "=" + value)
		return
	}
	this.println("Not Found " + args[0])
}

func (this *Console) regList(args []string) {
	this.section("DoGetRegeditList")
	entries := this.Registry.Entries()
	this.println("Regedit Count = " + strconv.Itoa(len(entries)))
	for _, e := range entries {
		this.println(e.Key + "=" + e.Value)
	}
}

func (this *Console) delReg(args []string) {
	this.section("DoDeleteRegedit")
	if len(args) != 1 {
		this.println("DelReg Param number is error!")
		return
	}
	if this.Registry.Delete(args[0]) {
		this.save()
	}
	this.println("Delete Regedit:" + args[0] + " OK!")
}

func (this *Console) clrReg(args []string) {
	this.section("DoCleanRegeditList")
	if len(this.Registry.Entries()) != 0 {
		this.Registry.Clear()
		this.save()
	}
	this.println("Clean Regedit OK!")
}

func (this *Console) setType(args []string) {
	this.println("DoSetType")
	if len(args) != 1 {
		this.println("Argc number Error!")
		return
	}
	var cubeType [8]byte
	copy(cubeType[:], args[0])
	if this.Device.SetType(cubeType) != nil {
		this.println("SetType Error!")
		return
	}
	this.reboot(nil)
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

func (this *Console) setSerial(args []string) {
	this.println("DoSetSerial")
	if len(args) != 1 || len(args[0]) != 16 {
		this.println("Argc number Error!")
		return
	}
	var serial [8]byte
	for i := range serial {
		serial[i] = hexDigit(args[0][2*i])<<4 | hexDigit(args[0][2*i+1])
	}
	if this.Device.SetSerial(serial) != nil {
		this.println("SetSerial Error!")
		return
	}
	this.reboot(nil)
}

func (this *Console) test(args []string) {
	this.println("DoTest")
}

// remoteArg parses the single numeric argument of a control command within [low, high].
func (this *Console) remoteArg(args []string, low, high int) (int, bool) {
	if this.Remote == nil {
		this.println("Not Ready!")
		return 0, false
	}
	if len(args) == 1 {
		if v, ok := parseNumber(args[0]); ok && int(v) >= low && int(v) <= high {
			return int(v), true
		}
	}
	this.println("Param Error!")
	return 0, false
}

func (this *Console) key(args []string) {
	if key, ok := this.remoteArg(args, 1, 7); ok {
		this.Remote.Key(key - 1)
	}
}

func (this *Console) action(args []string) {
	if action, ok := this.remoteArg(args, 1, 19); ok {
		this.Remote.Action(action)
	}
}

func (this *Console) cartoon(args []string) {
	if index, ok := this.remoteArg(args, 0, 255); ok && !this.Remote.Select(index) {
		this.println("Param Error!")
	}
}

func (this *Console) status(args []string) {
	if this.Remote == nil {
		this.println("Not Ready!")
		return
	}
	this.println("Status:" + this.Remote.Status())
}
