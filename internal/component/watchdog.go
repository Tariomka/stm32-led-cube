package component

import (
	"device/stm32"
	"machine"
)

// TinyGo assumes a 32.768 kHz watchdog clock; the F1 LSI runs at about 40 kHz, so 4 s becomes ~3.2 s.
const watchdogMillis = 4000

// StartWatchdog starts the independent watchdog (like the factory firmware's ~4 s one) and
// reports whether it caused the last reset. Feed it with FeedWatchdog.
func StartWatchdog() (causedReset bool) {
	causedReset = stm32.RCC.CSR.HasBits(stm32.RCC_CSR_IWDGRSTF)
	stm32.RCC.CSR.SetBits(stm32.RCC_CSR_RMVF)
	machine.Watchdog.Configure(machine.WatchdogConfig{TimeoutMillis: watchdogMillis})
	machine.Watchdog.Start()
	return causedReset
}

func FeedWatchdog() { machine.Watchdog.Update() }
