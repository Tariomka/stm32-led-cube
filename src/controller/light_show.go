package controller

import "time"

func InitialProgram(yb *YellowBoard, ll LedLayout) {
	ll.LedColor(1, 1, 7, Red)
	ll.LedColor(2, 3, 4, Green)
	yb.LightLeds(ll)
}

func DemoProgram(yb *YellowBoard, ll LedLayout) {
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 0, Red, 255)
		ll.LedRow(uint8(y), 1, Red, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 2, Red, 255)
		ll.LedRow(uint8(y), 2, Green, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 3, Green, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 4, Green, 255)
		ll.LedRow(uint8(y), 4, Blue, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 5, Blue, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 6, Blue, 255)
		ll.LedRow(uint8(y), 6, Red, 255)
	}
	for y := 0; y < 8; y++ {
		ll.LedRow(uint8(y), 7, Green, 255)
		ll.LedRow(uint8(y), 7, Blue, 255)
		ll.LedRow(uint8(y), 7, Red, 255)
	}

	for {
		yb.LedGreen.Pin.Low()
		yb.LightLeds(ll)
		yb.LedGreen.Pin.High()
		time.Sleep(100 * time.Millisecond)
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

func TrialAndError(yb *YellowBoard, ll LedLayout) {
	// Fucking try everything

	// Setting manually shit
	// yb.LedGreen.Pin.Low()
	// time.Sleep(50 * time.Millisecond)

	// yb.Demultiplexer.MultiEnable.Pin.High()
	// yb.Demultiplexer.MultiA0.Pin.High()
	// yb.Demultiplexer.MultiA1.Pin.Low()
	// yb.Demultiplexer.MultiA2.Pin.Low()

	// yb.LedDriver.blank.Pin.Low()
	// yb.LedDriver.latch.Pin.High()
	// yb.LedDriver.latch.Pin.Low()
	// yb.LedDriver.spi.Tx([]byte{255, 255}, nil)
	// yb.LedDriver.latch.Pin.High()

	// yb.LedGreen.Pin.High()

	//
	yb.LedGreen.Pin.Low()
	time.Sleep(50 * time.Millisecond)

	yb.Demultiplexer.EnableLayer(0)
	yb.LedDriver.ClearLayer()

	time.Sleep(100 * time.Millisecond)
	yb.LedDriver.LightLayer([]byte{255})
	time.Sleep(100 * time.Millisecond)
	yb.LedDriver.LightLayer([]byte{255, 255})
	time.Sleep(100 * time.Millisecond)
	yb.LedDriver.LightLayer([]byte{255, 255, 255})
	time.Sleep(100 * time.Millisecond)
	yb.LedDriver.LightLayer([]byte{
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
	})

	yb.LedGreen.Pin.High()
	time.Sleep(100 * time.Millisecond)
}

func GreenOnBoard(yb *YellowBoard, ll LedLayout) {
	yb.LedRed.Pin.High()
	yb.LedGreen.Pin.Low()

	// for {
	// 	yb.LedGreen.Pin.High()
	// 	time.Sleep(500 * time.Millisecond)

	// 	time.Sleep(500 * time.Millisecond)
	// }
}

func StaticLights(yb *YellowBoard, ll LedLayout) {
	yb.LedDriver.ClearLayer()
	yb.Demultiplexer.EnableLayer(7)

	prev := 23
	curr := 0
	data := make([]byte, 24)
	for {
		if curr > 23 {
			curr = 0
		}
		if prev > 23 {
			prev = 0
		}

		data[prev] = 0
		data[curr] = 255
		yb.LedDriver.LightLayer(data)
		prev++
		curr++

		time.Sleep(150 * time.Millisecond)
	}

	// yb.LedDriver.LightLayer([]byte{
	// 	255, 255, 255, 255, 255, 255, 255, 255,
	// 	255, 255, 255, 255, 255, 255, 255, 255,
	// 	255, 255, 255, 255, 255, 255, 255, 255,
	// })
	// yb.LedDriver.LightLayer([]byte{
	// 	0, 0, 0, 0, 255, 255,
	// 	0, 0, 0, 0, 255, 255,
	// 	0, 0, 0, 0, 255, 255,
	// 	0, 0, 0, 0, 255, 255,
	// })
	// yb.LedDriver.LightLayer([]byte{
	// 	255, 0, 0, 0, 255, 0,
	// 	0, 255, 0, 0, 0, 255,
	// 	0, 0, 255, 0, 0, 0,
	// 	0, 0, 0, 255, 0, 0,
	// })
}

func Empty(yb *YellowBoard, ll LedLayout) {

}

func NewLedShowList() []func(*YellowBoard, LedLayout) {
	return []func(*YellowBoard, LedLayout){
		InitialProgram,
		DemoProgram,
		RedGreenOnBoard,
		TrialAndError,
		GreenOnBoard,
		StaticLights,
		Empty,
	}
}
