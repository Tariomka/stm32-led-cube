package controller

import (
	"device/arm"
	"time"
	"unsafe"

	"github.com/Tariomka/stm32-led-cube/internal/component"
	"github.com/Tariomka/stm32-led-cube/internal/registry"
	"github.com/Tariomka/stm32-led-cube/internal/settings"
)

// Vendor records in the 24C02, see .agent/re/firmware-map.md
const (
	eepromSerial = 0xA0
	eepromType   = 0xA8
)

// The YellowBoard is the console.Device of the vendor-style UART console.

// rebootDelay lets the vendor PC tool's upload connect frame, sent right after its "reboot"
// command, cancel the reboot: the original rebooted into its bootloader, we take the upload.
const rebootDelay = 300 * time.Millisecond

func (this *YellowBoard) Reboot() {
	this.rebootAt = time.Now().Add(rebootDelay)
}

// ReadMemory reads flash, SRAM and system memory; peripherals are refused because reads
// can have side effects.
func (this *YellowBoard) ReadMemory(address uint32) (byte, bool) {
	valid := (address >= 0x08000000 && address < 0x08080000) ||
		(address >= 0x20000000 && address < 0x20010000) ||
		(address >= 0x1FFFF000 && address < 0x1FFFF810)
	if !valid {
		return 0, false
	}
	return *(*byte)(unsafe.Pointer(uintptr(address))), true
}

func (this *YellowBoard) SetType(cubeType [8]byte) error {
	return this.EEPROM.Write(eepromType, cubeType[:])
}

func (this *YellowBoard) SetSerial(serial [8]byte) error {
	return this.EEPROM.Write(eepromSerial, serial[:])
}

// SaveRegistry writes the registry to its flash page and the backup page.
func (this *YellowBoard) SaveRegistry(r *registry.Registry) error {
	image := r.Marshal()
	if err := component.WriteFlashPage(component.RegistryPage, image); err != nil {
		return err
	}
	return component.WriteFlashPage(component.RegistryBackup, image)
}

func loadRegistry() *registry.Registry {
	for _, page := range []uintptr{component.RegistryPage, component.RegistryBackup} {
		if r, err := registry.Unmarshal(component.FlashBytes(page, component.FlashPageSize)); err == nil {
			return r
		}
	}
	println("Regedit format error!")
	return &registry.Registry{}
}

// CubeType returns the vendor type string stored in the EEPROM (e.g. "8x8x8A").
func (this *YellowBoard) CubeType() string {
	var b [8]byte
	if this.EEPROM.Read(eepromType, b[:]) != nil {
		return ""
	}
	for i, c := range b {
		if c == 0 || c == 0xFF {
			return string(b[:i])
		}
	}
	return string(b[:])
}

func (this *YellowBoard) LoadSettings() (settings.Settings, bool) {
	var b [settings.Size]byte
	if this.EEPROM.Read(settings.Address, b[:]) != nil {
		return settings.Settings{}, false
	}
	return settings.Unmarshal(b)
}

func (this *YellowBoard) SaveSettings(s settings.Settings) error {
	b := s.Marshal()
	return this.EEPROM.Write(settings.Address, b[:])
}

func (this *YellowBoard) PollSerial() {
	this.Serial.Poll()
	if !this.rebootAt.IsZero() && time.Now().After(this.rebootAt) {
		arm.SystemReset()
	}
}

func (this *YellowBoard) OnUpload(started func(), finished func(ok bool)) {
	this.Upload.Started = func() {
		this.rebootAt = time.Time{}
		println("Show upload started")
		started()
	}
	this.Upload.Finished = finished
}

func (this *YellowBoard) UploadedShow() string {
	return component.ShowSlotData()
}
