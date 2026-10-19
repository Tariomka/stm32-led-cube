package show

import (
	"bytes"
	"testing"

	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe"
	"github.com/Tariomka/stm32-led-cube/internal/engine/player"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
	"github.com/Tariomka/stm32-led-cube/internal/show/original"
)

func newCanvas() *render.Canvas {
	return &render.Canvas{Displays: []render.Display{render.NewDisplay(0, 0, 0, 8, 8, 8)}}
}

// TestMatchesOriginal renders every frame of every Go cartoon and of the original show data.
func TestMatchesOriginal(t *testing.T) {
	b, err := blob.Parse(original.Blob)
	if err != nil {
		t.Fatal(err)
	}
	SetMissing(func(kind, name string) { t.Errorf("missing %s %q", kind, name) })
	for group, cartoons := range [][]func() *keyframe.Cartoon{Dynamic, Audio, Test} {
		want := player.NewBlobShow(b, group)
		got := New(cartoons)
		if got.Count() != want.Count() {
			t.Fatalf("group %d: %d cartoons, want %d", group, got.Count(), want.Count())
		}
		for key := 0; key < 7; key++ {
			if got.KeyFunction(key) != want.KeyFunction(key) {
				t.Errorf("KEY%d function = %d, want %d", key+1, got.KeyFunction(key), want.KeyFunction(key))
			}
		}
		for i := 0; i < want.Count(); i++ {
			wc, gc := want.Cartoon(i), got.Cartoon(i)
			if gc.Kind() != wc.Kind() || gc.Header() != wc.Header() {
				t.Errorf("group %d cartoon %d: kind/header %d %+v, want %d %+v", group, i, gc.Kind(), gc.Header(), wc.Kind(), wc.Header())
				continue
			}
			frames := wc.Header().FrameCount
			if wc.Kind() == blob.Audio {
				frames = 64
			}
			wantCanvas, gotCanvas := newCanvas(), newCanvas()
			for frame := 0; frame < frames; frame++ {
				wantCanvas.Clear(0)
				gotCanvas.Clear(0)
				wc.Render(frame, wantCanvas)
				gc.Render(frame, gotCanvas)
				if !bytes.Equal(gotCanvas.Displays[0].Voxels, wantCanvas.Displays[0].Voxels) {
					t.Errorf("group %d cartoon %d frame %d differs", group, i, frame)
					break
				}
			}
		}
	}
}
