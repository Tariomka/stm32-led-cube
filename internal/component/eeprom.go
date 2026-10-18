package component

import (
	"device/arm"
	"device/stm32"
	"errors"
	"time"
)

// EEPROM is the 24C02 (256 bytes, 8-byte pages) on a bit-banged I2C bus like the original
// firmware: SCL = PB6, SDA = PB7, open-drain outputs with external pull-ups.
type EEPROM struct{}

const (
	eepromAddress = 0xA0
	eepromPage    = 8
	sclBit        = 1 << 6
	sdaBit        = 1 << 7
)

var ErrEEPROM = errors.New("eeprom not responding")

func NewEEPROM() EEPROM {
	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_IOPBEN)
	stm32.GPIOB.BSRR.Set(sclBit | sdaBit)
	stm32.GPIOB.CRL.ReplaceBits(0x77, 0xff, 24) // PB6, PB7: open-drain output, 50 MHz
	return EEPROM{}
}

func i2cDelay() {
	for i := 0; i < 60; i++ { // about 5 µs at 72 MHz
		arm.Asm("nop")
	}
}

func setLine(bit uint32, high bool) {
	if high {
		stm32.GPIOB.BSRR.Set(bit)
	} else {
		stm32.GPIOB.BRR.Set(bit)
	}
	i2cDelay()
}

func sda() bool { return stm32.GPIOB.IDR.HasBits(sdaBit) }

func i2cStart() {
	setLine(sdaBit, true)
	setLine(sclBit, true)
	setLine(sdaBit, false)
	setLine(sclBit, false)
}

func i2cStop() {
	setLine(sdaBit, false)
	setLine(sclBit, true)
	setLine(sdaBit, true)
}

func i2cWrite(b byte) bool {
	for i := 0; i < 8; i++ {
		setLine(sdaBit, b&0x80 != 0)
		setLine(sclBit, true)
		setLine(sclBit, false)
		b <<= 1
	}
	setLine(sdaBit, true)
	setLine(sclBit, true)
	ack := !sda()
	setLine(sclBit, false)
	return ack
}

func i2cRead(ack bool) byte {
	var b byte
	setLine(sdaBit, true)
	for i := 0; i < 8; i++ {
		setLine(sclBit, true)
		b <<= 1
		if sda() {
			b |= 1
		}
		setLine(sclBit, false)
	}
	setLine(sdaBit, !ack)
	setLine(sclBit, true)
	setLine(sclBit, false)
	setLine(sdaBit, true)
	return b
}

func (EEPROM) Read(address byte, buffer []byte) error {
	defer i2cStop()
	i2cStart()
	if !i2cWrite(eepromAddress) || !i2cWrite(address) {
		return ErrEEPROM
	}

	i2cStart()
	if !i2cWrite(eepromAddress | 1) {
		return ErrEEPROM
	}

	for i := range buffer {
		buffer[i] = i2cRead(i < len(buffer)-1)
	}
	return nil
}

// Write stores data page by page and waits for each internal write cycle (max 5 ms).
func (EEPROM) Write(address byte, data []byte) error {
	for len(data) > 0 {
		n := min(len(data), eepromPage-int(address)%eepromPage)
		i2cStart()
		ok := i2cWrite(eepromAddress) && i2cWrite(address)
		for i := 0; ok && i < n; i++ {
			ok = i2cWrite(data[i])
		}
		i2cStop()
		if !ok || !waitWriteCycle() {
			return ErrEEPROM
		}
		address += byte(n)
		data = data[n:]
	}
	return nil
}

func waitWriteCycle() bool {
	deadline := time.Now().Add(10 * time.Millisecond)
	for time.Now().Before(deadline) {
		i2cStart()
		ack := i2cWrite(eepromAddress)
		i2cStop()
		if ack {
			return true
		}
	}
	return false
}
