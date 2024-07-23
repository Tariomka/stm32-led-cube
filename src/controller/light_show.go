package controller

import "time"

func DemoProgram(yb *YellowBoard, ll LedLayout) {
	yb.LedDriver.ClearLayer()
	first := func() {
		for y := uint8(0); y < 8; y++ {
			ll.LedRow(y, 0, Red, 255)
			ll.LedRow(y, 1, Red, 255)
			ll.LedRow(y, 2, Red, 255)
			ll.LedRow(y, 2, Green, 255)
			ll.LedRow(y, 3, Green, 255)
			ll.LedRow(y, 4, Green, 255)
			ll.LedRow(y, 4, Blue, 255)
			ll.LedRow(y, 5, Blue, 255)
			ll.LedRow(y, 6, Blue, 255)
			ll.LedRow(y, 6, Red, 255)
			ll.LedRow(y, 7, Green, 255)
			ll.LedRow(y, 7, Blue, 255)
			ll.LedRow(y, 7, Red, 255)
		}
	}
	second := func() {
		for y := uint8(0); y < 8; y++ {
			ll.LedRow(y, 0, Red, 255)
			ll.LedRow(y, 1, Red, 255)
			ll.LedRow(y, 1, Green, 255)
			ll.LedRow(y, 2, Red, 255)
			ll.LedRow(y, 2, Green, 255)
			ll.LedRow(y, 3, Green, 255)
			ll.LedRow(y, 4, Green, 255)
			ll.LedRow(y, 4, Blue, 255)
			ll.LedRow(y, 5, Blue, 255)
			ll.LedRow(y, 6, Blue, 255)
			ll.LedRow(y, 6, Red, 255)
			ll.LedRow(y, 7, Green, 255)
			ll.LedRow(y, 7, Blue, 255)
			ll.LedRow(y, 7, Red, 255)
		}
	}
	for {
		first()
		yb.LightLeds(ll)
		second()
		yb.LightLeds(ll)
	}
}

func RedGreenOnBoard(yb *YellowBoard, ll LedLayout) {
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

func StaticLights(yb *YellowBoard, ll LedLayout) {
	yb.Demultiplexer.MultiEnable.Pin.Low()
	yb.Demultiplexer.EnableLayer(5)
	yb.LedDriver.LightLayer([]byte{
		255, 255,
	})
	yb.Demultiplexer.MultiEnable.Pin.High()
}

func NewLedShowList() []func(*YellowBoard, LedLayout) {
	return []func(*YellowBoard, LedLayout){
		DemoProgram,
		RedGreenOnBoard,
		StaticLights,
	}
}
