package registers

import (
	"device/stm32"
	"machine"
	"strconv"
)

func PrintAndResetCrashLog() uint32 {
	println("=================================")
	println("------- Reset cause flags -------")
	println("=================================")

	resetFlags := stm32.RCC.CSR.Get()
	println(" - RCC CSR value:                  ", strconv.FormatUint(uint64(resetFlags), 2))
	println(" - Low-power reset flag:           ", stm32.RCC.GetCSR_LPWRRSTF())
	println(" - Window watchdog reset flag:     ", stm32.RCC.GetCSR_WWDGRSTF())
	println(" - Independent watchdog reset flag:", stm32.RCC.GetCSR_IWDGRSTF())
	println(" - Software reset flag:            ", stm32.RCC.GetCSR_SFTRSTF())
	println(" - POR/PDR reset flag:             ", stm32.RCC.GetCSR_PORRSTF())
	println(" - PIN reset flag:                 ", stm32.RCC.GetCSR_PINRSTF())
	println("=================================")

	stm32.RCC.SetCSR_RMVF(0x1)

	return resetFlags
}

// TODO: Remove eveerything below, when pins start working properly
func PrintButtonStates() {
	println("=================================")
	println("--------- Button states ---------")
	println("=================================")
	println(" - Button Previous:  ", !machine.PC0.Get())
	println(" - Button Next:      ", !machine.PC1.Get())
	println(" - Button Speed More:", !machine.PC2.Get())
	println(" - Button Speed Less:", !machine.PB3.Get())
	println(" - Button Run/Pause: ", !machine.PA14.Get())
	println(" - Button Cycle:     ", !machine.PA13.Get())
	println(" - Button On/Off:    ", !machine.PA11.Get())
	println("---------------------------------")
	println(" - Switch 1:", machine.PA1.Get())
	println(" - Switch 2:", machine.PC3.Get())
	println("=================================")
}

func PrintCheck() {
	println("=================================")
	println(" RCC APB2ENR AFIOEN value:", strconv.FormatUint(uint64(stm32.RCC.GetAPB2ENR_AFIOEN()), 2))
	println("---------------------------------")
	println(" AFIO MAPR SWJ_CFG value:", strconv.FormatUint(uint64(stm32.AFIO.GetMAPR_SWJ_CFG()), 2))
	println("---------------------------------")
	mode := strconv.FormatUint(uint64(stm32.GPIOA.GetCRH_MODE14()), 2)
	cnf := strconv.FormatUint(uint64(stm32.GPIOA.GetCRH_CNF14()), 2)
	println(" GPIOA Pin PA14 MODE value:", mode)
	println(" GPIOA Pin PA14 CNF value: ", cnf)
	println(" Pin PA14 value:", cnf+mode)
	println("=================================")
}
