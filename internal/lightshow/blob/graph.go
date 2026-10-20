package blob

import "errors"

// Compression types of graph data
const (
	NoCompress = 0
	Normal2    = 1 // 4-bit packed palette indexes
	Normal4    = 2 // 2-bit
	Normal8    = 3 // 1-bit
	Continue2  = 4 // RLE with 4-bit palette indexes
	Continue4  = 5 // 2-bit
	Continue8  = 6 // 1-bit
)

var ErrCompression = errors.New("blob: bad graph data")

type Graph struct {
	Name                  string
	Length, Height, Width int
	Compression           int
	Data                  string
	offset                int
}

// ID is unique per graph within a blob, usable as a cache key.
func (this Graph) ID() int { return this.offset }

func (this Graph) Voxels() int { return this.Length * this.Height * this.Width }

// Decode expands the graph voxels (palette indexes, x fastest, then y, then z) into dst.
func (this Graph) Decode(destination []byte) error {
	voxels := this.Voxels()
	if len(destination) < voxels {
		return ErrCompression
	}

	destination = destination[:voxels]
	switch this.Compression {
	case NoCompress:
		if len(this.Data) < voxels {
			return ErrCompression
		}

		copy(destination, this.Data[:voxels])
		return nil
	case Normal2, Normal4, Normal8:
		return decodePacked(this.Data, destination, bitsFor(this.Compression-Normal2))
	case Continue2, Continue4, Continue8:
		return decodeRLE(this.Data, destination, bitsFor(this.Compression-Continue2))
	}

	return ErrCompression
}

func bitsFor(index int) uint {
	return [3]uint{4, 2, 1}[index]
}

func decodePacked(source string, destination []byte, bits uint) error {
	if len(source) < 1 || len(source) < 1+int(source[0]) {
		return ErrCompression
	}

	lut := source[1 : 1+int(source[0])]
	mask := byte(1)<<bits - 1
	perByte := 8 / int(bits)
	out := 0
	for index := 1 + len(lut); index < len(source) && out < len(destination); index++ {
		b := source[index]
		for j := 0; j < perByte && out < len(destination); j++ {
			idx := int(b & mask)
			if idx >= len(lut) {
				return ErrCompression
			}
			destination[out] = lut[idx]
			out++
			b >>= bits
		}
	}
	return nil
}

func decodeRLE(source string, destination []byte, bits uint) error {
	if len(source) < 1 || len(source) < 1+int(source[0]) {
		return ErrCompression
	}

	lut := source[1 : 1+int(source[0])]
	mask := byte(1)<<bits - 1
	out := 0
	for index := 1 + len(lut); index < len(source); {
		b := source[index]
		var count int
		var idx byte
		if b&0x80 != 0 {
			if index+1 >= len(source) {
				return ErrCompression
			}

			b2 := source[index+1]
			count = int(b&0x7F)<<(8-bits) | int(b2>>bits)
			idx = b2 & mask
			index += 2
		} else {
			count = int(b >> bits)
			idx = b & mask
			index++
		}
		if int(idx) >= len(lut) {
			return ErrCompression
		}

		for ; count > 0 && out < len(destination); count-- {
			destination[out] = lut[idx]
			out++
		}
	}

	return nil
}
