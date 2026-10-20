package component

// Layer data layout: 24 bytes per layer, sent LSB first.
// Bytes 0-7 green, 8-15 blue, 16-23 red; byte = row, bit = column.
const (
	frameLayers   = 8
	bytesPerLayer = 24
	framePlanes   = 3 // channel level n (0..3) lights planes 0..n-1, shown 1:3:6 of the time like the factory firmware

	greenOffset = 0
	blueOffset  = 8
	redOffset   = 16
)

// frameBits holds the planes of every layer as shifted into the MBI5024 chain.
type frameBits [framePlanes][frameLayers][bytesPerLayer]byte

// Mapping orients the show's 8x8x8 coordinate space (x, y = up, z) on the physical cube.
type Mapping struct {
	FlipX     bool // mirror columns
	FlipZ     bool // mirror rows
	FlipLayer bool // show y=0 is the top layer instead of the bottom
	SwapXZ    bool // exchange columns and rows
}

// encodeFrame converts an 8x8x8 voxel buffer (index x + y*8 + z*64, colors 0bRRGGBB) into planes.
func encodeFrame(voxels []byte, mapping Mapping, out *frameBits) {
	*out = frameBits{}
	for z := range 8 {
		for y := range 8 {
			for x := range 8 {
				c := voxels[x+y*8+z*64]
				if c == 0 {
					continue
				}

				column, row, layer := x, z, y
				if mapping.SwapXZ {
					column, row = row, column
				}
				if mapping.FlipX {
					column = 7 - column
				}
				if mapping.FlipZ {
					row = 7 - row
				}
				if mapping.FlipLayer {
					layer = 7 - layer
				}
				bit := byte(1) << column
				red, green, blue := int(c>>4&3), int(c>>2&3), int(c&3)
				for plane := range framePlanes {
					if red > plane {
						out[plane][layer][redOffset+row] |= bit
					}
					if green > plane {
						out[plane][layer][greenOffset+row] |= bit
					}
					if blue > plane {
						out[plane][layer][blueOffset+row] |= bit
					}
				}
			}
		}
	}
}
