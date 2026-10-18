package cascade

import "testing"

func feedAll(p *Parser, frame []byte) (payload []byte, ok bool) {
	for _, b := range frame {
		if payload, ok = p.Feed(b); ok {
			return append([]byte{}, payload...), true
		}
	}
	return nil, false
}

func TestFrame(t *testing.T) {
	frame := Encode([]byte{MsgNumber, 1})
	if want := []byte{0xA5, 0x55, 2, 2, 1, 2 ^ 2 ^ 1}; string(frame) != string(want) {
		t.Fatalf("frame = % X, want % X", frame, want)
	}
	var p Parser
	payload, ok := feedAll(&p, append([]byte{0x00, 0xA5, 0x12}, frame...))
	if !ok || string(payload) != string([]byte{MsgNumber, 1}) {
		t.Errorf("payload = % X, %v", payload, ok)
	}
	frame[len(frame)-1] ^= 1
	if _, ok := feedAll(&p, frame); ok {
		t.Error("accepted a bad checksum")
	}
}

func TestChain(t *testing.T) {
	var sent [][]byte
	slave := NewLink(func(p []byte) { sent = append(sent, append([]byte{}, p...)) })

	slave.Handle([]byte{MsgHello}) // upstream cube says hello while discovering
	if slave.Master {
		t.Fatal("a cube that heard a Hello stays master")
	}
	slave.state = stateNumber
	slave.Handle([]byte{MsgNumber, 1})
	if slave.Number != 1 || string(sent[0]) != string([]byte{MsgNumber, 2}) {
		t.Errorf("number = %d, forwarded % X", slave.Number, sent[0])
	}

	slave.state = stateRunning
	master := NewLink(func(p []byte) {
		if !slave.Handle(p) {
			t.Error("sync not handled")
		}
	})
	master.SendSync(1234, 0b100)
	if slave.Clock != 1234 || slave.Keys != 0b100 || len(sent) != 2 {
		t.Errorf("clock = %d, keys = %b, sent = %d", slave.Clock, slave.Keys, len(sent))
	}
}

func TestKeyMask(t *testing.T) {
	var k KeyMask
	// Short tap: one press edge, never seen as held
	_, edges := k.Next(0, 1)
	if edges != 1 {
		t.Fatalf("tap edge = %b", edges)
	}
	for range 5 {
		if _, edges := k.Next(0, 0); edges != 0 {
			t.Fatalf("extra edge %b", edges)
		}
	}
	// Long press: one edge, bit stays set while held
	_, edges = k.Next(2, 2)
	var mask uint8
	for range 10 {
		var e uint8
		mask, e = k.Next(2, 0)
		edges |= e
	}
	if edges != 2 || mask != 2 {
		t.Errorf("hold edges = %b, mask = %b", edges, mask)
	}
}
