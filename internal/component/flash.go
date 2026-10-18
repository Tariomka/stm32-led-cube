package component

import (
	"device/stm32"
	"errors"
	"runtime/volatile"
	"unsafe"
)

// Internal flash page reserved by the linker script (targets/gd32f103re.ld) for the console
// registry, plus a backup copy in the next page, like the original firmware.
const (
	RegistryPage   = 0x0807F000
	RegistryBackup = 0x0807F800
	FlashPageSize  = 0x800
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

// WriteFlashPage erases the 2 KB page at address and programs data (padded to halfwords).
// The CPU stalls while the flash is busy, so the display blanks for a moment.
func WriteFlashPage(address uintptr, data []byte) error {
	if len(data) > FlashPageSize {
		return ErrFlash
	}

	if stm32.FLASH.CR.HasBits(stm32.FLASH_CR_LOCK) {
		stm32.FLASH.KEYR.Set(0x45670123)
		stm32.FLASH.KEYR.Set(0xCDEF89AB)
	}
	defer stm32.FLASH.CR.SetBits(stm32.FLASH_CR_LOCK)

	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_PER)
	stm32.FLASH.AR.Set(uint32(address))
	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_STRT)
	err := flashWait()
	stm32.FLASH.CR.ClearBits(stm32.FLASH_CR_PER)
	if err != nil {
		return err
	}

	stm32.FLASH.CR.SetBits(stm32.FLASH_CR_PG)
	defer stm32.FLASH.CR.ClearBits(stm32.FLASH_CR_PG)
	for i := 0; i < len(data); i += 2 {
		half := uint16(data[i])
		if i+1 < len(data) {
			half |= uint16(data[i+1]) << 8
		}
		(*volatile.Register16)(unsafe.Pointer(address + uintptr(i))).Set(half)
		if err := flashWait(); err != nil {
			return err
		}
	}
	return nil
}
