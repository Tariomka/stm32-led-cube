package controller

import "time"

func Demo(yb *YellowBoard, ll *LedLayout) {
	first := func() {
		ll.LedLayer(0, Red)
		ll.LedLayer(1, Red)
		ll.LedLayer(2, Red)
		ll.LedLayer(2, Green)
		ll.LedLayer(3, Green)
		ll.LedLayer(4, Green)
		ll.LedLayer(4, Blue)
		ll.LedLayer(5, Blue)
		ll.LedLayer(6, Blue)
		ll.LedLayer(6, Red)
		ll.LedLayer(7, Green)
		ll.LedLayer(7, Blue)
		ll.LedLayer(7, Red)
	}
	second := func() {
		ll.LedLayer(0, Red)
		ll.LedLayer(1, Red)
		ll.LedLayer(1, Green)
		ll.LedLayer(2, Red)
		ll.LedLayer(2, Green)
		ll.LedLayer(3, Green)
		ll.LedLayer(4, Green)
		ll.LedLayer(4, Blue)
		ll.LedLayer(5, Blue)
		ll.LedLayer(6, Blue)
		ll.LedLayer(6, Red)
		ll.LedLayer(7, Green)
		ll.LedLayer(7, Blue)
		ll.LedLayer(7, Red)
	}

	for {
		ll.LedBlockOff()
		first()
		yb.LightLeds(*ll)

		ll.LedBlockOff()
		second()
		yb.LightLeds(*ll)
	}
}

func DemoProgram(yb *YellowBoard, ll *LedLayout) {
	ll.LedBlockOff()
	for y := uint8(0); y < 8; y++ {
		ll.LedRowIndividual(y, 0, Red, 0b11111111)
		ll.LedRowIndividual(y, 1, Red, 0b11111111)
		ll.LedRowIndividual(y, 2, Red, 0b11111111)
		ll.LedRowIndividual(y, 2, Green, 0b11111111)
		ll.LedRowIndividual(y, 3, Green, 0b11111111)
		ll.LedRowIndividual(y, 4, Green, 0b11111111)
		ll.LedRowIndividual(y, 4, Blue, 0b11111111)
		ll.LedRowIndividual(y, 5, Blue, 0b11111111)
		ll.LedRowIndividual(y, 6, Blue, 0b11111111)
		ll.LedRowIndividual(y, 6, Red, 0b11111111)
		ll.LedRowIndividual(y, 7, Green, 0b11111000)
		ll.LedRowIndividual(y, 7, Blue, 0b11100111)
		ll.LedRowIndividual(y, 7, Red, 0b00011111)
	}

	for {
		yb.LightLeds(*ll)
	}
}

func SingledLeds(yb *YellowBoard, ll *LedLayout) {
	ll.LedBlockOff()
	ll.LedBlock(Red)
	ll.LedLayer(7, Blue)
	ll.LedRow(0, 3, Green)

	for {
		yb.LightLeds(*ll)
	}
}

func RedGreenOnBoard(yb *YellowBoard, ll *LedLayout) {
	for {
		yb.LedGreen.Pin.High()
		time.Sleep(500 * time.Millisecond)

		yb.LedRed.Pin.High()
		time.Sleep(250 * time.Millisecond)

		yb.LedGreen.Pin.Low()
		time.Sleep(500 * time.Millisecond)

		yb.LedRed.Pin.Low()
		time.Sleep(250 * time.Millisecond)
	}
}

func NewLedShowList() []func(*YellowBoard, *LedLayout) {
	return []func(*YellowBoard, *LedLayout){
		Demo,
		DemoProgram,
		SingledLeds,
		RedGreenOnBoard,
	}
}
