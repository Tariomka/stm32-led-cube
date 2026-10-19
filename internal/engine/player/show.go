package player

import (
	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

// Cartoon is one light show the player can run.
type Cartoon interface {
	Kind() int // blob.Normal, blob.Script or blob.Audio
	Header() blob.FrameHeader
	// Render draws frame (Normal) or audio level (Audio) into a cleared canvas.
	Render(frame int, canvas *render.Canvas)
}

// Show is the list of cartoons played in one mode, with its key configuration.
type Show interface {
	Count() int
	// Cartoon returns cartoon index; it may be built on each call.
	Cartoon(index int) Cartoon
	// KeyFunction returns the action of KEY1..KEY7 (key 0..6).
	KeyFunction(key int) int
}

// scriptCartoon is a cartoon run by the script VM.
type scriptCartoon interface {
	script() (body string, a *assets)
}

// BlobShow plays cartoons from show data made by the iCubeSmart PC software.
type BlobShow struct {
	blob     *blob.Blob
	assets   *assets
	cartoons []blob.Cartoon
}

// NewBlobShow lists the cartoons of the given groups (0..3), in order.
func NewBlobShow(b *blob.Blob, groups ...int) *BlobShow {
	show := &BlobShow{blob: b, assets: newAssets(b)}
	for _, group := range groups {
		show.cartoons = append(show.cartoons, b.Group(group)...)
	}
	return show
}

func (this *BlobShow) Count() int { return len(this.cartoons) }

func (this *BlobShow) Cartoon(index int) Cartoon {
	return blobCartoon{this.cartoons[index], this.assets}
}

func (this *BlobShow) KeyFunction(key int) int { return this.blob.KeyFunction(key) }

type blobCartoon struct {
	cartoon blob.Cartoon
	assets  *assets
}

func (this blobCartoon) Kind() int { return this.cartoon.Type }

func (this blobCartoon) Header() blob.FrameHeader { return this.cartoon.Header() }

func (this blobCartoon) Render(frame int, canvas *render.Canvas) {
	renderFrame(this.cartoon, frame, canvas, this.assets)
}

func (this blobCartoon) script() (string, *assets) { return this.cartoon.Body, this.assets }
