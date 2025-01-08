package main

import (
	"github.com/Tariomka/stm32-led-cube/internal/runner"
	"github.com/Tariomka/stm32-led-cube/internal/show"
)

func main() {
	runner := runner.NewRunner(runner.NewConfig())
	runner.Start(show.NewLedShowList())
}
