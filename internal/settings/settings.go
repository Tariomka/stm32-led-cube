// Package settings encodes the player settings kept in the 24C02 EEPROM across power cycles.
// The original firmware doesn't persist them; it only keeps its type/serial at 0xA0..0xAF.
package settings

const (
	Address = 0x00 // EEPROM address, clear of the vendor records at 0xA0..0xAF
	Size    = 8

	magic0, magic1 = 'L', 'C'
	version        = 1
)

type Settings struct {
	Speed      uint8
	PlayAll    bool
	Paused     bool
	DisplayOff bool
	Group      uint8 // DIP group the cartoon belongs to
	Cartoon    uint8
}

// Bits of the flags byte
const (
	flagPlayAll = 1 << iota
	flagPaused
	flagDisplayOff
)

func checksum(bytes []byte) byte {
	var sum byte
	for _, value := range bytes {
		sum += value
	}
	return ^sum
}

func (this Settings) Marshal() [Size]byte {
	bytes := [Size]byte{magic0, magic1, version, this.Speed, 0, this.Group, this.Cartoon}
	if this.PlayAll {
		bytes[4] |= flagPlayAll
	}
	if this.Paused {
		bytes[4] |= flagPaused
	}
	if this.DisplayOff {
		bytes[4] |= flagDisplayOff
	}
	bytes[Size-1] = checksum(bytes[:Size-1])
	return bytes
}

// Unmarshal returns false for blank or foreign EEPROM contents.
func Unmarshal(bytes [Size]byte) (Settings, bool) {
	if bytes[0] != magic0 || bytes[1] != magic1 || bytes[2] != version || bytes[Size-1] != checksum(bytes[:Size-1]) {
		return Settings{}, false
	}

	flags := bytes[4]
	return Settings{
		Speed:      bytes[3],
		PlayAll:    flags&flagPlayAll != 0,
		Paused:     flags&flagPaused != 0,
		DisplayOff: flags&flagDisplayOff != 0,
		Group:      bytes[5],
		Cartoon:    bytes[6],
	}, true
}
