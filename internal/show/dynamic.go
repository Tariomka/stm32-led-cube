// Generated from the original iCubeSmart show data by tools/showgen.

package show

import (
	"github.com/Tariomka/stm32-led-cube/internal/engine/keyframe"
	"github.com/Tariomka/stm32-led-cube/internal/engine/render"
)

// dynamic01X8x8 is cartoon 1 of the original Dynamic group.
func dynamic01X8x8() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 100, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 29, Shape: "8x8x8"},
				keyframe.Graph{Start: 30, End: 45, Shape: "HELLO", From: keyframe.Pose{Angle: 180}, Axis: render.AxisY},
				keyframe.Graph{Start: 46, End: 64, Shape: "HELLO"},
				keyframe.Graph{Start: 65, End: 79, Shape: "3D8", From: keyframe.Pose{Angle: 180}},
				keyframe.Graph{Start: 80, End: 99, Shape: "3D8"},
			},
		},
	}
}

// dynamic02Shape is cartoon 2 of the original Dynamic group.
func dynamic02Shape() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 100, Loop: 1, Delay: 20,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "3", From: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 8, End: 15, Shape: "3", Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 23, Shape: "2", From: keyframe.Pose{Z: 7}, Axis: render.AxisY},
				keyframe.Graph{Start: 24, End: 31, Shape: "2", Axis: render.AxisY},
				keyframe.Graph{Start: 32, End: 39, Shape: "1", From: keyframe.Pose{Z: 7}, Axis: render.AxisY},
				keyframe.Graph{Start: 40, End: 47, Shape: "1", Axis: render.AxisY},
				keyframe.Graph{Start: 48, End: 55, Shape: "GO", From: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 56, End: 63, Shape: "GO", To: keyframe.Pose{Z: 7}, Axis: render.AxisY},
				keyframe.Graph{Start: 64, End: 71, Shape: "GO", From: keyframe.Pose{Z: 7}, Axis: render.AxisY},
				keyframe.Graph{Start: 72, End: 79, Shape: "GO", To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 80, End: 87, Shape: "GO", From: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 88, End: 99, Shape: "GO"},
			},
		},
	}
}

// dynamic03Waterfall1 is cartoon 3 of the original Dynamic group.
func dynamic03Waterfall1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 95, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "waterfall1", Palette: "red-series", To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 8, End: 18, Shape: "waterfall2", Palette: "red-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -3}},
				keyframe.Graph{Start: 19, End: 26, Shape: "waterfall1", Palette: "green-series", To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 27, End: 37, Shape: "waterfall2", Palette: "green-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -3}},
				keyframe.Graph{Start: 38, End: 45, Shape: "waterfall1", Palette: "blue-series", To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 46, End: 56, Shape: "waterfall2", Palette: "blue-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -3}},
				keyframe.Graph{Start: 57, End: 64, Shape: "waterfall1", Palette: "yellow-series", To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 65, End: 75, Shape: "waterfall2", Palette: "yellow-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -3}},
				keyframe.Graph{Start: 76, End: 83, Shape: "waterfall1", Palette: "cyan-series", To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 84, End: 94, Shape: "waterfall2", Palette: "cyan-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -3}},
			},
		},
	}
}

// dynamic04Strip is cartoon 4 of the original Dynamic group.
func dynamic04Strip() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 85, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 32, End: 39, Shape: "4strip", From: keyframe.Pose{X: -7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 40, End: 84, Shape: "4strip"},
			},
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "vertical8dot", From: keyframe.Pose{Z: -7}},
				keyframe.Graph{Start: 8, End: 31, Shape: "vertical8dot", Mode: render.ModeMix},
				keyframe.Graph{Start: 40, End: 47, Shape: "vertical8dot", From: keyframe.Pose{X: 7, Z: -7}, To: keyframe.Pose{X: 7}},
				keyframe.Graph{Start: 48, End: 84, Shape: "vertical8dot", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 7}},
			},
			{
				keyframe.Graph{Start: 8, End: 15, Shape: "crosswise8dot", From: keyframe.Pose{Y: -7, Z: 7}, To: keyframe.Pose{Z: 7}},
				keyframe.Graph{Start: 16, End: 31, Shape: "crosswise8dot", From: keyframe.Pose{Z: 7}, To: keyframe.Pose{Z: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 48, End: 55, Shape: "crosswise8dot", From: keyframe.Pose{X: 7, Y: -7, Z: 7}, To: keyframe.Pose{X: 7, Z: 7}},
				keyframe.Graph{Start: 56, End: 84, Shape: "crosswise8dot", From: keyframe.Pose{X: 7, Z: 7}, To: keyframe.Pose{X: 7, Z: 7}},
			},
			{
				keyframe.Graph{Start: 16, End: 23, Shape: "vertical8dot", From: keyframe.Pose{Y: 7, Z: 7}, To: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 24, End: 31, Shape: "vertical8dot", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 56, End: 63, Shape: "vertical8dot", From: keyframe.Pose{X: 7, Y: 7, Z: 7}, To: keyframe.Pose{X: 7, Y: 7}},
				keyframe.Graph{Start: 64, End: 84, Shape: "vertical8dot", From: keyframe.Pose{X: 7, Y: 7}, To: keyframe.Pose{X: 7, Y: 7}},
			},
			{
				keyframe.Graph{Start: 24, End: 31, Shape: "crosswise8dot", From: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 64, End: 71, Shape: "crosswise8dot", From: keyframe.Pose{X: 7, Y: 7}, To: keyframe.Pose{X: 7}},
				keyframe.Graph{Start: 72, End: 84, Shape: "crosswise8dot", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 7}},
			},
			{
				keyframe.Graph{Start: 32, End: 84, Shape: "1x8x8", Mode: render.ModeMix},
			},
		},
	}
}

// dynamic05Cubeamp1 is cartoon 5 of the original Dynamic group.
func dynamic05Cubeamp1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "cubeamp1", From: keyframe.Pose{X: 3, Y: 3, Z: 3}, To: keyframe.Pose{X: 3, Y: 3, Z: 3}},
				keyframe.Graph{Start: 1, End: 1, Shape: "cubeamp2", From: keyframe.Pose{X: 2, Y: 2, Z: 2}, To: keyframe.Pose{X: 2, Y: 2, Z: 2}},
				keyframe.Graph{Start: 2, End: 2, Shape: "cubeamp3", From: keyframe.Pose{X: 1, Y: 1, Z: 1}, To: keyframe.Pose{X: 1, Y: 1, Z: 1}},
				keyframe.Graph{Start: 3, End: 3, Shape: "cubeamp4"},
				keyframe.Graph{Start: 4, End: 4, Shape: "cubeamp3", From: keyframe.Pose{X: -1, Y: -1, Z: -1}, To: keyframe.Pose{X: -1, Y: -1, Z: -1}},
				keyframe.Graph{Start: 5, End: 5, Shape: "cubeamp2", From: keyframe.Pose{X: -2, Y: -2, Z: -2}, To: keyframe.Pose{X: -2, Y: -2, Z: -2}},
				keyframe.Graph{Start: 6, End: 6, Shape: "cubeamp1", From: keyframe.Pose{X: -3, Y: -3, Z: -3}, To: keyframe.Pose{X: -3, Y: -3, Z: -3}},
				keyframe.Graph{Start: 7, End: 7, Shape: "empty"},
				keyframe.Graph{Start: 8, End: 8, Shape: "cubeamp1", Palette: "red", From: keyframe.Pose{X: -3, Y: -3, Z: -3}, To: keyframe.Pose{X: -3, Y: -3, Z: -3}},
				keyframe.Graph{Start: 9, End: 9, Shape: "cubeamp2", Palette: "red", From: keyframe.Pose{X: -2, Y: -2, Z: -2}, To: keyframe.Pose{X: -2, Y: -2, Z: -2}},
				keyframe.Graph{Start: 10, End: 10, Shape: "cubeamp3", Palette: "red", From: keyframe.Pose{X: -1, Y: -1, Z: -1}, To: keyframe.Pose{X: -1, Y: -1, Z: -1}},
				keyframe.Graph{Start: 11, End: 11, Shape: "cubeamp4", Palette: "red"},
				keyframe.Graph{Start: 12, End: 12, Shape: "cubeamp3", Palette: "red", From: keyframe.Pose{X: 1, Y: 1, Z: 1}, To: keyframe.Pose{X: 1, Y: 1, Z: 1}},
				keyframe.Graph{Start: 13, End: 13, Shape: "cubeamp2", Palette: "red", From: keyframe.Pose{X: 2, Y: 2, Z: 2}, To: keyframe.Pose{X: 2, Y: 2, Z: 2}},
				keyframe.Graph{Start: 14, End: 14, Shape: "cubeamp1", Palette: "red", From: keyframe.Pose{X: 1, Y: 1, Z: 1}, To: keyframe.Pose{X: 1, Y: 1, Z: 1}},
				keyframe.Graph{Start: 15, End: 15, Shape: "empty", Palette: "red"},
			},
		},
	}
}

