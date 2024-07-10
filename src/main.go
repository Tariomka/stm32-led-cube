package main

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/src/controller"
)

func main() {
	board := controller.NewYellowBoard()

	// go lightSwitch(board.LedGreen.Pin)
	// go secondLightSwitch(board.LedRed.Pin)
	for {
		board.LedGreen.Pin.High()
		time.Sleep(1 * time.Second)

		board.LedRed.Pin.High()
		time.Sleep(1 * time.Second)

		board.LedGreen.Pin.Low()
		time.Sleep(1 * time.Second)

		board.LedRed.Pin.Low()
		time.Sleep(1 * time.Second)
	}
}

// func lightSwitch(pin machine.Pin) {
// 	for {
// 		pin.High()
// 		time.Sleep(500 * time.Millisecond)

// 		pin.Low()
// 		time.Sleep(500 * time.Millisecond)

// 		runtime.Gosched()
// 	}
// }

// func secondLightSwitch(pin machine.Pin) {
// 	for {
// 		pin.High()
// 		time.Sleep(500 * time.Millisecond)

// 		runtime.Gosched()

// 		pin.Low()
// 		time.Sleep(500 * time.Millisecond)
// 	}
// }
