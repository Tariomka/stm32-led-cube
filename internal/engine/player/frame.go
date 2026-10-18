package player

import (
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

// assets resolves resource names and keeps the most recently decoded graph.
type assets struct {
	blob     *blob.Blob
	buf      []byte
	cachedID int
	cached   render.Voxels
	Missing  func(kind, name string) // optional, for diagnostics
}

func newAssets(b *blob.Blob) *assets {
	return &assets{blob: b, cachedID: -1}
}

func (a *assets) graph(name string) (render.Voxels, bool) {
	g, ok := a.blob.Graph(name)
	if !ok {
		a.missing("graph", name)
		return render.Voxels{}, false
	}
	if g.ID() == a.cachedID {
		return a.cached, true
	}
	n := g.Voxels()
	if cap(a.buf) < n {
		a.buf = make([]byte, n)
	}
	if err := g.Decode(a.buf[:n]); err != nil {
		a.missing("graph data", name)
		return render.Voxels{}, false
	}
	a.cachedID = g.ID()
	a.cached = render.Voxels{L: g.L, H: g.H, W: g.W, Data: a.buf[:n]}
	return a.cached, true
}

// palette returns the default palette for an empty name.
func (a *assets) palette(name string) (string, bool) {
	if name == "" {
		return render.DefaultPalette, true
	}
	p, ok := a.blob.Palette(name)
	if !ok {
		a.missing("palette", name)
	}
	return p, ok
}

func (a *assets) index(name string) (string, bool) {
	idx, ok := a.blob.Index(name)
	if !ok {
		a.missing("index", name)
	}
	return idx, ok
}

func (a *assets) missing(kind, name string) {
	if a.Missing != nil {
		a.Missing(kind, name)
	}
}

func lerp(v0, v1, frame, start, end int) int {
	return v0 + (v1-v0)*(frame-start)/(end-start)
}

// renderFrame draws all element properties of a Normal/Audio cartoon that are active at frame.
func renderFrame(c blob.Cartoon, frame int, canvas *render.Canvas, a *assets) {
	elementStart := -1
	f := frame
	c.ForEachProp(func(p blob.Prop) {
		// The vendor engine reuses the frame variable across the properties of one element and
		// subtracts the property start from it while interpolating. Reproduced as-is.
		if p.Element != elementStart {
			elementStart, f = p.Element, frame
		}
		if f < p.Start || f > p.End {
			return
		}
		switch p.Kind {
		case blob.SetGraph:
			f = drawGraph(p, f, canvas, a)
		case blob.DrawBall:
			f = drawBall(p, f, canvas)
		case blob.TitleStyle:
			f = drawTitle(p, f, canvas, a)
		}
	})
}

func drawGraph(p blob.Prop, f int, canvas *render.Canvas, a *assets) int {
	d := p.Data
	field := func(i int) int { return blob.I16(d, 32+2*i) }
	var x, y, z, angle int
	switch f {
	case p.Start:
		x, y, z, angle = field(0), field(2), field(4), field(6)
	case p.End:
		x, y, z, angle = field(1), field(3), field(5), field(7)
	default:
		x = lerp(field(0), field(1), f, p.Start, p.End)
		y = lerp(field(2), field(3), f, p.Start, p.End)
		z = lerp(field(4), field(5), f, p.Start, p.End)
		angle = lerp(field(6), field(7), f, p.Start, p.End)
		f -= p.Start
	}
	axis, mode := int(uint16(field(8))), int(uint16(field(9)))
	graphName := resolveName(blob.Name(d, 0), f, p.Start)
	paletteName := resolveName(blob.Name(d, 16), f, p.Start)
	g, ok := a.graph(graphName)
	if !ok {
		return f
	}
	pal, ok := a.palette(paletteName)
	if !ok {
		return f
	}
	canvas.Graph(g, pal, x, y, z, axis, angle, mode)
	return f
}

func drawBall(p blob.Prop, f int, canvas *render.Canvas) int {
	d := p.Data
	r0, r1 := blob.I16(d, 6), blob.I16(d, 8)
	var r int
	switch f {
	case p.Start:
		r = r0
	case p.End:
		r = r1
	default:
		r = lerp(r0, r1, f, p.Start, p.End)
		f -= p.Start
	}
	hollow := uint16(blob.I16(d, 10)) != 0
	canvas.Ball(blob.I16(d, 0), blob.I16(d, 2), blob.I16(d, 4), r, d[13], d[12], hollow)
	return f
}

func drawTitle(p blob.Prop, f int, canvas *render.Canvas, a *assets) int {
	d := p.Data
	f -= p.Start
	graphName := resolveName(blob.Name(d, 0), f, p.Start)
	paletteName := resolveName(blob.Name(d, 16), f, p.Start)
	indexName := resolveName(blob.Name(d, 32), f, p.Start)
	g, ok := a.graph(graphName)
	if !ok {
		return f
	}
	path, ok := a.index(indexName)
	if !ok {
		return f
	}
	pal, ok := a.palette(paletteName)
	if !ok {
		return f
	}
	canvas.Title(g, pal, path, blob.I16(d, 48), int(uint16(blob.I16(d, 50))), f)
	return f
}
