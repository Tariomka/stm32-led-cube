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

// Decode expands the graph voxels (palette indexes, x fastest, then y, then z) into dst.
func (g Graph) Decode(dst []byte) error {
	n := g.Voxels()
	if len(dst) < n {
		return ErrCompression
	}
	dst = dst[:n]
	switch g.Compression {
	case NoCompress:
		if len(g.Data) < n {
			return ErrCompression
		}
		copy(dst, g.Data[:n])
		return nil
	case Normal2, Normal4, Normal8:
		return decodePacked(g.Data, dst, bitsFor(g.Compression-Normal2))
	case Continue2, Continue4, Continue8:
		return decodeRLE(g.Data, dst, bitsFor(g.Compression-Continue2))
	}
	return ErrCompression
}

func bitsFor(i int) uint {
	return [3]uint{4, 2, 1}[i]
}

func decodePacked(src string, dst []byte, bits uint) error {
	if len(src) < 1 || len(src) < 1+int(src[0]) {
		return ErrCompression
	}
	lut := src[1 : 1+int(src[0])]
	mask := byte(1)<<bits - 1
	perByte := 8 / int(bits)
	out := 0
	for i := 1 + len(lut); i < len(src) && out < len(dst); i++ {
		b := src[i]
		for j := 0; j < perByte && out < len(dst); j++ {
			idx := int(b & mask)
			if idx >= len(lut) {
				return ErrCompression
			}
			dst[out] = lut[idx]
			out++
			b >>= bits
		}
	}
	return nil
}

func decodeRLE(src string, dst []byte, bits uint) error {
	if len(src) < 1 || len(src) < 1+int(src[0]) {
		return ErrCompression
	}
	lut := src[1 : 1+int(src[0])]
	mask := byte(1)<<bits - 1
	out := 0
	for i := 1 + len(lut); i < len(src); {
		b := src[i]
		var count int
		var idx byte
		if b&0x80 != 0 {
			if i+1 >= len(src) {
				return ErrCompression
			}
			b2 := src[i+1]
			count = int(b&0x7F)<<(8-bits) | int(b2>>bits)
			idx = b2 & mask
			i += 2
		} else {
			count = int(b >> bits)
			idx = b & mask
			i++
		}
		if int(idx) >= len(lut) {
			return ErrCompression
		}
		for ; count > 0 && out < len(dst); count-- {
			dst[out] = lut[idx]
			out++
		}
	}
	return nil
}
