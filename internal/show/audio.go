// Generated from the original iCubeSmart show data by tools/showgen.

package show

import "github.com/Tariomka/stm32-led-cube/internal/engine/keyframe"

// audio01AudioG1 is cartoon 1 of the original Audio group.
func audio01AudioG1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG1", Palette: "red-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio02AudioG2 is cartoon 2 of the original Audio group.
func audio02AudioG2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG2", Palette: "yellow-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio03AudioG1 is cartoon 3 of the original Audio group.
func audio03AudioG1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG1", Palette: "green-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio04AudioG2 is cartoon 4 of the original Audio group.
func audio04AudioG2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG2", Palette: "cyan-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio05AudioG1 is cartoon 5 of the original Audio group.
func audio05AudioG1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG1", Palette: "blue-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio06AudioG2 is cartoon 6 of the original Audio group.
func audio06AudioG2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG2", Palette: "red-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio07AudioG1 is cartoon 7 of the original Audio group.
func audio07AudioG1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG1", Palette: "yellow-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio08AudioG2 is cartoon 8 of the original Audio group.
func audio08AudioG2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG2", Palette: "green-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio09AudioG1 is cartoon 9 of the original Audio group.
func audio09AudioG1() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG1", Palette: "cyan-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// audio10AudioG2 is cartoon 10 of the original Audio group.
func audio10AudioG2() *keyframe.Cartoon {
	return &keyframe.Cartoon{
		Audio: true, Frames: 10, Loop: 0, Delay: 10, // 10 s, level mode 0
		Library: library,
		Elements: []keyframe.Element{
			{
				keyframe.Graph{Start: 0, End: 63, Shape: "AudioG2", Palette: "blue-series", From: keyframe.Pose{Y: -7}},
			},
		},
	}
}

// Audio lists the cartoons of the original Audio group, in order.
var Audio = []func() *keyframe.Cartoon{
	audio01AudioG1,
	audio02AudioG2,
	audio03AudioG1,
	audio04AudioG2,
	audio05AudioG1,
	audio06AudioG2,
	audio07AudioG1,
	audio08AudioG2,
	audio09AudioG1,
	audio10AudioG2,
}
