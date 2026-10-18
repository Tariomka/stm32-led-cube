// Package audio turns microphone samples into the 0..63 level that selects the frame of Audio
// cartoons, like the original firmware (FUN_08007048 sampler, FUN_08032584 level).
package audio

const (
	window = 20   // samples, taken every millisecond
	idle   = 1050 // 12-bit ADC value of a silent microphone

	MaxLevel = 63

	// Mode 2 lets the level fall by one step per read instead of dropping at once.
	ModeDecay = 2
)

type Meter struct {
	samples  [window]int32
	pos      int
	max, min int32
	held     int
}

func NewMeter() *Meter {
	meter := &Meter{}
	for index := range meter.samples {
		meter.samples[index] = idle
	}
	meter.max, meter.min = idle, idle
	return meter
}

// Add stores a 12-bit sample and updates the peaks of the last 20 samples.
func (this *Meter) Add(sample int) {
	this.samples[this.pos] = int32(sample)
	this.pos = (this.pos + 1) % window
	this.max, this.min = 0, 99999
	for _, s := range this.samples {
		this.max = max(this.max, s)
		this.min = min(this.min, s)
	}
}

// Level returns the loudness 0..63 from how far the signal swings above or below idle.
func (this *Meter) Level(mode int) int {
	up := 0
	switch {
	case this.max >= 2300:
		up = MaxLevel
	case this.max > 1250:
		up = int((this.max-1250)*62/1050) + 1
	}
	down := 0
	switch {
	case this.min < -199:
		down = MaxLevel
	case this.min < 850:
		down = int((850-this.min)*62/1050) + 1
	}
	level := max(up, down)

	if mode == ModeDecay {
		if level < this.held {
			this.held--
			return this.held
		}

		this.held = level
	}

	return level
}
