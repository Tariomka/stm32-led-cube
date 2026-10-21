// Package builtin holds the light shows built into the firmware: the original iCubeSmart cartoons
// as Go functions (generated once from the factory show data: Dynamic in dynamic.go, Audio in audio.go, Test in
// test.go, shapes, palettes and Keys in library.go).
package builtin

import "github.com/Tariomka/led-common-lib/pkg/lightshow/keyframe"

var library = &keyframe.Library{Shapes: shapes, Palettes: palettes}

// New makes a playable show of cartoons with the original key functions.
func New(cartoons []func() *keyframe.Cartoon) *keyframe.Show {
	return &keyframe.Show{Cartoons: cartoons, Keys: Keys, Library: library}
}
