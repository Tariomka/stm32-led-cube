// showgen turns the original show data (internal/show/original/original.bin) into the Go light
// shows of package internal/show: one function per cartoon plus the shapes and palettes they use.
// The original Text group (one scrolling-text cartoon) is not converted.
//
// Usage (from the repo root): go run ./tools/showgen [output dir, default internal/show]
package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/show/original"
)

type groupInfo struct {
	index      int
	file, list string
	label      string
}

var groups = []groupInfo{
	{0, "dynamic.go", "dynamicCartoons", "Dynamic"},
	{1, "audio.go", "audioCartoons", "Audio"},
	{2, "test.go", "testCartoons", "Test"},
}

const header = "// Generated from the original iCubeSmart show data by tools/showgen.\n\npackage show\n\n"

func main() {
	out := "internal/show"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	b, err := blob.Parse(original.Blob)
	if err != nil {
		fail(err)
	}

	shapes, palettes := map[string]bool{}, map[string]bool{}
	for _, g := range groups {
		var body strings.Builder
		var names []string
		for i, c := range b.Group(g.index) {
			name := cartoonName(g, i, c)
			names = append(names, name)
			writeCartoon(&body, g, i, name, c, shapes, palettes)
		}
		fmt.Fprintf(&body, "var %s = []func() *keyframe.Cartoon{\n", g.list)
		for _, name := range names {
			fmt.Fprintf(&body, "\t%s,\n", name)
		}
		body.WriteString("}\n")

		imports := "import \"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe\"\n\n"
		if strings.Contains(body.String(), "render.") {
			imports = "import (\n\t\"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe\"\n\t\"github.com/Tariomka/stm32-led-cube/internal/engine/render\"\n)\n\n"
		}
		write(filepath.Join(out, g.file), header+imports+body.String())
	}

	var src strings.Builder
	src.WriteString(header)
	src.WriteString("import \"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe\"\n\n")
	src.WriteString("// Keys are the original key functions of KEY1..KEY7 (player.Action* codes).\n")
	src.WriteString("var Keys = [7]int{")
	for key := 0; key < 7; key++ {
		if key > 0 {
			src.WriteString(", ")
		}
		fmt.Fprintf(&src, "%d", b.KeyFunction(key))
	}
	src.WriteString("}\n\n")
	src.WriteString("var palettes = []keyframe.Palette{\n")
	for _, name := range sorted(palettes) {
		colors, _ := b.Palette(name)
		fmt.Fprintf(&src, "\t{Name: %q, Colors: \"%s\"},\n", name, hexEscape(colors))
	}
	src.WriteString("}\n\n")
	src.WriteString("// Shape art: one line per layer from the top, z rows separated by spaces, x left to right.\n")
	src.WriteString("var shapes = []keyframe.Shape{\n")
	for _, name := range sorted(shapes) {
		g, _ := b.Graph(name)
		fmt.Fprintf(&src, "\t{Name: %q, L: %d, H: %d, W: %d, Art: `\n%s\t`},\n", name, g.Length, g.Height, g.Width, art(g))
	}
	src.WriteString("}\n")
	write(filepath.Join(out, "library.go"), src.String())
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "showgen:", err)
	os.Exit(1)
}

func write(path, src string) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		fail(fmt.Errorf("%s: %w", path, err))
	}
	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		fail(err)
	}
	fmt.Println("wrote", path)
}

