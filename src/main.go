package main

import (
	"github.com/Tariomka/stm32-led-cube/src/controller"
)

func main() {
	leds := controller.LedLayout{}
	board := controller.NewYellowBoard()
	shows := controller.NewLedShowList()

	board.BlinkStartup()
	board.Run(&leds, shows)
}
