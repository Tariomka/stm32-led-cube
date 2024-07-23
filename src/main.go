package main

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/src/controller"
)

func main() {
	board := controller.NewYellowBoard()
	leds := controller.LedLayout{}
	shows := controller.NewLedShowList()

	for i := 0; i < 3; i++ {
		board.LedRed.Pin.Low()
		board.LedGreen.Pin.Low()
		time.Sleep(100 * time.Millisecond)

		board.LedRed.Pin.High()
		board.LedGreen.Pin.High()
		time.Sleep(100 * time.Millisecond)
	}

	for {
		// board.Demultiplexer.EnableLayer(5)
		// board.LedDriver.LightLayer([]byte{
		// 	255,
		// })
		shows[0](board, leds)
	}
}

// func main() {
// 	red := machine.PB9
// 	red.Configure(machine.PinConfig{Mode: machine.PinOutput})

// 	green := machine.PB8
// 	green.Configure(machine.PinConfig{Mode: machine.PinOutput})

// 	a0 := machine.PB0
// 	a0.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	a1 := machine.PB1
// 	a1.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	a2 := machine.PB10
// 	a2.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	a3 := machine.PB11
// 	a3.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	aEnable := machine.PA8
// 	aEnable.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	aEnable.High()

// 	spi := machine.SPI0
// 	spi.Configure(machine.SPIConfig{
// 		Frequency: 100_000,
// 		SCK:       machine.SPI0_SCK_PIN,
// 		SDO:       machine.SPI0_SDO_PIN,
// 		LSBFirst:  true,
// 		Mode:      1,
// 	})
// 	latch := machine.PC4
// 	latch.Configure(machine.PinConfig{Mode: machine.PinOutput})
// 	blank := machine.PC5
// 	blank.Configure(machine.PinConfig{Mode: machine.PinOutput})

// 	for i := 0; i < 3; i++ {
// 		red.Low()
// 		green.Low()
// 		time.Sleep(100 * time.Millisecond)

// 		red.High()
// 		green.High()
// 		time.Sleep(100 * time.Millisecond)
// 	}

// 	for {
// 		green.Low()
// 	}
// }
