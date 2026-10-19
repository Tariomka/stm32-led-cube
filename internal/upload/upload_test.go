package upload

import (
	"bytes"
	"testing"

	"github.com/Tariomka/stm32-led-cube/internal/engine/blob"
	"github.com/Tariomka/stm32-led-cube/internal/show/original"
)

type memory struct {
	data   [ShowSize]byte
	begins int
}

func (m *memory) Begin() error { m.begins++; return nil }

func (m *memory) Write(offset int, data []byte) error {
	copy(m.data[offset:], data)
	return nil
}

func record(line, address int, data []byte) []byte {
	p := []byte{msgRecord, byte(len(data)), byte(line >> 8), byte(line),
		byte(address >> 24), byte(address >> 16), byte(address >> 8), byte(address)}
	return append(p, data...)
}

func TestUpload(t *testing.T) {
	store := &memory{}
	var replies [][]byte
	var result *bool
	r := Receiver{
		Storage:  store,
		Reply:    func(p []byte) { replies = append(replies, append([]byte{}, p...)) },
		Finished: func(ok bool) { result = &ok },
	}

	r.Frame(append([]byte{msgConnect}, make([]byte, 20)...))
	r.Frame(record(1, 0x08006000, []byte{1, 2, 3, 4}))              // firmware code: acknowledged, not stored
	r.Frame(record(2, ShowStart-2, []byte{9, 9, 0x65, 0x87}))       // straddles the show start
	r.Frame(record(3, ShowStart+2, []byte{0x34, 0x12, 0xAA, 0xBB})) // inside
	r.Frame([]byte{msgEnd, 0})
	r.Frame([]byte{msgEnd, 0})

	if result == nil || !*result {
		t.Fatalf("finished = %v", result)
	}
	if !bytes.Equal(store.data[:6], []byte{0x65, 0x87, 0x34, 0x12, 0xAA, 0xBB}) {
		t.Errorf("stored % X", store.data[:6])
	}
	want := [][]byte{{msgReady, 0, 0, 0, 0, 0}, {msgStored, 4, 0, 1}, {msgStored, 4, 0, 2}, {msgStored, 4, 0, 3}}
	if len(replies) != len(want) {
		t.Fatalf("replies = % X", replies)
	}
	for i := range want {
		if !bytes.Equal(replies[i], want[i]) {
			t.Errorf("reply %d = % X, want % X", i, replies[i], want[i])
		}
	}
}

// TestUploadOriginal sends the original show data as 16-byte hex records, like the PC tool.
func TestUploadOriginal(t *testing.T) {
	store := &memory{}
	acks := 0
	r := Receiver{Storage: store, Reply: func([]byte) { acks++ }}
	r.Frame(append([]byte{msgConnect}, make([]byte, 20)...))
	data := original.Blob
	line := 0
	for offset := 0; offset < len(data); offset += 16 {
		line++
		r.Frame(record(line, ShowStart+offset, []byte(data[offset:min(offset+16, len(data))])))
	}
	r.Frame([]byte{msgEnd, 0})

	if acks != line+1 {
		t.Errorf("acks = %d, want %d", acks, line+1)
	}
	if string(store.data[:len(data)]) != data {
		t.Error("stored show data differs from the original")
	}
	if _, err := blob.Parse(string(store.data[:])); err != nil {
		t.Error(err)
	}
}

func TestFrame(t *testing.T) {
	// Matches the PC tool's WriteSerialData for payload {37, 0}
	if got, want := Frame([]byte{37, 0}), []byte{0xA5, 0x05, 2, 37, 0, 0xA5 ^ 0x05 ^ 2 ^ 37 ^ 0x55, 0x55}; !bytes.Equal(got, want) {
		t.Errorf("frame = % X, want % X", got, want)
	}
}
