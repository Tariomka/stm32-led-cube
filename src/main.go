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
		shows[5](board, leds)

		board.LedGreen.Pin.Low()
		time.Sleep(500 * time.Millisecond)

		board.LedGreen.Pin.High()
		time.Sleep(500 * time.Millisecond)
	}
}
