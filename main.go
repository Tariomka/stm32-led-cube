package main

import (
	"time"

	"github.com/Tariomka/stm32-led-cube/src/controller"
)

func main() {
	board := controller.NewYellowBoard()

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
