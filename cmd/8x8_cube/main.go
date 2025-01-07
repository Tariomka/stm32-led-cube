package main

import (
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/show"
)

func main() {
	leds := controller.LedLayout{}
	board := controller.NewYellowBoard()
	shows := show.NewLedShowList()

	board.BlinkStartup()
	board.Run(&leds, shows)
}
