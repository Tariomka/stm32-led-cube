package component

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"
)

const (
	micWindow = 20   // samples, taken every millisecond
	micIdle   = 1050 // 12-bit ADC value of a silent microphone

	MaxAudioLevel = 63
	// AudioModeDecay lets the level fall by one step per read instead of dropping at once.
	AudioModeDecay = 2
)

// Microphone samples the sound sensor on PA0 (ADC1 channel 0) every millisecond and turns the
// samples into the 0..63 loudness of Audio cartoons, like the factory firmware: ADC in continuous
// mode, a TIM3 interrupt keeps the latest 20 samples and their peaks.
type Microphone struct {
	timer    *machine.TIM
	samples  [micWindow]int32
	pos      int
	max, min int32
	held     int
	last     int
}

func NewMicrophone() *Microphone {
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinInputModeAnalog})

	adc := stm32.ADC1
	stm32.RCC.SetAPB2ENR_ADC1EN(stm32.RCC_APB2ENR_ADC1EN_Enabled)
	adc.SetSMPR2_SMP0(5) // 55.5 cycles
	adc.SetSQR1_L(0)     // one conversion: channel 0
	adc.SetSQR3_SQ1(0)
	adc.CR2.Set(stm32.ADC_CR2_ADON | stm32.ADC_CR2_CONT | 7<<stm32.ADC_CR2_EXTSEL_Pos)

	adc.CR2.SetBits(stm32.ADC_CR2_RSTCAL)
	for adc.CR2.HasBits(stm32.ADC_CR2_RSTCAL) {
	}
	adc.CR2.SetBits(stm32.ADC_CR2_CAL)
	for adc.CR2.HasBits(stm32.ADC_CR2_CAL) {
	}
	adc.CR2.SetBits(stm32.ADC_CR2_EXTTRIG | stm32.ADC_CR2_SWSTART)

	microphone := &Microphone{timer: &machine.TIM3, max: micIdle, min: micIdle}
	for index := range microphone.samples {
		microphone.samples[index] = micIdle
	}
	microphone.timer.Configure(machine.PWMConfig{Period: 1_000_000}) // 1 ms
	microphone.timer.SetWraparoundInterrupt(microphone.sample)
	return microphone
}

func (this *Microphone) sample() {
	this.last = int(stm32.ADC1.DR.Get() & 0xfff)
	this.samples[this.pos] = int32(this.last)
	this.pos = (this.pos + 1) % micWindow
	this.max, this.min = 0, 99999
	for _, s := range this.samples {
		this.max = max(this.max, s)
		this.min = min(this.min, s)
	}
}

// Level returns the loudness 0..63 (how far the signal swings above or below idle) for an audio
// cartoon's level mode.
func (this *Microphone) Level(mode int) int {
	mask := interrupt.Disable()
	high, low := this.max, this.min
	interrupt.Restore(mask)

	up := 0
	switch {
	case high >= 2300:
		up = MaxAudioLevel
	case high > 1250:
		up = int((high-1250)*62/1050) + 1
	}
	down := 0
	switch {
	case low < -199:
		down = MaxAudioLevel
	case low < 850:
		down = int((850-low)*62/1050) + 1
	}
	level := max(up, down)

	if mode == AudioModeDecay {
		if level < this.held {
			this.held--
			return this.held
		}
		this.held = level
	}
	return level
}

// Raw returns the latest 12-bit sample.
func (this *Microphone) Raw() int {
	mask := interrupt.Disable()
	raw := this.last
	interrupt.Restore(mask)
	return raw
}
