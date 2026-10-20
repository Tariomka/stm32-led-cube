//go:build windows && 386

// Compares the Go show engine with the vendor's Simulator.dll (32-bit, from the iCubeSmart PC
// software) frame by frame. Run with: go run . oracle
package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"unsafe"

	"github.com/Tariomka/stm32-led-cube/internal/lightshow"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/blob"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/builtin/original"
	"github.com/Tariomka/stm32-led-cube/internal/lightshow/render"
)

const (
	dllPath = "Simulator.dll"
	// RVA of the player's speed factor (float32). Patched to make frame waits zero.
	speedRVA = 0x1c7b8
	cubeSize = 8 * 8 * 8
)

type simulator struct {
	dll                                          *syscall.DLL
	clean, add, get, setCartoon, run, end, doKey *syscall.Proc
}

func load(t *testing.T) *simulator {
	t.Helper()
	path, _ := filepath.Abs(dllPath)
	if _, err := os.Stat(path); err != nil {
		t.Skip("Simulator.dll not available: ", err)
	}
	dll, err := syscall.LoadDLL(path)
	if err != nil {
		t.Fatal(err)
	}
	s := &simulator{dll: dll}
	s.clean = dll.MustFindProc("CleanCube")
	s.add = dll.MustFindProc("AddCube")
	s.get = dll.MustFindProc("GetCube")
	s.setCartoon = dll.MustFindProc("SetCartoon")
	s.run = dll.MustFindProc("RunCartoon")
	s.end = dll.MustFindProc("EndCartoon")
	s.doKey = dll.MustFindProc("DoKeyFunction")
	*(*float32)(unsafe.Pointer(uintptr(dll.Handle) + speedRVA)) = 1e6
	return s
}

// record plays a group in the DLL and returns the first n frames.
// keys maps a frame number to a key (0..6) pressed right after that frame was drawn.
func (s *simulator) record(data []byte, group, n int, keys map[int]int) [][]byte {
	s.clean.Call()
	s.add.Call(0, 0, 0, 8, 8, 8)
	var frames [][]byte
	update := syscall.NewCallbackCDecl(func() uintptr {
		buf := make([]byte, cubeSize)
		s.get.Call(0, uintptr(unsafe.Pointer(&buf[0])), cubeSize)
		frames = append(frames, buf)
		if key, ok := keys[len(frames)-1]; ok {
			s.doKey.Call(uintptr(key))
		}
		if len(frames) >= n {
			s.end.Call()
		}
		return 0
	})
	s.setCartoon.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(group), update, 0)
	s.run.Call()
	return frames
}

func play(t *testing.T, group, n int, keys map[int]int) [][]byte {
	b, err := blob.Parse(original.Blob)
	if err != nil {
		t.Fatal(err)
	}
	canvas := &render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}}
	p := lightshow.New(lightshow.NewBlobShow(b, group), canvas)
	var frames [][]byte
	for len(frames) < n {
		rendered, _ := p.Step(0)
		if !rendered {
			t.Fatalf("player stalled after %d frames", len(frames))
		}
		frames = append(frames, append([]byte(nil), canvas.Displays[0].Voxels...))
		if key, ok := keys[len(frames)-1]; ok {
			p.Key(key)
		}
	}
	return frames
}

func groupFrames(group int) int {
	b, _ := blob.Parse(original.Blob)
	total := 0
	for _, c := range b.Group(group) {
		h := c.Header()
		if c.Type == blob.Normal {
			total += h.FrameCount * int(uint16(h.Loop))
		}
	}
	return total
}

func compare(t *testing.T, group, n int, keys map[int]int) {
	s := load(t)
	want := s.record([]byte(original.Blob), group, n, keys)
	got := play(t, group, n, keys)
	if len(want) < n {
		t.Fatalf("simulator produced %d frames, want %d", len(want), n)
	}
	bad, lit := 0, 0
	for i := 0; i < n; i++ {
		for _, v := range want[i] {
			if v != 0 {
				lit++
				break
			}
		}
		if string(want[i]) != string(got[i]) {
			if bad < 5 {
				t.Errorf("group %d frame %d differs:\n%s", group+1, i, diff(want[i], got[i]))
			}
			bad++
		}
	}
	t.Logf("group %d: compared %d frames, %d non-empty", group+1, n, lit)
	if lit < n/2 {
		t.Errorf("group %d: too many empty frames, the comparison is not meaningful", group+1)
	}
	if bad > 0 {
		t.Errorf("group %d: %d of %d frames differ", group+1, bad, n)
	}
}

func diff(want, got []byte) string {
	out := ""
	for i := range want {
		if want[i] != got[i] {
			out += fmt.Sprintf("  x=%d y=%d z=%d: dll=%02x go=%02x\n", i%8, i/8%8, i/64, want[i], got[i])
			if len(out) > 600 {
				return out + "  ...\n"
			}
		}
	}
	return out
}

func TestDynamicGroup(t *testing.T) {
	compare(t, 0, groupFrames(0)+20, nil)
}

func TestAudioGroup(t *testing.T) {
	compare(t, 1, 200, nil)
}

func TestTestGroup(t *testing.T) {
	compare(t, 2, groupFrames(2)+5, nil)
}

func TestTextGroup(t *testing.T) {
	compare(t, 3, groupFrames(3)+5, nil)
}

// Keys: KEY2 next, KEY1 previous, KEY6 play model (repeat one), KEY2 next again.
func TestKeys(t *testing.T) {
	compare(t, 0, 400, map[int]int{10: 1, 30: 1, 60: 0, 90: 5, 250: 1})
}
