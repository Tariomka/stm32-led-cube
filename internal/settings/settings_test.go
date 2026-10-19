package settings

import "testing"

func TestRoundTrip(t *testing.T) {
	s := Settings{Speed: 6, PlayAll: true, Paused: true, DisplayOff: true, Mode: 2, Cartoon: 3}
	got, ok := Unmarshal(s.Marshal())
	if !ok || got != s {
		t.Errorf("got %+v, %v", got, ok)
	}
}

func TestRejectsBlank(t *testing.T) {
	blank := [Size]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	if _, ok := Unmarshal(blank); ok {
		t.Error("accepted erased EEPROM")
	}
}
