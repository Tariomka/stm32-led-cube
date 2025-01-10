package main

import "github.com/Tariomka/stm32-led-cube/internal/runner"

func main() {
	if runner := runner.NewRunner(runner.NewConfig()); runner != nil {
		runner.Start()
	}
}
