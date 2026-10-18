// gpio_diag verifies the KEY5 (PA14) fix on hardware. Watch the output with scripts\monitor.cmd.
package main

import (
	"device/stm32"
	"machine"
	"time"

	"github.com/Tariomka/stm32-led-cube/internal/component"
)

type namedPin struct {
	name string
	pin  machine.Pin
}

var inputs = []namedPin{
	{"KEY1 PC0", machine.PC0}, {"KEY2 PC1", machine.PC1}, {"KEY3 PC2", machine.PC2}, {"KEY4 PB3", machine.PB3},
	{"KEY5 PA14", machine.PA14}, {"KEY6 PA13", machine.PA13}, {"KEY7 PA11", machine.PA11},
	{"SW1 PA1", machine.PA1}, {"SW2 PC3", machine.PC3}, {"SW3 PA15", machine.PA15}, {"SW4 PA12", machine.PA12},
}

func main() {
	time.Sleep(2 * time.Second)
	println("\r\n==== gpio_diag ====")
	dump("1. reset state")

	machine.PA14.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.PA13.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dump("2. TinyGo PinInputPullup on PA14 then PA13 (expect PA14=1001, PA15=1001: the bug)")

	for _, in := range inputs {
		component.NewInputPin(in.pin)
	}
	dump("3. workaround on all inputs, SWJ still enabled (expect all inputs 1000)")
	println("   press KEY5 / KEY6 within 5 s (KEY5 stuck low by the SWCLK pull-down, KEY6 works via the SWDIO pull-up)")
	watch(5 * time.Second)

	component.DisableDebugPort()
	dump("4. SWJ disabled (expect all inputs 1000, MAPR SWJ_CFG=100)")
	println("   press any key / flip any switch")
	watch(0)
}

func watch(d time.Duration) {
	var last uint32 = 0xFFFFFFFF
	start := time.Now()
	for d == 0 || time.Since(start) < d {
		var state uint32
		for i, in := range inputs {
			if !in.pin.Get() {
				state |= 1 << i
			}
		}
		if state != last {
			last = state
			print("   low (key pressed / switch ON):")
			for i, in := range inputs {
				if state&(1<<i) != 0 {
					print(" [", in.name, "]")
				}
			}
			println()
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func dump(title string) {
	println(title)
	println("   AFIO MAPR SWJ_CFG:", stm32.AFIO.GetMAPR_SWJ_CFG())
	port("A", stm32.GPIOA)
	port("B", stm32.GPIOB)
	port("C", stm32.GPIOC)
}

func port(name string, gpio *stm32.GPIO_Type) {
	crl, crh := gpio.CRL.Get(), gpio.CRH.Get()
	print("   P", name, " CNF|MODE 15..0:")
	for pin := 15; pin >= 0; pin-- {
		reg, shift := crh, uint32(pin-8)*4
		if pin < 8 {
			reg, shift = crl, uint32(pin)*4
		}
		nib := (reg >> shift) & 0xF
		print(" ", (nib>>3)&1, (nib>>2)&1, (nib>>1)&1, nib&1)
	}
	println("  IDR:", gpio.IDR.Get(), "ODR:", gpio.ODR.Get())
}
