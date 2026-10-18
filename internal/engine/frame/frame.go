// Package frame converts engine voxel buffers into the bit planes shifted into the cube's
// MBI5024 driver chain.
package frame

// Layer data layout (from led-common-lib): 24 bytes per layer, sent LSB first.
// Bytes 0-7 green, 8-15 blue, 16-23 red; byte = row, bit = column.
const (
	Layers        = 8
	BytesPerLayer = 24
	Planes        = 2 // 2-bit color depth per channel, shown with binary code modulation

	greenOffset = 0
	blueOffset  = 8
	redOffset   = 16
)

// Bits holds plane 0 (weight 1) and plane 1 (weight 2) for every layer.
type Bits [Planes][Layers][BytesPerLayer]byte

// Mapping orients the engine's 8x8x8 coordinate space (x, y = up, z) on the physical cube.
type Mapping struct {
	FlipX     bool // mirror columns
	FlipZ     bool // mirror rows
	FlipLayer bool // engine y=0 is the top layer instead of the bottom
	SwapXZ    bool // exchange columns and rows
}

// Encode converts an 8x8x8 voxel buffer (index x + y*8 + z*64, colors 0bRRGGBB) into bit planes.
func Encode(voxels []byte, m Mapping, out *Bits) {
	*out = Bits{}
	for z := 0; z < 8; z++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				c := voxels[x+y*8+z*64]
				if c == 0 {
					continue
				}
				col, row, layer := x, z, y
				if m.SwapXZ {
					col, row = row, col
				}
				if m.FlipX {
					col = 7 - col
				}
				if m.FlipZ {
					row = 7 - row
				}
				if m.FlipLayer {
					layer = 7 - layer
				}
				bit := byte(1) << col
				for p := 0; p < Planes; p++ {
					if c>>(4+p)&1 != 0 {
						out[p][layer][redOffset+row] |= bit
					}
					if c>>(2+p)&1 != 0 {
						out[p][layer][greenOffset+row] |= bit
					}
					if c>>p&1 != 0 {
						out[p][layer][blueOffset+row] |= bit
					}
				}
			}
		}
	}
}

// FromLayer fills both planes with the same 24-byte layer (full brightness, used for
// led-common-lib layouts which have no color depth).
func FromLayer(out *Bits, layer int, data []byte) {
	copy(out[0][layer][:], data)
	copy(out[1][layer][:], data)
}