// dynamic06Cubeamp1 is cartoon 6 of the original Dynamic group.
func dynamic06Cubeamp1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 36, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "cubeamp1"},
				keyframe.Graph{Start: 1, End: 1, Shape: "cubeamp2"},
				keyframe.Graph{Start: 2, End: 2, Shape: "cubeamp3"},
				keyframe.Graph{Start: 3, End: 3, Shape: "cubeamp4"},
				keyframe.Graph{Start: 4, End: 4, Shape: "cubeamp3"},
				keyframe.Graph{Start: 5, End: 5, Shape: "cubeamp2"},
				keyframe.Graph{Start: 6, End: 6, Shape: "cubeamp1", Palette: "cyan"},
				keyframe.Graph{Start: 7, End: 7, Shape: "cubeamp2", Palette: "cyan"},
				keyframe.Graph{Start: 8, End: 8, Shape: "cubeamp3", Palette: "cyan"},
				keyframe.Graph{Start: 9, End: 9, Shape: "cubeamp4", Palette: "cyan"},
				keyframe.Graph{Start: 10, End: 10, Shape: "cubeamp3", Palette: "cyan"},
				keyframe.Graph{Start: 11, End: 11, Shape: "cubeamp2", Palette: "cyan"},
				keyframe.Graph{Start: 12, End: 12, Shape: "cubeamp1", Palette: "green"},
				keyframe.Graph{Start: 13, End: 13, Shape: "cubeamp2", Palette: "green"},
				keyframe.Graph{Start: 14, End: 14, Shape: "cubeamp3", Palette: "green"},
				keyframe.Graph{Start: 15, End: 15, Shape: "cubeamp4", Palette: "green"},
				keyframe.Graph{Start: 16, End: 16, Shape: "cubeamp3", Palette: "green"},
				keyframe.Graph{Start: 17, End: 17, Shape: "cubeamp2", Palette: "green"},
				keyframe.Graph{Start: 18, End: 18, Shape: "cubeamp1", Palette: "yellow"},
				keyframe.Graph{Start: 19, End: 19, Shape: "cubeamp2", Palette: "yellow"},
				keyframe.Graph{Start: 20, End: 20, Shape: "cubeamp3", Palette: "yellow"},
				keyframe.Graph{Start: 21, End: 21, Shape: "cubeamp4", Palette: "yellow"},
				keyframe.Graph{Start: 22, End: 22, Shape: "cubeamp3", Palette: "yellow"},
				keyframe.Graph{Start: 23, End: 23, Shape: "cubeamp2", Palette: "yellow"},
				keyframe.Graph{Start: 24, End: 24, Shape: "cubeamp1", Palette: "red"},
				keyframe.Graph{Start: 25, End: 25, Shape: "cubeamp2", Palette: "red"},
				keyframe.Graph{Start: 26, End: 26, Shape: "cubeamp3", Palette: "red"},
				keyframe.Graph{Start: 27, End: 27, Shape: "cubeamp4", Palette: "red"},
				keyframe.Graph{Start: 28, End: 28, Shape: "cubeamp3", Palette: "red"},
				keyframe.Graph{Start: 29, End: 29, Shape: "cubeamp2", Palette: "red"},
				keyframe.Graph{Start: 30, End: 30, Shape: "cubeamp1", Palette: "purple"},
				keyframe.Graph{Start: 31, End: 31, Shape: "cubeamp2", Palette: "purple"},
				keyframe.Graph{Start: 32, End: 32, Shape: "cubeamp3", Palette: "purple"},
				keyframe.Graph{Start: 33, End: 33, Shape: "cubeamp4", Palette: "purple"},
				keyframe.Graph{Start: 34, End: 34, Shape: "cubeamp3", Palette: "purple"},
				keyframe.Graph{Start: 35, End: 35, Shape: "cubeamp2", Palette: "purple"},
			},
		},
	}
}

// dynamic07Wave1 is cartoon 7 of the original Dynamic group.
func dynamic07Wave1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 6, Delay: 60,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "Wave1", Palette: "red-series", To: keyframe.Pose{X: -15}},
			},
			{
				keyframe.Graph{Start: 9, End: 15, Shape: "Wave1", Palette: "red-series", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 1}},
			},
		},
	}
}

