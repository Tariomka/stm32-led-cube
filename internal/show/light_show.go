package show

import "github.com/Tariomka/led-common-lib/pkg/led"

func Demo() led.LightShow {
	firstFrame := func(lw led.LayoutWorker) {
		lw.SetLayer(0, led.Red)
		lw.SetLayer(1, led.Red)
		lw.SetLayer(2, led.Yellow)
		lw.SetLayer(3, led.Green)
		lw.SetLayer(4, led.Cyan)
		lw.SetLayer(5, led.Blue)
		lw.SetLayer(6, led.Violet)
		lw.SetLayer(7, led.White)
	}
	secondFrame := func(lw led.LayoutWorker) {
		lw.SetLayer(0, led.Red)
		lw.SetLayer(1, led.Yellow)
		lw.SetLayer(2, led.Yellow)
		lw.SetLayer(3, led.Green)
		lw.SetLayer(4, led.Cyan)
		lw.SetLayer(5, led.Blue)
		lw.SetLayer(6, led.Violet)
		lw.SetLayer(7, led.White)
	}

	return led.LightShow{firstFrame, secondFrame}
}

func Demo2() led.LightShow {
	firstFrame := func(lw led.LayoutWorker) {
		lw.ChangeLayer(0, led.Violet)
		lw.ChangeLayer(1, led.Violet)
		lw.ChangeLayer(2, led.Violet)
		lw.ChangeLayer(3, led.Violet)
		lw.ChangeLayer(4, led.Violet)
		lw.ChangeLayer(5, led.Violet)
		lw.ChangeLayer(6, led.Violet)
		lw.ChangeLayer(7, led.Violet)
	}
	secondFrame := func(lw led.LayoutWorker) {
		lw.ChangeLayer(0, led.Red)
		lw.ChangeLayer(1, led.Red)
		lw.ChangeLayer(2, led.Red)
		lw.ChangeLayer(3, led.Red)
		lw.ChangeLayer(4, led.Red)
		lw.ChangeLayer(5, led.Red)
		lw.ChangeLayer(6, led.Red)
		lw.ChangeLayer(7, led.Red)
	}

	return led.LightShow{firstFrame, secondFrame}
}

func DemoProgram() led.LightShow {
	return led.LightShow{func(lw led.LayoutWorker) {
		for y := uint8(0); y < 8; y++ {
			lw.SetRowIndividual(y, 0, led.Red, 0b11111111)
			lw.SetRowIndividual(y, 1, led.Red, 0b11111111)
			lw.SetRowIndividual(y, 2, led.Red, 0b11111111)
			lw.SetRowIndividual(y, 2, led.Green, 0b11111111)
			lw.SetRowIndividual(y, 3, led.Green, 0b11111111)
			lw.SetRowIndividual(y, 4, led.Green, 0b11111111)
			lw.SetRowIndividual(y, 4, led.Blue, 0b11111111)
			lw.SetRowIndividual(y, 5, led.Blue, 0b11111111)
			lw.SetRowIndividual(y, 6, led.Blue, 0b11111111)
			lw.SetRowIndividual(y, 6, led.Red, 0b11111111)
			lw.SetRowIndividual(y, 7, led.Green, 0b11111000)
			lw.SetRowIndividual(y, 7, led.Blue, 0b11100111)
			lw.SetRowIndividual(y, 7, led.Red, 0b00011111)
		}
	}}
}

func SingledLeds() led.LightShow {
	return led.LightShow{func(lw led.LayoutWorker) {
		lw.SetBlock(led.Red)
		lw.SetLayer(7, led.Blue)
		lw.SetRow(0, 3, led.Green)
	}}

}

func NewLedShowList() []led.LightShow {
	return []led.LightShow{
		Demo(),
		Demo2(),
		DemoProgram(),
		SingledLeds(),
	}
}
