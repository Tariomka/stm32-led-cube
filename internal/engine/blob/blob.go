// Package blob reads the show data ("CartoonDatas") produced by the iCubeSmart PC software.
// The data is kept as a string so the embedded original stays in flash on the microcontroller.
// Format: .agent/re/spec-blob-format.md
package blob

import "errors"

const (
	Magic      = 0x12348765
	headerSize = 68
	nameSize   = 16

	PaletteSize = 16
)

var (
	ErrMagic  = errors.New("blob: bad magic")
	ErrLayout = errors.New("blob: section outside data")
)

const (
	secConfig = iota
	secPalette
	secGraph
	secIndex
	secGroup1
	GroupCount = 4
)

// Cartoon types
const (
	Normal = 0
	Script = 1
	Audio  = 2
)

// Frame element kinds
const (
	SetGraph   = 0
	DrawBall   = 1
	TitleStyle = 2
)

// Prop sizes per element kind, including the start/end frame header
var propSize = [3]int{56, 20, 56}

type section struct{ off, len int }

type Blob struct {
	data string
	sec  [8]section
}

func U32(s string, off int) uint32 {
	return uint32(s[off]) | uint32(s[off+1])<<8 | uint32(s[off+2])<<16 | uint32(s[off+3])<<24
}

func I16(s string, off int) int {
	return int(int16(uint16(s[off]) | uint16(s[off+1])<<8))
}

// Name returns a NUL-padded 16-byte name field without allocating.
func Name(s string, off int) string {
	end := off
	for end < off+nameSize && s[end] != 0 {
		end++
	}
	return s[off:end]
}

func Parse(data string) (*Blob, error) {
	if len(data) < headerSize || U32(data, 0) != Magic {
		return nil, ErrMagic
	}
	b := &Blob{data: data}
	for i := range b.sec {
		off, n := int(U32(data, 4+8*i)), int(U32(data, 8+8*i))
		if off < headerSize || off+n > len(data) {
			return nil, ErrLayout
		}
		b.sec[i] = section{off, n}
	}
	return b, nil
}

// Size is the used length of the blob in bytes.
func (b *Blob) Size() int {
	last := b.sec[len(b.sec)-1]
	return last.off + last.len
}

// KeyFunction returns the configured action for KEY1..KEY7 (key 0..6).
func (b *Blob) KeyFunction(key int) int {
	if key < 0 || key > 6 {
		return 0
	}
	return int(b.data[b.sec[secConfig].off+key])
}

type Display struct {
	X, Y, Z  int
	Dir      int
	CubeType int // 1=8x8x8, 2=8x16x8, 3=16x16x16, 4=8x32x8
}

func (d Display) Size() (l, h, w int) {
	switch d.CubeType {
	case 2:
		return 8, 16, 8
	case 3:
		return 16, 16, 16
	case 4:
		return 8, 32, 8
	}
	return 8, 8, 8
}

func (b *Blob) DisplayCount() int {
	return int(b.data[b.sec[secConfig].off+7])
}

func (b *Blob) Display(i int) Display {
	o := b.sec[secConfig].off + 8 + 8*i
	return Display{
		X: I16(b.data, o), Y: I16(b.data, o+2), Z: I16(b.data, o+4),
		Dir: int(b.data[o+6]), CubeType: int(b.data[o+7]),
	}
}

// Palette returns the 16 color bytes of a named palette.
func (b *Blob) Palette(name string) (string, bool) {
	s := b.sec[secPalette]
	for o := s.off; o+32 <= s.off+s.len; o += 32 {
		if Name(b.data, o) == name {
			return b.data[o+nameSize : o+32], true
		}
	}
	return "", false
}

type Graph struct {
	Name        string
	L, H, W     int
	Compression int
	Data        string
	offset      int
}

// ID is unique per graph within a blob, usable as a cache key.
func (g Graph) ID() int { return g.offset }

func (g Graph) Voxels() int { return g.L * g.H * g.W }

func (b *Blob) graphAt(o int) Graph {
	dl := int(U32(b.data, o+20))
	return Graph{
		Name: Name(b.data, o),
		L:    int(b.data[o+24]), H: int(b.data[o+25]), W: int(b.data[o+26]),
		Compression: int(b.data[o+27]),
		Data:        b.data[o+28 : o+28+dl],
		offset:      o,
	}
}