// dynamic08II is cartoon 8 of the original Dynamic group.
func dynamic08II() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 144, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "II", Palette: "yellow-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 16, End: 31, Shape: "II", Palette: "yellow-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 32, End: 47, Shape: "II", Palette: "yellow-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 48, End: 63, Shape: "XX", Palette: "red-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 64, End: 79, Shape: "XX", Palette: "red-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 80, End: 95, Shape: "XX", Palette: "red-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 96, End: 111, Shape: "UU", Palette: "blue-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 112, End: 127, Shape: "UU", Palette: "blue-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
				keyframe.Graph{Start: 128, End: 143, Shape: "UU", Palette: "blue-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic09Arrows5 is cartoon 9 of the original Dynamic group.
func dynamic09Arrows5() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 72, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "Arrows5", Palette: "red-series", From: keyframe.Pose{X: -7}, To: keyframe.Pose{X: 8}},
				keyframe.Graph{Start: 32, End: 47, Shape: "Arrows5", Palette: "red-series", From: keyframe.Pose{X: 8, Z: -1, Angle: 180}, To: keyframe.Pose{X: -7, Z: -1, Angle: 180}, Axis: render.AxisY},
			},
			{
				keyframe.Graph{Start: 8, End: 23, Shape: "Arrows5", Palette: "green-series", From: keyframe.Pose{X: -7}, To: keyframe.Pose{X: 8}},
				keyframe.Graph{Start: 40, End: 55, Shape: "Arrows5", Palette: "green-series", From: keyframe.Pose{X: 8, Z: -1, Angle: 180}, To: keyframe.Pose{X: -7, Z: -1, Angle: 180}, Axis: render.AxisY},
			},
			{
				keyframe.Graph{Start: 16, End: 31, Shape: "Arrows5", Palette: "yellow-series", From: keyframe.Pose{X: -7}, To: keyframe.Pose{X: 8}},
				keyframe.Graph{Start: 48, End: 63, Shape: "Arrows5", Palette: "yellow-series", From: keyframe.Pose{X: 8, Z: -1, Angle: 180}, To: keyframe.Pose{X: -7, Z: -1, Angle: 180}, Axis: render.AxisY},
			},
			{
				keyframe.Graph{Start: 24, End: 39, Shape: "Arrows5", Palette: "cyan-series", From: keyframe.Pose{X: -7}, To: keyframe.Pose{X: 8}},
				keyframe.Graph{Start: 56, End: 71, Shape: "Arrows5", Palette: "cyan-series", From: keyframe.Pose{X: 8, Z: -1, Angle: 180}, To: keyframe.Pose{X: -7, Z: -1, Angle: 180}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic10Meteor is cartoon 10 of the original Dynamic group.
func dynamic10Meteor() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 79, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -7}},
				keyframe.Graph{Start: 14, End: 21, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 3, Y: 7}, To: keyframe.Pose{X: 3, Y: -7}},
				keyframe.Graph{Start: 27, End: 34, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 5, Y: 7, Z: 2}, To: keyframe.Pose{X: 5, Y: -7, Z: 2}},
				keyframe.Graph{Start: 41, End: 48, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 2, Y: 7, Z: 7}, To: keyframe.Pose{X: 2, Y: -7, Z: 7}},
				keyframe.Graph{Start: 53, End: 60, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 1, Y: 7, Z: 5}, To: keyframe.Pose{X: 1, Y: -7, Z: 5}},
			},
			{
				keyframe.Graph{Start: 16, End: 24, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 7, Y: 7, Z: 2}, To: keyframe.Pose{X: 7, Y: -7, Z: 2}},
				keyframe.Graph{Start: 31, End: 41, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 3, Y: 7, Z: 7}, To: keyframe.Pose{X: 3, Y: -7, Z: 7}},
				keyframe.Graph{Start: 47, End: 54, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 1, Y: 7, Z: 2}, To: keyframe.Pose{X: 1, Y: -7, Z: 2}},
				keyframe.Graph{Start: 59, End: 65, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 2, Y: 7, Z: 7}, To: keyframe.Pose{X: 2, Y: -7, Z: 7}},
				keyframe.Graph{Start: 70, End: 77, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 1, Y: 7, Z: 5}, To: keyframe.Pose{X: 1, Y: -7, Z: 5}},
			},
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 2, Y: 7, Z: 1}, To: keyframe.Pose{X: 2, Y: -7, Z: 1}},
				keyframe.Graph{Start: 11, End: 18, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 7, Y: 7, Z: 2}, To: keyframe.Pose{X: 7, Y: -7, Z: 2}},
				keyframe.Graph{Start: 23, End: 30, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 7, Y: 7, Z: 7}, To: keyframe.Pose{X: 7, Y: -7, Z: 7}},
				keyframe.Graph{Start: 36, End: 43, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 2, Y: 7, Z: 1}, To: keyframe.Pose{X: 2, Y: -7, Z: 1}},
				keyframe.Graph{Start: 63, End: 72, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 1, Y: 7}, To: keyframe.Pose{X: 1, Y: -7}},
			},
			{
				keyframe.Graph{Start: 3, End: 10, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 4, Y: 7, Z: 6}, To: keyframe.Pose{X: 4, Y: -7, Z: 6}},
				keyframe.Graph{Start: 14, End: 21, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 3, Y: 7, Z: 5}, To: keyframe.Pose{X: 3, Y: -7, Z: 5}},
				keyframe.Graph{Start: 27, End: 36, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 1, Y: 7, Z: 3}, To: keyframe.Pose{X: 1, Y: -7, Z: 3}},
				keyframe.Graph{Start: 53, End: 60, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 2, Y: 7, Z: 7}, To: keyframe.Pose{X: 2, Y: -7, Z: 7}},
				keyframe.Graph{Start: 68, End: 78, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 1, Y: 7, Z: 2}, To: keyframe.Pose{X: 1, Y: -7, Z: 2}},
			},
			{
				keyframe.Graph{Start: 16, End: 23, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 5, Y: 7, Z: 4}, To: keyframe.Pose{X: 5, Y: -7, Z: 4}},
				keyframe.Graph{Start: 27, End: 34, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 3, Y: 7, Z: 3}, To: keyframe.Pose{X: 3, Y: -7, Z: 3}},
				keyframe.Graph{Start: 35, End: 42, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 1, Y: 7}, To: keyframe.Pose{X: 1, Y: -7}},
				keyframe.Graph{Start: 47, End: 55, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 7, Y: 7, Z: 2}, To: keyframe.Pose{X: 7, Y: -7, Z: 2}},
				keyframe.Graph{Start: 61, End: 70, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 4, Y: 7, Z: 5}, To: keyframe.Pose{X: 4, Y: -7, Z: 5}},
			},
			{
				keyframe.Graph{Start: 9, End: 17, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 7, Y: 7, Z: 7}, To: keyframe.Pose{X: 7, Y: -7, Z: 7}},
				keyframe.Graph{Start: 21, End: 28, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 7, Y: 7}, To: keyframe.Pose{X: 7, Y: -7}},
				keyframe.Graph{Start: 37, End: 45, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 2, Y: 7, Z: 7}, To: keyframe.Pose{X: 2, Y: -7, Z: 7}},
				keyframe.Graph{Start: 50, End: 65, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 5, Y: 7, Z: 3}, To: keyframe.Pose{X: 5, Y: -7, Z: 3}},
				keyframe.Graph{Start: 67, End: 74, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 5, Y: 7, Z: 1}, To: keyframe.Pose{X: 5, Y: -7, Z: 1}},
			},
			{
				keyframe.Graph{Start: 3, End: 10, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 3, Y: 7}, To: keyframe.Pose{X: 3, Y: -7}},
				keyframe.Graph{Start: 14, End: 22, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 5, Y: 7, Z: 4}, To: keyframe.Pose{X: 5, Y: -7, Z: 4}},
				keyframe.Graph{Start: 27, End: 34, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 5, Y: 7, Z: 2}, To: keyframe.Pose{X: 5, Y: -7, Z: 2}},
				keyframe.Graph{Start: 40, End: 49, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 7, Y: 7}, To: keyframe.Pose{X: 7, Y: -7}},
				keyframe.Graph{Start: 53, End: 61, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 4, Y: 7, Z: 3}, To: keyframe.Pose{X: 4, Y: -7, Z: 3}},
			},
			{
				keyframe.Graph{Start: 9, End: 16, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 3, Y: 7, Z: 7}, To: keyframe.Pose{X: 3, Y: -7, Z: 7}},
				keyframe.Graph{Start: 22, End: 29, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{Y: 7, Z: 6}, To: keyframe.Pose{Y: -7, Z: 6}},
				keyframe.Graph{Start: 32, End: 39, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: -7}},
				keyframe.Graph{Start: 48, End: 58, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 1, Y: 7}, To: keyframe.Pose{X: 1, Y: -7}},
				keyframe.Graph{Start: 63, End: 72, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 6, Y: 7, Z: 7}, To: keyframe.Pose{X: 6, Y: -7, Z: 7}},
			},
			{
				keyframe.Graph{Start: 18, End: 25, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 1, Y: 7, Z: 4}, To: keyframe.Pose{X: 1, Y: -7, Z: 4}},
				keyframe.Graph{Start: 35, End: 42, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 4, Y: 7, Z: 1}, To: keyframe.Pose{X: 4, Y: -7, Z: 1}},
				keyframe.Graph{Start: 45, End: 56, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 4, Y: 7, Z: 6}, To: keyframe.Pose{X: 4, Y: -7, Z: 6}},
				keyframe.Graph{Start: 59, End: 66, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 6, Y: 7, Z: 3}, To: keyframe.Pose{X: 6, Y: -7, Z: 3}},
				keyframe.Graph{Start: 71, End: 78, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 5, Y: 7, Z: 5}, To: keyframe.Pose{X: 5, Y: -7, Z: 5}},
			},
			{
				keyframe.Graph{Start: 5, End: 12, Shape: "meteor", Palette: "red-series", From: keyframe.Pose{X: 7, Y: 7, Z: 2}, To: keyframe.Pose{X: 7, Y: -7, Z: 2}},
				keyframe.Graph{Start: 17, End: 24, Shape: "meteor", Palette: "green-series", From: keyframe.Pose{X: 3, Y: 7, Z: 7}, To: keyframe.Pose{X: 3, Y: -7, Z: 7}},
				keyframe.Graph{Start: 32, End: 39, Shape: "meteor", Palette: "blue-series", From: keyframe.Pose{X: 1, Y: 7, Z: 2}, To: keyframe.Pose{X: 1, Y: -7, Z: 2}},
				keyframe.Graph{Start: 46, End: 54, Shape: "meteor", Palette: "yellow-series", From: keyframe.Pose{X: 2, Y: 7, Z: 7}, To: keyframe.Pose{X: 2, Y: -7, Z: 7}},
				keyframe.Graph{Start: 65, End: 74, Shape: "meteor", Palette: "cyan-series", From: keyframe.Pose{X: 1, Y: 7, Z: 5}, To: keyframe.Pose{X: 1, Y: -7, Z: 5}},
			},
		},
	}
}

