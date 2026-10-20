package controller

import (
	"encoding/binary"
	"errors"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// Registry limits and flash image format of the factory firmware ("Regedit")
const (
	regMaxKey     = 11
	regMaxValue   = 251 // value + NUL, rounded up to 4, must fit in 256
	regMaxEntries = 64
	regImageSize  = component.FlashPageSize

	regHeader     = "Registry for 2.00"
	regHeaderSize = 0x20 // count (4) + header text (28)
	regEntrySize  = 0x10 // key[12] + value size (4)
)

var (
	errRegInvalid = errors.New("invalid key or value")
	errRegFull    = errors.New("registry full")
	errRegFormat  = errors.New("registry format error")
)

type registryEntry struct {
	key, value string
}

// registry is the console's key/value store (SetReg/GetReg/RegList/DelReg/ClrReg), kept in a
// flash page plus a backup page. The firmware itself never reads it.
type registry struct {
	entries []registryEntry
}

// loadRegistry reads the first valid image of the registry page and its backup.
func loadRegistry() *registry {
	for _, page := range []uintptr{component.RegistryPage, component.RegistryBackup} {
		if r, err := parseRegistry(component.FlashBytes(page, component.FlashPageSize)); err == nil {
			return r
		}
	}

	println("Regedit format error!")
	return &registry{}
}

// save writes the image to the registry page and then the backup page.
func (this *registry) save() error {
	image := this.marshal()
	if err := component.WriteFlashPage(component.RegistryPage, image); err != nil {
		return err
	}
	return component.WriteFlashPage(component.RegistryBackup, image)
}

func regValueSize(value string) int { return (len(value) + 4) &^ 3 }

func (this *registry) used() int {
	size := regHeaderSize
	for _, entry := range this.entries {
		size += regEntrySize + regValueSize(entry.value)
	}
	return size
}

func (this *registry) find(key string) int {
	for index, entry := range this.entries {
		if entry.key == key {
			return index
		}
	}
	return -1
}

// set adds or replaces a value.
func (this *registry) set(key, value string) error {
	if key == "" || len(key) > regMaxKey || value == "" || len(value) > regMaxValue {
		return errRegInvalid
	}

	if index := this.find(key); index >= 0 {
		if this.used()-regValueSize(this.entries[index].value)+regValueSize(value) > regImageSize {
			return errRegFull
		}
		this.entries[index].value = value
		return nil
	}

	if len(this.entries) >= regMaxEntries || this.used()+regEntrySize+regValueSize(value) > regImageSize {
		return errRegFull
	}
	this.entries = append(this.entries, registryEntry{key, value})
	return nil
}

func (this *registry) get(key string) (string, bool) {
	if index := this.find(key); index >= 0 {
		return this.entries[index].value, true
	}
	return "", false
}

// delete removes a key and reports whether it existed.
func (this *registry) delete(key string) bool {
	index := this.find(key)
	if index < 0 {
		return false
	}

	this.entries = append(this.entries[:index], this.entries[index+1:]...)
	return true
}

// marshal builds the flash image: count, header text, key/size table, then the values.
func (this *registry) marshal() []byte {
	image := make([]byte, this.used())
	binary.LittleEndian.PutUint32(image, uint32(len(this.entries)))
	copy(image[4:regHeaderSize], regHeader)
	pos := regHeaderSize + len(this.entries)*regEntrySize
	for index, entry := range this.entries {
		imageEntry := image[regHeaderSize+index*regEntrySize:]
		copy(imageEntry[:12], entry.key)
		size := regValueSize(entry.value)
		binary.LittleEndian.PutUint32(imageEntry[12:], uint32(size))
		copy(image[pos:], entry.value)
		pos += size
	}
	return image
}

// parseRegistry loads an image written by marshal or by the factory firmware.
func parseRegistry(image []byte) (*registry, error) {
	if len(image) < regHeaderSize || string(image[4:4+len(regHeader)]) != regHeader || image[4+len(regHeader)] != 0 {
		return nil, errRegFormat
	}

	count := int(binary.LittleEndian.Uint32(image))
	if count > regMaxEntries || regHeaderSize+count*regEntrySize > len(image) {
		return nil, errRegFormat
	}

	registry := &registry{}
	pos := regHeaderSize + count*regEntrySize
	for i := range count {
		entry := image[regHeaderSize+i*regEntrySize:]
		size := int(binary.LittleEndian.Uint32(entry[12:]))
		if size == 0 || pos+size > len(image) {
			return nil, errRegFormat
		}

		registry.entries = append(registry.entries, registryEntry{key: cString(entry[:12]), value: cString(image[pos : pos+size])})
		pos += size
	}
	return registry, nil
}

func cString(bytes []byte) string {
	for index, char := range bytes {
		if char == 0 {
			return string(bytes[:index])
		}
	}
	return string(bytes)
}
