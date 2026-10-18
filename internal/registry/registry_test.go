package registry

import (
	"strings"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	r := &Registry{}
	if err := r.Set("name", "cube"); err != nil {
		t.Fatal(err)
	}
	if err := r.Set("speed", "1.5"); err != nil {
		t.Fatal(err)
	}
	if err := r.Set("name", "my cube"); err != nil {
		t.Fatal(err)
	}

	image := r.Marshal()
	if len(image) != 0x20+2*0x10+8+4 {
		t.Errorf("image size = %d", len(image))
	}
	loaded, err := Unmarshal(append(image, 0xFF, 0xFF))
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := loaded.Get("name"); !ok || v != "my cube" {
		t.Errorf("name = %q, %v", v, ok)
	}
	if v, _ := loaded.Get("speed"); v != "1.5" {
		t.Errorf("speed = %q", v)
	}
}

func TestLimits(t *testing.T) {
	r := &Registry{}
	if r.Set("twelve_chars", "x") != ErrInvalid {
		t.Error("accepted a 12 character key")
	}
	if r.Set("k", strings.Repeat("v", MaxValueLength+1)) != ErrInvalid {
		t.Error("accepted a value over 256 bytes")
	}
	for i := 0; ; i++ {
		if err := r.Set(string(rune('A'+i/26))+string(rune('a'+i%26)), strings.Repeat("v", 100)); err != nil {
			if err != ErrFull || r.used() > ImageSize {
				t.Fatalf("err = %v, used = %d", err, r.used())
			}
			break
		}
	}
}

func TestErasedFlash(t *testing.T) {
	erased := make([]byte, ImageSize)
	for i := range erased {
		erased[i] = 0xFF
	}
	if _, err := Unmarshal(erased); err != ErrFormat {
		t.Errorf("err = %v, want ErrFormat", err)
	}
}