// dynamic11Balls is cartoon 11 of the original Dynamic group.
func dynamic11Balls() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 37, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Ball{Start: 0, End: 6, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x33},
				keyframe.Ball{Start: 10, End: 16, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x3C},
				keyframe.Ball{Start: 20, End: 26, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x0F},
				keyframe.Ball{Start: 30, End: 36, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x3F},
			},
			{
				keyframe.Ball{Start: 4, End: 10, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x30},
				keyframe.Ball{Start: 15, End: 21, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x0C},
				keyframe.Ball{Start: 25, End: 31, X: 350, Y: 350, Z: 350, From: 100, To: 700, Hollow: true, Out: 0x00, In: 0x03},
			},
		},
	}
}

// dynamic12X4x4 is cartoon 12 of the original Dynamic group.
func dynamic12X4x4() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 60, Loop: 2, Delay: 70,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 4, Shape: "4x4x4", Palette: "purple", To: keyframe.Pose{X: 4}},
				keyframe.Graph{Start: 5, End: 9, Shape: "4x4x4", Palette: "purple", From: keyframe.Pose{X: 4}, To: keyframe.Pose{X: 4, Y: 4}},
				keyframe.Graph{Start: 10, End: 14, Shape: "4x4x4", Palette: "purple", From: keyframe.Pose{X: 4, Y: 4}, To: keyframe.Pose{X: 4, Y: 4, Z: 4}},
				keyframe.Graph{Start: 15, End: 19, Shape: "4x4x4", Palette: "purple", From: keyframe.Pose{X: 4, Y: 4, Z: 4}, To: keyframe.Pose{Y: 4, Z: 4}},
				keyframe.Graph{Start: 20, End: 24, Shape: "4x4x4", Palette: "purple", From: keyframe.Pose{Y: 4, Z: 4}, To: keyframe.Pose{Z: 4}},
				keyframe.Graph{Start: 25, End: 29, Shape: "4x4x4", Palette: "purple", From: keyframe.Pose{Z: 4}},
				keyframe.Graph{Start: 30, End: 34, Shape: "4x4x4", Palette: "yellow", To: keyframe.Pose{Y: 4}},
				keyframe.Graph{Start: 35, End: 39, Shape: "4x4x4", Palette: "yellow", From: keyframe.Pose{Y: 4}, To: keyframe.Pose{Y: 4, Z: 4}},
				keyframe.Graph{Start: 40, End: 44, Shape: "4x4x4", Palette: "yellow", From: keyframe.Pose{Y: 4, Z: 4}, To: keyframe.Pose{X: 4, Y: 4, Z: 4}},
				keyframe.Graph{Start: 45, End: 49, Shape: "4x4x4", Palette: "yellow", From: keyframe.Pose{X: 4, Y: 4, Z: 4}, To: keyframe.Pose{X: 4, Y: 4}},
				keyframe.Graph{Start: 50, End: 54, Shape: "4x4x4", Palette: "yellow", From: keyframe.Pose{X: 4, Y: 4}, To: keyframe.Pose{Y: 4}},
				keyframe.Graph{Start: 55, End: 59, Shape: "4x4x4", Palette: "yellow", From: keyframe.Pose{Y: 4}},
			},
		},
	}
}

// dynamic13X4x4 is cartoon 13 of the original Dynamic group.
func dynamic13X4x4() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 40, Loop: 3, Delay: 70,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 4, Shape: "4x4x4", To: keyframe.Pose{X: 4}},
				keyframe.Graph{Start: 5, End: 9, Shape: "4x4x4", From: keyframe.Pose{X: 4}, To: keyframe.Pose{X: 4, Z: 4}},
				keyframe.Graph{Start: 10, End: 14, Shape: "4x4x4", From: keyframe.Pose{X: 4, Z: 4}, To: keyframe.Pose{Z: 4}},
				keyframe.Graph{Start: 15, End: 19, Shape: "4x4x4", From: keyframe.Pose{Z: 4}},
				keyframe.Graph{Start: 20, End: 24, Shape: "4x4x4", To: keyframe.Pose{Y: 4}},
				keyframe.Graph{Start: 25, End: 29, Shape: "4x4x4", From: keyframe.Pose{Y: 4}, To: keyframe.Pose{Y: 4, Z: 4}},
				keyframe.Graph{Start: 30, End: 34, Shape: "4x4x4", From: keyframe.Pose{Y: 4, Z: 4}, To: keyframe.Pose{Z: 4}},
				keyframe.Graph{Start: 35, End: 39, Shape: "4x4x4", From: keyframe.Pose{Z: 4}},
			},
			{
				keyframe.Graph{Start: 0, End: 4, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4, Y: 4, Z: 4}, To: keyframe.Pose{Y: 4, Z: 4}},
				keyframe.Graph{Start: 5, End: 9, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{Y: 4, Z: 4}, To: keyframe.Pose{Y: 4}},
				keyframe.Graph{Start: 10, End: 14, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{Y: 4}, To: keyframe.Pose{X: 4, Y: 4}},
				keyframe.Graph{Start: 15, End: 19, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4, Y: 4}, To: keyframe.Pose{X: 4, Y: 4, Z: 4}},
				keyframe.Graph{Start: 20, End: 24, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4, Y: 4, Z: 4}, To: keyframe.Pose{X: 4, Z: 4}},
				keyframe.Graph{Start: 25, End: 29, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4, Z: 4}, To: keyframe.Pose{X: 4}},
				keyframe.Graph{Start: 30, End: 34, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4}, To: keyframe.Pose{X: 4, Y: 4}},
				keyframe.Graph{Start: 35, End: 39, Shape: "4x4x4", Palette: "red", From: keyframe.Pose{X: 4, Y: 4}, To: keyframe.Pose{X: 4, Y: 4, Z: 4}},
			},
		},
	}
}