func (b *Blob) Graph(name string) (Graph, bool) {
	s := b.sec[secGraph]
	for o := s.off; o+28 <= s.off+s.len; {
		size := int(U32(b.data, o+16))
		if size <= 0 {
			break
		}
		if Name(b.data, o) == name {
			return b.graphAt(o), true
		}
		o += size
	}
	return Graph{}, false
}

// ForEachGraph visits all graphs in blob order.
func (b *Blob) ForEachGraph(fn func(Graph)) {
	s := b.sec[secGraph]
	for o := s.off; o+28 <= s.off+s.len; {
		size := int(U32(b.data, o+16))
		if size <= 0 {
			return
		}
		fn(b.graphAt(o))
		o += size
	}
}

// Index returns the raw data of a named index table.
func (b *Blob) Index(name string) (string, bool) {
	s := b.sec[secIndex]
	for o := s.off; o+20 <= s.off+s.len; {
		n := int(U32(b.data, o+16))
		if Name(b.data, o) == name {
			return b.data[o+20 : o+20+n], true
		}
		o += 20 + (n+3)&^3
	}
	return "", false
}

// IndexValue reads element i of an index table as the given type
// (0=u32, 1=i16, 2=u16, 3=i8, 4=u8), like the script GetIndex function.
func IndexValue(index string, i, typ int) int {
	if i < 0 {
		return 0
	}
	switch typ {
	case 0:
		if i < len(index)/4 {
			return int(int32(U32(index, 4*i)))
		}
	case 1:
		if i < len(index)/2 {
			return I16(index, 2*i)
		}
	case 2:
		if i < len(index)/2 {
			return int(uint16(index[2*i]) | uint16(index[2*i+1])<<8)
		}
	case 3:
		if i < len(index) {
			return int(int8(index[i]))
		}
	case 4:
		if i < len(index) {
			return int(index[i])
		}
	}
	return 0
}

// Cartoon is one show of a group.
type Cartoon struct {
	Type int
	Body string
}

// Group returns all cartoons of group g (0..3). Allocates once per call.
func (b *Blob) Group(g int) []Cartoon {
	if g < 0 || g >= GroupCount {
		return nil
	}
	s := b.sec[secGroup1+g]
	var list []Cartoon
	for o := s.off; o+8 <= s.off+s.len; {
		size := int(U32(b.data, o))
		if size < 8 || o+size > s.off+s.len {
			break
		}
		list = append(list, Cartoon{Type: int(b.data[o+4]), Body: b.data[o+8 : o+size]})
		o += size
	}
	return list
}

// Frame header of Normal and Audio cartoons.
type FrameHeader struct {
	FrameCount int // Normal: frames; Audio: duration in seconds (0 = forever)
	Loop       int // Normal: repeat count; Audio: audio type
	Delay      int // ms per frame before speed scaling
}

func (c Cartoon) Header() FrameHeader {
	if len(c.Body) < 8 {
		return FrameHeader{}
	}
	return FrameHeader{
		FrameCount: int(int32(U32(c.Body, 0))),
		Loop:       I16(c.Body, 4),
		Delay:      int(uint16(c.Body[6]) | uint16(c.Body[7])<<8),
	}
}

// Prop is one keyframed property of an element; Data starts after the start/end frame fields.
type Prop struct {
	Kind       int
	Element    int // offset of the owning element within the cartoon body
	Start, End int
	Data       string
}

// ForEachProp visits every property of every element of a Normal/Audio cartoon.
func (c Cartoon) ForEachProp(fn func(Prop)) {
	for e := 8; e+8 <= len(c.Body); {
		size := int(U32(c.Body, e))
		kind := int(U32(c.Body, e+4))
		if size < 8 || e+size > len(c.Body) {
			return
		}
		if kind >= 0 && kind < len(propSize) {
			ps := propSize[kind]
			for p := e + 8; p+ps <= e+size; p += ps {
				fn(Prop{Kind: kind, Element: e, Start: I16(c.Body, p), End: I16(c.Body, p+2), Data: c.Body[p+4 : p+ps]})
			}
		}
		e += size
	}
}
