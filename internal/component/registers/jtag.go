package registers

import (
	"device/stm32"
	"time"
)

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

func DisableJTAG() {
	stm32.RCC.SetAPB2ENR_AFIOEN(stm32.RCC_APB2ENR_AFIOEN) // Enable clock for AFIO peripheral
	time.Sleep(1 * time.Microsecond)
	stm32.AFIO.SetMAPR_SWJ_CFG(disableJTAG_SW) // JTAG-DP and SW-DP disabled
}