func sorted(set map[string]bool) []string {
	var list []string
	for name := range set {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func hexEscape(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(&sb, "\\x%02X", s[i])
	}
	return sb.String()
}

// art writes the shape as text, see keyframe.Shape.
func art(g blob.Graph) string {
	data := make([]byte, g.Voxels())
	if err := g.Decode(data); err != nil {
		fail(fmt.Errorf("graph %s: %w", g.Name, err))
	}
	const digits = ".123456789ABCDEF"
	var sb strings.Builder
	for y := g.Height - 1; y >= 0; y-- {
		sb.WriteString("\t\t")
		for z := 0; z < g.Width; z++ {
			if z > 0 {
				sb.WriteByte(' ')
			}
			for x := 0; x < g.Length; x++ {
				sb.WriteByte(digits[data[x+y*g.Length+z*g.Length*g.Height]&0xF])
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type element struct {
	props []blob.Prop
}

func elements(c blob.Cartoon) []element {
	var list []element
	last := -1
	c.ForEachProp(func(p blob.Prop) {
		if p.Element != last {
			list = append(list, element{})
			last = p.Element
		}
		list[len(list)-1].props = append(list[len(list)-1].props, p)
	})
	return list
}

func field(p blob.Prop, i int) int { return blob.I16(p.Data, 32+2*i) }

// usedNames collects the shape and palette names a graph prop resolves to on any frame,
// following the engine's per-element frame shift.
func usedNames(c blob.Cartoon, shapes, palettes map[string]bool) {
	h := c.Header()
	frames := h.FrameCount
	if c.Type == blob.Audio {
		frames = 64
	}
	for _, e := range elements(c) {
		for frame := 0; frame < frames; frame++ {
			inner := frame
			for _, p := range e.props {
				if inner < p.Start || inner > p.End {
					continue
				}
				nameFrame := inner
				if inner != p.Start && inner != p.End {
					nameFrame = inner - p.Start
				}
				if p.Kind == blob.SetGraph {
					shapes[blob.ResolveName(blob.Name(p.Data, 0), nameFrame, p.Start)] = true
					if pal := blob.ResolveName(blob.Name(p.Data, 16), nameFrame, p.Start); pal != "" {
						palettes[pal] = true
					}
				}
				inner = nameFrame
			}
		}
	}
}

func cartoonName(g groupInfo, i int, c blob.Cartoon) string {
	label := "Balls"
	c.ForEachProp(func(p blob.Prop) {
		if label == "Balls" && p.Kind == blob.SetGraph {
			label = identifier(blob.Name(p.Data, 0))
		}
	})
	return fmt.Sprintf("%s%02d%s", strings.ToLower(g.label), i+1, label)
}

// identifier keeps the letters and digits of a shape name before any "$" expression.
func identifier(name string) string {
	if i := strings.IndexByte(name, '$'); i >= 0 {
		name = name[:i]
	}
	var sb strings.Builder
	upper := true
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9' && sb.Len() > 0:
			if upper && r >= 'a' && r <= 'z' {
				r -= 'a' - 'A'
			}
			sb.WriteRune(r)
			upper = false
		default:
			upper = true
		}
	}
	if sb.Len() == 0 {
		return "Shape"
	}
	return sb.String()
}

func pose(x, y, z, angle int) string {
	if x == 0 && y == 0 && z == 0 && angle == 0 {
		return ""
	}
	var parts []string
	for _, f := range []struct {
		name  string
		value int
	}{{"X", x}, {"Y", y}, {"Z", z}, {"Angle", angle}} {
		if f.value != 0 {
			parts = append(parts, fmt.Sprintf("%s: %d", f.name, f.value))
		}
	}
	return "keyframe.Pose{" + strings.Join(parts, ", ") + "}"
}

var axes = []string{"render.AxisX", "render.AxisY", "render.AxisZ"}

func graphProp(p blob.Prop) string {
	fields := []string{fmt.Sprintf("Start: %d, End: %d", p.Start, p.End), fmt.Sprintf("Shape: %q", blob.Name(p.Data, 0))}
	if pal := blob.Name(p.Data, 16); pal != "" {
		fields = append(fields, fmt.Sprintf("Palette: %q", pal))
	}
	if from := pose(field(p, 0), field(p, 2), field(p, 4), field(p, 6)); from != "" {
		fields = append(fields, "From: "+from)
	}
	if to := pose(field(p, 1), field(p, 3), field(p, 5), field(p, 7)); to != "" {
		fields = append(fields, "To: "+to)
	}
	if axis := int(uint16(field(p, 8))); axis != 0 {
		if axis < len(axes) {
			fields = append(fields, "Axis: "+axes[axis])
		} else {
			fields = append(fields, fmt.Sprintf("Axis: %d", axis))
		}
	}
	switch mode := int(uint16(field(p, 9))); mode {
	case 0:
	case 1:
		fields = append(fields, "Mode: render.ModeMix")
	default:
		fields = append(fields, fmt.Sprintf("Mode: %d", mode))
	}
	return "keyframe.Graph{" + strings.Join(fields, ", ") + "}"
}

func ballProp(p blob.Prop) string {
	d := p.Data
	fields := []string{
		fmt.Sprintf("Start: %d, End: %d", p.Start, p.End),
		fmt.Sprintf("X: %d, Y: %d, Z: %d", blob.I16(d, 0), blob.I16(d, 2), blob.I16(d, 4)),
		fmt.Sprintf("From: %d, To: %d", blob.I16(d, 6), blob.I16(d, 8)),
	}
	if uint16(blob.I16(d, 10)) != 0 {
		fields = append(fields, "Hollow: true")
	}
	fields = append(fields, fmt.Sprintf("Out: 0x%02X, In: 0x%02X", d[13], d[12]))
	return "keyframe.Ball{" + strings.Join(fields, ", ") + "}"
}

func writeCartoon(src *strings.Builder, g groupInfo, i int, name string, c blob.Cartoon, shapes, palettes map[string]bool) {
	if c.Type != blob.Normal && c.Type != blob.Audio {
		fail(fmt.Errorf("%s: cartoon type %d is not supported", name, c.Type))
	}
	usedNames(c, shapes, palettes)
	h := c.Header()
	fmt.Fprintf(src, "// %s is cartoon %d of the original %s group.\n", name, i+1, g.label)
	fmt.Fprintf(src, "func %s() *keyframe.Cartoon {\n\treturn &keyframe.Cartoon{\n", name)
	if c.Type == blob.Audio {
		fmt.Fprintf(src, "\t\tAudio: true, Frames: %d, Loop: %d, Delay: %d, // %d s, level mode %d\n", h.FrameCount, h.Loop, h.Delay, h.FrameCount, h.Loop)
	} else {
		fmt.Fprintf(src, "\t\tFrames: %d, Loop: %d, Delay: %d,\n", h.FrameCount, h.Loop, h.Delay)
	}
	src.WriteString("\t\tLibrary: library,\n\t\tElements: []keyframe.Element{\n")
	for _, e := range elements(c) {
		src.WriteString("\t\t\t{\n")
		for _, p := range e.props {
			switch p.Kind {
			case blob.SetGraph:
				fmt.Fprintf(src, "\t\t\t\t%s,\n", graphProp(p))
			case blob.DrawBall:
				fmt.Fprintf(src, "\t\t\t\t%s,\n", ballProp(p))
			default:
				fail(fmt.Errorf("%s: element kind %d is not supported", name, p.Kind))
			}
		}
		src.WriteString("\t\t\t},\n")
	}
	src.WriteString("\t\t},\n\t}\n}\n\n")
}