// dynamic14X8x1R is cartoon 14 of the original Dynamic group.
func dynamic14X8x1R() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 32, Loop: 3, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 31, Shape: "8x8x1R", Palette: "red-series", From: keyframe.Pose{Angle: 350}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic15X8x1 is cartoon 15 of the original Dynamic group.
func dynamic15X8x1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 60, Loop: 2, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "8x8x1", Palette: "blue", From: keyframe.Pose{Z: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 20, End: 27, Shape: "8x8x1", Palette: "red", From: keyframe.Pose{Z: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 40, End: 47, Shape: "8x8x1", Palette: "green", From: keyframe.Pose{Z: 7}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 4, End: 19, Shape: "I", Palette: "blue", From: keyframe.Pose{Z: 3}, To: keyframe.Pose{Z: 3}, Mode: render.ModeMix},
				keyframe.Graph{Start: 24, End: 39, Shape: "X", Palette: "red", From: keyframe.Pose{Z: 3}, To: keyframe.Pose{Z: 3}, Mode: render.ModeMix},
				keyframe.Graph{Start: 44, End: 59, Shape: "U", Palette: "green", From: keyframe.Pose{Z: 3}, To: keyframe.Pose{Z: 3}, Mode: render.ModeMix},
			},
		},
	}
}

// dynamic16X2 is cartoon 16 of the original Dynamic group.
func dynamic16X2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 50, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 4, Shape: "2x2", Palette: "blue", From: keyframe.Pose{X: 3, Z: 3}, To: keyframe.Pose{X: 3, Y: 4, Z: 3}},
				keyframe.Graph{Start: 10, End: 14, Shape: "2x2", Palette: "green", From: keyframe.Pose{X: 2, Z: 2}, To: keyframe.Pose{X: 2, Y: 5, Z: 2}},
				keyframe.Graph{Start: 20, End: 24, Shape: "2x2", Palette: "red", From: keyframe.Pose{X: 6, Z: 6}, To: keyframe.Pose{X: 6, Y: 4, Z: 6}},
				keyframe.Graph{Start: 30, End: 34, Shape: "2x2", Palette: "blue", From: keyframe.Pose{X: 2, Z: 5}, To: keyframe.Pose{X: 2, Y: 5, Z: 5}},
				keyframe.Graph{Start: 40, End: 44, Shape: "2x2", Palette: "yellow", From: keyframe.Pose{X: 6, Z: 2}, To: keyframe.Pose{X: 6, Y: 5, Z: 2}},
			},
			{
				keyframe.Ball{Start: 4, End: 9, X: 300, Y: 500, Z: 300, From: 0, To: 300, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 14, End: 19, X: 200, Y: 500, Z: 200, From: 0, To: 400, Hollow: true, Out: 0x0C, In: 0x3F},
				keyframe.Ball{Start: 24, End: 29, X: 600, Y: 500, Z: 600, From: 0, To: 300, Hollow: true, Out: 0x30, In: 0x3F},
				keyframe.Ball{Start: 34, End: 39, X: 200, Y: 500, Z: 500, From: 0, To: 300, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 45, End: 49, X: 600, Y: 500, Z: 200, From: 0, To: 400, Hollow: true, Out: 0x3C, In: 0x3F},
			},
		},
	}
}

// dynamic17XX is cartoon 17 of the original Dynamic group.
func dynamic17XX() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 4, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "XX", Palette: "red-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic18Circle5 is cartoon 18 of the original Dynamic group.
func dynamic18Circle5() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 30, Loop: 2, Delay: 200,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "Circle5", Palette: "red"},
				keyframe.Graph{Start: 1, End: 1, Shape: "Circle4", Palette: "red"},
				keyframe.Graph{Start: 2, End: 2, Shape: "Circle3", Palette: "red"},
				keyframe.Graph{Start: 3, End: 3, Shape: "Circle2", Palette: "red"},
				keyframe.Graph{Start: 4, End: 4, Shape: "Circle1", Palette: "red"},
				keyframe.Graph{Start: 5, End: 5, Shape: "Circle5", Palette: "yellow"},
				keyframe.Graph{Start: 6, End: 6, Shape: "Circle4", Palette: "yellow"},
				keyframe.Graph{Start: 7, End: 7, Shape: "Circle3", Palette: "yellow"},
				keyframe.Graph{Start: 8, End: 8, Shape: "Circle2", Palette: "yellow"},
				keyframe.Graph{Start: 9, End: 9, Shape: "Circle1", Palette: "yellow"},
				keyframe.Graph{Start: 10, End: 10, Shape: "Circle5", Palette: "green"},
				keyframe.Graph{Start: 11, End: 11, Shape: "Circle4", Palette: "green"},
				keyframe.Graph{Start: 12, End: 12, Shape: "Circle3", Palette: "green"},
				keyframe.Graph{Start: 13, End: 13, Shape: "Circle2", Palette: "green"},
				keyframe.Graph{Start: 14, End: 14, Shape: "Circle1", Palette: "green"},
				keyframe.Graph{Start: 15, End: 15, Shape: "Circle5", Palette: "cyan"},
				keyframe.Graph{Start: 16, End: 16, Shape: "Circle4", Palette: "cyan"},
				keyframe.Graph{Start: 17, End: 17, Shape: "Circle3", Palette: "cyan"},
				keyframe.Graph{Start: 18, End: 18, Shape: "Circle2", Palette: "cyan"},
				keyframe.Graph{Start: 19, End: 19, Shape: "Circle1", Palette: "cyan"},
				keyframe.Graph{Start: 20, End: 20, Shape: "Circle5", Palette: "blue"},
				keyframe.Graph{Start: 21, End: 21, Shape: "Circle4", Palette: "blue"},
				keyframe.Graph{Start: 22, End: 22, Shape: "Circle3", Palette: "blue"},
				keyframe.Graph{Start: 23, End: 23, Shape: "Circle2", Palette: "blue"},
				keyframe.Graph{Start: 24, End: 24, Shape: "Circle1", Palette: "blue"},
				keyframe.Graph{Start: 25, End: 25, Shape: "Circle5", Palette: "purple"},
				keyframe.Graph{Start: 26, End: 26, Shape: "Circle4", Palette: "purple"},
				keyframe.Graph{Start: 27, End: 27, Shape: "Circle3", Palette: "purple"},
				keyframe.Graph{Start: 28, End: 28, Shape: "Circle2", Palette: "purple"},
				keyframe.Graph{Start: 29, End: 29, Shape: "Circle1", Palette: "purple"},
			},
		},
	}
}

