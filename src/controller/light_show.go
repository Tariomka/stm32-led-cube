package controller

func Demo(b Board, lw LayoutWorker) {
	first := func() {
		lw.SetLayer(0, Red)
		lw.SetLayer(1, Red)
		lw.SetLayer(2, Red)
		lw.SetLayer(2, Green)
		lw.SetLayer(3, Green)
		lw.SetLayer(4, Green)
		lw.SetLayer(4, Blue)
		lw.SetLayer(5, Blue)
		lw.SetLayer(6, Blue)
		lw.SetLayer(6, Red)
		lw.SetLayer(7, Green)
		lw.SetLayer(7, Blue)
		lw.SetLayer(7, Red)
	}
	second := func() {
		lw.SetLayer(0, Red)
		lw.SetLayer(1, Red)
		lw.SetLayer(1, Green)
		lw.SetLayer(2, Red)
		lw.SetLayer(2, Green)
		lw.SetLayer(3, Green)
		lw.SetLayer(4, Green)
		lw.SetLayer(4, Blue)
		lw.SetLayer(5, Blue)
		lw.SetLayer(6, Blue)
		lw.SetLayer(6, Red)
		lw.SetLayer(7, Green)
		lw.SetLayer(7, Blue)
		lw.SetLayer(7, Red)
	}

	for {
		lw.ResetBlock()
		first()
		b.LightLeds(lw)

		lw.ResetBlock()
		second()
		b.LightLeds(lw)
	}
}

func DemoProgram(b Board, lw LayoutWorker) {
	lw.ResetBlock()
	for y := uint8(0); y < 8; y++ {
		lw.SetRowIndividual(y, 0, Red, 0b11111111)
		lw.SetRowIndividual(y, 1, Red, 0b11111111)
		lw.SetRowIndividual(y, 2, Red, 0b11111111)
		lw.SetRowIndividual(y, 2, Green, 0b11111111)
		lw.SetRowIndividual(y, 3, Green, 0b11111111)
		lw.SetRowIndividual(y, 4, Green, 0b11111111)
		lw.SetRowIndividual(y, 4, Blue, 0b11111111)
		lw.SetRowIndividual(y, 5, Blue, 0b11111111)
		lw.SetRowIndividual(y, 6, Blue, 0b11111111)
		lw.SetRowIndividual(y, 6, Red, 0b11111111)
		lw.SetRowIndividual(y, 7, Green, 0b11111000)
		lw.SetRowIndividual(y, 7, Blue, 0b11100111)
		lw.SetRowIndividual(y, 7, Red, 0b00011111)
	}

	for {
		b.LightLeds(lw)
	}
}

func SingledLeds(b Board, lw LayoutWorker) {
	lw.ResetBlock()
	lw.SetBlock(Red)
	lw.SetLayer(7, Blue)
	lw.SetRow(0, 3, Green)

	for {
		b.LightLeds(lw)
	}
}

func NewLedShowList() []func(Board, LayoutWorker) {
	return []func(Board, LayoutWorker){
		Demo,
		DemoProgram,
		SingledLeds,
	}
}
