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
		shows[3](board, leds)
	}
}
