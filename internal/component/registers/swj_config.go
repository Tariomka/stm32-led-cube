package registers

import (
	"device/stm32"
)

// TODO: Does this even work? Maybe just delet dis?

func UpdateRegisters() {
	DisableJTAG()
	CANRemap()
	stm32.RCC.SetAPB2ENR_AFIOEN(stm32.RCC_APB2ENR_AFIOEN_Disabled) // Disable clock for AFIO peripheral
}

// Default Pin assignments
// PA15	: JTDI
// PA14	: JTCK/SWCLK
// PA13	: JTMS/SWDIO
// PB4	: NJTRST
// PB3	: JTDO

// SWJ_CFG field (Serial wire JTAG configuration) in AFIO port configuration register 0 (AFIO_PCF0)
// It contains 3 bits -> [26-24]. SWJ_CFG position can be found 'stm32.AFIO_MAPR_SWJ_CFG_Pos'

// SWJ_CFG values

const (
	// resetState     = 0b000 // Full SWJ (JTAG-DP + SW-DP) (Reset state)
	// noNJTRST       = 0b001 // Full SWJ (JTAG-DP + SW-DP) but without NJTRST
	// disableJTAG    = 0b010 // JTAG-DP Disabled and SW-DP Enabled

	disableJTAG_SW = 0b100 // JTAG-DP Disabled and SW-DP Disabled
)

// DisableJTAG disables Serial Wire JTAG interface on the board. This function frees up these pins for GPIO usage:
//
// - PA13
//
// - PA14
//
// - PA15
//
// - PB3
//
// - PB4
func DisableJTAG() {
	stm32.RCC.SetAPB2ENR_AFIOEN(stm32.RCC_APB2ENR_AFIOEN_Enabled) // Enable clock for AFIO peripheral
	stm32.AFIO.SetMAPR_SWJ_CFG(disableJTAG_SW)                    // JTAG-DP and SW-DP disabled

	// stm32.GPIOA.SetCRH_MODE14(stm32.GPIO_CRH_MODE14_Input) // Set PA14 to Input mode
	// stm32.GPIOA.SetCRH_MODE15(stm32.GPIO_CRH_MODE15_Input) // Set PA15 to Input mode
}

// Default Pin assignments
// PA11	: CAN_RX
// PA12	: CAN_TX

// CAN_REMAP values

const (
	// canPA      = 0b00 // RX -> PA11, TX -> PA12

	canNotUsed = 0b01
	// canPB      = 0b10 // RX -> PB8, TX -> PB9
	// canPD      = 0b11 // RX -> PD0, TX -> PD1
)

func CANRemap() {
	stm32.RCC.SetAPB2ENR_AFIOEN(stm32.RCC_APB2ENR_AFIOEN_Enabled) // Enable clock for AFIO peripheral
	stm32.AFIO.SetMAPR_CAN_REMAP(canNotUsed)                      // Remap CAN pins
}
