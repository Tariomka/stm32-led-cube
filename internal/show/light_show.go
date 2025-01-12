package show

import "github.com/Tariomka/stm32-led-cube/internal/controller"

func Demo() controller.LightShow {
	firstFrame := func(lw controller.LayoutWorker) {
		lw.SetLayer(0, controller.Red)
		lw.SetLayer(1, controller.Red)
		lw.SetLayer(2, controller.Yellow)
		lw.SetLayer(3, controller.Green)
		lw.SetLayer(4, controller.Cyan)
		lw.SetLayer(5, controller.Blue)
		lw.SetLayer(6, controller.Violet)
		lw.SetLayer(7, controller.White)
	}
	secondFrame := func(lw controller.LayoutWorker) {
		lw.SetLayer(0, controller.Red)
		lw.SetLayer(1, controller.Yellow)
		lw.SetLayer(2, controller.Yellow)
		lw.SetLayer(3, controller.Green)
		lw.SetLayer(4, controller.Cyan)
		lw.SetLayer(5, controller.Blue)
		lw.SetLayer(6, controller.Violet)
		lw.SetLayer(7, controller.White)
	}

	return controller.LightShow{firstFrame, secondFrame}
}

func Demo2() controller.LightShow {
	firstFrame := func(lw controller.LayoutWorker) {
		lw.ChangeLayer(0, controller.Violet)
		lw.ChangeLayer(1, controller.Violet)
		lw.ChangeLayer(2, controller.Violet)
		lw.ChangeLayer(3, controller.Violet)
		lw.ChangeLayer(4, controller.Violet)
		lw.ChangeLayer(5, controller.Violet)
		lw.ChangeLayer(6, controller.Violet)
		lw.ChangeLayer(7, controller.Violet)
	}
	secondFrame := func(lw controller.LayoutWorker) {
		lw.ChangeLayer(0, controller.Red)
		lw.ChangeLayer(1, controller.Red)
		lw.ChangeLayer(2, controller.Red)
		lw.ChangeLayer(3, controller.Red)
		lw.ChangeLayer(4, controller.Red)
		lw.ChangeLayer(5, controller.Red)
		lw.ChangeLayer(6, controller.Red)
		lw.ChangeLayer(7, controller.Red)
	}

	return controller.LightShow{firstFrame, secondFrame}
}

func DemoProgram() controller.LightShow {
	return controller.LightShow{func(lw controller.LayoutWorker) {
		for y := uint8(0); y < 8; y++ {
			lw.SetRowIndividual(y, 0, controller.Red, 0b11111111)
			lw.SetRowIndividual(y, 1, controller.Red, 0b11111111)
			lw.SetRowIndividual(y, 2, controller.Red, 0b11111111)
			lw.SetRowIndividual(y, 2, controller.Green, 0b11111111)
			lw.SetRowIndividual(y, 3, controller.Green, 0b11111111)
			lw.SetRowIndividual(y, 4, controller.Green, 0b11111111)
			lw.SetRowIndividual(y, 4, controller.Blue, 0b11111111)
			lw.SetRowIndividual(y, 5, controller.Blue, 0b11111111)
			lw.SetRowIndividual(y, 6, controller.Blue, 0b11111111)
			lw.SetRowIndividual(y, 6, controller.Red, 0b11111111)
			lw.SetRowIndividual(y, 7, controller.Green, 0b11111000)
			lw.SetRowIndividual(y, 7, controller.Blue, 0b11100111)
			lw.SetRowIndividual(y, 7, controller.Red, 0b00011111)
		}
	}}
}

func SingledLeds() controller.LightShow {
	return controller.LightShow{func(lw controller.LayoutWorker) {
		lw.SetBlock(controller.Red)
		lw.SetLayer(7, controller.Blue)
		lw.SetRow(0, 3, controller.Green)
	}}

}

func NewLedShowList() []controller.LightShow {
	return []controller.LightShow{
		Demo(),
		Demo2(),
		DemoProgram(),
		SingledLeds(),
	}
}
