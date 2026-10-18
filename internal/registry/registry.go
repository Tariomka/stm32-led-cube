// Package registry is the original firmware's key/value store ("Regedit", console commands
// SetReg/GetReg/RegList/DelReg/ClrReg) with its limits and flash image format.
package registry

import (
	"encoding/binary"
	"errors"
)

const (
	MaxKeyLength   = 11
	MaxValueLength = 251 // value + NUL, rounded up to 4, must fit in 256
	MaxEntries     = 64
	ImageSize      = 0x800 // one 2 KB flash page

	header     = "Registry for 2.00"
	headerSize = 0x20 // count (4) + header text (28)
	entrySize  = 0x10 // key[12] + value size (4)
)

var (
	ErrInvalid = errors.New("invalid key or value")
	ErrFull    = errors.New("registry full")
	ErrFormat  = errors.New("registry format error")
)

type Entry struct {
	Key, Value string
}

type Registry struct {
	entries []Entry
}

func valueSize(value string) int { return (len(value) + 4) &^ 3 }

func (this *Registry) used() int {
	sizeUsed := headerSize
	for _, entry := range this.entries {
		sizeUsed += entrySize + valueSize(entry.Value)
	}
	return sizeUsed
}

func (this *Registry) find(key string) int {
	for index, entry := range this.entries {
		if entry.Key == key {
			return index
		}
	}

	return -1
}

func validKey(key string) bool { return key != "" && len(key) <= MaxKeyLength }

// Set adds or replaces a value.
func (this *Registry) Set(key, value string) error {
	if !validKey(key) || value == "" || len(value) > MaxValueLength {
		return ErrInvalid
	}

	if index := this.find(key); index >= 0 {
		if this.used()-valueSize(this.entries[index].Value)+valueSize(value) > ImageSize {
			return ErrFull
		}

		this.entries[index].Value = value
		return nil
	}

	if len(this.entries) >= MaxEntries || this.used()+entrySize+valueSize(value) > ImageSize {
		return ErrFull
	}

	this.entries = append(this.entries, Entry{key, value})
	return nil
}

func (this *Registry) Get(key string) (string, bool) {
	if index := this.find(key); index >= 0 {
		return this.entries[index].Value, true
	}

	return "", false
}

// Delete removes a key and reports whether it existed.
func (this *Registry) Delete(key string) bool {
	index := this.find(key)
	if index < 0 {
		return false
	}

	this.entries = append(this.entries[:index], this.entries[index+1:]...)
	return true
}

func (this *Registry) Clear() { this.entries = nil }

func (this *Registry) Entries() []Entry { return this.entries }

// Marshal builds the flash image: count, header text, key/size table, then the values.
func (this *Registry) Marshal() []byte {
	image := make([]byte, this.used())
	binary.LittleEndian.PutUint32(image, uint32(len(this.entries)))
	copy(image[4:headerSize], header)
	pos := headerSize + len(this.entries)*entrySize
	for index, entry := range this.entries {
		imageEntry := image[headerSize+index*entrySize:]
		copy(imageEntry[:12], entry.Key)
		size := valueSize(entry.Value)
		binary.LittleEndian.PutUint32(imageEntry[12:], uint32(size))
		copy(image[pos:], entry.Value)
		pos += size
	}
	return image
}

// Unmarshal loads a flash image written by Marshal or by the original firmware.
func Unmarshal(image []byte) (*Registry, error) {
	if len(image) < headerSize || string(image[4:4+len(header)]) != header || image[4+len(header)] != 0 {
		return nil, ErrFormat
	}

	count := int(binary.LittleEndian.Uint32(image))
	if count > MaxEntries || headerSize+count*entrySize > len(image) {
		return nil, ErrFormat
	}

	registry := &Registry{}
	pos := headerSize + count*entrySize
	for i := range count {
		entry := image[headerSize+i*entrySize:]
		size := int(binary.LittleEndian.Uint32(entry[12:]))
		if size == 0 || pos+size > len(image) {
			return nil, ErrFormat
		}

		registry.entries = append(registry.entries, Entry{Key: cstring(entry[:12]), Value: cstring(image[pos : pos+size])})
		pos += size
	}

	return registry, nil
}

func cstring(bytes []byte) string {
	for index, char := range bytes {
		if char == 0 {
			return string(bytes[:index])
		}
	}

	return string(bytes)
}
