// fwtool extracts the show-data blob from an iCubeSmart firmware image (vendor *.hex or a full flash dump *.bin).
//
//	go run ./tools/fwtool <firmware.hex|dump.bin> <out.bin>
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Tariomka/stm32-led-cube/internal/lightshow/blob"
)

const (
	flashBase = 0x08000000
	flashSize = 0x80000
	blobAddr  = 0x0800B000
	blobMax   = 160000
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: fwtool <firmware.hex|dump.bin> <out.bin>")
		os.Exit(2)
	}
	img, err := load(os.Args[1])
	if err != nil {
		fail(err)
	}
	raw := img[blobAddr-flashBase : blobAddr-flashBase+blobMax]
	b, err := blob.Parse(string(raw))
	if err != nil {
		fail(err)
	}
	out := raw[:b.Size()]
	if err := os.WriteFile(os.Args[2], out, 0o644); err != nil {
		fail(err)
	}
	fmt.Printf("wrote %d bytes to %s\n", len(out), os.Args[2])
}

func load(path string) ([]byte, error) {
	if strings.EqualFold(path[max(0, len(path)-4):], ".bin") {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		img := make([]byte, flashSize)
		copy(img, data)
		return img, nil
	}
	return loadHex(path)
}

func loadHex(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img := make([]byte, flashSize)
	for i := range img {
		img[i] = 0xFF
	}
	var base uint64
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		l := strings.TrimSpace(sc.Text())
		if len(l) < 11 || l[0] != ':' {
			continue
		}
		n, _ := strconv.ParseUint(l[1:3], 16, 8)
		addr, _ := strconv.ParseUint(l[3:7], 16, 16)
		typ, _ := strconv.ParseUint(l[7:9], 16, 8)
		if len(l) < 11+2*int(n) {
			return nil, fmt.Errorf("%s:%d: truncated record", path, line)
		}
		switch typ {
		case 0:
			for k := 0; k < int(n); k++ {
				v, _ := strconv.ParseUint(l[9+2*k:11+2*k], 16, 8)
				off := int(base+addr) - flashBase + k
				if off >= 0 && off < len(img) {
					img[off] = byte(v)
				}
			}
		case 4:
			v, _ := strconv.ParseUint(l[9:13], 16, 16)
			base = v << 16
		}
	}
	return img, sc.Err()
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "fwtool:", err)
	os.Exit(1)
}
