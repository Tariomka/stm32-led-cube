package controller

func Demo() LightShow {
	firstFrame := func(lw LayoutWorker) {
		lw.SetLayer(0, Red)
		lw.SetLayer(1, Red)
		lw.SetLayer(2, Yellow)
		lw.SetLayer(3, Green)
		lw.SetLayer(4, Cyan)
		lw.SetLayer(5, Blue)
		lw.SetLayer(6, Violet)
		lw.SetLayer(7, White)
	}
	secondFrame := func(lw LayoutWorker) {
		lw.SetLayer(0, Red)
		lw.SetLayer(1, Yellow)
		lw.SetLayer(2, Yellow)
		lw.SetLayer(3, Green)
		lw.SetLayer(4, Cyan)
		lw.SetLayer(5, Blue)
		lw.SetLayer(6, Violet)
		lw.SetLayer(7, White)
	}

	return LightShow{firstFrame, secondFrame}
}

func Demo2() LightShow {
	firstFrame := func(lw LayoutWorker) {
		lw.ChangeLayer(0, Violet)
		lw.ChangeLayer(1, Violet)
		lw.ChangeLayer(2, Violet)
		lw.ChangeLayer(3, Violet)
		lw.ChangeLayer(4, Violet)
		lw.ChangeLayer(5, Violet)
		lw.ChangeLayer(6, Violet)
		lw.ChangeLayer(7, Violet)
	}
	secondFrame := func(lw LayoutWorker) {
		lw.ChangeLayer(0, Red)
		lw.ChangeLayer(1, Red)
		lw.ChangeLayer(2, Red)
		lw.ChangeLayer(3, Red)
		lw.ChangeLayer(4, Red)
		lw.ChangeLayer(5, Red)
		lw.ChangeLayer(6, Red)
		lw.ChangeLayer(7, Red)
	}

	return LightShow{firstFrame, secondFrame}
}

func DemoProgram() LightShow {
	return LightShow{func(lw LayoutWorker) {
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
	}}
}

func SingledLeds() LightShow {
	return LightShow{func(lw LayoutWorker) {
		lw.SetBlock(Red)
		lw.SetLayer(7, Blue)
		lw.SetRow(0, 3, Green)
	}}

}

func NewLedShowList() []LightShow {
	return []LightShow{
		Demo2(),
		Demo(),
		DemoProgram(),
		SingledLeds(),
	}
}
