// Generated from the original iCubeSmart show data by tools/showgen.

package show

import (
	"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

// test01X8x1 is cartoon 1 of the original Test group.
func test01X8x1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 24, Loop: 1, Delay: 400,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "8x8x1", Palette: "red", To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 8, End: 15, Shape: "8x8x1", Palette: "red", From: keyframe.Pose{X: -4, Z: 3, Angle: 90}, To: keyframe.Pose{X: 3, Z: 3, Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 23, Shape: "8x8x1", Palette: "red", From: keyframe.Pose{Y: -4, Z: 4, Angle: 90}, To: keyframe.Pose{Y: 3, Z: 4, Angle: 90}},
			},
		},
	}
}

// test02X8x1 is cartoon 2 of the original Test group.
func test02X8x1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 24, Loop: 1, Delay: 400,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "8x8x1", Palette: "green", To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 8, End: 15, Shape: "8x8x1", Palette: "green", From: keyframe.Pose{X: -4, Z: 3, Angle: 90}, To: keyframe.Pose{X: 3, Z: 3, Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 23, Shape: "8x8x1", Palette: "green", From: keyframe.Pose{Y: -4, Z: 4, Angle: 90}, To: keyframe.Pose{Y: 3, Z: 4, Angle: 90}},
			},
		},
	}
}

// test03X8x1 is cartoon 3 of the original Test group.
func test03X8x1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 24, Loop: 1, Delay: 400,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "8x8x1", To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 8, End: 15, Shape: "8x8x1", From: keyframe.Pose{X: -4, Z: 3, Angle: 90}, To: keyframe.Pose{X: 3, Z: 3, Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 23, Shape: "8x8x1", From: keyframe.Pose{Y: -4, Z: 4, Angle: 90}, To: keyframe.Pose{Y: 3, Z: 4, Angle: 90}},
			},
		},
	}
}

// test04X8x1 is cartoon 4 of the original Test group.
func test04X8x1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 24, Loop: 1, Delay: 400,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "8x8x1", Palette: "white", To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 8, End: 15, Shape: "8x8x1", Palette: "white", From: keyframe.Pose{X: -4, Z: 3, Angle: 90}, To: keyframe.Pose{X: 3, Z: 3, Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 23, Shape: "8x8x1", Palette: "white", From: keyframe.Pose{Y: -4, Z: 4, Angle: 90}, To: keyframe.Pose{Y: 3, Z: 4, Angle: 90}},
			},
		},
	}
}

// Test lists the cartoons of the original Test group, in order.
var Test = []func() *keyframe.Cartoon{
	test01X8x1,
	test02X8x1,
	test03X8x1,
	test04X8x1,
}
