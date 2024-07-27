package main

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/src/controller"
)

func main() {
	leds := controller.LedLayout{}
	board := controller.NewYellowBoard()
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
		shows[2](board, &leds)
		// controller.DemoProgram(board, leds)
	}
}
