package component

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"

	"github.com/Tariomka/stm32-led-cube/internal/engine/audio"
)

// Microphone samples the sound sensor on PA0 (ADC1 channel 0) every millisecond, like the
// original firmware: ADC in continuous mode, TIM3 interrupt reads the latest conversion.
type Microphone struct {
	meter *audio.Meter
	timer *machine.TIM
	last  int
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

	return &Microphone{meter: audio.NewMeter(), timer: &machine.TIM3}
}

func (this *Microphone) Start() {
	this.timer.Configure(machine.PWMConfig{Period: 1_000_000}) // 1 ms
	this.timer.SetWraparoundInterrupt(this.sample)
}

func (this *Microphone) sample() {
	this.last = int(stm32.ADC1.DR.Get() & 0xfff)
	this.meter.Add(this.last)
}

// Level returns the loudness 0..63 for an audio cartoon's level mode.
func (this *Microphone) Level(mode int) int {
	mask := interrupt.Disable()
	level := this.meter.Level(mode)
	interrupt.Restore(mask)
	return level
}

// Raw returns the latest 12-bit sample.
func (this *Microphone) Raw() int {
	mask := interrupt.Disable()
	raw := this.last
	interrupt.Restore(mask)
	return raw
}
