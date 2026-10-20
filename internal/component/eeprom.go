package component

import (
	"machine"
	"time"
)

const (
	eepromAddress = 0x50
	pageSizeB     = 8 // Page size in bytes
)

// EEPROM is the 24C02 (256 bytes, 8-byte pages) on I2C1: SCL = PB6, SDA = PB7.
type EEPROM struct{ bus *machine.I2C }

func NewEEPROM() EEPROM {
	return EEPROM{bus: NewOnBoardI2C()}
}

func (this EEPROM) Read(address byte, buffer []byte) error {
	return this.retry(func() error { return this.bus.ReadRegister(eepromAddress, address, buffer) })
}

// Write stores data page by page and waits for each internal write cycle (max 5 ms).
func (this EEPROM) Write(address byte, data []byte) error {
	var page [pageSizeB + 1]byte
	for len(data) > 0 {
		n := min(len(data), pageSizeB-int(address)%pageSizeB)
		page[0] = address
		copy(page[1:], data[:n])
		if err := this.retry(func() error { return this.bus.Tx(eepromAddress, page[:n+1], nil) }); err != nil {
			return err
		}

		if err := this.waitWriteCycle(); err != nil {
			return err
		}

		address += byte(n)
		data = data[n:]
	}
	return nil
}

// retry repeats a failed transfer once after resetting the I2C peripheral (a stuck BUSY flag or a
// transfer disturbed by the display interrupt).
func (this EEPROM) retry(transfer func() error) error {
	err := transfer()
	if err != nil {
		NewOnBoardI2C()
		err = transfer()
	}
	return err
}

// waitWriteCycle polls until the chip acknowledges its address again.
func (this EEPROM) waitWriteCycle() error {
	deadline := time.Now().Add(10 * time.Millisecond)
	for {
		err := this.bus.Tx(eepromAddress, nil, nil)
		if err == nil || time.Now().After(deadline) {
			return err
		}
	}
}