// dynamic19Balls is cartoon 19 of the original Dynamic group.
func dynamic19Balls() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 176, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Ball{Start: 0, End: 15, X: 0, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 16, End: 31, X: 700, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x0F, In: 0x3F},
				keyframe.Ball{Start: 32, End: 47, X: 700, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x0C, In: 0x3F},
				keyframe.Ball{Start: 48, End: 63, X: 0, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x3C, In: 0x3F},
				keyframe.Ball{Start: 64, End: 79, X: 0, Y: 700, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x30, In: 0x3F},
				keyframe.Ball{Start: 80, End: 95, X: 0, Y: 700, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x32, In: 0x3F},
				keyframe.Ball{Start: 96, End: 111, X: 0, Y: 0, Z: 0, From: 0, To: 1300, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 112, End: 127, X: 700, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x0B, In: 0x3F},
				keyframe.Ball{Start: 128, End: 143, X: 0, Y: 700, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x0C, In: 0x3F},
				keyframe.Ball{Start: 144, End: 159, X: 700, Y: 700, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x3C, In: 0x3F},
				keyframe.Ball{Start: 160, End: 175, X: 0, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x30, In: 0x3F},
			},
			{
				keyframe.Ball{Start: 96, End: 111, X: 700, Y: 700, Z: 700, From: 0, To: 1300, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 112, End: 127, X: 0, Y: 700, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x0B, In: 0x3F},
				keyframe.Ball{Start: 128, End: 143, X: 700, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x0C, In: 0x3F},
				keyframe.Ball{Start: 144, End: 159, X: 0, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x3C, In: 0x3F},
				keyframe.Ball{Start: 160, End: 175, X: 700, Y: 700, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x30, In: 0x3F},
			},
		},
	}
}

// dynamic20Balls is cartoon 20 of the original Dynamic group.
func dynamic20Balls() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 32, Loop: 3, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Ball{Start: 0, End: 15, X: 0, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x30, In: 0x3F},
				keyframe.Ball{Start: 16, End: 31, X: 700, Y: 700, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x38, In: 0x3F},
			},
			{
				keyframe.Ball{Start: 0, End: 15, X: 700, Y: 0, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x03, In: 0x3F},
				keyframe.Ball{Start: 16, End: 31, X: 0, Y: 700, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x0F, In: 0x3F},
			},
			{
				keyframe.Ball{Start: 0, End: 15, X: 0, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x0C, In: 0x3F},
				keyframe.Ball{Start: 16, End: 31, X: 700, Y: 700, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x32, In: 0x3F},
			},
			{
				keyframe.Ball{Start: 0, End: 15, X: 700, Y: 0, Z: 700, From: 0, To: 1400, Hollow: true, Out: 0x3C, In: 0x3F},
				keyframe.Ball{Start: 16, End: 31, X: 0, Y: 700, Z: 0, From: 0, To: 1400, Hollow: true, Out: 0x13, In: 0x3F},
			},
		},
	}
}

// dynamic21Frame1 is cartoon 21 of the original Dynamic group.
func dynamic21Frame1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 30,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "frame1", Palette: "green-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic22Hui is cartoon 22 of the original Dynamic group.
func dynamic22Hui() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 8, Delay: 60,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "Hui", Palette: "cyan-series", To: keyframe.Pose{X: -15}},
			},
			{
				keyframe.Graph{Start: 9, End: 15, Shape: "Hui", Palette: "cyan-series", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 1}},
			},
		},
	}
}

// dynamic23I is cartoon 23 of the original Dynamic group.
func dynamic23I() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 60, Loop: 2, Delay: 30,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "I", Palette: "purple", From: keyframe.Pose{Z: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 20, End: 27, Shape: "X", Palette: "cyan", From: keyframe.Pose{Z: 7}, To: keyframe.Pose{Z: 1}, Mode: render.ModeMix},
				keyframe.Graph{Start: 40, End: 47, Shape: "U", Palette: "yellow", From: keyframe.Pose{Z: 7}, To: keyframe.Pose{Z: 2}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 8, End: 19, Shape: "I", Palette: "purple", Mode: render.ModeMix},
				keyframe.Graph{Start: 28, End: 39, Shape: "X", Palette: "cyan", From: keyframe.Pose{Z: 1}, To: keyframe.Pose{Z: 1}, Mode: render.ModeMix},
				keyframe.Graph{Start: 48, End: 59, Shape: "U", Palette: "yellow", From: keyframe.Pose{Z: 2}, To: keyframe.Pose{Z: 2}, Mode: render.ModeMix},
			},
		},
	}
}

// dynamic24Vertical8dot is cartoon 24 of the original Dynamic group.
func dynamic24Vertical8dot() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 71, Loop: 1, Delay: 15,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "vertical8dot", From: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 8, End: 15, Shape: "vertical8dot", From: keyframe.Pose{X: 1, Y: 7}, To: keyframe.Pose{X: 1}},
				keyframe.Graph{Start: 16, End: 23, Shape: "crosswise8dot", From: keyframe.Pose{X: 2, Y: 7}, To: keyframe.Pose{X: 2}},
				keyframe.Graph{Start: 24, End: 31, Shape: "vertical8dot", From: keyframe.Pose{X: 3, Y: 7}, To: keyframe.Pose{X: 3}},
				keyframe.Graph{Start: 32, End: 39, Shape: "vertical8dot", From: keyframe.Pose{X: 4, Y: 7}, To: keyframe.Pose{X: 4}},
				keyframe.Graph{Start: 40, End: 47, Shape: "vertical8dot", From: keyframe.Pose{X: 5, Y: 7}, To: keyframe.Pose{X: 5}},
				keyframe.Graph{Start: 48, End: 55, Shape: "vertical8dot", From: keyframe.Pose{X: 6, Y: 7}, To: keyframe.Pose{X: 6}},
				keyframe.Graph{Start: 56, End: 63, Shape: "vertical8dot", From: keyframe.Pose{X: 7, Y: 7}, To: keyframe.Pose{X: 7}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 7, End: 14, Shape: "8x1x8", From: keyframe.Pose{X: -7}, To: keyframe.Pose{X: -7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 15, End: 22, Shape: "8x1x8", From: keyframe.Pose{X: -6}, To: keyframe.Pose{X: -6}},
				keyframe.Graph{Start: 23, End: 30, Shape: "8x1x8", From: keyframe.Pose{X: -5}, To: keyframe.Pose{X: -5}},
				keyframe.Graph{Start: 31, End: 38, Shape: "8x1x8", From: keyframe.Pose{X: -4}, To: keyframe.Pose{X: -4}},
				keyframe.Graph{Start: 39, End: 46, Shape: "8x1x8", From: keyframe.Pose{X: -3}, To: keyframe.Pose{X: -3}},
				keyframe.Graph{Start: 47, End: 54, Shape: "8x1x8", From: keyframe.Pose{X: -2}, To: keyframe.Pose{X: -2}},
				keyframe.Graph{Start: 55, End: 62, Shape: "8x1x8", From: keyframe.Pose{X: -1}, To: keyframe.Pose{X: -1}},
				keyframe.Graph{Start: 63, End: 70, Shape: "8x1x8", Mode: render.ModeMix},
			},
		},
	}
}

