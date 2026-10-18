package console

import (
	"strings"
	"testing"

	"github.com/Tariomka/stm32-led-cube/internal/registry"
)

type fakeDevice struct {
	rebooted bool
	saves    int
	cubeType [8]byte
	serial   [8]byte
}

func (d *fakeDevice) Reboot() { d.rebooted = true }

func (d *fakeDevice) ReadMemory(address uint32) (byte, bool) {
	return byte(address), address < 0x100
}

func (d *fakeDevice) SetType(t [8]byte) error   { d.cubeType = t; return nil }
func (d *fakeDevice) SetSerial(s [8]byte) error { d.serial = s; return nil }

func (d *fakeDevice) SaveRegistry(*registry.Registry) error { d.saves++; return nil }

func newConsole() (*Console, *fakeDevice, *strings.Builder) {
	dev := &fakeDevice{}
	out := &strings.Builder{}
	return &Console{Device: dev, Registry: &registry.Registry{}, Write: func(s string) { out.WriteString(s) }}, dev, out
}

func TestRegistryCommands(t *testing.T) {
	c, dev, out := newConsole()
	c.Run("setreg name cube")
	c.Run("GetReg name")
	c.Run("RegList")
	c.Run("DelReg name")
	c.Run("GetReg name")

	want := "DoWriteRegedit:-----------------\r\nWriteReg:name OK!\r\n" +
		"DoReadRegedit:-----------------\r\nReadReg:name=cube\r\n" +
		"DoGetRegeditList:-----------------\r\nRegedit Count = 1\r\nname=cube\r\n" +
		"DoDeleteRegedit:-----------------\r\nDelete Regedit:name OK!\r\n" +
		"DoReadRegedit:-----------------\r\nNot Found name\r\n"
	if out.String() != want {
		t.Errorf("output:\n%q\nwant:\n%q", out.String(), want)
	}
	if dev.saves != 2 {
		t.Errorf("saves = %d, want 2", dev.saves)
	}
}

func TestReadMem(t *testing.T) {
	c, _, out := newConsole()
	c.Run("ReadMem 0x10 2 18")
	want := "DoReadMemery:-----------------\r\nAddress=0x00000010,offset=2,count=18\r\n" +
		"12 13 14 15 16 17 18 19 1A 1B 1C 1D 1E 1F 20 21 \r\n22 23 \r\nReadMemery End!\r\n"
	if out.String() != want {
		t.Errorf("output:\n%q\nwant:\n%q", out.String(), want)
	}
}

func TestSetSerialReboots(t *testing.T) {
	c, dev, _ := newConsole()
	c.Run("SetSerial 0123456789abcdef")
	if dev.serial != [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF} || !dev.rebooted {
		t.Errorf("serial = %X, rebooted = %v", dev.serial, dev.rebooted)
	}
}

func TestDemux(t *testing.T) {
	var lines []string
	d := Demux{Command: func(line string) { lines = append(lines, line) }}
	var passed []byte
	feed := func(data []byte, now uint32) {
		for _, b := range data {
			if d.Feed(b, now) {
				passed = append(passed, b)
			}
		}
	}

	frame := []byte{0xAA, 1, 0, 3, '#', '?', 'x', 1, 2, 3, 4, 0x55}
	feed([]byte("SGFuZHNoYWtl"), 0)
	feed([]byte("#?reboot?!"), 10)
	feed(frame, 20)
	feed([]byte("#?help"), 30)
	feed([]byte("?!"), 2000) // too late: dropped

	if len(lines) != 1 || lines[0] != "reboot" {
		t.Errorf("lines = %q", lines)
	}
	if want := append([]byte("SGFuZHNoYWtl"), frame...); string(passed) != string(want)+"?!" {
		t.Errorf("passed = %q", passed)
	}
}
