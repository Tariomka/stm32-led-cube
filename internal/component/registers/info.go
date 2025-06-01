package registers

import "device/stm32"

func PrintAndResetCrashLog() {
	println("=================================")
	println("---- Reset and clock control ----")
	println("=================================")
	println("RCC CSR value: ", stm32.RCC.CSR.Get())
	println("Low-power reset flag: ", stm32.RCC.GetCSR_LPWRRSTF())
	println("Window watchdog reset flag: ", stm32.RCC.GetCSR_WWDGRSTF())
	println("Independent watchdog reset flag: ", stm32.RCC.GetCSR_IWDGRSTF())
	println("Software reset flag: ", stm32.RCC.GetCSR_SFTRSTF())
	println("POR/PDR reset flag: ", stm32.RCC.GetCSR_PORRSTF())
	println("PIN reset flag: ", stm32.RCC.GetCSR_PINRSTF())
	println("=================================")

	stm32.RCC.SetCSR_RMVF(0x1)
}

func PrintPVDLog() {
	println("=================================")
	println("- Programmable voltage detector -")
	println("=================================")
	print("PVD level selection value: (", stm32.PWR.GetCR_PLS(), ") ")
	switch stm32.PWR.GetCR_PLS() {
	case 0b000:
		println("2.2V")
	case 0b001:
		println("2.3V")
	case 0b010:
		println("2.4V")
	case 0b011:
		println("2.5V")
	case 0b100:
		println("2.6V")
	case 0b101:
		println("2.7V")
	case 0b110:
		println("2.8V")
	case 0b111:
		println("2.9V")
	}
	println("PVD enable value: ", stm32.PWR.GetCR_PVDE())
	println("PVD output value: ", stm32.PWR.GetCSR_PVDO())
	println("=================================")
}
