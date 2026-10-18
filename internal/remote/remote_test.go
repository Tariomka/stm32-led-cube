package remote

import "testing"

// send feeds a NEC frame (leader + 32 bits, MSB first per byte) and returns the decoder result.
func send(d *Decoder, bytes [4]byte) (int, bool) {
	d.Pulse(9000, 4500)
	var code int
	var ok bool
	for i := 0; i < 32; i++ {
		space := 560
		if bytes[i/8]&(0x80>>(i%8)) != 0 {
			space = 1690
		}
		code, ok = d.Pulse(560, space)
		if ok && i < 31 {
			return code, false
		}
	}
	return code, ok
}

func TestDecode(t *testing.T) {
	d := NewDecoder()
	code, ok := send(d, [4]byte{0x00, 0xFF, 0x68, 0x97})
	if !ok || code != 0x68 {
		t.Fatalf("code = %#x, %v; want 0x68, true", code, ok)
	}
	if Action(code) != 10 {
		t.Errorf("action = %d, want 10", Action(code))
	}
}

func TestRejectsBadComplement(t *testing.T) {
	d := NewDecoder()
	if code, ok := send(d, [4]byte{0x00, 0xFF, 0x68, 0x96}); ok {
		t.Errorf("accepted corrupted frame %#x", code)
	}
}

func TestIgnoresRepeat(t *testing.T) {
	d := NewDecoder()
	d.Pulse(9000, 2250)
	if _, ok := d.Pulse(560, 560); ok {
		t.Error("repeat frame produced a code")
	}
}
