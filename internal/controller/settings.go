package controller

import (
	"github.com/Tariomka/stm32-led-cube/internal/common"
	"github.com/Tariomka/stm32-led-cube/internal/component"
)

// EEPROM records: the player settings (this firmware) and the factory serial number and cube type.
const (
	settingsAddress = 0x00
	settingsSize    = 8
	serialAddress   = 0xA0
	cubeTypeAddress = 0xA8

	settingsMagic0, settingsMagic1 = 'L', 'C'
	settingsVersion                = 1
)

// Bits of the settings flags byte
const (
	flagPlayAll = 1 << iota
	flagPaused
	flagDisplayOff
)

// Factory cube type strings, indexed like the factory firmware's scan variants.
var cubeTypes = [...]string{"8x8x8A", "8x8x8B", "8x16x8A", "161616A", "161616B", "161616C", "8x32x8A", "8x32x8B"}

// Settings are the player settings kept across power cycles.
type Settings struct {
	Speed      uint8
	PlayAll    bool
	Paused     bool
	DisplayOff bool
	Mode       Mode // DIP switch mode the cartoon belongs to
	Cartoon    uint8
}

// SettingsController keeps the player settings and the factory cube type and serial number in
// the 24C02 EEPROM (I2C1: SCL - PB6, SDA - PB7).
type SettingsController struct {
	eeprom   component.EEPROM
	cubeType string
	typeRead bool
}

func NewSettingsController() *SettingsController {
	return &SettingsController{eeprom: component.NewEEPROM()}
}

// Load returns the saved player settings; false for a blank or foreign EEPROM.
func (this *SettingsController) Load() (Settings, bool) {
	var bytes [settingsSize]byte
	if this.eeprom.Read(settingsAddress, bytes[:]) != nil ||
		bytes[0] != settingsMagic0 || bytes[1] != settingsMagic1 || bytes[2] != settingsVersion ||
		bytes[settingsSize-1] != settingsChecksum(bytes[:settingsSize-1]) {
		return Settings{}, false
	}

	flags := bytes[4]
	return Settings{
		Speed:      bytes[3],
		PlayAll:    flags&flagPlayAll != 0,
		Paused:     flags&flagPaused != 0,
		DisplayOff: flags&flagDisplayOff != 0,
		Mode:       Mode(bytes[5]),
		Cartoon:    bytes[6],
	}, true
}

func (this *SettingsController) Save(settings Settings) error {
	bytes := [settingsSize]byte{settingsMagic0, settingsMagic1, settingsVersion, settings.Speed, 0, byte(settings.Mode), settings.Cartoon}
	if settings.PlayAll {
		bytes[4] |= flagPlayAll
	}
	if settings.Paused {
		bytes[4] |= flagPaused
	}
	if settings.DisplayOff {
		bytes[4] |= flagDisplayOff
	}
	bytes[settingsSize-1] = settingsChecksum(bytes[:settingsSize-1])
	return this.eeprom.Write(settingsAddress, bytes[:])
}

func settingsChecksum(bytes []byte) byte {
	var sum byte
	for _, value := range bytes {
		sum += value
	}
	return ^sum
}

// CubeType returns the factory cube type string (this kit: "8x32x8b"), "" if unreadable.
func (this *SettingsController) CubeType() string {
	if this.typeRead {
		return this.cubeType
	}

	var bytes [8]byte
	if this.eeprom.Read(cubeTypeAddress, bytes[:]) != nil {
		return ""
	}
	this.cubeType, this.typeRead = string(bytes[:]), true
	for index, char := range bytes {
		if char == 0 || char == 0xFF {
			this.cubeType = string(bytes[:index])
			break
		}
	}
	return this.cubeType
}

// CubeTypeIndex is the factory scan variant of the cube type (compared case-insensitively), -1
// if unknown. This firmware always drives one 8x8x8 cube; it only decides how display dir 1 is
// shown.
func (this *SettingsController) CubeTypeIndex() int {
	cubeType := this.CubeType()
	for index, name := range cubeTypes {
		if common.EqualFold(name, cubeType) {
			return index
		}
	}

	return -1
}

// IgnoresDisplayDir reports the cube types whose factory dir 1 conversion is the dir 0 one
// (161616B, 161616C, 8x32x8B, which includes this kit); the others show dir 1 upside down.
func (this *SettingsController) IgnoresDisplayDir() bool {
	index := this.CubeTypeIndex()
	return index == 4 || index == 5 || index == 7
}

func (this *SettingsController) SetCubeType(cubeType [8]byte) error {
	this.typeRead = false
	return this.eeprom.Write(cubeTypeAddress, cubeType[:])
}

func (this *SettingsController) SetSerial(serial [8]byte) error {
	return this.eeprom.Write(serialAddress, serial[:])
}
