package runner

import (
	"github.com/Tariomka/led-common-lib/pkg/led"
	"github.com/Tariomka/stm32-led-cube/internal/show"
)

type Size uint8

const (
	Size8 Size = iota
	Size16
	Size32
)

type LedType uint8

const (
	Monochrome LedType = iota
	RGB
)

type RunnerConfig struct {
	BaseSize   Size
	Height     Size
	LedType    LedType
	LightShows []led.LightShow
}

func NewConfig() RunnerConfig {
	return RunnerConfig{
		BaseSize:   Size8,
		Height:     Size8,
		LedType:    RGB,
		LightShows: show.NewLedShowList(),
	}
}
