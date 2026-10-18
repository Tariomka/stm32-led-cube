package audio

import "testing"

func TestLevel(t *testing.T) {
	m := NewMeter()
	if got := m.Level(0); got != 0 {
		t.Fatalf("silent level = %d, want 0", got)
	}

	m.Add(1775) // (1775-1250)*62/1050+1 = 32
	if got := m.Level(0); got != 32 {
		t.Errorf("level = %d, want 32", got)
	}

	m.Add(325) // (850-325)*62/1050+1 = 32, the swing below idle counts too
	m.Add(4095)
	if got := m.Level(0); got != MaxLevel {
		t.Errorf("loud level = %d, want %d", got, MaxLevel)
	}
}

func TestDecay(t *testing.T) {
	m := NewMeter()
	m.Add(4095)
	if got := m.Level(ModeDecay); got != MaxLevel {
		t.Fatalf("level = %d, want %d", got, MaxLevel)
	}
	for range window {
		m.Add(idle)
	}
	if got := m.Level(0); got != 0 {
		t.Errorf("level without decay = %d, want 0", got)
	}
	if got := m.Level(ModeDecay); got != MaxLevel-1 {
		t.Errorf("decayed level = %d, want %d", got, MaxLevel-1)
	}
}
