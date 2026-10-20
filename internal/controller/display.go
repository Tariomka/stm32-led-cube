package controller

import "github.com/Tariomka/stm32-led-cube/internal/component"

// DisplayController shows frames on the LED cube: the MBI5024 column drivers and the 74HC154
// layer decoder, refreshed by the TIM2 interrupt.
type DisplayController struct {
	scanner    *component.Scanner
	mapping    component.Mapping // orientation of the show coordinates on this kit (none needed)
	upsideDown bool
}

func NewDisplayController() *DisplayController {
	return &DisplayController{scanner: component.NewScanner()}
}

// Show presents an 8x8x8 frame (index x + y*8 + z*64, colors 0bRRGGBB).
func (this *DisplayController) Show(voxels []byte) {
	mapping := this.mapping
	if this.upsideDown {
		mapping.FlipX, mapping.FlipLayer = !mapping.FlipX, !mapping.FlipLayer
	}
	this.scanner.Show(voxels, mapping)
}

// SetUpsideDown rotates the frames 180° about the z axis (show display dir 1).
func (this *DisplayController) SetUpsideDown(upsideDown bool) { this.upsideDown = upsideDown }

// SetBlank turns all LEDs off (true) or back on (false).
func (this *DisplayController) SetBlank(blank bool) { this.scanner.SetBlank(blank) }