// dynamic25Vertical8dot is cartoon 25 of the original Dynamic group.
func dynamic25Vertical8dot() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 71, Loop: 2, Delay: 15,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "vertical8dot", To: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 8, End: 15, Shape: "vertical8dot", From: keyframe.Pose{X: 1}, To: keyframe.Pose{X: 1, Y: 7}},
				keyframe.Graph{Start: 16, End: 23, Shape: "vertical8dot", From: keyframe.Pose{X: 2}, To: keyframe.Pose{X: 2, Y: 7}},
				keyframe.Graph{Start: 24, End: 31, Shape: "vertical8dot", From: keyframe.Pose{X: 3}, To: keyframe.Pose{X: 3, Y: 7}},
				keyframe.Graph{Start: 32, End: 39, Shape: "vertical8dot", From: keyframe.Pose{X: 4}, To: keyframe.Pose{X: 4, Y: 7}},
				keyframe.Graph{Start: 40, End: 47, Shape: "vertical8dot", From: keyframe.Pose{X: 5}, To: keyframe.Pose{X: 5, Y: 7}},
				keyframe.Graph{Start: 48, End: 55, Shape: "vertical8dot", From: keyframe.Pose{X: 6}, To: keyframe.Pose{X: 6, Y: 7}},
				keyframe.Graph{Start: 56, End: 63, Shape: "vertical8dot", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 7, Y: 7}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 7, End: 14, Shape: "8x1x8", From: keyframe.Pose{X: -7, Y: 7}, To: keyframe.Pose{X: -7, Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 15, End: 22, Shape: "8x1x8", From: keyframe.Pose{X: -6, Y: 7}, To: keyframe.Pose{X: -6, Y: 7}},
				keyframe.Graph{Start: 23, End: 30, Shape: "8x1x8", From: keyframe.Pose{X: -5, Y: 7}, To: keyframe.Pose{X: -5, Y: 7}},
				keyframe.Graph{Start: 31, End: 38, Shape: "8x1x8", From: keyframe.Pose{X: -4, Y: 7}, To: keyframe.Pose{X: -4, Y: 7}},
				keyframe.Graph{Start: 39, End: 46, Shape: "8x1x8", From: keyframe.Pose{X: -3, Y: 7}, To: keyframe.Pose{X: -3, Y: 7}},
				keyframe.Graph{Start: 47, End: 54, Shape: "8x1x8", From: keyframe.Pose{X: -2, Y: 7}, To: keyframe.Pose{X: -2, Y: 7}},
				keyframe.Graph{Start: 55, End: 62, Shape: "8x1x8", From: keyframe.Pose{X: -1, Y: 7}, To: keyframe.Pose{X: -1, Y: 7}},
				keyframe.Graph{Start: 63, End: 70, Shape: "8x1x8", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 1, End: 8, Shape: "8x1x8", From: keyframe.Pose{X: 1}, To: keyframe.Pose{X: 1}},
				keyframe.Graph{Start: 9, End: 16, Shape: "8x1x8", From: keyframe.Pose{X: 2}, To: keyframe.Pose{X: 2}},
				keyframe.Graph{Start: 17, End: 24, Shape: "8x1x8", From: keyframe.Pose{X: 3}, To: keyframe.Pose{X: 3}},
				keyframe.Graph{Start: 25, End: 32, Shape: "8x1x8", From: keyframe.Pose{X: 4}, To: keyframe.Pose{X: 4}},
				keyframe.Graph{Start: 33, End: 40, Shape: "8x1x8", From: keyframe.Pose{X: 5}, To: keyframe.Pose{X: 5}},
				keyframe.Graph{Start: 41, End: 48, Shape: "8x1x8", From: keyframe.Pose{X: 6}, To: keyframe.Pose{X: 6}},
				keyframe.Graph{Start: 49, End: 56, Shape: "8x1x8", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 7}},
			},
		},
	}
}

// dynamic26X2 is cartoon 26 of the original Dynamic group.
func dynamic26X2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 17, Loop: 3, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "2x2", From: keyframe.Pose{X: 3, Y: 7, Z: 3}, To: keyframe.Pose{X: 3, Z: 3}},
				keyframe.Graph{Start: 8, End: 9, Shape: "Circle-2", Palette: "cyan"},
				keyframe.Graph{Start: 10, End: 11, Shape: "Circle-3", Palette: "green"},
				keyframe.Graph{Start: 12, End: 13, Shape: "Circle-4", Palette: "yellow"},
				keyframe.Graph{Start: 14, End: 15, Shape: "Circle-5", Palette: "red"},
			},
		},
	}
}

// dynamic27X2 is cartoon 27 of the original Dynamic group.
func dynamic27X2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 17, Loop: 3, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 7, Shape: "2x2", Palette: "blue", From: keyframe.Pose{X: 3, Z: 3}, To: keyframe.Pose{X: 3, Y: 7, Z: 3}},
				keyframe.Graph{Start: 8, End: 9, Shape: "Circle-2", Palette: "purple", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 10, End: 11, Shape: "Circle-3", Palette: "red", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 12, End: 13, Shape: "Circle-4", Palette: "yellow", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}},
				keyframe.Graph{Start: 14, End: 15, Shape: "Circle-5", Palette: "green", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}},
			},
		},
	}
}

// dynamic28Frame2 is cartoon 28 of the original Dynamic group.
func dynamic28Frame2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 30,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "frame2", Palette: "yellow-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic29Fly1 is cartoon 29 of the original Dynamic group.
func dynamic29Fly1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 7, Loop: 10, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 1, Shape: "Fly_1", Palette: "yellow-series"},
				keyframe.Graph{Start: 2, End: 2, Shape: "Fly_2", Palette: "yellow-series"},
				keyframe.Graph{Start: 3, End: 3, Shape: "Fly_3", Palette: "yellow-series"},
				keyframe.Graph{Start: 4, End: 4, Shape: "Fly_4", Palette: "yellow-series"},
				keyframe.Graph{Start: 5, End: 5, Shape: "Fly_3", Palette: "yellow-series"},
				keyframe.Graph{Start: 6, End: 6, Shape: "Fly_2", Palette: "yellow-series"},
			},
		},
	}
}

// dynamic30Wave2 is cartoon 30 of the original Dynamic group.
func dynamic30Wave2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 8, Delay: 60,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "Wave2", Palette: "blue-series", To: keyframe.Pose{X: -15}},
			},
			{
				keyframe.Graph{Start: 9, End: 15, Shape: "Wave2", Palette: "blue-series", From: keyframe.Pose{X: 7}, To: keyframe.Pose{X: 1}},
			},
		},
	}
}

// dynamic31Frame4 is cartoon 31 of the original Dynamic group.
func dynamic31Frame4() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 6, Delay: 30,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "frame4", Palette: "red-series", From: keyframe.Pose{Z: -1}, To: keyframe.Pose{Z: -1, Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic32TurnRight1 is cartoon 32 of the original Dynamic group.
func dynamic32TurnRight1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "TurnRight1", Palette: "red-series", To: keyframe.Pose{Angle: 170}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic33TurnRight1 is cartoon 33 of the original Dynamic group.
func dynamic33TurnRight1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "TurnRight1", Palette: "green-series", To: keyframe.Pose{Angle: 170}, Axis: render.AxisZ},
			},
		},
	}
}

// dynamic34TurnRight2 is cartoon 34 of the original Dynamic group.
func dynamic34TurnRight2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "TurnRight2", Palette: "blue-series", To: keyframe.Pose{Angle: 170}, Axis: render.AxisZ},
			},
		},
	}
}

// dynamic35II is cartoon 35 of the original Dynamic group.
func dynamic35II() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 33, Loop: 2, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 16, Shape: "II", Palette: "red-series", From: keyframe.Pose{Y: -8}, To: keyframe.Pose{Y: 8}},
			},
			{
				keyframe.Graph{Start: 8, End: 24, Shape: "XX", Palette: "cyan-series", From: keyframe.Pose{Y: -8}, To: keyframe.Pose{Y: 8}},
			},
			{
				keyframe.Graph{Start: 16, End: 32, Shape: "UU", Palette: "green-series", From: keyframe.Pose{Y: -8}, To: keyframe.Pose{Y: 8}},
			},
		},
	}
}

// dynamic36X8x1R is cartoon 36 of the original Dynamic group.
func dynamic36X8x1R() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 32, Loop: 3, Delay: 50,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 31, Shape: "8x8x1R", Palette: "blue-series", From: keyframe.Pose{Angle: 350}},
			},
		},
	}
}

