package blob

import (
	"bytes"
	"testing"

	"github.com/Tariomka/stm32-led-cube/internal/lightshow/builtin/original"
)

func mustParse(t *testing.T) *Blob {
	t.Helper()
	b, err := Parse(original.Blob)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestOriginalLayout(t *testing.T) {
	b := mustParse(t)
	if b.Size() != 32100 {
		t.Errorf("size = %d, want 32100", b.Size())
	}
	wantKeys := []int{2, 3, 4, 6, 8, 7, 9}
	for k, want := range wantKeys {
		if got := b.KeyFunction(k); got != want {
			t.Errorf("key %d = %d, want %d", k+1, got, want)
		}
	}
	if b.DisplayCount() != 1 || b.Display(0) != (Display{CubeType: 1}) {
		t.Errorf("displays = %d %+v", b.DisplayCount(), b.Display(0))
	}
	if p, ok := b.Palette("red-series"); !ok || p != "\x00\x30\x31\x32\x33\x3B\x3F\x00\x00\x00\x00\x00\x00\x00\x00\x00" {
		t.Errorf("red-series = %q %v", p, ok)
	}
	graphs := 0
	b.ForEachGraph(func(Graph) { graphs++ })
	if graphs != 79 {
		t.Errorf("graphs = %d, want 79", graphs)
	}
	if idx, ok := b.Index("TwoSide"); !ok || len(idx) != 60 || IndexValue(idx, 0, 1) != 7 {
		t.Errorf("TwoSide index = %d bytes, ok=%v", len(idx), ok)
	}
	for g, want := range []int{41, 10, 4, 1} {
		if got := len(b.Group(g)); got != want {
			t.Errorf("group %d has %d cartoons, want %d", g+1, got, want)
		}
	}
	if h := b.Group(3)[0].Header(); h != (FrameHeader{FrameCount: 102, Loop: 1, Delay: 150}) {
		t.Errorf("text cartoon header = %+v", h)
	}
}

// Re-encoding the decoded voxels with the PC software's encoder must reproduce the stored bytes.
func TestDecodeRoundTrip(t *testing.T) {
	b := mustParse(t)
	b.ForEachGraph(func(g Graph) {
		vox := make([]byte, g.Voxels())
		if err := g.Decode(vox); err != nil {
			t.Fatalf("%s: %v", g.Name, err)
		}
		data, comp := encode(vox)
		if comp != g.Compression || !bytes.Equal(data[:len(g.Data)], []byte(g.Data)) {
			t.Errorf("%s: round trip mismatch (comp %d vs %d)", g.Name, comp, g.Compression)
		}
	})
}

// encode is a port of _3D_Guan.Compile.Compress.Encode.
func encode(vox []byte) ([]byte, int) {
	var palette []byte
	seen := map[byte]bool{}
	for _, v := range vox {
		if !seen[v] {
			seen[v] = true
			palette = append(palette, v)
		}
	}
	for i := 1; i < len(palette); i++ {
		for j := i; j > 0 && palette[j] < palette[j-1]; j-- {
			palette[j], palette[j-1] = palette[j-1], palette[j]
		}
	}
	if len(palette) > 16 {
		out := make([]byte, (len(vox)+3)&^3)
		copy(out, vox)
		return out, NoCompress
	}
	bits, comp := 4, Normal2
	switch {
	case len(palette) <= 2:
		bits, comp = 1, Normal8
	case len(palette) <= 4:
		bits, comp = 2, Normal4
	}
	index := func(v byte) byte {
		for i, p := range palette {
			if p == v {
				return byte(i)
			}
		}
		return 0
	}
	perByte := 8 / bits
	c := append([]byte{byte(len(palette))}, palette...)
	add := func(color byte, count int) {
		if count <= 0 {
			return
		}
		if count < 1<<(7-bits) {
			c = append(c, byte(count<<bits)|color)
			return
		}
		c = append(c, byte(count>>(8-bits))|0x80, byte(count<<bits)|color)
	}
	last, count := vox[0], 1
	for _, v := range vox[1:] {
		if v == last {
			count++
			if count == 1<<(15-bits)-1 {
				add(index(last), count)
				count = 0
			}
		} else {
			add(index(last), count)
			last, count = v, 1
		}
	}
	add(index(last), count)
	if len(c) < len(vox)/perByte {
		return c, comp + 3
	}
	n := append([]byte{byte(len(palette))}, palette...)
	for i := 0; i < len(vox); i += perByte {
		var b byte
		for j := 0; j < perByte && i+j < len(vox); j++ {
			b |= index(vox[i+j]) << (j * bits)
		}
		n = append(n, b)
	}
	return n, comp
}
