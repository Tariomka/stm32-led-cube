package component

import (
	"device/stm32"
	"errors"
	"runtime/volatile"
	"unsafe"
)

// Internal flash areas reserved by the linker script (targets/gd32f103re.ld):
// uploaded show data, then the console registry page and its backup like the original firmware.
const (
	FlashPageSize  = 0x800
	ShowSlot       = 0x08057800
	ShowSlotSize   = 0x27800 // 79 pages, room for the vendor's whole show region (0x27100 bytes)
	RegistryPage   = 0x0807F000
	RegistryBackup = 0x0807F800
)

var ErrFlash = errors.New("flash program error")

// FlashBytes returns a read-only view of size bytes of flash at address.
func FlashBytes(address uintptr, size int) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(address)), size)
}

func flashWait() error {
	for stm32.FLASH.SR.HasBits(stm32.FLASH_SR_BSY) {
	}
	sr := stm32.FLASH.SR.Get()
	stm32.FLASH.SR.Set(stm32.FLASH_SR_EOP | stm32.FLASH_SR_PGERR | stm32.FLASH_SR_WRPRTERR)
	if sr&(stm32.FLASH_SR_PGERR|stm32.FLASH_SR_WRPRTERR) != 0 {
		return ErrFlash
	}
	return nil
}

func unlockFlash() {
	if stm32.FLASH.CR.HasBits(stm32.FLASH_CR_LOCK) {
		stm32.FLASH.KEYR.Set(0x45670123)
		stm32.FLASH.KEYR.Set(0xCDEF89AB)
	}
}

// eraseFlashPage erases the 2 KB page at address. The CPU stalls while the flash is busy,
// so the display blanks for a moment.
func eraseFlashPage(address uintptr) error {
	unlockFlash()
	defer stm32.FLASH.CR.SetBits(stm32.FLASH_CR_LOCK)

	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_PER)
	stm32.FLASH.AR.Set(uint32(address))
	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_STRT)
	err := flashWait()
	stm32.FLASH.CR.ClearBits(stm32.FLASH_CR_PER)
	return err
}

// programFlash writes data at an even address of erased flash; halfwords that already hold
// the value are skipped, so a repeated write is harmless.
func programFlash(address uintptr, data []byte) error {
	if address%2 != 0 {
		return ErrFlash
	}
	unlockFlash()
	defer stm32.FLASH.CR.SetBits(stm32.FLASH_CR_LOCK)

	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_PG)
	defer stm32.FLASH.CR.ClearBits(stm32.FLASH_CR_PG)
	for i := 0; i < len(data); i += 2 {
		half := uint16(data[i]) | 0xFF00
		if i+1 < len(data) {
			half = uint16(data[i]) | uint16(data[i+1])<<8
		}
		target := (*volatile.Register16)(unsafe.Pointer(address + uintptr(i)))
		if target.Get() == half {
			continue
		}
		target.Set(half)
		if err := flashWait(); err != nil {
			return err
		}
	}
	return nil
}

// WriteFlashPage erases the 2 KB page at address and programs data.
func WriteFlashPage(address uintptr, data []byte) error {
	if len(data) > FlashPageSize {
		return ErrFlash
	}
	if err := eraseFlashPage(address); err != nil {
		return err
	}
	return programFlash(address, data)
}

// ShowSlotStorage stores uploaded show data in the show slot (an upload.Storage); each page is
// erased the first time an upload writes to it.
type ShowSlotStorage struct {
	erased [ShowSlotSize / FlashPageSize]bool
}

func (this *ShowSlotStorage) Begin() error {
	this.erased = [len(this.erased)]bool{}
	return nil
}

func (this *ShowSlotStorage) Write(offset int, data []byte) error {
	if offset < 0 || offset+len(data) > ShowSlotSize || len(data) == 0 {
		return ErrFlash
	}
	for page := offset / FlashPageSize; page <= (offset+len(data)-1)/FlashPageSize; page++ {
		if !this.erased[page] {
			if err := eraseFlashPage(ShowSlot + uintptr(page*FlashPageSize)); err != nil {
				return err
			}
			this.erased[page] = true
		}
	}
	return programFlash(ShowSlot+uintptr(offset), data)
}

// ShowSlotData is a read-only view of the show slot.
func ShowSlotData() string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ShowSlot))), ShowSlotSize)
}