// dynamic37Square1 is cartoon 37 of the original Dynamic group.
func dynamic37Square1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 6, Loop: 10, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "square1", Palette: "red-series"},
				keyframe.Graph{Start: 1, End: 1, Shape: "square2", Palette: "red-series"},
				keyframe.Graph{Start: 2, End: 2, Shape: "square3", Palette: "red-series"},
				keyframe.Graph{Start: 3, End: 3, Shape: "square4", Palette: "red-series"},
				keyframe.Graph{Start: 4, End: 4, Shape: "square5", Palette: "red-series"},
				keyframe.Graph{Start: 5, End: 5, Shape: "square6", Palette: "red-series"},
			},
		},
	}
}

// dynamic38Square1 is cartoon 38 of the original Dynamic group.
func dynamic38Square1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 6, Loop: 10, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "square1", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
				keyframe.Graph{Start: 1, End: 1, Shape: "square2", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
				keyframe.Graph{Start: 2, End: 2, Shape: "square3", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
				keyframe.Graph{Start: 3, End: 3, Shape: "square4", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
				keyframe.Graph{Start: 4, End: 4, Shape: "square5", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
				keyframe.Graph{Start: 5, End: 5, Shape: "square6", Palette: "green-series", From: keyframe.Pose{Y: -1, Angle: 90}, To: keyframe.Pose{Y: -1, Angle: 90}},
			},
		},
	}
}

// dynamic39Square1 is cartoon 39 of the original Dynamic group.
func dynamic39Square1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 6, Loop: 10, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "square1", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 1, End: 1, Shape: "square2", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 2, End: 2, Shape: "square3", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 3, End: 3, Shape: "square4", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 4, End: 4, Shape: "square5", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
				keyframe.Graph{Start: 5, End: 5, Shape: "square6", Palette: "blue-series", From: keyframe.Pose{Angle: 90}, To: keyframe.Pose{Angle: 90}, Axis: render.AxisY},
			},
		},
	}
}

// dynamic40L is cartoon 40 of the original Dynamic group.
func dynamic40L() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 44, Loop: 1, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 15, Shape: "L", Palette: "blue", From: keyframe.Pose{X: 8}, To: keyframe.Pose{X: -8}, Mode: render.ModeMix},
				keyframe.Graph{Start: 21, End: 36, Shape: "E", Palette: "yellow", From: keyframe.Pose{X: 8}, To: keyframe.Pose{X: -8}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 7, End: 22, Shape: "L", Palette: "blue", From: keyframe.Pose{Z: -8, Angle: 90}, To: keyframe.Pose{Z: 8, Angle: 90}, Axis: render.AxisY, Mode: render.ModeMix},
				keyframe.Graph{Start: 28, End: 43, Shape: "E", Palette: "yellow", From: keyframe.Pose{Z: -8, Angle: 90}, To: keyframe.Pose{Z: 8, Angle: 90}, Axis: render.AxisY, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 7, End: 22, Shape: "O", Palette: "green", From: keyframe.Pose{X: 8}, To: keyframe.Pose{X: -8}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 14, End: 29, Shape: "O", Palette: "green", From: keyframe.Pose{Z: -8, Angle: 90}, To: keyframe.Pose{Z: 8, Angle: 90}, Axis: render.AxisY, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 14, End: 29, Shape: "V", Palette: "red", From: keyframe.Pose{X: 8}, To: keyframe.Pose{X: -8}, Mode: render.ModeMix},
			},
			{
				keyframe.Graph{Start: 21, End: 36, Shape: "V", Palette: "red", From: keyframe.Pose{Z: -8, Angle: 90}, To: keyframe.Pose{Z: 8, Angle: 90}, Axis: render.AxisY, Mode: render.ModeMix},
			},
		},
	}
}

// dynamic41Rain1 is cartoon 41 of the original Dynamic group.
func dynamic41Rain1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Frames: 16, Loop: 5, Delay: 100,
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 0, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: 7}, To: keyframe.Pose{Y: 7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 1, End: 1, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: 6}, To: keyframe.Pose{Y: 6}, Mode: render.ModeMix},
				keyframe.Graph{Start: 2, End: 2, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: 5}, To: keyframe.Pose{Y: 5}, Mode: render.ModeMix},
				keyframe.Graph{Start: 3, End: 3, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: 4}, To: keyframe.Pose{Y: 4}, Mode: render.ModeMix},
				keyframe.Graph{Start: 4, End: 4, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: 3}, To: keyframe.Pose{Y: 3}, Mode: render.ModeMix},
				keyframe.Graph{Start: 5, End: 5, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: 2}, To: keyframe.Pose{Y: 2}, Mode: render.ModeMix},
				keyframe.Graph{Start: 6, End: 6, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: 1}, To: keyframe.Pose{Y: 1}, Mode: render.ModeMix},
				keyframe.Graph{Start: 7, End: 7, Shape: "rain2", Palette: "red-series", Mode: render.ModeMix},
				keyframe.Graph{Start: 8, End: 8, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: -1}, To: keyframe.Pose{Y: -1}, Mode: render.ModeMix},
				keyframe.Graph{Start: 9, End: 9, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: -2}, To: keyframe.Pose{Y: -2}, Mode: render.ModeMix},
				keyframe.Graph{Start: 10, End: 10, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: -3}, To: keyframe.Pose{Y: -3}, Mode: render.ModeMix},
				keyframe.Graph{Start: 11, End: 11, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: -4}, To: keyframe.Pose{Y: -4}, Mode: render.ModeMix},
				keyframe.Graph{Start: 12, End: 12, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: -5}, To: keyframe.Pose{Y: -5}, Mode: render.ModeMix},
				keyframe.Graph{Start: 13, End: 13, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: -6}, To: keyframe.Pose{Y: -6}, Mode: render.ModeMix},
				keyframe.Graph{Start: 14, End: 14, Shape: "rain1", Palette: "red-series", From: keyframe.Pose{Y: -7}, To: keyframe.Pose{Y: -7}, Mode: render.ModeMix},
				keyframe.Graph{Start: 15, End: 15, Shape: "rain2", Palette: "red-series", From: keyframe.Pose{Y: -8}, To: keyframe.Pose{Y: -8}, Mode: render.ModeMix},
			},
		},
	}
}

var dynamicCartoons = []func() *keyframe.Cartoon{
	dynamic01X8x8,
	dynamic02Shape,
	dynamic03Waterfall1,
	dynamic04Strip,
	dynamic05Cubeamp1,
	dynamic06Cubeamp1,
	dynamic07Wave1,
	dynamic08II,
	dynamic09Arrows5,
	dynamic10Meteor,
	dynamic11Balls,
	dynamic12X4x4,
	dynamic13X4x4,
	dynamic14X8x1R,
	dynamic15X8x1,
	dynamic16X2,
	dynamic17XX,
	dynamic18Circle5,
	dynamic19Balls,
	dynamic20Balls,
	dynamic21Frame1,
	dynamic22Hui,
	dynamic23I,
	dynamic24Vertical8dot,
	dynamic25Vertical8dot,
	dynamic26X2,
	dynamic27X2,
	dynamic28Frame2,
	dynamic29Fly1,
	dynamic30Wave2,
	dynamic31Frame4,
	dynamic32TurnRight1,
	dynamic33TurnRight1,
	dynamic34TurnRight2,
	dynamic35II,
	dynamic36X8x1R,
	dynamic37Square1,
	dynamic38Square1,
	dynamic39Square1,
	dynamic40L,
	dynamic41Rain1,
}
